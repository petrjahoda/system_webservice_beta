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
            locale: sessionStorage.getItem("locale")
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


function DrawTimeLine(productionDataset, downtimeDataset, powerOffDataset, data) {
    // access data
    const xAccessor = d => d["Date"]
    const yAccessor = d => d["Value"]
    //set dimensions
    let dimensions = {
        width: screen.width - 500,
        height: 100,
        margin: {top: 0, right: 40, bottom: 20, left: 40,},
    }
    //draw canvas
    const wrapper = d3.select("#navbar-charts-standard-2-timeline")
        .append("svg")
        .attr("viewBox", "0 0 " + dimensions.width + " " + dimensions.height)
    const bounds = wrapper.append("g")


    //set scales
    const chartStartsAt = productionDataset[0]["Date"]
    const chartEndsAt = productionDataset[productionDataset.length - 1]["Date"]
    const xScale = d3.scaleTime()
        .domain([chartStartsAt, chartEndsAt])
        .range([dimensions.margin.left, dimensions.width - dimensions.margin.right])
    const yScale = d3.scaleLinear()
        .domain(d3.extent(productionDataset, yAccessor))
        .range([dimensions.height - dimensions.margin.bottom, 0])


    // prepare data
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

    // prepare tooltip
    let div = d3.select("#navbar-charts-standard-2-timeline").append("div")
        .attr("class", "tooltip")
        .style("opacity", 0.7)
        .style("visibility", "hidden");
    const timeScale = d3.scaleTime()
        .domain([new Date(chartStartsAt * 1000), new Date(chartEndsAt * 1000)])
        .range([dimensions.margin.left, dimensions.width - dimensions.margin.right])
    bounds.append("g")
        .attr("transform", "translate(0," + (dimensions.height - dimensions.margin.bottom) + ")")
        .call(d3.axisBottom(timeScale))

    // draw data with tooltip
    bounds.append("path")
        .attr("d", productionAreaGenerator(productionDataset))
        .attr("fill", data["ProductionColor"])
        .on('mousemove', (event) => {
                let coords = d3.pointer(event);
                let timeEntered = timeScale.invert(coords[0]) / 1000
                let now = new Date(timeEntered * 1000).toLocaleString()
                let start = new Date(productionDataset.filter(i => i["Date"] < timeEntered).pop()["Date"] * 1000).toLocaleString()
                let end = new Date(productionDataset.filter(i => i["Date"] > timeEntered)[0]["Date"] * 1000).toLocaleString()
                div.html(now + "<br/>" + data["ProductionLocale"].toUpperCase() + "<br/>" + start + "<br/>" + end)
                    .style("visibility", "visible")
                    .style("top", (event.pageY) - 60 + "px")
                    .style("left", (event.pageX) - 60 + "px")
            }
        )
        .on('mouseout', () => {
                div.transition()
                div.style("visibility", "hidden")
            }
        )
    bounds.append("path")
        .attr("d", downtimeAreaGenerator(downtimeDataset))
        .attr("fill", data["DowntimeColor"])
        .on('mousemove', (event) => {
                let coords = d3.pointer(event);
                let timeEntered = timeScale.invert(coords[0]) / 1000
                let now = new Date(timeEntered * 1000).toLocaleString()
                let start = new Date(downtimeDataset.filter(i => i["Date"] < timeEntered).pop()["Date"] * 1000).toLocaleString()
                let end = new Date(downtimeDataset.filter(i => i["Date"] > timeEntered)[0]["Date"] * 1000).toLocaleString()
                div.html(now + "<br/>" + data["DowntimeLocale"].toUpperCase() + "<br/>" + start + "<br/>" + end)
                    .style("visibility", "visible")
                    .style("top", (event.pageY) - 60 + "px")
                    .style("left", (event.pageX) - 60 + "px")
            }
        )
        .on('mouseout', () => {
                div.style("visibility", "hidden")
            }
        )
    bounds.append("path")
        .attr("d", powerOffAreaGenerator(powerOffDataset))
        .attr("fill", data["PoweroffColor"])
        .on('mousemove', (event) => {
                let coords = d3.pointer(event);
                let timeEntered = timeScale.invert(coords[0]) / 1000
                let now = new Date(timeEntered * 1000).toLocaleString()
                let start = new Date(powerOffDataset.filter(i => i["Date"] < timeEntered).pop()["Date"] * 1000).toLocaleString()
                let end = new Date(powerOffDataset.filter(i => i["Date"] > timeEntered)[0]["Date"] * 1000).toLocaleString()
                div.html(now + "<br/>" + data["PoweroffLocale"].toUpperCase() + "<br/>" + start + "<br/>" + end)
                    .style("visibility", "visible")
                    .style("top", (event.pageY) - 60 + "px")
                    .style("left", (event.pageX) - 60 + "px")
            }
        )
        .on('mouseout', () => {
                div.style("visibility", "hidden")
            }
        )
}

