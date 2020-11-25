const sessionLogin = sessionStorage.getItem("user")
const navbarLogin = document.getElementById("navbar")
const navbarMenu = document.getElementById("navbar-menu")
const login = document.getElementById("login")
const userEmail = document.getElementById("user-email")
const userPassword = document.getElementById("user-password")
const userLoginButton = document.getElementById("user-login-button")


if (sessionLogin === null || sessionLogin === "") {
    console.log("session is not logged in")
    navbarLogin.classList.add("hidden")
    navbarMenu.classList.add("hidden")
    login.classList.remove("hidden")
} else {
    console.log("session is logged in")
    verifyUser();
}

userLoginButton.addEventListener("click", () => {
    console.log(userEmail.value)
    console.log(userPassword.value)
    verifyUser();
})

userEmail.addEventListener("keydown", (key) => {
    if (/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(userEmail.value)) {
        userEmail.classList.remove("false")
    } else {
        userEmail.classList.add("false")
    }
    if (userEmail.value === "") {
        userEmail.classList.remove("false")
    }
    if (key.key === "Enter") {
        verifyUser();
    }
})

userPassword.addEventListener("keydown", (key) => {
    if (key.key === "Enter") {
        verifyUser();
    }
})

function verifyUser() {
    let data = {
        UserEmail: userEmail.value,
        UserPassword: userPassword.value,
        UserSessionStorage: sessionLogin,
    };
    console.log(data)
    fetch("/verify_user", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                navbarLogin.classList.remove("hidden")
                navbarMenu.classList.remove("hidden")
                login.classList.add("hidden")
                sessionStorage.setItem("user", result.Data)
                sessionStorage.setItem("locale", result.Locale)
                for (menuLocale of result.MenuLocales) {
                    const menuItem = document.getElementById(menuLocale.Name)
                    try {
                        menuItem.textContent = menuLocale[result.Locale]
                    } catch (e) {
                        console.log(menuLocale.Name + " not loaded for actual page")
                    }
                }
            } else {
                userEmail.classList.add("false")
                userPassword.classList.add("false")
            }
        });
    }).catch((error) => {
        console.log(error)
    });
}