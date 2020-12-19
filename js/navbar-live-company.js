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
            displayData(result, "navbar-live-company-1-data-day", result["Today"], result["Yesterday"]);
            displayData(result, "navbar-live-company-1-data-week", result["ThisWeek"], result["LastWeek"]);
            displayData(result, "navbar-live-company-1-data-month", result["ThisMonth"], result["LastMonth"]);
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
    lastColor.style.background = "rgba(137,171,15," + (+last) / 100
    lastColor.style.border = "0.2px solid black"
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
    actualColor.style.background = "rgba(137,171,15," + (+actual) / 100
    actualColor.style.border = "0.2px solid black"
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
    let data = {
        input: input,
    };
    fetch("/get_calendar_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            drawCalendar(result["Data"]);
        });
    }).catch((error) => {
        console.log(error)
    });
}

function drawCalendar(data) {
    let thisYear = new Date().getFullYear()
    const width = 960, height = 136, cellSize = 17;
    const color = d3.scaleQuantize()
        .domain([0, 100])
        .range(["#f3f6e7", "#e7eecf", "#dbe5b7", "#d0dd9f", "#c4d587", "#b8cd6f", "#acc457", "#a1bc3f", "#94b327", "#89ab0f"]);
    const svg = d3.select("#navbar-live-company-2-calendar-chart")
        .selectAll("svg")
        .data(d3.range(thisYear - 1, thisYear + 1))
        .enter().append("svg")
        .attr("width", width)
        .attr("height", height)
        .append("g")
        .attr("transform", "translate(" + ((width - cellSize * 53) / 2) + "," + (height - cellSize * 7 - 1) + ")");
    svg.append("text")
        .attr("transform", "translate(-6," + cellSize * 3.5 + ")rotate(-90)")
        .attr("font-family", "sans-serif")
        .attr("font-size", "1em")
        .attr("text-anchor", "middle")
        .text(function (d) {
            return d;
        });
    svg.append("g")
        .attr("fill", "none")
        .attr("stroke", "#000")
        .attr("stroke-width", "0.1px")
        .selectAll("rect")
        .data(function (d) {
            return d3.timeDays(new Date(d, 0, 1), new Date(d + 1, 0, 1));
        })
        .enter().append("rect")
        .attr("width", cellSize)
        .attr("height", cellSize)
        .attr("x", function (d) {
            return d3.timeMonday.count(d3.timeYear(d), d) * cellSize;
        })
        .attr("y", function (d) {
            return d.getUTCDay() * cellSize;
        })
        .datum(d3.timeFormat("%Y-%m-%d"))
        .attr('fill', function (d) {
            for (let oneDate of data) {
                if (oneDate["Date"] === d) {
                    return color(oneDate["ProductionValue"])
                }
            }
        })
        .on("mouseover", function () {
            d3.select(this).attr('stroke-width', "1px");
        })
        .on("mouseout", function () {
            d3.select(this).attr('stroke-width', "0.1px");
        })
        .append("title")
        .text(function (d) {
            for (let oneDate of data) {
                if (oneDate["Date"] === d) {
                    return d + ": " + oneDate["ProductionValue"] + "%"
                }
            }
        })

    svg.append("g")
        .attr("fill", "none")
        .attr("stroke", "#000")
        .attr("stroke-width", "1.5px")
        .selectAll("path")
        .data(function (d) {
            return d3.timeMonths(new Date(d, 0, 1), new Date(d + 1, 0, 1));
        })
        .enter().append("path")
        .attr("d", function (d) {
            const t1 = new Date(d.getFullYear(), d.getMonth() + 1, 0),
                d0 = d.getUTCDay(), w0 = d3.timeMonday.count(d3.timeYear(d), d),
                d1 = t1.getUTCDay(), w1 = d3.timeMonday.count(d3.timeYear(t1), t1);
            return "M" + (w0 + 1) * cellSize + "," + d0 * cellSize
                + "H" + w0 * cellSize + "V" + 7 * cellSize
                + "H" + w1 * cellSize + "V" + (d1 + 1) * cellSize
                + "H" + (w1 + 1) * cellSize + "V" + 0
                + "H" + (w0 + 1) * cellSize + "Z";
        });
}


function displayOverview(input) {
    console.log("displaying overview for " + input)
    let data = {
        input: input,
    };
    fetch("/get_live_overview_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            displayOverViewData(result, "navbar-live-company-3-data-production", result["Production"], "#89ab0f");
            displayOverViewData(result, "navbar-live-company-3-data-downtime", result["Downtime"], "#e6ad3c");
            displayOverViewData(result, "navbar-live-company-3-data-poweroff", result["Poweroff"], "#de6b59");
        });
    }).catch((error) => {
        console.log(error)
    });
}

function displayOverViewData(result, elementId, resultElement, elementColor) {
    const data = document.getElementById(elementId)
    let color = document.createElement("div");
    color.style.width = "0.9em"
    color.style.height = "0.9em"
    color.style.background = elementColor
    color.style.border = "0.2px solid black"
    data.appendChild(color)
    let actualPercent = document.createElement("div");
    actualPercent.textContent = resultElement;
    actualPercent.style.width = "30px";
    actualPercent.style.textAlign = "center";
    actualPercent.style.paddingLeft = "5px";
    actualPercent.style.fontWeight = "bold"
    data.appendChild(actualPercent)
}