function updateProductionDataset(productionDataset, startBrush, endBrush) {
    let productionDatasetUpdated
    productionDatasetUpdated = []
    productionDatasetUpdated.push({Date: startBrush * 1000, Value: 1})
    let startProduction = true;
    let lastProduction = 0
    for (const element of productionDataset) {
        if ((+(element["Date"] / 1000) <= +startBrush) && element["Value"] === 1) {
            lastProduction = element["Date"]
        }
        if ((+(element["Date"] / 1000) >= +startBrush) && (+(element["Date"] / 1000) <= +endBrush)) {
            if (startProduction) {
                if (element["Value"] === 0) {
                    productionDatasetUpdated.push({Date: startBrush * 1000, Value: 1})
                }
                startProduction = false;
            }
            productionDatasetUpdated.push(element)
        }
    }
    productionDatasetUpdated.push({Date: endBrush * 1000, Value: 0})
    return {productionDatasetUpdated, lastProduction};
}

function updateDowntimeDataset(downtimeDataset, startBrush, endBrush) {
    let downtimeDatasetUpdated
    downtimeDatasetUpdated = []
    let startDowntime = true;
    let lastDowntime = 0
    for (const element of downtimeDataset) {
        if ((+(element["Date"] / 1000) <= +startBrush) && element["Value"] === 1) {
            lastDowntime = element["Date"]
        }
        if ((+(element["Date"] / 1000) >= +startBrush) && (+(element["Date"] / 1000) <= +endBrush)) {
            if (startDowntime) {
                if (element["Value"] === 0) {
                    downtimeDatasetUpdated.push({Date: startBrush * 1000, Value: 1})
                }
                startDowntime = false;
            }
            downtimeDatasetUpdated.push(element)
        }
    }
    if ((downtimeDatasetUpdated.length > 0) && (downtimeDatasetUpdated[downtimeDatasetUpdated.length - 1]["Value"] === 1)) {
        downtimeDatasetUpdated.push({Date: endBrush * 1000, Value: 0})
    }
    return {downtimeDatasetUpdated, lastDowntime};
}

function updatePowerOffDataset(powerOffDataset, startBrush, endBrush) {
    let powerOffDatasetUpdated
    powerOffDatasetUpdated = []
    let startPowerOff = true;
    let lastPowerOff = 0
    for (const element of powerOffDataset) {
        if ((+(element["Date"] / 1000) <= +startBrush) && element["Value"] === 1) {
            lastPowerOff = element["Date"]
        }
        if ((+(element["Date"] / 1000) >= +startBrush) && (+(element["Date"] / 1000) <= +endBrush)) {
            if (startPowerOff) {
                if (element["Value"] === 0) {
                    powerOffDatasetUpdated.push({Date: startBrush * 1000, Value: 1})
                }
                startPowerOff = false;
            }
            powerOffDatasetUpdated.push(element)
        }
    }
    if ((powerOffDatasetUpdated.length > 0) && (powerOffDatasetUpdated[powerOffDatasetUpdated.length - 1]["Value"] === 1)) {
        powerOffDatasetUpdated.push({Date: endBrush * 1000, Value: 0})
    }
    return {powerOffDatasetUpdated, lastPowerOff};
}

