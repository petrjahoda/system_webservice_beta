package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type LoginData struct {
	Email    string
	Password string
}

type LoginResponseData struct {
	Result string
	Data   string
}

func checkLogin(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Checking login")
	var data LoginData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LoginResponseData
		responseData.Result = "nok"
		responseData.Data = "Error parsing data"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	if len(data.Email) == 0 || len(data.Password) == 0 {
		logError("MAIN", "Bad input data")
		var responseData LoginResponseData
		responseData.Result = "nok"
		responseData.Data = "Bad input data"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
		logInfo("MAIN", "Checking login started for " + data.Email + ":" + data.Password)
	var responseData LoginResponseData
	responseData.Result = "ok"
	responseData.Data = "54231hdgs36ss4js"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("MAIN", "Checking login ended")
	return
}