function displayBestWorkplaces(input) {
    console.log("displaying best workplaces for " + input)
    let data = {
        input: input,
    };
    fetch("/get_live_best_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            displayBestWorstPoweroffData(result, "navbar-live-company-4-best-data-first", result["Data"][0], "rgba(137,171,15,");
            displayBestWorstPoweroffData(result, "navbar-live-company-4-best-data-second", result["Data"][1], "rgba(137,171,15,");
            displayBestWorstPoweroffData(result, "navbar-live-company-4-best-data-third", result["Data"][2], "rgba(137,171,15,");
        });
    }).catch((error) => {
        console.log(error)
    });
}

function displayBestWorstPoweroffData(result, elementId, resultElement, workplaceColor) {
    const data = document.getElementById(elementId)
    // data.style.border = "1px solid black"
    let workplaceName = document.createElement("div");
    workplaceName.textContent = resultElement["WorkplaceName"].replace(/(.{20})..+/, "$1…");
    workplaceName.title = resultElement["WorkplaceName"];
    workplaceName.style.width = "220px";
    workplaceName.style.textAlign = "right";
    workplaceName.style.paddingRight = "5px";
    data.appendChild(workplaceName)

    let color = document.createElement("div");
    color.style.width = "0.9em"
    color.style.height = "0.9em"
    if (elementId.includes("poweroff")) {
        color.style.background = workplaceColor
    } else {
        color.style.background = workplaceColor + (+resultElement["WorkplaceProduction"]) / 100
    }
    color.style.border = "0.2px solid black"
    color.style.marginTop = "1px"
    data.appendChild(color)

    let actualPercent = document.createElement("div");
    if (elementId.includes("poweroff")) {
        actualPercent.textContent = resultElement["WorkplaceProduction"];
    } else {
        actualPercent.textContent = resultElement["WorkplaceProduction"] + "%";
    }
    actualPercent.style.width = "60px";
    actualPercent.style.textAlign = "left";
    actualPercent.style.paddingLeft = "5px";
    data.appendChild(actualPercent)



}

function displayWorstWorkplaces(input) {
    console.log("displaying worst workplaces for " + input)
    let data = {
        input: input,
    };
    fetch("/get_live_worst_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            displayBestWorstPoweroffData(result, "navbar-live-company-4-worst-data-first", result["Data"][0], "rgba(230,173,60,");
            displayBestWorstPoweroffData(result, "navbar-live-company-4-worst-data-second", result["Data"][1], "rgba(230,173,60,");
            displayBestWorstPoweroffData(result, "navbar-live-company-4-worst-data-third", result["Data"][2], "rgba(230,173,60,");
        });
    }).catch((error) => {
        console.log(error)
    });
}

function displayPowerOffWorkplaces(input) {
    console.log("displaying power off workplaces for " + input)
    let data = {
        input: input,
    };
    fetch("/get_live_poweroff_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            displayBestWorstPoweroffData(result, "navbar-live-company-4-poweroff-data-first", result["Data"][0], "#de6b59");
            displayBestWorstPoweroffData(result, "navbar-live-company-4-poweroff-data-second", result["Data"][1], "#de6b59");
            displayBestWorstPoweroffData(result, "navbar-live-company-4-poweroff-data-third", result["Data"][2], "#de6b59");
        });
    }).catch((error) => {
        console.log(error)
    });
}


function displayWorkplaceOverview(input) {
    console.log("displaying workplace overview for " + input)
    let data = {
        input: input,
    };
    fetch("/get_live_all_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            for (const workplace of result["Data"]) {
                console.log(workplace)
                displayWorkplaceData("navbar-live-company-5-all-data", workplace);
            }
        });
    }).catch((error) => {
        console.log(error)
    });
}

function displayWorkplaceData(elementId, workplace) {
    const data = document.getElementById(elementId)
    let container = document.createElement("div")
    container.style.display = "inline-block"

    let workplaceName = document.createElement("div");
    workplaceName.textContent = workplace["WorkplaceName"].replace(/(.{20})..+/, "$1…");
    workplaceName.title = workplace["WorkplaceName"];
    workplaceName.style.width = "220px";
    workplaceName.style.display = "flex"
    workplaceName.style.textAlign = "right";
    workplaceName.style.paddingRight = "5px";
    workplaceName.style.display = "inline-block"
    container.appendChild(workplaceName)

    let color = document.createElement("div");
    color.style.width = "0.9em"
    color.style.height = "0.9em"
    color.style.verticalAlign = "middle"
    if (workplace["WorkplaceProduction"].includes("production")) {
        color.style.background = "#89ab0f"
    } else if (workplace["WorkplaceProduction"].includes("downtime")) {
        color.style.background = "#e6ad3c"
    } else {
        color.style.background = "#de6b59"
    }
    color.style.border = "0.2px solid black"
    color.style.display = "inline-block"
    color.style.marginRight = "2px"
    container.appendChild(color)


    data.appendChild(container)
}