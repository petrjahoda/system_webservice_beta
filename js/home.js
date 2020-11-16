
//
// const mainContainer = document.getElementById("main-container");
// const chartsContainer = document.getElementById("charts-container");
// const statisticsContainer = document.getElementById("statistics-container");
// const dataContainer = document.getElementById("data-container");
// const settingsContainer = document.getElementById("settings-container");
//
//
// liveSelect.addEventListener("click", () => {
//     let data = {
//         ContainerName: "live"
//     };
//     sessionStorage.setItem("container", "live")
//     fetch("/get_container", {
//         method: "POST",
//         body: JSON.stringify(data)
//     }).then((response) => {
//         response.text().then(function (data) {
//             let result = JSON.parse(data);
//             if (result.Result === "ok") {
//                 mainContainer.innerHTML = result.Data
//             } else {
//                 console.log(result.Result)
//             }
//         });
//     }).catch((error) => {
//         console.log(error.toString())
//     });
//
//     liveContainer.hidden = false
//     chartsContainer.hidden = true
//     statisticsContainer.hidden = true
//     dataContainer.hidden = true
//     settingsContainer.hidden = true
//     console.log("chart selected")
// })
//
// chartsSelect.addEventListener("click", () => {
//     let data = {
//         ContainerName: "charts"
//     };
//     sessionStorage.setItem("container", "charts")
//     fetch("/get_container", {
//         method: "POST",
//         body: JSON.stringify(data)
//     }).then((response) => {
//         response.text().then(function (data) {
//             let result = JSON.parse(data);
//             if (result.Result === "ok") {
//                 mainContainer.innerHTML = result.Data
//             } else {
//                 console.log(result.Result)
//             }
//         });
//     }).catch((error) => {
//         console.log(error.toString())
//     });
//     liveContainer.hidden = true
//     chartsContainer.hidden = false
//     statisticsContainer.hidden = true
//     dataContainer.hidden = true
//     settingsContainer.hidden = true
//     console.log("chart selected")
// })
//
//
// statisticsSelect.addEventListener("click", () => {
//     let data = {
//         ContainerName: "statistics"
//     };
//     sessionStorage.setItem("container", "statistics")
//     fetch("/get_container", {
//         method: "POST",
//         body: JSON.stringify(data)
//     }).then((response) => {
//         response.text().then(function (data) {
//             let result = JSON.parse(data);
//             if (result.Result === "ok") {
//                 mainContainer.innerHTML = result.Data
//             } else {
//                 console.log(result.Result)
//             }
//         });
//     }).catch((error) => {
//         console.log(error.toString())
//     });
//     liveContainer.hidden = true
//     chartsContainer.hidden = true
//     statisticsContainer.hidden = false
//     dataContainer.hidden = true
//     settingsContainer.hidden = true
//     console.log("chart selected")
// })
//
// dataSelect.addEventListener("click", () => {
//     let data = {
//         ContainerName: "data"
//     };
//     sessionStorage.setItem("container", "data")
//     fetch("/get_container", {
//         method: "POST",
//         body: JSON.stringify(data)
//     }).then((response) => {
//         response.text().then(function (data) {
//             let result = JSON.parse(data);
//             if (result.Result === "ok") {
//                 mainContainer.innerHTML = result.Data
//             } else {
//                 console.log(result.Result)
//             }
//         });
//     }).catch((error) => {
//         console.log(error.toString())
//     });
//     liveContainer.hidden = true
//     chartsContainer.hidden = true
//     statisticsContainer.hidden = true
//     dataContainer.hidden = false
//     settingsContainer.hidden = true
//     console.log("chart selected")
// })
//
// settingsSelect.addEventListener("click", () => {
//     let data = {
//         ContainerName: "settings"
//     };
//     sessionStorage.setItem("container", "settings")
//     fetch("/get_container", {
//         method: "POST",
//         body: JSON.stringify(data)
//     }).then((response) => {
//         response.text().then(function (data) {
//             let result = JSON.parse(data);
//             if (result.Result === "ok") {
//                 mainContainer.innerHTML = result.Data
//             } else {
//                 console.log(result.Result)
//             }
//         });
//     }).catch((error) => {
//         console.log(error.toString())
//     });
//     liveContainer.hidden = true
//     chartsContainer.hidden = true
//     statisticsContainer.hidden = true
//     dataContainer.hidden = true
//     settingsContainer.hidden = false
//     console.log("chart selected")
// })