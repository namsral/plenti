package build

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"rogchap.com/v8go"
)

// SSRctx is a v8go context for loaded with components needed to render HTML.
var SSRctx *v8go.Context

// Client builds the SPA.
func Client(buildPath string) {

	defer Benchmark(time.Now(), "Compiling client SPA with Svelte")

	Log("\nCompiling client SPA with svelte")

	stylePath := buildPath + "/spa/bundle.css"

	// Set up counter for logging output.
	compiledComponentCounter := 0

	// Get svelte compiler code from node_modules.
	compiler, err := ioutil.ReadFile("node_modules/svelte/compiler.js")
	if err != nil {
		fmt.Printf("Can't read node_modules/svelte/compiler.js: %v", err)
	}
	// Remove reference to 'self' that breaks v8go on line 19 of node_modules/svelte/compiler.js.
	compilerStr := strings.Replace(string(compiler), "self.performance.now();", "'';", 1)
	ctx, _ := v8go.NewContext(nil)
	_, addCompilerErr := ctx.RunScript(compilerStr, "compile_svelte")
	if addCompilerErr != nil {
		fmt.Printf("Could not add svelte compiler: %v\n", addCompilerErr)
	}

	// Use v8go and add create_ssr_component() function.
	createSsrComponent, npmReadErr := ioutil.ReadFile("node_modules/svelte/internal/index.js")
	if npmReadErr != nil {
		fmt.Printf("Can't read node_modules/svelte/internal/index.js: %v", npmReadErr)
	}
	// Fix "Cannot access 'on_destroy' before initialization" errors on line 1320 & line 1337 of node_modules/svelte/internal/index.js.
	createSsrStr := strings.ReplaceAll(string(createSsrComponent), "function create_ssr_component(fn) {", "function create_ssr_component(fn) {var on_destroy= {};")
	SSRctx, _ = v8go.NewContext(nil)
	_, createFuncErr := SSRctx.RunScript(createSsrStr, "create_ssr")
	if err != nil {
		fmt.Printf("Could not add create_ssr_component() func from svelte/internal: %v", createFuncErr)
	}
	// Fix "ReferenceError: exports is not defined" errors on line 1319 (exports.current_component;).
	SSRctx.RunScript("var exports = {};", "create_ssr")

	compileSvelte(ctx, SSRctx, "ejected/router.svelte", buildPath+"/spa/ejected/router.js", stylePath)

	// Go through all file paths in the "/layout" folder.
	layoutFilesErr := filepath.Walk("layout", func(layoutPath string, layoutFileInfo os.FileInfo, err error) error {
		// Create destination path.
		destFile := buildPath + strings.Replace(layoutPath, "layout", "/spa", 1)
		// Make sure path is a directory
		if layoutFileInfo.IsDir() {
			// Create any sub directories need for filepath.
			os.MkdirAll(destFile, os.ModePerm)
		} else {
			// If the file is in .svelte format, compile it to .js
			if filepath.Ext(layoutPath) == ".svelte" {

				// Replace .svelte file extension with .js.
				destFile = strings.TrimSuffix(destFile, filepath.Ext(destFile)) + ".js"

				compileSvelte(ctx, SSRctx, layoutPath, destFile, stylePath)

				compiledComponentCounter++

			}
		}
		return nil
	})
	if layoutFilesErr != nil {
		fmt.Printf("Could not get layout file: %s", layoutFilesErr)
	}

	Log("Number of components compiled: " + strconv.Itoa(compiledComponentCounter))

}

