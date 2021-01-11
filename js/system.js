const navbar = document.getElementById("navbar");

const content = document.getElementById("content")

navbar.addEventListener("click", (event) => {
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
                    console.log(result.MenuLocales)
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
        displaySelection("group")
        let savedSelection = sessionStorage.getItem("groupName")
        displayLiveProductivity("group", savedSelection)
        displayCalendar("group", savedSelection)
        displayOverview("group", savedSelection);
    }
    if (menuData.includes("navbar-live-workplace")) {
        console.log("Processing page " + menuData)
        displaySelection("workplace")
        let savedSelection = sessionStorage.getItem("workplaceName")
        displayLiveProductivity("workplace", savedSelection)
        displayCalendar("workplace", savedSelection)
        displayWorkplaceData("workplace", savedSelection)
    }
}



