const content = document.getElementById("content");
const navbar = document.getElementById("navbar");

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
        fetch("/get_content", {
            method: "POST",
            body: JSON.stringify(data)
        }).then((response) => {
            response.text().then(function (data) {
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
            });
        }).catch((error) => {
            console.log(error)
        });
    }
})


content.addEventListener("click", (event) => {
    let contentData = null
    if (event.target.id !== "navbar-menu") {
        const menuItems = document.getElementsByClassName("navbar-menu-item")
        for (let i = 0; i < menuItems.length; i++) {
            if (event.target.id === menuItems[i].id) {
                menuItems[i].classList.add("underlined")
                contentData = menuItems[i].id
            } else {
                menuItems[i].classList.remove("underlined")
            }
        }
    }
    console.log("Downloading data content for " + contentData)
    let data = {Content: contentData};
    fetch("/get_content_data", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            const contentData = document.getElementById("content-data");
            contentData.innerHTML = result.Html
        });
    }).catch((error) => {
        console.log(error)
    });
})