func compileSvelte(ctx *v8go.Context, SSRctx *v8go.Context, layoutPath string, destFile string, stylePath string) {

	component, err := ioutil.ReadFile(layoutPath)
	if err != nil {
		fmt.Printf("Can't read component: %v", err)
	}
	componentStr := string(component)

	// Compile component with Svelte.
	ctx.RunScript("var { js, css } = svelte.compile(`"+componentStr+"`, {css: false});", "compile_svelte")

	// Get the JS code from the compiled result.
	jsCode, err := ctx.RunScript("js.code;", "compile_svelte")
	if err != nil {
		fmt.Printf("V8go could not execute js.code: %v", err)
	}
	jsBytes := []byte(jsCode.String())
	jsWriteErr := ioutil.WriteFile(destFile, jsBytes, 0755)
	if jsWriteErr != nil {
		fmt.Printf("Unable to write file: %v", jsWriteErr)
	}

	// Get the CSS code from the compiled result.
	cssCode, err := ctx.RunScript("css.code;", "compile_svelte")
	if err != nil {
		fmt.Printf("V8go could not execute css.code: %v", err)
	}
	cssStr := strings.TrimSpace(cssCode.String())
	// If there is CSS, write it into the bundle.css file.
	if cssStr != "null" {
		cssFile, WriteStyleErr := os.OpenFile(stylePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if WriteStyleErr != nil {
			fmt.Printf("Could not open bundle.css for writing: %s", WriteStyleErr)
		}
		defer cssFile.Close()
		if _, err := cssFile.WriteString(cssStr); err != nil {
			log.Println(err)
		}
	}

	// Get Server Side Rendered (SSR) JS.
	ctx.RunScript("var { js: ssrJs, css: ssrCss } = svelte.compile(`"+componentStr+"`, {generate: 'ssr'});", "compile_svelte")
	ssrJsCode, err := ctx.RunScript("ssrJs.code;", "compile_svelte")
	if err != nil {
		fmt.Printf("V8go could not execute ssrJs.code: %v", err)
	}
	// Regex match static import statements.
	reStaticImport := regexp.MustCompile(`import\s((.*)\sfrom(.*);|(((.*)\n){0,})\}\sfrom(.*);)`)
	// Regex match static export statements.
	reStaticExport := regexp.MustCompile(`export\s(.*);`)
	// Remove static import statements.
	ssrStr := reStaticImport.ReplaceAllString(ssrJsCode.String(), `/*$0*/`)
	reStaticImportPath := regexp.MustCompile(`(?:'|").*(?:'|")`)
	namedImports := reStaticImport.FindAllString(ssrStr, -1)
	for _, namedImport := range namedImports {
		// Get path only from static import statement.
		importPath := reStaticImportPath.FindString(namedImport)
		// Remove quotes around path.
		importPath = strings.Trim(importPath, `'"`)
		// Get individual path arguments.
		layoutParts := strings.Split(layoutPath, "/")
		layoutFileName := layoutParts[len(layoutParts)-1]
		layoutRootPath := strings.TrimSuffix(layoutPath, layoutFileName)

		importParts := strings.Split(importPath, "/")
		importSignature := ""
		// Check if the path ends with a file extension, e.g. ".svelte".
		if len(filepath.Ext(importParts[len(importParts)-1])) > 1 {
			for _, importPart := range importParts {
				if importPart == "." {
					importPath = strings.TrimPrefix(importPath, "./")
				}
				if importPart == ".." {
					layoutParts = strings.Split(layoutRootPath, "/")
					importPath = strings.TrimPrefix(importPath, importPart+"/")
					layoutRootPath = strings.TrimSuffix(layoutRootPath, layoutParts[len(layoutParts)-2]+"/")
				}
			}
			importSignature = strings.ReplaceAll(strings.ReplaceAll((layoutRootPath+importPath), "/", "_"), ".", "_")
		} else {
			// No file ext was found, e.g. "svelte/internal".
			importSignature = strings.ReplaceAll(strings.ReplaceAll(importPath, "/", "_"), ".", "_")
		}
		fmt.Println(importSignature)
		//fmt.Println(layoutPath)
		//fmt.Println(importPath)
	}
	// Remove static export statements.
	ssrStr = reStaticExport.ReplaceAllString(ssrStr, `/*$0*/`)
	// Use var instead of const so it can be redeclared multiple times.
	ssrStr = strings.ReplaceAll(ssrStr, "const", "var")
	// Create custom variable name for component based on the file path for the layout.
	componentSignature := strings.ReplaceAll(strings.ReplaceAll(layoutPath, "/", "_"), ".", "_")

	// TODO: Need to account for imports using name not based on layout filename,
	// e.g. "Uses" instead of "Template" - for now must manually change in project.
	ssrStr = strings.ReplaceAll(ssrStr, "Component", componentSignature)
	// Add component to context so it can be used to render HTML in data_source.go.
	_, addSSRCompErr := SSRctx.RunScript(ssrStr, "create_ssr")
	if addSSRCompErr != nil {
		fmt.Printf("Could not add SSR Component: %v\n", addSSRCompErr)
	}

}
