package generated

// Do not edit, this file is automatically generated.

// Defaults: scaffolding used in 'build' command
var Defaults = map[string][]byte{
	"/.gitignore": []byte(`public
node_modules`),
	"/content/blog/_blueprint.json": []byte(`{
    "title": "text",
    "body": ["text"],
    "author": "text",
    "date": "date"
}`),
	"/content/blog/adding_pletiform.json": []byte(`{
    "title": "Build sites with good form",
    "body": [
        "Need an easy webform solution?",
        "Try adding a <a href='https://plentiform.com' target='blank' rel='noopener noreferrer'>plentiform</a>! (Coming soon)"
    ],
    "author": "Jim Fisk",
    "date": "1/26/2020"
}`),
	"/content/blog/post1.json": []byte(`{
    "title": "Post 1",
    "body": [
        "The first of the posts"
    ],
    "author": "Jim Fisk",
    "date": "1/24/2020"
}`),
	"/content/blog/post2.json": []byte(`{
    "title": "Post 2",
    "body": [
        "This is the second post."
    ],
    "author": "Jim Fisk",
    "date": "1/25/2020"
}`),
	"/content/index.json": []byte(`{
	"title": "My Plenti Site",
	"intro": {
		"slogan": "Visit the <a href=\"https://svelte.dev/tutorial\">Svelte tutorial</a> to learn how to build Svelte apps.",
		"color": "red"
	},
	"components": [
		{
			"component": "uses",
			"fields": {"type": "index"}
		}
	]
}`),
	"/content/pages/_blueprint.json": []byte(`{
	"title": "text",
	"description": ["text"],
	"author": "text"
}`),
	"/content/pages/about.json": []byte(`{
	"title": "About Plenti",
	"description": [
		"Plenti is <a href='https://jamstack.org/' target='blank' rel='noopener noreferrer'>JAMstack</a> framework with a modern frontend for creating dynamic experiences. We've cut out as many dependencies as possible so you can focus on being productive instead of wrestling with a complicated toolchain.",
		"The <a href='https://svelte.dev/' target='blank' rel='noopener noreferrer'>Svelte</a> frontend <em>cuts weight</em> so users get a snappy experience, even with bad internet connections or underpowered devices.",
		"The <a href='https://golang.org/' target='blank' rel='noopener noreferrer'>Go</a> backend <em>cuts wait</em> so apps build faster allowing devs to get more done and editors to get realtime feedback on content changes.",
		"Thanks for taking a look!"
	],
	"author": "Jim Fisk"
}`),
	"/content/pages/contact.json": []byte(`{
	"title": "Contact",
	"description": [
		"The project is 100% open source, so if you'd like to fork it for your own purposes, or help us out by reporting bugs / contributing code: <a href=\"https://github.com/plentico/plenti\">https://github.com/plentico/plenti</a>"
	],
	"author": "Jim Fisk"
}`),
	"/layout/components/grid.svelte": []byte(`<script>
  import { sortByDate } from '../scripts/sort_by_date.svelte';
  export let items, filter;
</script>

<div class="grid">
  {#each sortByDate(items) as item}
		{#if item.type == filter}
      <a class="grid-item" href="{item.path}">{item.fields.title}</a>
		{/if}
  {/each}
</div>

<style>
  .grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    grid-column-gap: 10px;
    grid-row-gap: 10px;
  }
  .grid-item {
    display: flex;
    flex-grow: 1;
    height: 200px;
    background: var(--base);
    align-items: center;
    justify-content: center;
  }
</style>
`),
	"/layout/components/uses.svelte": []byte(`<script>
  export let type;

  let path;
  let copyText = "Copy";
  const copy = async () => {
    if (!navigator.clipboard) {
      return
    }
    try {
      copyText = "Copied";
      await navigator.clipboard.writeText(path.innerHTML);
      setTimeout(() => copyText = "Copy", 800);
    } catch (err) {
      console.error('Failed to copy!', err)
    }
  }
</script>

<details>
  <summary>Uses the "{type}" template</summary>
  <pre>
    <code bind:this={path}>layout/content/{type}.svelte</code><button on:click={copy}>{copyText}</button>
  </pre>
</details>

<style>
  summary {
      cursor: pointer;
  }
  code {
      background-color: var(--base);
      padding: 5px 10px;
  }
  button {
    border: 1px solid rgba(0,0,0,.1);
    background: white;
    padding: 4px;
    border-top-right-radius: 5px;
    border-bottom-right-radius: 5px;
    cursor: pointer;
  }
</style>`),
	"/layout/content/404.svelte": []byte(`<h1>Oops... 404 not found</h1>
<a href="/">Go home?</a>`),
	"/layout/content/blog.svelte": []byte(`<script>
	export let title, body, author, date;
  import Uses from "../components/uses.svelte";
</script>

<h1>{title}</h1>

<p><em>{#if author}Written by {author}{/if}{#if date}&nbsp;on {date}{/if}</em></p>

<div>
  {#each body as paragraph}
    <p>{@html paragraph}</p>
  {/each}
</div>

<Uses type="blog" />

<p><a href="/">Back home</a></p>
`),
	"/layout/content/index.svelte": []byte(`<script>
	export let title, intro, components, allNodes;
	import Grid from '../components/grid.svelte';
	import { loadComponent } from '../scripts/load_component.svelte';
</script>

<h1>{title}</h1>

<section id="intro">
	<p>{@html intro.slogan}</p>
</section>

<div>
	<h3>Recent blog posts:</h3>
	<Grid items={allNodes} filter="blog" />
	<br />
</div>

{#if components}
	{#each components as { component, fields }}
		{#await loadComponent(component)}
		{:then compClass}
			<svelte:component this="{compClass}" {...fields} />
		{:catch error}
		{/await}
	{/each}
{/if}`),
	"/layout/content/pages.svelte": []byte(`<script>
  export let title, description;
  import Uses from "../components/uses.svelte";
</script>

<h1>{title}</h1>

<div>
  {#each description as paragraph}
    <p>{@html paragraph}</p>
  {/each}
</div>

<Uses type="pages" />

<p><a href="/">Back home</a></p>`),
	"/layout/global/footer.svelte": []byte(`<script>
  export let allNodes;
  import { makeTitle } from '../scripts/make_title.svelte';
</script>
<footer>
  <div class="container">
    <span>All nodes:</span>
    {#each allNodes as node}
      <a href="{node.path}">{makeTitle(node.filename)}</a>&nbsp;
    {/each}
  </div>
</footer>

<style>
  footer {
    min-height: 200px;
    display: flex;
    align-items: center;
    background-color: var(--base);
    margin-top: 100px;
  }
</style>
`),
	"/layout/global/head.svelte": []byte(`<script>
  export let title;
</script>

<head>
  <meta charset='utf-8'>
  <meta name='viewport' content='width=device-width,initial-scale=1'>

  <title>{ title }</title>

  <link href="https://fonts.googleapis.com/css2?family=Rubik:ital,wght@0,300;0,700;1,300&display=swap" rel="stylesheet">
  <link rel='icon' type='image/png' href='/favicon.png'>
  <link rel='stylesheet' href='/spa/bundle.css'>
</head>
`),
	"/layout/global/html.svelte": []byte(`<script>
  import Head from './head.svelte';
  import Nav from './nav.svelte';
  import Footer from './footer.svelte';
  import { makeTitle } from '../scripts/make_title.svelte';

  export let route, node, allNodes;
</script>

<html lang="en">
<Head title={makeTitle(node.filename)} />
<body>
  <Nav />
  <main>
    <div class="container">
      <svelte:component this={route} {...node.fields} {allNodes} />
      <br />
    </div>
  </main>
  <Footer {allNodes} />
</body>
</html>

<style>
  body {
    font-family: 'Rubik', sans-serif;
    display: flex;
    flex-direction: column;
    margin: 0;
  }
  main {
    flex-grow: 1;
  }
  :global(.container) {
    max-width: 1024px;
    margin: 0 auto;
    flex-grow: 1;
    padding: 0 20px;
  }
  :global(:root) {
    --primary: (50, 50, 50);
    --accent: rgb(1, 1, 1);
    --base: rgb(245, 245, 245);
  }
</style>
`),
	"/layout/global/nav.svelte": []byte(`<nav>
  <div class="container">
    <span id="brand"><a href="/">Home</a></span>
    <a href="/about">About</a>&nbsp;
    <a href="/contact">Contact</a>
  </div>
</nav>

<style>
  nav {
    min-height: 60px;
    display: flex;
    align-items: center;
    box-shadow: 0px 2px 3px var(--base);
  }
  .container {
    display: flex;
  }
  #brand {
    flex: 1;
  }
</style>
`),
	"/layout/scripts/load_component.svelte": []byte(`<script context="module">
  export const loadComponent = component => {
    let compClassPromise = import("../components/" + component + ".svelte").then(res => res.default);
    // Fix "Unhandled promise rejection" error.
    // See: https://github.com/sveltejs/sapper/issues/487#issuecomment-529145749
    $: compClassPromise.catch(err => null)
    return compClassPromise;
  }
</script>
`),
	"/layout/scripts/make_title.svelte": []byte(`<script context="module">
  export const makeTitle = filename => {
  if (filename == 'index.json') {
    return 'Home';
  } else if (filename) {
    // Remove file extension.
    filename = filename.split('.').slice(0, -1).join('.');
    // Convert underscores and hyphens to spaces.
    filename = filename.replace(/_|-/g,' ');
    // Capitalize first letter of each word.
    filename = filename.split(' ').map((s) => s.charAt(0).toUpperCase() + s.substring(1)).join(' ');
  } 
  return filename;
  }
</script>
`),
	"/layout/scripts/sort_by_date.svelte": []byte(`<script context="module">
  export const sortByDate = (items, order) => {
    items.sort((a, b) => { 
      // Must have a field specifically named "date" to work.
      // Feel free to extend to other custom named date fields.
      if (a.fields.hasOwnProperty("date") && b.fields.hasOwnProperty("date")) {
        let aDate = new Date(a.fields.date);
        let bDate = new Date(b.fields.date);
        if (order == "oldest") {
            return aDate - bDate;
        }
        return bDate - aDate;
      }
    });
    return items;
  };
</script>`),
	"/package.json": []byte(`{
  "name": "my-plenti-app",
  "version": "1.0.0",
  "type": "module",
  "private": true,
  "devDependencies": {
    "require-relative": "^0.8.7"
  },
  "dependencies": {
    "navaid": "^1.0.5",
    "regexparam": "^1.3.0",
    "svelte": "^3.21.0"
  }
}
`),
	"/plenti.json": []byte(`{
	"baseurl": "http://example.org/",
	"title": "My New Plenti Site",
	"types": {
		"pages": "/:filename"
	},
	"build": "public",
	"local": {
		"port": 3000
	}
}`),
}
