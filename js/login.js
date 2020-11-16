const email = document.getElementById("email");
const password = document.getElementById("password");
const login = document.getElementById("login");
const info = document.getElementById("info");

login.addEventListener("click", () => {
    let data = {
        Email: email.value,
        Password: password.value,
    };
    fetch("/check_login", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                console.log(result.Data)
                sessionStorage.setItem("user", result.Data)
                window.location.replace('live');
            } else {
                info.textContent = result.Data+ ", please enter again"
                info.classList.add("uk-text-danger")
            }
        });
    }).catch((error) => {
        console.log(error.toString())
    });
    console.log(email.value)
    console.log(password.value)
})