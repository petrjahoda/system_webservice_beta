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
            drawCalendar(result.Data);
        });
    }).catch((error) => {
        console.log(error)
    });
}

function drawCalendar(data) {
    const width = 960, height = 136, cellSize = 17;
    const color = d3.scaleQuantize()
        .domain([0, 100])
        .range(["#ffffff", "#f3f6e7", "#e7eecf", "#dbe5b7", "#d0dd9f", "#c4d587", "#b8cd6f", "#acc457", "#a1bc3f", "#94b327", "#89ab0f"]);
    const svg = d3.select("#navbar-live-company-2-calendar-chart")
        .selectAll("svg")
        .data(d3.range(2019, 2021))
        .enter().append("svg")
        .attr("width", width)
        .attr("height", height)
        .append("g")
        .attr("transform", "translate(" + ((width - cellSize * 53) / 2) + "," + (height - cellSize * 7 - 1) + ")");
    // svg.append("text")
    //     .attr("transform", "translate(-6," + cellSize * 3.5 + ")rotate(-90)")
    //     .attr("font-family", "sans-serif")
    //     .attr("font-size", 10)
    //     .attr("text-anchor", "middle")
    //     .text(function (d) {
    //         return d;
    //     });
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
                if (oneDate.Date === d) {
                    return color(oneDate.ProductionValue)
                }
            }
        })
        .append("title")
        .text(function (d) {
            for (let oneDate of data) {
                if (oneDate.Date === d) {
                    return d + ": " + oneDate.ProductionValue + "%"
                }
            }
        });
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