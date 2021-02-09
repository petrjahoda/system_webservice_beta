const navbar = document.getElementById("navbar");

const content = document.getElementById("content")

navbar.addEventListener("click", (event) => {
    console.log("mouse on main menu")
    if (event.target.id !== "navbar") {
        const menuItems = document.getElementsByClassName("navbar-item")
        for (let i = 0; i < menuItems.length; i++) {
            if (menuItems[i].id === event.target.id) {
                menuItems[i].classList.add("underlined")
            } else {
                menuItems[i].classList.remove("underlined")
            }
        }
        console.log("Downloading data for " + event.target.id)
        let data = {Content: event.target.id};
        fetch("/get_menu", {
            method: "POST",
            body: JSON.stringify(data)
        }).then((response) => {
            response.text().then(function (data) {
                let result = JSON.parse(data);
                navbarMenu.innerHTML = result.Html
                if (result.MenuLocales !== null) {
                    for (menuLocale of result.MenuLocales) {
                        const menuItem = document.getElementById(menuLocale.Name)
                        try {
                            menuItem.textContent = menuLocale[sessionStorage.getItem("locale")]
                        } catch (e) {
                            console.log(menuLocale.Name + " not loaded for actual page")
                        }
                    }
                } else {
                    console.log("no data")
                }
            });

        }).catch((error) => {
            console.log(error)
        });
    }
})

navbarMenu.addEventListener("click", (event) => {
    let menuData = null
    if (event.target.id !== "navbar-menu") {
        const menuItems = document.getElementsByClassName("navbar-menu-item")
        for (let i = 0; i < menuItems.length; i++) {
            if (event.target.id === menuItems[i].id) {
                menuItems[i].classList.add("underlined")
                menuData = menuItems[i].id
            } else {
                menuItems[i].classList.remove("underlined")
            }
        }
        console.log("Downloading data content for " + menuData)
        let data = {Content: menuData,};
        fetch("/get_content", {
            method: "POST",
            body: JSON.stringify(data)
        }).then((response) => {
            response.text().then(function (data) {
                const content = document.getElementById("content")
                let result = JSON.parse(data);
                content.innerHTML = result.Html
                if (result.MenuLocales !== null) {
                    for (menuLocale of result.MenuLocales) {
                        const menuItem = document.getElementById(menuLocale.Name)
                        try {
                            menuItem.textContent = menuLocale[sessionStorage.getItem("locale")]
                        } catch (e) {
                            console.log(menuLocale.Name + " not loaded for actual page")
                        }
                    }
                } else {
                    console.log("no data")
                }
                Process(menuData)
            });

        }).catch((error) => {
            console.log(error)
        });
    }
})


function displayAllStandardCharts(chartsStart, chartsEnd) {
    let workplace = sessionStorage.getItem("workplaceName")
    document.getElementById("navbar-charts-standard-2-timeline").innerHTML = ""
    document.getElementById("navbar-charts-standard-2-timeline-panel").innerHTML = ""
    document.getElementById("navbar-charts-standard-3-digital").innerHTML = ""
    displayAllCharts(chartsStart.value, chartsEnd.value, workplace)
}

