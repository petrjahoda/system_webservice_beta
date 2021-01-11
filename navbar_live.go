package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

type LiveDataInput struct {
	Input     string
	Selection string
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
	Production []WorkplaceData
	Downtime   []WorkplaceData
	Poweroff   []WorkplaceData
	Result     string
}

type CalendarDataOutput struct {
	Data []CalendarData
}

type CalendarData struct {
	Date            string
	ProductionValue int
}

type WorkplaceData struct {
	WorkplaceName string
	State         string
	StateDuration string
}

type SelectionDataOutput struct {
	Result        string
	SelectionData []string
}

type WorkplaceDataOutput struct {
	Result            string
	Order             string
	OrderDuration     string
	User              string
	Downtime          string
	DowntimeDuration  string
	Breakdown         string
	BreakdownDuration string
	Alarm             string
	AlarmDuration     string
}

type CompanyOutput struct {
	Result      string
	CompanyName string
}

func getCalendarData(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing calendar data started for "+data.Input+" and "+data.Selection)
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

func getLiveProductivityData(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok: " + err.Error()
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

func getLiveOverviewData(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing live overview data started for "+data.Input+" and "+data.Selection)
	//todo: process real live data from database
	var productionWorkplaces []WorkplaceData
	productionWorkplaces = append(productionWorkplaces, WorkplaceData{WorkplaceName: "CNC-1", State: "production", StateDuration: "4h 1m"})
	productionWorkplaces = append(productionWorkplaces, WorkplaceData{WorkplaceName: "CNC-4", State: "production", StateDuration: "4h 0m"})
	productionWorkplaces = append(productionWorkplaces, WorkplaceData{WorkplaceName: "CNC-5", State: "production", StateDuration: "3h 51m"})
	productionWorkplaces = append(productionWorkplaces, WorkplaceData{WorkplaceName: "CNC-9", State: "production", StateDuration: "3h 43m"})
	productionWorkplaces = append(productionWorkplaces, WorkplaceData{WorkplaceName: "CNC-7", State: "production", StateDuration: "2h 43m"})
	productionWorkplaces = append(productionWorkplaces, WorkplaceData{WorkplaceName: "CNC-8", State: "production", StateDuration: "1h 17m"})

	var downtimeWorkplaces []WorkplaceData
	downtimeWorkplaces = append(downtimeWorkplaces, WorkplaceData{WorkplaceName: "MATSUSHITA 620C LEFT ORDERED PRODUCTION", State: "downtime", StateDuration: "12 h1m"})
	downtimeWorkplaces = append(downtimeWorkplaces, WorkplaceData{WorkplaceName: "MATSUSHITA 620D", State: "downtime", StateDuration: "8h 12m"})
	downtimeWorkplaces = append(downtimeWorkplaces, WorkplaceData{WorkplaceName: "CNC-2", State: "downtime", StateDuration: "3h 11m"})

	var poweroffWorkplaces []WorkplaceData
	poweroffWorkplaces = append(poweroffWorkplaces, WorkplaceData{WorkplaceName: "CNC-3", State: "poweroff", StateDuration: "3d 11h 1m"})
	poweroffWorkplaces = append(poweroffWorkplaces, WorkplaceData{WorkplaceName: "CNC-6", State: "poweroff", StateDuration: "12h 12m"})

	var outputData LiveOverviewDataOutput
	outputData.Result = "ok"
	outputData.Production = productionWorkplaces
	outputData.Downtime = downtimeWorkplaces
	outputData.Poweroff = poweroffWorkplaces
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}

func getLiveSelection(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing selection started for "+data.Input)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData SelectionDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Verifying user ended")
		return
	}
	var selectionData []string
	if data.Input == "group" {
		var groups []database.WorkplaceSection
		db.Find(&groups)
		for _, group := range groups {
			selectionData = append(selectionData, group.Name)
		}
	}
	if data.Input == "workplace" {
		var workplaces []database.Workplace
		db.Find(&workplaces)
		for _, workplace := range workplaces {
			selectionData = append(selectionData, workplace.Name)
		}
	}
	sort.Strings(selectionData)
	var outputData SelectionDataOutput
	outputData.Result = "ok"
	outputData.SelectionData = selectionData
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}
func getCompanyName(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Downloading company name")

	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData CompanyOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Verifying user ended")
		return
	}
	var settings database.Setting
	db.Where("name = 'company'").First(&settings)
	var outputData CompanyOutput
	outputData.Result = "ok"
	outputData.CompanyName = settings.Value
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}
func getWorkplaceData(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData WorkplaceDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing workplace data for "+data.Input)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData WorkplaceDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Verifying user ended")
		return
	}

	var workplace database.Workplace
	db.Where("Name = ?", data.Input).Find(&workplace)
	var orderRecord database.OrderRecord
	db.Where("workplace_id = ?", workplace.ID).Where("date_time_end is null").Find(&orderRecord)
	logInfo("MAIN", orderRecord.DateTimeStart.String())
	var order database.Order
	db.Where("id = ?", orderRecord.OrderID).Find(&order)
	logInfo("MAIN", "Order: "+order.Name)

	var userRecord database.UserRecord
	db.Where("workplace_id = ?", workplace.ID).Where("date_time_end is null").Where("order_record_id = ?", orderRecord.ID).Find(&userRecord)
	var user database.User
	db.Where("id = ?", userRecord.UserID).Find(&user)
	logInfo("MAIN", "User: "+user.FirstName+" "+user.SecondName)

	var downtimeRecord database.DowntimeRecord
	db.Where("workplace_id = ?", workplace.ID).Where("date_time_end is null").Find(&downtimeRecord)
	var downtime database.Downtime
	db.Where("id = ?", downtimeRecord.DowntimeID).Find(&downtime)
	logInfo("MAIN", "Downtime: "+downtime.Name)

	var breakdownRecord database.BreakdownRecord
	db.Where("workplace_id = ?", workplace.ID).Where("date_time_end is null").Find(&breakdownRecord)
	var breakdown database.Breakdown
	db.Where("id = ?", breakdownRecord.BreakdownID).Find(&breakdown)
	logInfo("MAIN", "Breakdown: "+breakdown.Name)

	var alarmRecord database.AlarmRecord
	db.Where("workplace_id = ?", workplace.ID).Where("date_time_end is null").Find(&alarmRecord)
	var alarm database.Alarm
	db.Where("id = ?", alarmRecord.AlarmID).Find(&alarm)
	logInfo("MAIN", "Alarm: "+alarm.Name)

	var outputData WorkplaceDataOutput
	outputData.Result = "ok"
	outputData.Order = order.Name
	outputData.OrderDuration = time.Now().Sub(orderRecord.DateTimeStart).Round(1 * time.Second).String()
	outputData.User = user.FirstName + " " + user.SecondName
	outputData.Downtime = downtime.Name
	outputData.DowntimeDuration = time.Now().Sub(downtimeRecord.DateTimeStart).Round(1 * time.Second).String()
	outputData.Breakdown = breakdown.Name
	outputData.BreakdownDuration = time.Now().Sub(breakdownRecord.DateTimeStart).Round(1 * time.Second).String()
	outputData.Alarm = alarm.Name
	outputData.AlarmDuration = time.Now().Sub(alarmRecord.DateTimeStart).Round(1 * time.Second).String()
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}
