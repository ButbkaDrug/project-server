
window.onload = (event) => {
const tag = "active-line"
const activeLine = [tag, "text-success" ]
const activeSection = ["font-bold", "animate-pulse"]
const scrollSettings = {behavior:"instant", block:"center", inline:"center"}
const navbar = document.getElementById("navbar")
const body = document.getElementById("body")


let line = document.createElement("div");
let e = document.createElement("div");
let section = document.createElement("div");
let scroll = true
let playing = false

const alert = (message, duration) => {
    const toast = document.getElementById("toast")
    const d = duration || 2000

    if(!toast){
        console.log("toast not found")
        return
    }
    try {

    setTimeout(function() { toast.firstChild.classList.add("hidden") }, d);

    toast.innerHTML = `<div id="message" class="alert alert-info animate-bounce transition duration-300 ease-in-out">
        <span>${message}</span>
    </div>`
    }catch (e) {
        console.error(e)
    }
};


section.innerHTML = "<div></div>"


const eventSource = new EventSource('/events');
eventSource.onopen = (e) => {
    body.classList.remove("border-red-500")
};

eventSource.onerror = (err) => {
    body.classList.add("border-red-500")
};

const playbackStopped = (element) => {
    body.classList.remove("border-green-500")
    if (playing === true) {
        alert("connection is broken?", 3000)
    }
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
        if (!body.classList.contains("border-green-500")) {
            body.classList.add("border-green-500")
        }
    }
    lastEvent = Date.now()


    const lineID = `line-${event.data}`
    const newLine = document.getElementById(lineID)

    try {

    if (newLine != null && newLine != line) {
        if (line.classList.contains(tag)) {
                line.classList.remove(...activeLine)
            }
        line = newLine
        line.classList.add(...activeLine)
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


eventSource.addEventListener("play", (event) => {
    const message = event.data

        alert(message, 1000)
        playing = true
});

eventSource.addEventListener("stop", (event) => {
    const message = event.data

        alert(message, 1000)
        playing = false
});

window.addEventListener('keydown', function (e) {
    switch (e.key) {
    case "Escape":
        try {
                    console.log("escape activated")
            checkbox = document.getElementById("my-drawer-3")
            checkbox.checked = false
        } catch (e) {
            console.error(e)
        }
            case "1":
            console.log("one have been pressed")
            default:
            console.log(e.key)

    }
});

};

