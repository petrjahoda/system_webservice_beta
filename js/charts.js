const liveSelect = document.getElementById("live");
const chartsSelect = document.getElementById("charts");
const statisticsSelect = document.getElementById("statistics");
const dataSelect = document.getElementById("data");
const settingsSelect = document.getElementById("settings");

const user = sessionStorage.getItem("user")
if (user === null) {
    window.location.replace('login');
}

liveSelect.addEventListener("click", () => {
    console.log("live")
    window.location.replace('live');
})
chartsSelect.addEventListener("click", () => {
    console.log("charts")
    window.location.replace('charts');
})
statisticsSelect.addEventListener("click", () => {
    console.log("statistics")
    window.location.replace('statistics');
})
dataSelect.addEventListener("click", () => {
    console.log("data")
    window.location.replace('data');
})
settingsSelect.addEventListener("click", () => {
    console.log("settings")
    window.location.replace('settings');
})