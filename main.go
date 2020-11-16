package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/data/*filepath", http.Dir("data"))
	router.ServeFiles("/fonts/*filepath", http.Dir("fonts"))
	router.ServeFiles("/mif/*filepath", http.Dir("mif"))
	router.GET("/", live)
	router.GET("/live", live)
	router.GET("/charts", charts)
	router.GET("/statistics", statistics)
	router.GET("/data", data)
	router.GET("/settings", settings)
	router.GET("/login", login)
	router.POST("/check_login", checkLogin)
	_ = http.ListenAndServe(":80", router)
}



func live(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "./html/live.html")
}

func charts(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "./html/charts.html")
}
func statistics(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "./html/statistics.html")
}
func data(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "./html/data.html")
}
func settings(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "./html/settings.html")
}

func login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "./html/login.html")
}

