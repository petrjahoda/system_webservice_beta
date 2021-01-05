function callNavbarLiveCompanyJs() {
    console.log("Starting processing data for navbar-live-company")
    displayLiveProductivityData("company")
    displayCalendar("company")
    displayOverview("company")
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

function displayCalendar(input) {
    console.log("displaying calendar for " + input)
    d3.json("/get_calendar_data", {
        method: "POST",
        body: JSON.stringify({input: input})
    }).then((data) => {
        drawCalendar(data["Data"]);
    }).catch((error) => {
        console.error("Error loading the data: " + error);
    });
}

function drawCalendar(data) {
    let result = new Map(data.map(value => [value["Date"], value["ProductionValue"]]));
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
            displayOverViewData("navbar-live-company-3-production", result["Production"], "#89ab0f");
            displayOverViewData("navbar-live-company-3-downtime", result["Downtime"], "#e6ad3c");
            displayOverViewData("navbar-live-company-3-poweroff", result["Poweroff"], "#de6b59");
        });
    }).catch((error) => {
        console.log(error)
    });
}

function displayOverViewData(elementId, resultElement, elementColor) {
    console.log("Processing data with size of " + resultElement.length + " for " + elementId)
    const data = document.getElementById(elementId)
    data.textContent = resultElement.length + " " + data.textContent.toUpperCase()
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
        workplaceName.textContent = workplace["WorkplaceName"].substring(0,25) + " . " + workplace["StateDuration"];
        workplaceName.title = workplace["WorkplaceName"];
        workplaceName.style.fontWeight = "normal"
        workplaceName.style.fontSize = "0.9em"
        
        workplaceName.id = workplace["WorkplaceName"]
        workplaceContent.appendChild(workplaceName)
        let elementData = document.getElementById(workplace["WorkplaceName"])
        console.log(elementData.clientWidth)
        while (elementData.clientWidth <= 275) {
            elementData.textContent = elementData.textContent.replace(".", "..")
        }
        elementData.style.marginRight = "" + 278-elementData.clientWidth + "px"
    }
}