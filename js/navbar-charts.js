function displaySelectionCharts(input, name) {
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
            displaySelectionChartsData("charts-select-" + input, result["SelectionData"], input);
            let selectionElement = document.getElementById("charts-select-" + input);
            if (sessionStorage.getItem(input + "Name") === null) {
                sessionStorage.setItem(input + "Name", result["SelectionData"][0])
            }
            selectionElement.addEventListener("change", (event) => {
                sessionStorage.setItem(input + "Name", selectionElement.value)
            })
        });

    }).catch((error) => {
        console.log(error)
    });
}

function displaySelectionChartsData(element, selectionData, input) {
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

function displayTimelineChart(chartsStart, chartsEnd, workplace) {
    console.log("Displaying timeline chart for " + workplace)
    d3.json("/get_timeline_chart", {
        method: "POST",
        body: JSON.stringify({
            chartsStart: new Date(chartsStart).toISOString(),
            chartsEnd: new Date(chartsEnd).toISOString(),
            workplace: workplace,
        })
    }).then((data) => {
        drawTimelineChart(data);
    }).catch((error) => {
        console.error("Error loading the data: " + error);
    });
}

function displayAnalogChart() {
    console.log("Displaying analog chart")
}

function displayDigitalChart() {
    console.log("Displaying digital chart")
}

function displayProductionRateChart() {
    console.log("Displaying production rate chart")
}


function drawTimelineChart(data) {
    // data
    const productionDataset = data["ProductionData"]
    const downtimeDataset = data["DowntimeData"]
    const powerOffDataset = data["PowerOffData"]
    const xAccessor = d => d["Date"]
    const yAccessor = d => d["Value"]

    //dimensions
    let dimensions = {
        width: screen.width - 500,
        height: 70,
        margin: {top: 0, right: 40, bottom: 20, left: 40,},
    }

    //canvas
    const wrapper = d3.select("#navbar-charts-standard-2-timeline")
        .append("svg")
        .attr("viewBox", "0 0 " + dimensions.width + " 100")
        .attr("preserveAspectRatio", "xMinYMin meet")
    const bounds = wrapper.append("g")

    //scales
    const chartStartsAt = productionDataset[0]["Date"]
    const chartEndsAt = productionDataset[productionDataset.length - 1]["Date"]
    const xScale = d3.scaleTime()
        .domain([chartStartsAt, chartEndsAt])
        .range([dimensions.margin.left, dimensions.width - dimensions.margin.right])
    const yScale = d3.scaleLinear()
        .domain(d3.extent(productionDataset, yAccessor))
        .range([dimensions.height - dimensions.margin.bottom, 0])

    // draw data
    const productionAreaGenerator = d3.area()
        .x(d => xScale(xAccessor(d)))
        .y0(dimensions.height - dimensions.margin.bottom)
        .y1(d => yScale(yAccessor(d)))
        .curve(d3.curveStepAfter);
    const downtimeAreaGenerator = d3.area()
        .x(d => xScale(xAccessor(d)))
        .y0(dimensions.height - dimensions.margin.bottom)
        .y1(d => yScale(yAccessor(d)))
        .curve(d3.curveStepAfter);
    const powerOffAreaGenerator = d3.area()
        .x(d => xScale(xAccessor(d)))
        .y0(dimensions.height - dimensions.margin.bottom)
        .y1(d => yScale(yAccessor(d)))
        .curve(d3.curveStepAfter);

    bounds.append("path")
        .attr("d", productionAreaGenerator(productionDataset))
        .attr("fill", data["ProductionColor"])
        .on('mouseenter', (event) => {
                let coords = d3.pointer(event);
                let timeEntered = timeScale.invert(coords[0]) / 1000
                let start = new Date(productionDataset.filter(i => i["Date"] < timeEntered).pop()["Date"] * 1000).toLocaleString()
                let end = new Date(productionDataset.filter(i => i["Date"] > timeEntered)[0]["Date"] * 1000).toLocaleString()
                div.transition()
                    .duration(200)
                    .style("opacity", .9);
                div.html(start + "<br/>" + end)
                    .style("visibility", "visible")
                    .style("top", (event.pageY) + "px")
                    .style("left", (event.pageX) - 60 + "px")
            }
        )
        .on('mouseout', () => {
                div.style("visibility", "hidden")
            }
        )
    bounds.append("path")
        .attr("d", downtimeAreaGenerator(downtimeDataset))
        .attr("fill", data["DowntimeColor"])
        .on('mouseenter', (event) => {
                let coords = d3.pointer(event);
                let timeEntered = timeScale.invert(coords[0]) / 1000
                let start = new Date(downtimeDataset.filter(i => i["Date"] < timeEntered).pop()["Date"] * 1000).toLocaleString()
                let end = new Date(downtimeDataset.filter(i => i["Date"] > timeEntered)[0]["Date"] * 1000).toLocaleString()
                div.transition()
                    .duration(200)
                    .style("opacity", .9);
                div.html(start + "<br/>" + end)
                    .style("visibility", "visible")
                    .style("top", (event.pageY) + "px")
                    .style("left", (event.pageX) - 60 + "px")
            }
        )
        .on('mouseout', () => {
                div.style("visibility", "hidden")
            }
        )

    let div = d3.select("#navbar-charts-standard-2-timeline").append("div")
        .attr("class", "tooltip")
        .style("opacity", 0);

    bounds.append("path")
        .attr("d", powerOffAreaGenerator(powerOffDataset))
        .attr("fill", data["PoweroffColor"])
        .on('mouseenter', (event) => {
                let coords = d3.pointer(event);
                let timeEntered = timeScale.invert(coords[0]) / 1000
                let start = new Date(powerOffDataset.filter(i => i["Date"] < timeEntered).pop()["Date"] * 1000).toLocaleString()
                let end = new Date(powerOffDataset.filter(i => i["Date"] > timeEntered)[0]["Date"] * 1000).toLocaleString()
                div.transition()
                    .duration(200)
                    .style("opacity", .9);
                div.html(start + "<br/>" + end)
                    .style("visibility", "visible")
                    .style("top", (event.pageY) + "px")
                    .style("left", (event.pageX) - 60 + "px")
            }
        )
        .on('mouseout', () => {
                div.style("visibility", "hidden")
            }
        )


// Define the axes
    const timeScale = d3.scaleTime()
        .domain([new Date(chartStartsAt * 1000), new Date(chartEndsAt * 1000)])
        .nice()
        .range([dimensions.margin.left, dimensions.width - dimensions.margin.right])
    bounds.append("g")
        .attr("transform", "translate(0," + (dimensions.height - dimensions.margin.bottom) + ")")
        .call(d3.axisBottom(timeScale))


}