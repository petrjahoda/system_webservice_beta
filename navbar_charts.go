package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ChartDataInput struct {
	ChartsStart string
	ChartsEnd   string
	Workplace   string
	Locale      string
}

type ChartDataOutput struct {
	Result           string
	TimelineDataList []TimelineDataList
	DigitalDataList  []DigitalDataList
	AnalogDataList   []AnalogDataList
}

type TimelineDataList struct {
	Name  string
	Color string
	Data  []TimelineData
}

type TimelineData struct {
	Date  int64
	Value int
}

type AnalogDataList struct {
	Name string
	Data []AnalogData
}

type AnalogData struct {
	Date  int64
	Value float32
}

type DigitalDataList struct {
	Name string
	Data []DigitalData
}

type DigitalData struct {
	Date  int64
	Value int
}

func processDownloadedDigitalRecords(listOfDigitalRecords [][]database.DevicePortDigitalRecord, db *gorm.DB) []DigitalDataList {
	var responseDigitalData []DigitalDataList
	for _, digitalRecords := range listOfDigitalRecords {
		var portDigitalData DigitalDataList
		var workplacePort database.WorkplacePort
		db.Where("device_port_id = ?", digitalRecords[0].DevicePortID).Find(&workplacePort)
		portDigitalData.Name = workplacePort.Name
		for _, record := range digitalRecords {
			portDigitalData.Data = append(portDigitalData.Data, DigitalData{
				Date:  record.DateTime.Unix(),
				Value: record.Data,
			})
		}
		responseDigitalData = append(responseDigitalData, portDigitalData)
	}
	return responseDigitalData
}

func getDigitalRecords(listOfDigitalDevicePorts []database.DevicePort, db *gorm.DB, recordsFrom time.Time, recordsTo time.Time) [][]database.DevicePortDigitalRecord {
	var listOfDigitalRecords [][]database.DevicePortDigitalRecord
	for _, port := range listOfDigitalDevicePorts {
		var digitalRecords []database.DevicePortDigitalRecord
		db.Where("device_port_id = ?", port.ID).Where("date_time >= ?", recordsFrom).Where("date_time <= ?", recordsTo).Order("date_time").Find(&digitalRecords)
		if len(digitalRecords) > 0 {
			digitalRecords = prependFirstRecord(digitalRecords, recordsFrom, int(port.ID))
			digitalRecords = appendLastRecord(recordsTo, digitalRecords, int(port.ID))
			listOfDigitalRecords = append(listOfDigitalRecords, digitalRecords)
		}
	}
	return listOfDigitalRecords
}

func appendLastRecord(recordsTo time.Time, digitalRecords []database.DevicePortDigitalRecord, portId int) []database.DevicePortDigitalRecord {
	var lastRecord database.DevicePortDigitalRecord
	lastRecord.DateTime = recordsTo
	lastRecord.Data = 0
	lastRecord.DevicePortID = portId
	digitalRecords = append(digitalRecords, lastRecord)
	return digitalRecords
}

func prependFirstRecord(digitalRecords []database.DevicePortDigitalRecord, recordsFrom time.Time, portId int) []database.DevicePortDigitalRecord {
	var firstRecord database.DevicePortDigitalRecord
	if digitalRecords[0].Data == 1 {
		firstRecord.DateTime = recordsFrom
		firstRecord.Data = 0
		firstRecord.DevicePortID = portId
	} else {
		firstRecord.DateTime = recordsFrom
		firstRecord.Data = 1
		firstRecord.DevicePortID = portId
	}
	digitalRecords = append([]database.DevicePortDigitalRecord{firstRecord}, digitalRecords...)
	return digitalRecords
}

func getListOdDeviceDigitalPorts(db *gorm.DB, data ChartDataInput) []database.DevicePort {
	var workplace database.Workplace
	db.Where("name = ?", data.Workplace).Find(&workplace)
	var listOfWorkplacePorts []database.WorkplacePort
	db.Where("workplace_id = ?", workplace.ID).Find(&listOfWorkplacePorts)
	var listOfDigitalDevicePorts []database.DevicePort
	for _, port := range listOfWorkplacePorts {
		var digitalDevicePort database.DevicePort
		db.Where("id = ?", port.DevicePortID).Where("device_port_type_id = 1").Find(&digitalDevicePort)
		if digitalDevicePort.ID > 0 {
			listOfDigitalDevicePorts = append(listOfDigitalDevicePorts, digitalDevicePort)
		}
	}
	return listOfDigitalDevicePorts
}

