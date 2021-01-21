function displayLiveProductivity(input, selectionData) {
    console.log("displaying live productivity data for " + input + " and " +selectionData)
    let data = {
        input: input,
        selection: selectionData
    };
    fetch("/get_live_productivity_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            displayLiveProductivityData(result, "navbar-live-" + input + "-1-data-day", result["Today"], result["Yesterday"]);
            displayLiveProductivityData(result, "navbar-live-" + input + "-1-data-week", result["ThisWeek"], result["LastWeek"]);
            displayLiveProductivityData(result, "navbar-live-" + input + "-1-data-month", result["ThisMonth"], result["LastMonth"]);
        });
    }).catch((error) => {
        console.log(error)
    });
}

function displayLiveProductivityData(result, elementId, actual, last) {
    const data = document.getElementById(elementId)
    data.innerHTML = ""
    let lastPercent = document.createElement("div");
    lastPercent.textContent = last + "%";
    lastPercent.style.width = "55px";
    lastPercent.style.textAlign = "right";
    lastPercent.style.paddingRight = "5px";
    lastPercent.style.fontWeight = "bold"
    data.appendChild(lastPercent)
    let lastColor = document.createElement("div");
    lastColor.style.width = "15px"
    lastColor.style.height = "15px"
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
    actualColor.style.width = "15px"
    actualColor.style.height = "15px"
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

function displayCalendar(input, selectionData) {
    console.log("displaying calendar for " + input)
    d3.json("/get_calendar_data", {
        method: "POST",
        body: JSON.stringify({input: input, selection: selectionData})
    }).then((data) => {
        drawCalendar(data["Data"], input);
    }).catch((error) => {
        console.error("Error loading the data: " + error);
    });
}

function drawCalendar(data, input) {
    if ((input === "group") || (input === "workplace")) {
        document.getElementById("navbar-live-" + input + "-2-calendar-chart").innerHTML = ""
    }
    let result = new Map(data.map(value => [value["Date"], value["ProductionValue"]]));
    const width = 960, height = 136, cellSize = 17;
    const color = d3.scaleQuantize()
        .domain([0, 100])
        .range(["#f3f6e7", "#e7eecf", "#dbe5b7", "#d0dd9f", "#c4d587", "#b8cd6f", "#acc457", "#a1bc3f", "#94b327", "#89ab0f"]);

    let today = new Date();
    let startYear = today.getFullYear()
    let endYear = startYear+1
    for (const dayData of data) {
        let yearOfDay = parseInt(dayData["Date"])
        if (yearOfDay < startYear)
            startYear = yearOfDay
    }
    const svg = d3.select("#navbar-live-" + input + "-2-calendar-chart")
        .selectAll("svg")
        .data(d3.range(startYear, endYear))
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
        .text(d => d)

    svg.append("g")
        .attr("fill", "none")
        .attr("stroke", "#000")
        .attr("stroke-width", "0.1px")
        .selectAll("rect")
        .data(d => d3.timeDays(new Date(d, 0, 1), new Date(d + 1, 0, 1)))
        .enter().append("rect")
        .attr("width", cellSize)
        .attr("height", cellSize)
        .attr("x", d => d3.timeMonday.count(d3.timeYear(d), d) * cellSize)
        .attr("y", d => d.getUTCDay() * cellSize)
        .datum(d3.timeFormat("%Y-%m-%d"))
        .attr('fill', d => color(result.get(d)))
        .on("mouseover", function () {
            d3.select(this).attr('stroke-width', "1px");
        })
        .on("mouseout", function () {
            d3.select(this).attr('stroke-width', "0.1px");
        })
        .append("title")
        .text(d => d + ": " + result.get(d) + "%")

    svg.append("g")
        .attr("fill", "none")
        .attr("stroke", "#000")
        .attr("stroke-width", "1.5px")
        .selectAll("path")
        .data(d => d3.timeMonths(new Date(d, 0, 1), new Date(d + 1, 0, 1)))
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


function displayOverview(input, selectionData) {
    console.log("displaying overview for " + input)
    let data = {
        input: input,
        selection: selectionData,
    };
    fetch("/get_live_overview_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            displayOverViewData("navbar-live-" + input + "-3-production", result["Production"], "#89ab0f");
            displayOverViewData("navbar-live-" + input + "-3-downtime", result["Downtime"], "#e6ad3c");
            displayOverViewData("navbar-live-" + input + "-3-poweroff", result["Poweroff"], "#de6b59");
        });
    }).catch((error) => {
        console.log(error)
    });
}

