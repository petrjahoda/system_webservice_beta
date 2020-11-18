const navbarLive = document.getElementById("navbar-live");
const navbarCharts = document.getElementById("navbar-charts");
const navbarStatistics = document.getElementById("navbar-statistics");
const navbarData = document.getElementById("navbar-data");
const navbarSettings = document.getElementById("navbar-settings");


navbarLive.addEventListener("click", () => {
    console.log("live clicked")
    navbarLive.classList.add("underlined")
    navbarCharts.classList.remove("underlined")
    navbarStatistics.classList.remove("underlined")
    navbarData.classList.remove("underlined")
    navbarSettings.classList.remove("underlined")

})

navbarCharts.addEventListener("click", () => {
    console.log("charts clicked")
    navbarLive.classList.remove("underlined")
    navbarCharts.classList.add("underlined")
    navbarStatistics.classList.remove("underlined")
    navbarData.classList.remove("underlined")
    navbarSettings.classList.remove("underlined")
})

navbarStatistics.addEventListener("click", () => {
    console.log("statistics clicked")
    navbarLive.classList.remove("underlined")
    navbarCharts.classList.remove("underlined")
    navbarStatistics.classList.add("underlined")
    navbarData.classList.remove("underlined")
    navbarSettings.classList.remove("underlined")
})

navbarData.addEventListener("click", () => {
    console.log("data clicked")
    navbarLive.classList.remove("underlined")
    navbarCharts.classList.remove("underlined")
    navbarStatistics.classList.remove("underlined")
    navbarData.classList.add("underlined")
    navbarSettings.classList.remove("underlined")
})

navbarSettings.addEventListener("click", () => {
    console.log("settings clicked")
    navbarLive.classList.remove("underlined")
    navbarCharts.classList.remove("underlined")
    navbarStatistics.classList.remove("underlined")
    navbarData.classList.remove("underlined")
    navbarSettings.classList.add("underlined")
})