func getChartData(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data ChartDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData ChartDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing chart data ended")
		return
	}
	logInfo("MAIN", "Processing chart data for "+data.Workplace+" from "+data.ChartsStart+" to "+data.ChartsEnd)
	recordsFrom, err := time.Parse("2006-01-02T15:04:05.000Z", data.ChartsStart)
	if err != nil {
		logError("MAIN", "Problem parsing date: "+data.ChartsStart)
		var responseData ChartDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing chart data ended")
		return
	}
	recordsTo, err := time.Parse("2006-01-02T15:04:05.000Z", data.ChartsEnd)
	if err != nil {
		logError("MAIN", "Problem parsing date: "+data.ChartsEnd)
		var responseData ChartDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing chart data ended")
		return
	}
	logInfo("MAIN", "Processing chart data for "+data.Workplace+" from "+strconv.Itoa(int(recordsFrom.Unix()))+" to "+strconv.Itoa(int(recordsTo.Unix())))
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData ChartDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Processing chart data ended")
		return
	}

	productionData, downtimeData, powerOffData := getTimelineData(db, data, recordsFrom.In(time.UTC), recordsTo.In(time.UTC))
	productionState, downtimeState, poweroffState := getStates(db)
	productionLocale, downtimeLocale, poweroffLocale := getLocales(data, db)
	timelineDataList := processAllStates(productionData, productionState, productionLocale, downtimeData, downtimeState, downtimeLocale, powerOffData, poweroffState, poweroffLocale)

	listOfDigitalDevicePorts := getListOdDeviceDigitalPorts(db, data)
	listOfDigitalRecords := getDigitalRecords(listOfDigitalDevicePorts, db, recordsFrom, recordsTo)
	digitalDataList := processDownloadedDigitalRecords(listOfDigitalRecords, db)

	var outputData ChartDataOutput
	outputData.Result = "ok"
	outputData.TimelineDataList = timelineDataList
	outputData.DigitalDataList = digitalDataList
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Processing chart data ended")
}

func processAllStates(productionData []TimelineData, productionState database.State, productionLocale database.Locale, downtimeData []TimelineData, downtimeState database.State, downtimeLocale database.Locale, powerOffData []TimelineData, poweroffState database.State, poweroffLocale database.Locale) []TimelineDataList {
	var timelineDataList []TimelineDataList
	productionDataList := processProductionStates(productionData, productionState, productionLocale)
	timelineDataList = append(timelineDataList, productionDataList)
	downtimeDataList := processDowntimeStates(downtimeData, downtimeState, downtimeLocale)
	timelineDataList = append(timelineDataList, downtimeDataList)
	powerOffDataList := processPowerOffStates(powerOffData, poweroffState, poweroffLocale)
	timelineDataList = append(timelineDataList, powerOffDataList)
	return timelineDataList
}

func processPowerOffStates(powerOffData []TimelineData, poweroffState database.State, poweroffLocale database.Locale) TimelineDataList {
	var powerOffDataList TimelineDataList
	powerOffDataList.Data = powerOffData
	powerOffDataList.Color = poweroffState.Color
	value := reflect.ValueOf(poweroffLocale)
	for i := 0; i < value.NumField(); i++ {
		temp := value.Field(i).String()
		if len(temp) > 0 && !strings.Contains(temp, "{") {
			powerOffDataList.Name = temp
		}
	}
	return powerOffDataList
}

func processDowntimeStates(downtimeData []TimelineData, downtimeState database.State, downtimeLocale database.Locale) TimelineDataList {
	var downtimeDataList TimelineDataList
	downtimeDataList.Data = downtimeData
	downtimeDataList.Color = downtimeState.Color
	value := reflect.ValueOf(downtimeLocale)
	for i := 0; i < value.NumField(); i++ {
		temp := value.Field(i).String()
		if len(temp) > 0 && !strings.Contains(temp, "{") {
			downtimeDataList.Name = temp
		}
	}
	return downtimeDataList
}

func processProductionStates(productionData []TimelineData, productionState database.State, productionLocale database.Locale) TimelineDataList {
	var productionDataList TimelineDataList
	productionDataList.Data = productionData
	productionDataList.Color = productionState.Color
	value := reflect.ValueOf(productionLocale)
	for i := 0; i < value.NumField(); i++ {
		temp := value.Field(i).String()
		if len(temp) > 0 && !strings.Contains(temp, "{") {
			productionDataList.Name = temp
		}
	}
	return productionDataList
}

func getLocales(data ChartDataInput, db *gorm.DB) (database.Locale, database.Locale, database.Locale) {
	logInfo("MAIN", "Downloading locales for "+data.Locale)
	var productionLocale database.Locale
	var downtimeLocale database.Locale
	var poweroffLocale database.Locale
	db.Where("name = 'production'").Select(data.Locale).Find(&productionLocale)
	db.Where("name = 'downtime'").Select(data.Locale).Find(&downtimeLocale)
	db.Where("name = 'poweroff'").Select(data.Locale).Find(&poweroffLocale)
	return productionLocale, downtimeLocale, poweroffLocale
}

func getStates(db *gorm.DB) (database.State, database.State, database.State) {
	var productionState database.State
	var downtimeState database.State
	var poweroffState database.State
	db.Where("name = ?", "Production").Find(&productionState)
	db.Where("name = ?", "Downtime").Find(&downtimeState)
	db.Where("name = ?", "Poweroff").Find(&poweroffState)
	return productionState, downtimeState, poweroffState
}

func getTimelineData(db *gorm.DB, data ChartDataInput, recordsFrom time.Time, recordsTo time.Time) ([]TimelineData, []TimelineData, []TimelineData) {
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

func readTimeZoneFromDatabase() string {
	logInfo("MAIN", "Reading timezone from database")
	timer := time.Now()
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return ""
	}
	var settings database.Setting
	db.Where("name=?", "timezone").Find(&settings)
	logInfo("MAIN", "Timezone read in "+time.Since(timer).String())
	return settings.Value
}
