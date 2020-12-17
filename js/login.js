const sessionLogin = sessionStorage.getItem("user")
const navbarLogin = document.getElementById("navbar")
const login = document.getElementById("login")
const userEmail = document.getElementById("user-email")
const userPassword = document.getElementById("user-password")
const userLoginButton = document.getElementById("user-login-button")
const navbarMenu = document.getElementById("navbar-menu")
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
    if (userEmail.value === "email@email.com" && userPassword.value === "pa33w0rd") {
        document.getElementById('canvas').classList.remove("none");
        document.getElementById('content').classList.add("hidden");
        login.classList.add("hidden")
        setTimeout(displayMatrix, 500)
        return;
    }

    let data = {
        UserEmail: userEmail.value,
        UserPassword: userPassword.value,
        UserSessionStorage: sessionLogin,
    };
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
                document.getElementById('canvas').classList.add("none");
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


function displayMatrix() {
    const canvas = document.getElementById('canvas');
    canvas.height = window.screen.height;
    canvas.width = window.screen.width;
    canvas.getContext('2d').translate(canvas.width, 0);
    canvas.getContext('2d').scale(-1, 1);
    let columns = [];
    for (let i = 0; i < 256; columns[i++] = 1) ;

    function step() {
        canvas.getContext('2d').fillStyle = 'rgba(0,0,0,0.05)';
        canvas.getContext('2d').fillRect(0, 0, canvas.width, canvas.height);
        canvas.getContext('2d').fillStyle = '#0F0';
        columns.map(function (value, index) {
                let character = String.fromCharCode(12448 + Math.random() * 96);
                canvas.getContext('2d').fillText(character, index * 10, value);
                columns[index] = value > 758 + Math.random() * 1e4 ? 0 : value + 10
            }
        )
    }

    setInterval(step, 33)
}