function DrawTimeLinePanel(productionDataset, downtimeDataset, powerOffDataset, data) {
    const xAccessor = d => d["Date"]
    const yAccessor = d => d["Value"]
    let dimensionsPanel = {
        width: screen.width - 500,
        height: 20,
        margin: {top: 0, right: 40, bottom: 0, left: 40,},
    }

    const chartStartsAt = productionDataset[0]["Date"]
    const chartEndsAt = productionDataset[productionDataset.length - 1]["Date"]

    const wrapperPanel = d3.select("#navbar-charts-standard-2-timeline-panel")
        .append("svg")
        .attr("viewBox", "0 0 " + dimensionsPanel.width + " " + dimensionsPanel.height)
    const boundsPanel = wrapperPanel.append("g")

    const xScalePanel = d3.scaleTime()
        .domain([chartStartsAt, chartEndsAt])
        .range([dimensionsPanel.margin.left, dimensionsPanel.width - dimensionsPanel.margin.right])
    const yScalePanel = d3.scaleLinear()
        .domain(d3.extent(productionDataset, yAccessor))
        .range([dimensionsPanel.height - dimensionsPanel.margin.bottom, 0])

    const productionAreaGeneratorPanel = d3.area()
        .x(d => xScalePanel(xAccessor(d)))
        .y0(dimensionsPanel.height - dimensionsPanel.margin.bottom)
        .y1(d => yScalePanel(yAccessor(d)))
        .curve(d3.curveStepAfter)
    const downtimeAreaGeneratorPanel = d3.area()
        .x(d => xScalePanel(xAccessor(d)))
        .y0(dimensionsPanel.height - dimensionsPanel.margin.bottom)
        .y1(d => yScalePanel(yAccessor(d)))
        .curve(d3.curveStepAfter);
    const powerOffAreaGeneratorPanel = d3.area()
        .x(d => xScalePanel(xAccessor(d)))
        .y0(dimensionsPanel.height - dimensionsPanel.margin.bottom)
        .y1(d => yScalePanel(yAccessor(d)))
        .curve(d3.curveStepAfter);

    boundsPanel.append("path")
        .attr("d", productionAreaGeneratorPanel(productionDataset))
        .attr("fill", data["ProductionColor"])
    boundsPanel.append("path")
        .attr("d", downtimeAreaGeneratorPanel(downtimeDataset))
        .attr("fill", data["DowntimeColor"])
    boundsPanel.append("path")
        .attr("d", powerOffAreaGeneratorPanel(powerOffDataset))
        .attr("fill", data["PoweroffColor"])


    const brush = d3.brushX()
        .extent([[dimensionsPanel.margin.left, 0.5], [dimensionsPanel.width - dimensionsPanel.margin.right, dimensionsPanel.height - dimensionsPanel.margin.bottom + 0.5]])
        .on("brush", brushEnd)
        .on("end", brushEnd);
    wrapperPanel.append("g")
        .call(brush)

    function CheckDatasetsForZeroValues(productionDatasetUpdated, downtimeDatasetUpdated, powerOffDatasetUpdated, lastDowntime, lastPowerOff, startBrush, endBrush, lastProduction) {
        console.log("")
        console.log(productionDatasetUpdated.length)
        console.log(downtimeDatasetUpdated.length)
        console.log(powerOffDatasetUpdated.length)
        if ((productionDatasetUpdated.length === 2) && (downtimeDatasetUpdated.length === 0) && (powerOffDatasetUpdated.length === 0)) {
            if (lastDowntime > lastPowerOff) {
                downtimeDatasetUpdated = []
                downtimeDatasetUpdated.push({Date: startBrush * 1000, Value: 1})
                downtimeDatasetUpdated.push({Date: endBrush * 1000, Value: 0})
                powerOffDatasetUpdated = []
            } else if (lastDowntime < lastPowerOff) {
                powerOffDatasetUpdated = []
                powerOffDatasetUpdated.push({Date: startBrush * 1000, Value: 1})
                powerOffDatasetUpdated.push({Date: endBrush * 1000, Value: 0})
                downtimeDatasetUpdated = []
            }
            if ((lastProduction > lastDowntime) && (lastProduction>lastPowerOff)) {
                powerOffDatasetUpdated = []
                downtimeDatasetUpdated = []
            }
        }
        return {downtimeDatasetUpdated, powerOffDatasetUpdated};
    }

    function brushEnd({selection}) {
        if (selection) {
            let startBrush = xScalePanel.invert(selection[0]) / 1000
            let endBrush = xScalePanel.invert(selection[1]) / 1000
            if (selection[0] > 0) {
                document.getElementById("navbar-charts-standard-2-timeline").innerHTML = ""
                let {productionDatasetUpdated, lastProduction} = updateProductionDataset(productionDataset, startBrush, endBrush);
                let {downtimeDatasetUpdated, lastDowntime} = updateDowntimeDataset(downtimeDataset, startBrush, endBrush);
                let {powerOffDatasetUpdated, lastPowerOff} = updatePowerOffDataset(powerOffDataset, startBrush, endBrush);
                const updatedDatasets = CheckDatasetsForZeroValues(productionDatasetUpdated, downtimeDatasetUpdated, powerOffDatasetUpdated, lastDowntime, lastPowerOff, startBrush, endBrush, lastProduction);
                downtimeDatasetUpdated = updatedDatasets.downtimeDatasetUpdated;
                powerOffDatasetUpdated = updatedDatasets.powerOffDatasetUpdated;
                DrawTimeLine(productionDatasetUpdated, downtimeDatasetUpdated, powerOffDatasetUpdated, data);
            }
        }

    }
}

function drawTimelineChart(data) {
    const productionDataset = data["ProductionData"]
    const downtimeDataset = data["DowntimeData"]
    const powerOffDataset = data["PowerOffData"]
    DrawTimeLine(productionDataset, downtimeDataset, powerOffDataset, data);
    DrawTimeLinePanel(productionDataset, downtimeDataset, powerOffDataset, data);
}