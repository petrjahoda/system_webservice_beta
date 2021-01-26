package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type ChartDataInput struct {
	ChartsStart string
	ChartsEnd   string
	Workplace   string
}

type TimelineChartData struct {
	Result         string
	ProductionData []TimelineData
	DowntimeData   []TimelineData
	PowerOffData   []TimelineData
}

type TimelineData struct {
	Date  int64
	Value int
}

func getTimelineChart(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data ChartDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData TimelineChartData
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing timeline charts data for "+data.Workplace+" from "+data.ChartsStart+" to "+data.ChartsEnd)
	recordsFrom, err := time.Parse("2006-01-02T15:04", data.ChartsStart)
	if err != nil {
		logError("MAIN", "Problem parsing date: "+data.ChartsStart)
		var responseData TimelineChartData
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Verifying user ended")
		return
	}
	recordsTo, err := time.Parse("2006-01-02T15:04", data.ChartsEnd)
	if err != nil {
		logError("MAIN", "Problem parsing date: "+data.ChartsEnd)
		var responseData TimelineChartData
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Verifying user ended")
		return
	}
	logInfo("MAIN", "Processing timeline charts data for "+data.Workplace+" from "+strconv.Itoa(int(recordsFrom.Unix()))+" to "+strconv.Itoa(int(recordsTo.Unix())))
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData TimelineChartData
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Verifying user ended")
		return
	}
	productionData, downtimeData, powerOffData := downloadTimelineData(db, data, recordsFrom, recordsTo)
	logInfo("MAIN", "There are "+strconv.Itoa(len(productionData))+" production records")
	logInfo("MAIN", "There are "+strconv.Itoa(len(downtimeData))+" downtime records")
	logInfo("MAIN", "There are "+strconv.Itoa(len(powerOffData))+" poweroff records")
	for _, datum := range productionData {
		fmt.Println(datum.Date, datum.Value)
	}
	fmt.Println()
	for _, datum := range downtimeData {
		fmt.Println(datum.Date, datum.Value)
	}
	fmt.Println()
	for _, datum := range powerOffData {
		fmt.Println(datum.Date, datum.Value)
	}
	var outputData TimelineChartData
	outputData.Result = "ok"
	outputData.ProductionData = productionData
	outputData.DowntimeData = downtimeData
	outputData.PowerOffData = powerOffData
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}

func downloadTimelineData(db *gorm.DB, data ChartDataInput, recordsFrom time.Time, recordsTo time.Time) ([]TimelineData, []TimelineData, []TimelineData) {
	var workplace database.Workplace
	db.Where("Name = ?", data.Workplace).Find(&workplace)
	var lastStateRecord database.StateRecord
	db.Where("workplace_id = ?", workplace.ID).Where("date_time_start <= ?", recordsFrom).Last(&lastStateRecord)
	var allStateRecords []database.StateRecord
	db.Where("workplace_id = ?", workplace.ID).Where("date_time_start >= ?", recordsFrom).Where("date_time_start < ?", recordsTo).Order("date_time_start ASC").Find(&allStateRecords)
	var productionData []TimelineData
	var downtimeData []TimelineData
	var powerOffData []TimelineData
	previousRecordStateId := lastStateRecord.StateID
	lastStateIdToInsert := lastStateRecord.StateID
	productionData, downtimeData, powerOffData = addFirstRecordToData(lastStateRecord, productionData, recordsFrom, downtimeData, powerOffData)
	lastStateIdToInsert, productionData, downtimeData, powerOffData = addAllRecordsToData(allStateRecords, previousRecordStateId, productionData, lastStateIdToInsert, downtimeData, powerOffData)
	productionData, downtimeData, powerOffData = addLastRecordToData(lastStateIdToInsert, productionData, recordsTo, downtimeData, powerOffData)
	return productionData, downtimeData, powerOffData
}

