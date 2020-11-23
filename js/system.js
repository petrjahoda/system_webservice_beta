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
        let data = {Content: event.target.id};
        fetch("/get_content", {
            method: "POST",
            body: JSON.stringify(data)
        }).then((response) => {
            response.text().then(function (data) {
                let result = JSON.parse(data);
                content.innerHTML = result.Html
            });
        }).catch((error) => {
            console.log(error)
        });
    }
})


content.addEventListener("click", (event) => {
    if (event.target.id !== "navbar-menu") {
        const menuItems = document.getElementsByClassName("navbar-menu-item")
        for (let i = 0; i < menuItems.length; i++) {
            if (event.target.id === menuItems[i].id) {
                menuItems[i].classList.add("underlined")
            } else {
                menuItems[i].classList.remove("underlined")
            }
        }
    }
})



