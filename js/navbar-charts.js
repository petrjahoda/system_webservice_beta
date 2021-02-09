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
            selectionData.selected = "selected"
        }
        content.appendChild(selectionData)
    }
}

function displayAllCharts(chartsStart, chartsEnd, workplace) {
    console.log("Displaying charts for " + workplace)
    document.getElementById("charts-select-workplace").disabled = true
    document.getElementById("charts-standard-start").disabled = true
    document.getElementById("charts-standard-end").disabled = true
    d3.json("/get_chart_data", {
        method: "POST",
        body: JSON.stringify({
            chartsStart: new Date(chartsStart).toISOString(),
            chartsEnd: new Date(chartsEnd).toISOString(),
            workplace: workplace,
            locale: sessionStorage.getItem("locale")
        })
    }).then((data) => {
        drawCharts(data);
    }).catch((error) => {
        console.error("Error loading the data: " + error);
    });
}

function displayAnalogChart() {
    console.log("Displaying analog chart")
}


function displayProductionRateChart() {
    console.log("Displaying production rate chart")
}


function DrawTimeLine(productionDataset, downtimeDataset, powerOffDataset, productionColor, downtimeColor, powerOffColor, productionName, downtimeName, powerOffName) {
    const xAccessor = d => d["Date"]
    const yAccessor = d => d["Value"]
    let dimensions = {
        width: screen.width - 500,
        height: 60,
        margin: {top: 0, right: 40, bottom: 20, left: 40,},
    }

    const wrapper = d3.select("#navbar-charts-standard-2-timeline")
        .append("svg")
        .attr("viewBox", "0 0 " + dimensions.width + " " + dimensions.height)
    const bounds = wrapper.append("g")

    const chartStartsAt = productionDataset[0]["Date"]
    const chartEndsAt = productionDataset[productionDataset.length - 1]["Date"]

    const xScale = d3.scaleTime()
        .domain([chartStartsAt, chartEndsAt])
        .range([dimensions.margin.left, dimensions.width - dimensions.margin.right])
    const yScale = d3.scaleLinear()
        .domain(d3.extent(productionDataset, yAccessor))
        .range([dimensions.height - dimensions.margin.bottom, 0])

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

    const timeScale = d3.scaleTime()
        .domain([new Date(chartStartsAt * 1000), new Date(chartEndsAt * 1000)])
        .range([dimensions.margin.left, dimensions.width - dimensions.margin.right])
    bounds.append("g")
        .attr("transform", "translate(0," + (dimensions.height - dimensions.margin.bottom) + ")")
        .call(d3.axisBottom(timeScale))

    let div = d3.select("#navbar-charts-standard-2-timeline").append("div")
        .attr("class", "tooltip")
        .style("opacity", 0.7)
        .style("visibility", "hidden");

    bounds.append("path")
        .attr("d", productionAreaGenerator(productionDataset))
        .attr("fill", productionColor)
        .on('mousemove', (event) => {
                let coords = d3.pointer(event);
                let timeEntered = timeScale.invert(coords[0]) / 1000
                let now = new Date(timeEntered * 1000).toLocaleString()
                let start = new Date(productionDataset.filter(i => i["Date"] < timeEntered).pop()["Date"] * 1000).toLocaleString()
                let end = new Date(productionDataset.filter(i => i["Date"] > timeEntered)[0]["Date"] * 1000).toLocaleString()
                div.html(now + "<br/>" + productionName.toUpperCase() + "<br/>" + start + "<br/>" + end)
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
        .attr("fill", downtimeColor)
        .on('mousemove', (event) => {
                let coords = d3.pointer(event);
                let timeEntered = timeScale.invert(coords[0]) / 1000
                let now = new Date(timeEntered * 1000).toLocaleString()
                let start = new Date(downtimeDataset.filter(i => i["Date"] < timeEntered).pop()["Date"] * 1000).toLocaleString()
                let end = new Date(downtimeDataset.filter(i => i["Date"] > timeEntered)[0]["Date"] * 1000).toLocaleString()
                div.html(now + "<br/>" + downtimeName.toUpperCase() + "<br/>" + start + "<br/>" + end)
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
        .attr("fill", powerOffColor)
        .on('mousemove', (event) => {
                let coords = d3.pointer(event);
                let timeEntered = timeScale.invert(coords[0]) / 1000
                let now = new Date(timeEntered * 1000).toLocaleString()
                let start = new Date(powerOffDataset.filter(i => i["Date"] < timeEntered).pop()["Date"] * 1000).toLocaleString()
                let end = new Date(powerOffDataset.filter(i => i["Date"] > timeEntered)[0]["Date"] * 1000).toLocaleString()
                div.html(now + "<br/>" + powerOffName.toUpperCase() + "<br/>" + start + "<br/>" + end)
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

function UpdateDigitalDataset(digitalDataset, startBrush, endBrush) {
    let updatedDigitalDataset
    updatedDigitalDataset = []
    for (const dataset of digitalDataset) {let datasetUpdated
        datasetUpdated = {Name: dataset["Name"], Data:[]}
        let data
        data = []
        let startValue = true;
        let lastDataValue = 0
        for (const element of dataset["Data"]) {
            if ((+(element["Date"] / 1000) <= +startBrush)) {
                lastDataValue = element["Value"]
            }
            if ((+(element["Date"] / 1000) >= +startBrush) && (+(element["Date"] / 1000) <= +endBrush)) {
                if (startValue) {
                    if (element["Value"] === 0) {
                        data.push({Date: startBrush * 1000, Value: 1})
                    }
                    if (element["Value"] === 1) {
                        data.push({Date: startBrush * 1000, Value: 0})
                    }
                    startValue = false;
                }
                data.push(element)
            }
        }
        if (data.length === 0) {
            if (lastDataValue === 1) {
                data.push({Date: startBrush * 1000, Value: 1})
                data.push({Date: endBrush * 1000, Value: 0})
            } else {
                data.push({Date: startBrush * 1000, Value: 0})
                data.push({Date: endBrush * 1000, Value: 1})
            }
        } else {
            data.push({Date: endBrush * 1000, Value: 1})
        }
        datasetUpdated["Data"] = data
        updatedDigitalDataset.push(datasetUpdated)
    }
    return updatedDigitalDataset
}

function DrawTimeLinePanel(productionDataset, downtimeDataset, powerOffDataset, productionColor, downtimeColor, powerOffColor, productionName, downtimeName, powerOffName, digitalDataset) {
    const xAccessor = d => d["Date"]
    const yAccessor = d => d["Value"]
    let dimensionsPanel = {
        width: screen.width - 500,
        height: 40,
        margin: {top: 0, right: 40, bottom: 20, left: 40,},
    }

    const wrapperPanel = d3.select("#navbar-charts-standard-2-timeline-panel")
        .append("svg")
        .attr("viewBox", "0 0 " + dimensionsPanel.width + " " + dimensionsPanel.height)
    const boundsPanel = wrapperPanel.append("g")

    const chartStartsAt = productionDataset[0]["Date"]
    const chartEndsAt = productionDataset[productionDataset.length - 1]["Date"]

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

    const timeScalePanel = d3.scaleTime()
        .domain([new Date(chartStartsAt * 1000), new Date(chartEndsAt * 1000)])
        .range([dimensionsPanel.margin.left, dimensionsPanel.width - dimensionsPanel.margin.right])
    boundsPanel.append("g")
        .attr("transform", "translate(0," + (dimensionsPanel.height - dimensionsPanel.margin.bottom) + ")")
        .call(d3.axisBottom(timeScalePanel))

    boundsPanel.append("path")
        .attr("d", productionAreaGeneratorPanel(productionDataset))
        .attr("fill", productionColor)
    boundsPanel.append("path")
        .attr("d", downtimeAreaGeneratorPanel(downtimeDataset))
        .attr("fill", downtimeColor)
    boundsPanel.append("path")
        .attr("d", powerOffAreaGeneratorPanel(powerOffDataset))
        .attr("fill", powerOffColor)


    const brush = d3.brushX()
        .extent([[dimensionsPanel.margin.left, 0.5], [dimensionsPanel.width - dimensionsPanel.margin.right, dimensionsPanel.height - dimensionsPanel.margin.bottom + 0.5]])
        .on("brush", brushEnd)
        .on("end", brushEnd);
    wrapperPanel.append("g")
        .call(brush)

    //
    function CheckDatasetsForZeroValues(productionDatasetUpdated, downtimeDatasetUpdated, powerOffDatasetUpdated, lastDowntime, lastPowerOff, startBrush, endBrush, lastProduction) {
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
            if ((lastProduction > lastDowntime) && (lastProduction > lastPowerOff)) {
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
                DrawTimeLine(productionDatasetUpdated, downtimeDatasetUpdated, powerOffDatasetUpdated, productionColor, downtimeColor, powerOffColor, productionName, downtimeName, powerOffName);
                document.getElementById("navbar-charts-standard-3-digital").innerHTML = ""
                const updatedDigitalDataset = UpdateDigitalDataset(digitalDataset, startBrush, endBrush)
                DrawDigital(updatedDigitalDataset)
            }
        }

    }
}

function DrawDigital(digitalDataset) {
    for (const dataset of digitalDataset) {
        // access data
        const xAccessor = d => d["Date"]
        const yAccessor = d => d["Value"]
        //set dimensions
        let dimensions = {
            width: screen.width - 500,
            height: 60,
            margin: {top: 0, right: 40, bottom: 20, left: 40,},
        }

        //draw canvas
        const wrapper = d3.select("#navbar-charts-standard-3-digital")
            .append("svg")
            .attr("viewBox", "0 0 " + dimensions.width + " " + dimensions.height)
        const bounds = wrapper.append("g")

        //set scales
        const chartStartsAt = dataset["Data"][0]["Date"]
        const chartEndsAt = dataset["Data"][dataset["Data"].length - 1]["Date"]
        const xScale = d3.scaleTime()
            .domain([chartStartsAt, chartEndsAt])
            .range([dimensions.margin.left, dimensions.width - dimensions.margin.right])
        const yScale = d3.scaleLinear()
            .domain(d3.extent(dataset["Data"], yAccessor))
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
        let div = d3.select("#navbar-charts-standard-3-digital").append("div")
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
            .attr("d", productionAreaGenerator(dataset["Data"]))
            .attr("fill", 'steelblue')
            .on('mousemove', (event) => {
                    let coords = d3.pointer(event);
                    let timeEntered = timeScale.invert(coords[0]) / 1000
                    let now = new Date(timeEntered * 1000).toLocaleString()
                    let start = new Date(dataset["Data"].filter(i => i["Date"] < timeEntered).pop()["Date"] * 1000).toLocaleString()
                    let end = new Date(dataset["Data"].filter(i => i["Date"] > timeEntered)[0]["Date"] * 1000).toLocaleString()
                    div.html(now + "<br/>" + dataset["Name"].toUpperCase() + "<br/>" + start + "<br/>" + end)
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
    }
}

function drawCharts(data) {
    const productionDataset = data["TimelineDataList"][0]
    const downtimeDataset = data["TimelineDataList"][1]
    const powerOffDataset = data["TimelineDataList"][2]
    const digitalDataset = data["DigitalDataList"]
    const analogDataset = data["AnalogDataList"]
    DrawTimeLine(productionDataset["Data"], downtimeDataset["Data"], powerOffDataset["Data"], productionDataset["Color"], downtimeDataset["Color"], powerOffDataset["Color"], productionDataset["Name"], downtimeDataset["Name"], powerOffDataset["Name"]);
    DrawTimeLinePanel(productionDataset["Data"], downtimeDataset["Data"], powerOffDataset["Data"], productionDataset["Color"], downtimeDataset["Color"], powerOffDataset["Color"], productionDataset["Name"], downtimeDataset["Name"], powerOffDataset["Name"], digitalDataset);
    DrawDigital(digitalDataset)
    document.getElementById("charts-select-workplace").disabled = false
    document.getElementById("charts-standard-start").disabled = false
    document.getElementById("charts-standard-end").disabled = false
    // DrawAnalog(analogDataset)
}


function DrawDigitalDataToChart(name, dataset) {
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
    const wrapper = d3.select("#navbar-charts-standard-3-digital")
        .append("svg")
        .attr("viewBox", "0 0 " + dimensions.width + " " + dimensions.height)
    const bounds = wrapper.append("g")


    //set scales
    const chartStartsAt = dataset[0]["Date"]
    const chartEndsAt = dataset[dataset.length - 1]["Date"]
    const xScale = d3.scaleTime()
        .domain([chartStartsAt, chartEndsAt])
        .range([dimensions.margin.left, dimensions.width - dimensions.margin.right])
    const yScale = d3.scaleLinear()
        .domain(d3.extent(dataset, yAccessor))
        .range([dimensions.height - dimensions.margin.bottom, 0])


    // prepare data
    const productionAreaGenerator = d3.area()
        .x(d => xScale(xAccessor(d)))
        .y0(dimensions.height - dimensions.margin.bottom)
        .y1(d => yScale(yAccessor(d)))
        .curve(d3.curveStepAfter);

    // prepare tooltip
    let div = d3.select("#navbar-charts-standard-3-digital").append("div")
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
        .attr("d", productionAreaGenerator(dataset))
        .attr("fill", "grey")
}

function DrawDigitalChart(digitalDataList) {
    for (const digitalData of digitalDataList) {
        const dataset = digitalData["Data"]
        DrawDigitalDataToChart(digitalData["Name"], dataset)
    }
}