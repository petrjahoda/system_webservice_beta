package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"math/rand"
	"net/http"
	"time"
)

type LiveDataInput struct {
	Input string
}

type LiveDataOutput struct {
	Yesterday string
	Today     string
	LastWeek  string
	ThisWeek  string
	LastMonth string
	ThisMonth string
	Result    string
}

type LiveOverviewDataOutput struct {
	Production string
	Downtime   string
	Poweroff   string
	Result     string
}

type CalendarDataOutput struct {
	Data []CalendarData
}

type CalendarData struct {
	Date            string
	ProductionValue int
}

type BestWorstPoweroffOutputData struct {
	Data   []WorkplaceData
	Result string
}

type WorkplaceData struct {
	WorkplaceName       string
	WorkplaceProduction string
}

func getCalendarData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing calendar data started for "+data.Input)
	//todo: process real calendar data from database

	var calendarData []CalendarData
	for i := 0; i < 365; i++ {
		var oneCalendarData CalendarData
		oneCalendarData.Date = time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		randomProduction := rand.Intn(100-0) + 0
		oneCalendarData.ProductionValue = randomProduction
		calendarData = append(calendarData, oneCalendarData)
	}
	var outputData CalendarDataOutput
	outputData.Data = calendarData
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}

func getLiveProductivityData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing live productivity data started for "+data.Input)
	//todo: process real live data from database
	var outputData LiveDataOutput
	outputData.Result = "ok"
	outputData.Today = "67.5"
	outputData.Yesterday = "17.5"
	outputData.LastWeek = "15.3"
	outputData.ThisWeek = "12.7"
	outputData.LastMonth = "6.1"
	outputData.ThisMonth = "43.4"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}

func getLiveOverviewData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing live overview data started for "+data.Input)
	//todo: process real live data from database
	var outputData LiveOverviewDataOutput
	outputData.Result = "ok"
	outputData.Production = "999"
	outputData.Downtime = "17"
	outputData.Poweroff = "4"
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}

func getLiveBestData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing live best data started for "+data.Input)
	//todo: process real live data from database

	var workplacesData []WorkplaceData
	var first WorkplaceData
	first.WorkplaceName = "CNC-7"
	first.WorkplaceProduction = "99.7"
	workplacesData = append(workplacesData, first)
	var second WorkplaceData
	second.WorkplaceName = "MATSUSHITA 620C LEFT ORDERED PRODUCTION"
	second.WorkplaceProduction = "97.1"
	workplacesData = append(workplacesData, second)
	var third WorkplaceData
	third.WorkplaceName = "MATSUSHITA 120-I"
	third.WorkplaceProduction = "96.9"
	workplacesData = append(workplacesData, third)
	var outputData BestWorstPoweroffOutputData
	outputData.Result = "ok"
	outputData.Data = workplacesData
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}

func getLiveWorstData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing live worst data started for "+data.Input)
	//todo: process real live data from database

	var workplacesData []WorkplaceData
	var first WorkplaceData
	first.WorkplaceName = "CNC-4"
	first.WorkplaceProduction = "100.0"
	workplacesData = append(workplacesData, first)
	var second WorkplaceData
	second.WorkplaceName = "CNC-5"
	second.WorkplaceProduction = "99.1"
	workplacesData = append(workplacesData, second)
	var third WorkplaceData
	third.WorkplaceName = "MATSUSHITA 120-II"
	third.WorkplaceProduction = "96.4"
	workplacesData = append(workplacesData, third)
	var outputData BestWorstPoweroffOutputData
	outputData.Result = "ok"
	outputData.Data = workplacesData
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}

func getLivePoweroffData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing live poweroff data started for "+data.Input)
	//todo: process real live data from database

	var workplacesData []WorkplaceData
	var first WorkplaceData
	first.WorkplaceName = "CNC-14"
	first.WorkplaceProduction = "24d 3h"
	workplacesData = append(workplacesData, first)
	var second WorkplaceData
	second.WorkplaceName = "CNC-13"
	second.WorkplaceProduction = "22d 12h"
	workplacesData = append(workplacesData, second)
	var third WorkplaceData
	third.WorkplaceName = "CNC-11"
	third.WorkplaceProduction = "6d 23h"
	workplacesData = append(workplacesData, third)
	var outputData BestWorstPoweroffOutputData
	outputData.Result = "ok"
	outputData.Data = workplacesData
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}

func getLiveAllData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing live all data started for "+data.Input)
	//todo: process real live data from database

	var workplacesData []WorkplaceData
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-1", WorkplaceProduction: "production"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-2", WorkplaceProduction: "downtime"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-3", WorkplaceProduction: "production"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-4", WorkplaceProduction: "downtime"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-5", WorkplaceProduction: "downtime"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-6", WorkplaceProduction: "downtime"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-7", WorkplaceProduction: "production"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-8", WorkplaceProduction: "production"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-9", WorkplaceProduction: "production"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-10", WorkplaceProduction: "downtime"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-11", WorkplaceProduction: "poweroff"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-12", WorkplaceProduction: "production"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-13", WorkplaceProduction: "poweroff"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "CNC-14", WorkplaceProduction: "poweroff"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "MATSUSHITA-120-I", WorkplaceProduction: "production"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "MATSUSHITA-120-II", WorkplaceProduction: "downtime"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "MATSUSHITA 620C LEFT ORDERED PRODUCTION", WorkplaceProduction: "production"})
	workplacesData = append(workplacesData, WorkplaceData{WorkplaceName: "MATSUSHITA-620D", WorkplaceProduction: "production"})
	var outputData BestWorstPoweroffOutputData
	outputData.Result = "ok"
	outputData.Data = workplacesData
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}
