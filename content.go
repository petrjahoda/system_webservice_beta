package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

type ContentInput struct {
	Content string
}

type ContentOutput struct {
	Html string
}

func getContent(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var data ContentInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData ContentOutput
		responseData.Html = ""
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	file, err := ioutil.ReadFile("html/" + data.Content + ".html")
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData ContentOutput
		responseData.Html = ""
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	logInfo("MAIN", "Html read from file")
	var responseData ContentOutput
	responseData.Html = string(file)
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("MAIN", "Parsing data from page ended")
	return

}