function displayOverViewData(elementId, resultElement, elementColor) {
    console.log("Clearing previous data")
    let node = document.getElementById(elementId);
    node.querySelectorAll('*').forEach(n => n.remove());
    console.log("Updating text for " + elementId)
    const data = document.getElementById(elementId)
    if (resultElement != null) {
        console.log(resultElement.length)
        if (resultElement.length > 0) {
            data.textContent = resultElement.length + "" + data.textContent.replace(/\d/g, "")
        } else {
            data.textContent = "0" + data.textContent.replace(/\d/g, "")
        }
        console.log("Adding new data with size of " + resultElement.length + " for " + elementId)
        let content = document.createElement("div")
        content.style.display = "block"
        content.style.marginTop = "20px"
        data.appendChild(content)
        for (const workplace of resultElement) {
            let workplaceContent = document.createElement("div")
            workplaceContent.style.marginTop = "5px"
            workplaceContent.style.marginLeft = "5px"
            workplaceContent.style.marginRight = "5px"
            workplaceContent.style.display = "flex"
            workplaceContent.style.justifyContent = "space-evenly"
            content.appendChild(workplaceContent)

            let color = document.createElement("div");
            color.style.width = "15px"
            color.style.height = "15px"
            color.style.background = elementColor
            color.style.border = "0.2px solid black"
            color.style.display = "flex"
            workplaceContent.appendChild(color)

            let workplaceName = document.createElement("div");
            workplaceName.textContent = workplace["WorkplaceName"].substring(0, 25) + " . " + workplace["StateDuration"];
            workplaceName.title = workplace["WorkplaceName"];
            workplaceName.style.fontWeight = "normal"
            workplaceName.style.fontSize = "0.9em"

            workplaceName.id = workplace["WorkplaceName"]
            workplaceContent.appendChild(workplaceName)
            let elementData = document.getElementById(workplace["WorkplaceName"])
            while (elementData.clientWidth <= 270) {
                elementData.textContent = elementData.textContent.replace(".", "..")
            }
            elementData.style.marginRight = "" + 276 - elementData.clientWidth + "px"
        }
    } else {
        data.textContent = "0" + data.textContent.replace(/\d/g, "")
    }

}


function displaySelection(input) {
    console.log("downloading selection for " + input)
    let data = {
        input: input,
    };
    fetch("/get_live_selection", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            displaySelectionData("live-select-" + input, result["SelectionData"], input);
            let selectionElement = document.getElementById("live-select-" + input);
            if (sessionStorage.getItem(input + "Name") === null) {
                sessionStorage.setItem(input + "Name", result["SelectionData"][0])
            }
            selectionElement.addEventListener("change", (event) => {
                sessionStorage.setItem(input + "Name", selectionElement.value)
                displayLiveProductivity(input, selectionElement.value)
                displayCalendar(input, selectionElement.value)
                if (input === "group") {
                    displayOverview(input, selectionElement.value)
                }
                if (input === "workplace") {
                    displayWorkplaceData("workplace", selectionElement.value)
                }

            })
        });

    }).catch((error) => {
        console.log(error)
    });
}

function displaySelectionData(element, selectionData, input) {
    let content = document.getElementById(element)
    for (let selection of selectionData) {
        let selectionData = document.createElement("option")
        selectionData.value = selection
        selectionData.appendChild(document.createTextNode(selection));
        if (selection === sessionStorage.getItem(input + "Name")) {
            console.log("match for " + selection)
            selectionData.selected = "selected"
        }
        content.appendChild(selectionData)
    }
}

function displayCompanyName(company) {
    console.log("downloading name for " + company)
    let companyName = sessionStorage.getItem("CompanyName")
    console.log(companyName)
    if (companyName != null) {
        document.getElementById("navbar-live-company-name").textContent = companyName
        return
    }
    fetch("/get_company_name", {
        method: "POST",
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            document.getElementById("navbar-live-company-name").textContent = result["CompanyName"]
            sessionStorage.setItem("CompanyName", result["CompanyName"])
        });

    }).catch((error) => {
        console.log(error)
    });
}


function displayWorkplaceData(workplace, savedSelection) {
    console.log("downloading workplace data for  for " + savedSelection)
    let data = {
        input: savedSelection,
    };
    fetch("/get_workplace_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            console.log(result)

            if (result["Order"].length > 0) {
                document.getElementById("navbar-live-workplace-3-order-data").textContent = result["Order"] + " (" + result["OrderDuration"] + ")"
                console.log(result["OrderDuration"])
            } else {
                document.getElementById("navbar-live-workplace-3-order-data").textContent = "-"
            }
            if (result["User"].length > 1) {
                document.getElementById("navbar-live-workplace-3-user-data").textContent = result["User"]
            } else {
                document.getElementById("navbar-live-workplace-3-user-data").textContent = "-"
            }
            if (result["Downtime"].length > 0) {
                document.getElementById("navbar-live-workplace-3-downtime-data").textContent = result["Downtime"] + " (" + result["DowntimeDuration"] + ")"
            } else {
                document.getElementById("navbar-live-workplace-3-downtime-data").textContent = "-"
            }
            if (result["Breakdown"].length > 0) {
                document.getElementById("navbar-live-workplace-3-breakdown-data").textContent = result["Breakdown"] + " (" + result["BreakdownDuration"] + ")"
            } else {
                document.getElementById("navbar-live-workplace-3-breakdown-data").textContent = "-"
            }
            if (result["Alarm"].length > 0) {
                document.getElementById("navbar-live-workplace-3-alarm-data").textContent = result["Alarm"] + " (" + result["AlarmDuration"] + ")"
            } else {
                document.getElementById("navbar-live-workplace-3-alarm-data").textContent = "-"
            }
        });
    }).catch((error) => {
        console.log(error)
    });
}