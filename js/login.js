const sessionLogin = sessionStorage.getItem("user")
const navbarLogin = document.getElementById("navbar")
const navbarContent = document.getElementById("content")
const login = document.getElementById("login")
const userEmail = document.getElementById("user-email")
const userPassword = document.getElementById("user-password")
const userLoginButton = document.getElementById("user-login-button")


if (sessionLogin === null || sessionLogin === "") {
    console.log("session is not logged in")
    navbarLogin.classList.add("hidden")
    navbarContent.classList.add("hidden")
    login.classList.remove("hidden")
} else {
    console.log("session is logged in")
    verifyUser();
}

function verifyUser() {
    let data = {
        UserEmail: userEmail.value,
        UserPassword: userPassword.value,
        UserSessionStorage: sessionLogin
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
                navbarContent.classList.remove("hidden")
                login.classList.add("hidden")
                sessionStorage.setItem("user", result.Data)
                sessionStorage.setItem("locale", result.Locale)
            }
        });
    }).catch((error) => {
        console.log(error)
    });
}


userEmail.addEventListener("keydown", (event) => {
    if (/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(userEmail.value)) {
        userEmail.classList.remove("false")
    } else {
        userEmail.classList.add("false")
    }
    if (userEmail.value === "") {
        userEmail.classList.remove("false")
    }
})

userLoginButton.addEventListener("click", (event) => {
    console.log(userEmail.value)
    console.log(userPassword.value)
    verifyUser();
})