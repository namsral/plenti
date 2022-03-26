export const isDate = date => (new Date(date) !== "Invalid Date") && !isNaN(new Date(date));
export const makeDate = date => new Date(date).toISOString().split('T')[0];
export const formatDate = (date, format) => {
    let parts = date.split('-');
    let year = parts[0];
    let month = parts[1];
    let day = parts[2];

    let days = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];
    let daysFull = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];
    let months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    let monthsFull = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];

    /**
     * 6/7/2008
     * 6/7/08
     * 6-7-2008
     * 6-7-08
     * 6.7.2008
     * 6.7.08
     */
    let re = new RegExp("^(0?[1-9]|1[0-2])(/|-|\.)(0?[1-9]|1[0-9]|2[0-9]|3[0-1])(/|-|\.)(([0-9][0-9])?[0-9][0-9])$");
    if (re.test(format)) {
        let replacements = format.match(re);
        let delimeter = replacements[2];
        let yearFull = replacements[6];
        year = yearFull !== undefined ? year : year.slice(-2);
        return Number(month) + delimeter + Number(day) + delimeter + year;
    }
    /**
     * Saturday, June 7, 2008
     * Saturday, Jun 7, 2008
     * Sat, June 7, 2008
     * Sat, Jun 7, 2008
     * 
     * (with or without commas and case insensitive)
     */
    re = new RegExp("^\\b((Monday|Mon)|(Tuesday|Tue)|(Wednesday|Wed)|(Thursday|Thu)|(Friday|Fri)|(Saturday|Sat)|(Sunday|Sun))\\b(,?) \\b((January|Jan)|(February|Feb)|(March|Mar)|(April|Apr)|May|(June|Jun)|(July|Jul)|(August|Aug)|(September|Sep)|(October|Oct)|(November|Nov)|(December|Dec))\\b (0?[1-9]|1[0-9]|2[0-9]|3[0-1])(,?) ([0-9][0-9][0-9][0-9])$", "i");
    if (re.test(format)) {
        let replacements = format.match(re);
        let delimeterDayOfWeek = replacements[9] !== undefined ? replacements[9] : '';
        let dayPos = new Date(date).getDay();
        let dayOfWeek = replacements[1].length > 3 ? daysFull[dayPos] : days[dayPos];
        let monthName = replacements[10].length > 3 ? monthsFull[month -1] : months[month -1];
        let delimeterDay = replacements[23] !== undefined ? replacements[23] : '';
        return dayOfWeek + delimeterDayOfWeek + ' ' + monthName + ' ' + Number(day) + delimeterDay + ' ' + year;
    }
    // Can't find format
    return date;
}