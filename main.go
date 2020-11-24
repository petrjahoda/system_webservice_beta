package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const config = "user=postgres password=Zps05..... dbname=version3 host=localhost port=5432 sslmode=disable"

func main() {
	router := httprouter.New()
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/fonts/*filepath", http.Dir("fonts"))
	router.GET("/", system)
	router.POST("/get_content", getContent)
	router.POST("/get_content_data", getContentData)
	router.POST("/verify_user", verifyUser)
	_ = http.ListenAndServe(":80", router)
}

func system(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "./html/system.html")
}
