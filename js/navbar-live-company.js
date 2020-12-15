function callNavbarLiveCompanyJs() {
    console.log("Starting processing data for navbar-live-company")
    displayLiveProductivityData("company")
    displayCalendar("company")
    displayOverview("company")
    displayBestWorkplaces("company")
    displayWorstWorkplaces("company")
    displayPowerOffWorkplaces("company")
    displayWorkplaceOverview("company")
}

function displayLiveProductivityData(input) {
    console.log("displaying live productivity data for " + input)
    let data = {
        input: input,
    };
    fetch("/get_live_productivity_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            console.log(result)
            displayData(result, "navbar-live-company-1-data-day", result.Today, result.Yesterday);
            displayData(result, "navbar-live-company-1-data-week", result.ThisWeek, result.LastWeek);
            displayData(result, "navbar-live-company-1-data-month", result.ThisMonth, result.LastMonth);
        });
    }).catch((error) => {
        console.log(error)
    });
}

function displayData(result, elementId, actual, last) {
    const data = document.getElementById(elementId)
    let lastPercent = document.createElement("div");
    lastPercent.textContent = last + "%";
    lastPercent.style.width = "55px";
    lastPercent.style.textAlign = "right";
    lastPercent.style.paddingRight = "5px";
    lastPercent.style.fontWeight = "bold"
    data.appendChild(lastPercent)
    let lastColor = document.createElement("div");
    lastColor.style.width = "0.9em"
    lastColor.style.height = "0.9em"
    lastColor.style.background = "rgba(0,128,9," + (+last) / 100
    data.appendChild(lastColor)
    let arrow = document.createElement("div");
    arrow.style.width = "45px"
    arrow.style.textAlign = "center"
    if (Math.abs(+last - +actual) <= 5) {
        arrow.textContent = "→"
    } else if (+last < +actual) {
        arrow.textContent = "↑"
    } else if (+last > +actual) {
        arrow.textContent = "↓"
    }
    data.appendChild(arrow)
    let actualColor = document.createElement("div");
    actualColor.style.width = "0.9em"
    actualColor.style.height = "0.9em"
    actualColor.style.background = "rgba(0,128,9," + (+actual) / 100
    data.appendChild(actualColor)
    let actualPercent = document.createElement("div");
    actualPercent.textContent = actual + "%";
    actualPercent.style.width = "55px";
    actualPercent.style.textAlign = "left";
    actualPercent.style.paddingLeft = "5px";
    actualPercent.style.fontWeight = "bold"
    data.appendChild(actualPercent)
}

function displayCalendar(input) {
    console.log("displaying calendar for " + input)
}

function displayOverview(input) {
    console.log("displaying overview for " + input)
}

function displayBestWorkplaces(input) {
    console.log("displaying best workplaces for " + input)
}

function displayWorstWorkplaces(input) {
    console.log("displaying worst workplaces for " + input)
}

function displayPowerOffWorkplaces(input) {
    console.log("displaying power off workplaces for " + input)
}

function displayWorkplaceOverview(input) {
    console.log("displaying workplace overview for " + input)
}