package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
)

type ContentDataInput struct {
	Content string
}

type ContentDataOutput struct {
	Html string
}

func getContentData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var data ContentDataInput
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
		var responseData ContentDataOutput
		responseData.Html = ""
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData ContentDataOutput
		responseData.Html = ""
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Verifying user ended")
		return
	}

	logInfo("MAIN", "Html read from file")
	var responseData ContentDataOutput
	responseData.Html = string(file)
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("MAIN", "Parsing data from page ended")
	return
}
