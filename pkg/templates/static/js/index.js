const activeLine = "text-2xl font-bold"
const activeSection = "font-bold"
const scrollSettings = {behavior:"instant", block:"center", inline:"center"}

let line = document.createElement("div");
let section = document.createElement("div");

const eventSource = new EventSource('/events');

// we shell receive a json object with two fields
// section and line. Which refer to current section id and approximate
// current line number.
eventSource.addEventListener("cursor", (event) => {
    console.log(event.data)
    const data = JSON.parse(event.data)
    try {
    const lineID = `${data.section}-${data.line}`
    const newLine = document.getElementById(lineID)
    const newSection = document.getElementById(data.section)

    if (newSection && newSection != section) {
        section.classList.remove(activeSection)
        section = newSection
        section.classList.add(activeSection)
    }

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
        toast.append("<div class=\"alert alert-warning\"><span>Project file was updated. Please, refresh the page!</span></div>")
    }catch (e) {
        console.error(e)
    }
});
