const activeLine = "font-bold"
const activeSection = ["font-bold", "animate-pulse"]
const scrollSettings = {behavior:"instant", block:"center", inline:"center"}
const navbar = document.getElementById("navbar")

let line = document.createElement("div");
let e = document.createElement("div");
let section = document.createElement("div");
let scroll = true

const scrollToggle = document.getElementById("scroll")

try {
    scrollToggle.addEventListener("click", (e)=>{
        scroll = !scroll
    })
} catch (e) {
    console.log(e)
}


section.innerHTML = "<div></div>"


const eventSource = new EventSource('/events');
eventSource.onopen = (e) => {
    navbar.classList.remove("bg-error")
};

eventSource.onerror = (err) => {
    navbar.classList.add("bg-error")
};

const playbackStopped = (element) => {
    navbar.classList.remove("bg-success")
}


var lastEvent = Date.now()
var timeoutID = 0

eventSource.addEventListener("line", (event) => {
    if (!scroll) {
        return 
    }

    if ((Date.now() - lastEvent) < 200 ){
        clearTimeout(timeoutID)
        timeoutID = setTimeout(playbackStopped, 200)
        if (!navbar.classList.contains("bg-success")) {
            navbar.classList.add("bg-success")
        }
    }
    lastEvent = Date.now()


    const lineID = `line-${event.data}`
    const newLine = document.getElementById(lineID)

    try {

    if (newLine != null && newLine != line) {
        if (line.classList.contains(activeLine)) {
                line.classList.remove(activeLine)
            }
        line = newLine
        line.classList.add(activeLine)
        line.scrollIntoView(scrollSettings)
        if (scrollSettings.behavior !== "smooth") scrollSettings.behavior = "smooth"
    } else { return }

    if (newLine && newLine.parentElement && newLine.parentElement != section){
            section.firstElementChild.classList.remove(...activeSection)
            newLine.parentElement.firstElementChild.classList.add(...activeSection)
            section.classList.toggle("active")
            section = newLine.parentElement
            section.classList.toggle("active")
    }

    }catch (e) {
        console.log(e)
    }
});

eventSource.addEventListener("map", (event) => {
    const eID = `map-${event.data}`
    const newE = document.getElementById(eID)

    try {

    if (newE && newE != e) {
        e = newE
        line.scrollIntoView(scrollSettings)
        if (scrollSettings.behavior !== "smooth") scrollSettings.behavior = "smooth"
    }


    }catch (e) {
        console.log(e)
    }
});

eventSource.addEventListener("update", (event) => {
    console.log(event.data)
    window.location.reload()
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
