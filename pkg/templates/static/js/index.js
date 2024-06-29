const activeLine = "font-bold"
const activeSection = "font-bold"
const scrollSettings = {behavior:"instant", block:"center", inline:"center"}

let line = document.createElement("div");
let section = document.createElement("div");

const eventSource = new EventSource('/events');

// we shell receive a json object with two fields
// section and line. Which refer to current section id and approximate
// current line number.
eventSource.addEventListener("line", (event) => {
    const data = event.data
    try {
    const lineID = `line-${data}`
    const newLine = document.getElementById(lineID)



    if (newLine && newLine != line) {
        line.classList = ""
        line = newLine
        line.classList = activeLine
        line.scrollIntoView(scrollSettings)
        if (scrollSettings.behavior !== "smooth") scrollSettings.behavior = "smooth"
    }

    }catch (e) {
        console.error(e)
    }
});

eventSource.addEventListener("update", (event) => {
    console.log(event.data)
    const toast = document.getElementById("toast")


    if(!toast){
        console.log("toast not found")
        return
    }
    try {
        toast.innerHTML = "<div class=\"alert alert-warning\"><span>Project file was updated. Please, refresh the page!</span></div>"
    }catch (e) {
        console.error(e)
    }
});


eventSource.addEventListener("alert", (event) => {
    console.log(event.data)
    const toast = document.getElementById("toast")
    const data = JSON.parse(event.data)


    if(!toast){
        console.log("toast not found")
        return
    }
    try {
        toast.innerHTML = `<div class="alert alert-warning">
                <span>${data.message}</span>
        </div>`
    }catch (e) {
        console.error(e)
    }
});

eventSource.addEventListener("cowntdown", (event) => {
    console.log(event.data)
    const min = document.getElementById("cowntdown-min")
    const sec = document.getElementById("cowntdown-sec")
    const title = document.getElementById("countdown-title")


    try {
        const data = JSON.parse(event.data)
        title.innerHTML = `${data.name} starts in `
        min.style=`--value:${data.min};`
        sec.style=`--value:${data.sec};`
    }catch (e) {
        console.error(e)
    }
});
