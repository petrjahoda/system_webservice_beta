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
            chartsStart: chartsStart,
            chartsEnd: chartsEnd,
            workplace: workplace,
        })
    }).then((data) => {
        drawTimelineChart(data, chartsStart, chartsEnd);
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

function drawTimelineChart(data, chartsStart, chartsEnd) {
    // data
    const productionDataset = data["ProductionData"]

    let result = new Map(productionDataset.map(value => [value["Date"], value["Date"]]));
    const downtimeDataset = data["DowntimeData"]
    const powerOffDataset = data["PowerOffData"]
    const xAccessor = d => d["Date"]
    const yAccessor = d => d["Value"]

    //dimensions
    let dimensions = {
        width: 960,
        height: 100,
        // margin: {top: 15, right: 15, bottom: 40, left: 60,},
        margin: {top: 0, right: 0, bottom: 0, left: 0,},
    }
    dimensions.boundedWidth = dimensions.width - dimensions.margin.left - dimensions.margin.right
    dimensions.boundedHeight = dimensions.height - dimensions.margin.top - dimensions.margin.bottom

    //canvas
    const wrapper = d3.select("#navbar-charts-standard-2-timeline")
        .append("svg")
        .attr("width", dimensions.width)
        .attr("height", dimensions.height)
    const bounds = wrapper.append("g")
        .style("transform", `translate(${dimensions.margin.left}px, ${dimensions.margin.top}px)`);

    //scales

    const chartStartsAt = new Date(chartsStart).valueOf()/1000
    const chartEndsAt = new Date(chartsEnd).valueOf()/1000
    const productionxScale = d3.scaleTime()
        .domain([chartStartsAt, chartEndsAt])
        .range([0, dimensions.boundedWidth])
        .nice()
    const productionyScale = d3.scaleLinear()
        .domain(d3.extent(productionDataset, yAccessor))
        .range([dimensions.boundedHeight, 0])
        .nice()

    const downtimexScale = d3.scaleTime()
        .domain([chartStartsAt, chartEndsAt])
        .range([0, dimensions.boundedWidth])
        .nice()
    const downtimeyScale = d3.scaleLinear()
        .domain(d3.extent(downtimeDataset, yAccessor))
        .range([dimensions.boundedHeight, 0])
        .nice()

    const poweroffxScale = d3.scaleTime()
        .domain([chartStartsAt, chartEndsAt])
        .range([0, dimensions.boundedWidth])
        .nice()
    const poweroffyScale = d3.scaleLinear()
        .domain(d3.extent(powerOffDataset, yAccessor))
        .range([dimensions.boundedHeight, 0])
        .nice()




    // draw data
    const productionLineGenerator = d3.line()
        .x(d => productionxScale(xAccessor(d)))
        .y(d => productionyScale(yAccessor(d)))
        // .curve(d3.curveStep);
    const productionAreaGenerator = d3.area()
        .x(d => productionxScale(xAccessor(d)))
        .y0(dimensions.boundedHeight)
        .y1(d => productionyScale(yAccessor(d)))
        .curve(d3.curveStepAfter);

    const downtimeLineGenerator = d3.line()
        .x(d => downtimexScale(xAccessor(d)))
        .y(d => downtimeyScale(yAccessor(d)))
        // .curve(d3.curveStep);
    const downtimeAreaGenerator = d3.area()
        .x(d => downtimexScale(xAccessor(d)))
        .y0(dimensions.boundedHeight)
        .y1(d => downtimeyScale(yAccessor(d)))
        .curve(d3.curveStepAfter);

    const powerOffLineGenerator = d3.line()
        .x(d => poweroffxScale(xAccessor(d)))
        .y(d => poweroffyScale(yAccessor(d)))
        // .curve(d3.curveStep);
    const powerOffAreaGenerator = d3.area()
        .x(d => poweroffxScale(xAccessor(d)))
        .y0(dimensions.boundedHeight)
        .y1(d => poweroffyScale(yAccessor(d)))
        .curve(d3.curveStepAfter);

    //
    bounds.append("path")
        .attr("d", productionAreaGenerator(productionDataset))
        .attr("fill", "green")
        .on("mouseover", function (e) {
            console.log(e)
        })


    bounds.append("path")
        .attr("d", downtimeAreaGenerator(downtimeDataset))
        .attr("fill", "yellow")

    bounds.append("path")
        .attr("d", powerOffAreaGenerator(powerOffDataset))
        .attr("fill", "red")

    bounds.append("path")
        .attr("d", productionLineGenerator(productionDataset))
        .attr("fill", "none")
        .attr("stroke", "black")
        .attr("stroke-width", 1)
    bounds.append("path")
        .attr("d", downtimeLineGenerator(downtimeDataset))
        .attr("fill", "none")
        .attr("stroke", "black")
        .attr("stroke-width", 1)
        .attr("stroke-dasharray", ("3, 3"))
    bounds.append("path")
        .attr("d", powerOffLineGenerator(powerOffDataset))
        .attr("fill", "none")
        .attr("stroke", "black")
        .attr("stroke-width", 1)
        .attr("stroke-dasharray", ("1, 1"))


}