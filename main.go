package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/fonts/*filepath", http.Dir("fonts"))
	router.GET("/", system)
	router.POST("/check_login", checkLogin)
	_ = http.ListenAndServe(":80", router)
}

func system(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "./html/system.html")
}