function Process(menuData) {
    if (menuData.includes("navbar-live-company")) {
        console.log("Processing page " + menuData)
        displayCompanyName("company")
        displayLiveProductivity("company", "")
        displayCalendar("company", "")
        displayOverview("company", "");
    }
    if (menuData.includes("navbar-live-group")) {
        console.log("Processing page " + menuData)
        displaySelectionLive("group", "groupName")
    }
    if (menuData.includes("navbar-live-workplace")) {
        console.log("Processing page " + menuData)
        displaySelectionLive("workplace", "workplaceName")
    }
    if (menuData.includes("navbar-charts-standard")) {
        console.log("Processing page " + menuData)
        displaySelectionCharts("workplace", "workplaceName")

        let chartsStart = document.getElementById("charts-standard-start")
        let chartsEnd = document.getElementById("charts-standard-end")
        let workplaceSelect = document.getElementById("charts-select-workplace")

        if (sessionStorage.getItem("chartsStart")!== null) {
            chartsStart.value = sessionStorage.getItem("chartsStart")
        }
        if (sessionStorage.getItem("chartsEnd")!== null) {
            chartsEnd.value = sessionStorage.getItem("chartsEnd")
        }

        if (sessionStorage.getItem("chartsStart")!== null && sessionStorage.getItem("chartsEnd")!== null) {
            console.log("Displaying charts right away")
            if (chartsStart.value < chartsEnd.value) {
                chartsStart.style.border = "1px solid black"
                chartsEnd.style.border = "1px solid black"
                displayAllStandardCharts(chartsStart, chartsEnd);

            } else {
                chartsStart.style.border = "1px solid lightcoral"
                chartsEnd.style.border = "1px solid lightcoral"
            }

        }

        workplaceSelect.addEventListener("change", (event) => {
            sessionStorage.setItem("workplaceName", event.target.value)
            if (chartsStart.value < chartsEnd.value) {
                chartsStart.style.border = "1px solid black"
                chartsEnd.style.border = "1px solid black"
                displayAllStandardCharts(chartsStart, chartsEnd);

            } else {
                chartsStart.style.border = "1px solid lightcoral"
                chartsEnd.style.border = "1px solid lightcoral"
            }
        })

        chartsStart.addEventListener("change", (event) => {
            sessionStorage.setItem("chartsStart", chartsStart.value)
            if (chartsEnd.value.length>0) {
                if (chartsStart.value < chartsEnd.value) {
                    chartsStart.style.border = "1px solid black"
                    chartsEnd.style.border = "1px solid black"
                    displayAllStandardCharts(chartsStart, chartsEnd);

                } else {
                    chartsStart.style.border = "1px solid lightcoral"
                    chartsEnd.style.border = "1px solid lightcoral"
                }
            }
        })

        chartsEnd.addEventListener("change", (event) => {
            sessionStorage.setItem("chartsEnd", chartsEnd.value)
            if (chartsStart.value.length>0) {
                if (chartsStart.value < chartsEnd.value) {
                    chartsStart.style.border = "1px solid black"
                    chartsEnd.style.border = "1px solid black"
                    displayAllStandardCharts(chartsStart, chartsEnd);

                } else {
                    chartsStart.style.border = "1px solid lightcoral"
                    chartsEnd.style.border = "1px solid lightcoral"
                }
            }
        })



    }
    if (menuData.includes("navbar-charts-special")) {
        console.log("Processing page " + menuData)
        displaySelectionCharts("workplace", "workplaceName")
        let chartsStart = document.getElementById("charts-special-start")
        let chartsEnd = document.getElementById("charts-special-end")
        let workplace = document.getElementById("charts-select-workplace")
        if (sessionStorage.getItem("chartsStart")!== null) {
            chartsStart.value = sessionStorage.getItem("chartsStart")
        }
        if (sessionStorage.getItem("chartsEnd")!== null) {
            chartsEnd.value = sessionStorage.getItem("chartsEnd")
        }
        if (sessionStorage.getItem("chartsStart")!== null && sessionStorage.getItem("chartsEnd")!== null) {
            console.log("Displaying charts right away")
            if (chartsStart.value < chartsEnd.value) {
                chartsStart.style.border = "1px solid black"
                chartsEnd.style.border = "1px solid black"
                displayProductionRateChart(chartsStart.value, chartsEnd.value, workplace)
            } else {
                chartsStart.style.border = "1px solid lightcoral"
                chartsEnd.style.border = "1px solid lightcoral"
            }
        }
        chartsStart.addEventListener("change", (event) => {
            sessionStorage.setItem("chartsStart", chartsStart.value)
            if (chartsEnd.value.length>0) {
                if (chartsStart.value < chartsEnd.value) {
                    chartsStart.style.border = "1px solid black"
                    chartsEnd.style.border = "1px solid black"
                    displayProductionRateChart(chartsStart.value, chartsEnd.value, workplace)
                } else {
                    chartsStart.style.border = "1px solid lightcoral"
                    chartsEnd.style.border = "1px solid lightcoral"

                }
            }
        })
        chartsEnd.addEventListener("change", (event) => {
            sessionStorage.setItem("chartsEnd", chartsEnd.value)
            if (chartsStart.value.length>0) {
                if (chartsStart.value < chartsEnd.value) {
                    chartsStart.style.border = "1px solid black"
                    chartsEnd.style.border = "1px solid black"
                    displayProductionRateChart(chartsStart.value, chartsEnd.value, workplace)

                } else {
                    chartsStart.style.border = "1px solid lightcoral"
                    chartsEnd.style.border = "1px solid lightcoral"
                }
            }
        })

    }
}