func addAllRecordsToData(allStateRecords []database.StateRecord, previousRecordStateId int, productionData []TimelineData, lastStateIdToInsert int, downtimeData []TimelineData, powerOffData []TimelineData) (int, []TimelineData, []TimelineData, []TimelineData) {
	for _, record := range allStateRecords {
		if record.StateID == previousRecordStateId {
			continue
		}
		if record.StateID == 1 {
			productionData = append(productionData, TimelineData{
				Date:  record.DateTimeStart.Unix(),
				Value: 1,
			})
			lastStateIdToInsert = 1
		} else if record.StateID == 2 {
			downtimeData = append(downtimeData, TimelineData{
				Date:  record.DateTimeStart.Unix(),
				Value: 1,
			})
			lastStateIdToInsert = 2
		} else if record.StateID == 3 {
			powerOffData = append(powerOffData, TimelineData{
				Date:  record.DateTimeStart.Unix(),
				Value: 1,
			})
			lastStateIdToInsert = 3
		}
		if previousRecordStateId == 1 {
			productionData = append(productionData, TimelineData{
				Date:  record.DateTimeStart.Unix(),
				Value: 0,
			})

		} else if previousRecordStateId == 2 {
			downtimeData = append(downtimeData, TimelineData{
				Date:  record.DateTimeStart.Unix(),
				Value: 0,
			})
		} else if previousRecordStateId == 3 {
			powerOffData = append(powerOffData, TimelineData{
				Date:  record.DateTimeStart.Unix(),
				Value: 0,
			})
		}
		previousRecordStateId = record.StateID
	}
	return lastStateIdToInsert, productionData, downtimeData, powerOffData
}

func addLastRecordToData(lastStateIdToInsert int, productionData []TimelineData, recordsTo time.Time, downtimeData []TimelineData, powerOffData []TimelineData) ([]TimelineData, []TimelineData, []TimelineData) {
	if lastStateIdToInsert == 1 {
		productionData = append(productionData, TimelineData{
			Date:  recordsTo.Unix(),
			Value: 0,
		})
		downtimeData = append(downtimeData, TimelineData{
			Date:  recordsTo.Unix(),
			Value: 1,
		})
		powerOffData = append(powerOffData, TimelineData{
			Date:  recordsTo.Unix(),
			Value: 1,
		})
	} else if lastStateIdToInsert == 2 {
		downtimeData = append(downtimeData, TimelineData{
			Date:  recordsTo.Unix(),
			Value: 0,
		})
		productionData = append(productionData, TimelineData{
			Date:  recordsTo.Unix(),
			Value: 1,
		})
		powerOffData = append(powerOffData, TimelineData{
			Date:  recordsTo.Unix(),
			Value: 1,
		})
	} else if lastStateIdToInsert == 3 {
		powerOffData = append(powerOffData, TimelineData{
			Date:  recordsTo.Unix(),
			Value: 0,
		})
		productionData = append(productionData, TimelineData{
			Date:  recordsTo.Unix(),
			Value: 1,
		})
		downtimeData = append(downtimeData, TimelineData{
			Date:  recordsTo.Unix(),
			Value: 1,
		})
	}
	return productionData, downtimeData, powerOffData
}

func addFirstRecordToData(lastStateRecord database.StateRecord, productionData []TimelineData, recordsFrom time.Time, downtimeData []TimelineData, powerOffData []TimelineData) ([]TimelineData, []TimelineData, []TimelineData) {
	if lastStateRecord.StateID == 1 {
		productionData = append(productionData, TimelineData{
			Date:  recordsFrom.Unix(),
			Value: 1,
		})
		downtimeData = append(downtimeData, TimelineData{
			Date:  recordsFrom.Unix(),
			Value: 0,
		})
		powerOffData = append(powerOffData, TimelineData{
			Date:  recordsFrom.Unix(),
			Value: 0,
		})

	} else if lastStateRecord.StateID == 2 {
		downtimeData = append(downtimeData, TimelineData{
			Date:  recordsFrom.Unix(),
			Value: 1,
		})
		productionData = append(productionData, TimelineData{
			Date:  recordsFrom.Unix(),
			Value: 0,
		})
		powerOffData = append(powerOffData, TimelineData{
			Date:  recordsFrom.Unix(),
			Value: 0,
		})
	} else if lastStateRecord.StateID == 3 {
		powerOffData = append(powerOffData, TimelineData{
			Date:  recordsFrom.Unix(),
			Value: 1,
		})
		productionData = append(productionData, TimelineData{
			Date:  recordsFrom.Unix(),
			Value: 0,
		})
		downtimeData = append(downtimeData, TimelineData{
			Date:  recordsFrom.Unix(),
			Value: 0,
		})
	}
	return productionData, downtimeData, powerOffData
}
