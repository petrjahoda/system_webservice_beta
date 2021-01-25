package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"github.com/snabb/isoweek"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"sort"
	"strconv"
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
	Production []WorkplaceDataString
	Downtime   []WorkplaceDataString
	Poweroff   []WorkplaceDataString
	Result     string
}

type CalendarDataOutput struct {
	Result string
	Data   []CalendarData
}

type CalendarData struct {
	Date            string
	ProductionValue string
}

type WorkplaceData struct {
	WorkplaceName string
	State         string
	StateDuration time.Duration
}

type WorkplaceDataString struct {
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
	start := time.Now()
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData CalendarDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Processing calendar data started for "+data.Input+" and "+data.Selection)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData CalendarDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended in "+time.Since(start).String())
		return
	}
	lastYear := time.Now().AddDate(-1, 0, 0)
	lastYearStart := time.Date(lastYear.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	workplaceStateRecords := downloadWorkplaceStateRecords(data, db, lastYearStart)
	productionRate := calculateProductionRate(workplaceStateRecords)
	var calendarData []CalendarData
	secondInOneDay := 86400.0
	todayStart := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	secondToday := time.Now().Sub(todayStart).Seconds()
	today := time.Now()
	todayFormatted := today.Format("2006-01-02")
	for datetime, totalDuration := range productionRate {
		averageDayDurationInSeconds := totalDuration / float64(len(workplaceStateRecords))
		var oneCalendarData CalendarData
		oneCalendarData.Date = datetime
		if oneCalendarData.Date == todayFormatted {
			oneCalendarData.ProductionValue = strconv.FormatFloat(averageDayDurationInSeconds*100/secondToday, 'f', 1, 64)
		} else {
			oneCalendarData.ProductionValue = strconv.FormatFloat(averageDayDurationInSeconds*100/secondInOneDay, 'f', 1, 64)
		}
		calendarData = append(calendarData, oneCalendarData)
	}
	var outputData CalendarDataOutput
	outputData.Data = calendarData
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended")
	return
}

func calculateProductionRate(workplaceStateRecords map[database.Workplace][]database.StateRecord) map[string]float64 {
	var productionRate map[string]float64
	productionRate = make(map[string]float64)
	for _, records := range workplaceStateRecords {
		var previousRecord database.StateRecord
		for _, record := range records {
			if record.StateID != 1 {
				format := record.DateTimeStart.Format("2006-01-02")
				if previousRecord.DateTimeStart.YearDay() == record.DateTimeStart.YearDay() {
					productionRate[format] = productionRate[format] + record.DateTimeStart.Sub(previousRecord.DateTimeStart).Seconds()
				} else {
					if previousRecord.DateTimeStart.IsZero() {
					} else {
						endOfDay := time.Date(previousRecord.DateTimeStart.Year(), previousRecord.DateTimeStart.Month(), previousRecord.DateTimeStart.Day(), 24, 0, 0, 0, time.UTC)
						dayFormat := endOfDay.Format("2006-01-02")
						productionRate[dayFormat] = productionRate[dayFormat] + endOfDay.Sub(previousRecord.DateTimeStart).Seconds()
						for endOfDay.YearDay() != record.DateTimeStart.YearDay() {
							previousDay := endOfDay
							previousDayFormat := endOfDay.Format("2006-01-02")
							endOfDay = endOfDay.Add(24 * time.Hour)
							if productionRate[previousDayFormat] == 0 {
								productionRate[previousDayFormat] = productionRate[previousDayFormat] + endOfDay.Sub(previousDay).Seconds()
							}
						}
						productionRate[format] = productionRate[format] + endOfDay.Sub(record.DateTimeStart).Seconds()
					}
				}
			}
			previousRecord = record
		}
	}
	return productionRate
}

func getLiveProductivityData(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	start := time.Now()
	logInfo("MAIN", "Parsing data")
	var data LiveDataInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData LiveDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended in "+time.Since(start).String())
		return
	}
	logInfo("MAIN", "Processing live productivity data started for "+data.Input+" "+data.Selection)

	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData LiveDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended in "+time.Since(start).String())
		return
	}
	lastMonth := time.Now().AddDate(0, -1, 0)
	lastMonthStart := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	thisMonthStarts := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	year, week := time.Now().ISOWeek()
	thisWeekStart := isoweek.StartTime(year, week, time.UTC)
	lastWeekStart := thisWeekStart.AddDate(0, 0, -7)
	yesterday := time.Now().AddDate(0, 0, -1)
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	todayStart := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

	workplaceStateRecords := downloadWorkplaceStateRecords(data, db, lastMonthStart)
	productionRate := calculateProductionRate(workplaceStateRecords)

	secondInOneDay := 86400.0
	var outputData LiveDataOutput
	outputData.Result = "ok"
	lastMonthDays := 0
	thisMonthDays := 0
	lastWeekDays := 0
	thisWeekDays := 0
	lastMonthTotal := 0.0
	thisMonthTotal := 0.0
	lastWeekTotal := 0.0
	thisWeekTotal := 0.0

	for datetime, totalDuration := range productionRate {
		layout := "2006-01-02"
		dateTimeConverted, _ := time.Parse(layout, datetime)
		if dateTimeConverted == lastMonthStart || (lastMonthStart.Before(dateTimeConverted) && dateTimeConverted.Before(thisMonthStarts)) {
			lastMonthTotal = lastMonthTotal + totalDuration
			lastMonthDays++
		}
		if dateTimeConverted == thisMonthStarts || thisMonthStarts.Before(dateTimeConverted) {
			thisMonthTotal = thisMonthTotal + totalDuration
			thisMonthDays++
		}
		if dateTimeConverted == lastWeekStart || (lastWeekStart.Before(dateTimeConverted) && dateTimeConverted.Before(thisWeekStart)) {
			lastWeekTotal = lastWeekTotal + totalDuration
			lastWeekDays++
		}
		if dateTimeConverted == thisWeekStart || thisWeekStart.Before(dateTimeConverted) {
			thisWeekTotal = thisWeekTotal + totalDuration
			thisWeekDays++
		}
		if dateTimeConverted == yesterdayStart || (yesterdayStart.Before(dateTimeConverted) && dateTimeConverted.Before(todayStart)) {
			outputData.Yesterday = strconv.FormatFloat(totalDuration/float64(len(workplaceStateRecords))/secondInOneDay*100, 'f', 1, 64)
		}
		if dateTimeConverted == todayStart || todayStart.Before(dateTimeConverted) {
			secondToday := time.Now().Sub(todayStart).Seconds()
			outputData.Today = strconv.FormatFloat(totalDuration/float64(len(workplaceStateRecords))/secondToday*100, 'f', 1, 64)
		}
	}
	if lastMonthDays == 0 {
		outputData.LastMonth = "0"
	} else {
		outputData.LastMonth = strconv.FormatFloat(lastMonthTotal/float64(len(workplaceStateRecords))/(secondInOneDay*float64(lastMonthDays))*100, 'f', 1, 64)
	}
	if thisMonthDays == 0 {
		outputData.ThisMonth = "0"
	} else {
		outputData.ThisMonth = strconv.FormatFloat(thisMonthTotal/float64(len(workplaceStateRecords))/(secondInOneDay*float64(thisMonthDays))*100, 'f', 1, 64)
	}
	if lastWeekDays == 0 {
		outputData.LastWeek = "0"
	} else {
		outputData.LastWeek = strconv.FormatFloat(lastWeekTotal/float64(len(workplaceStateRecords))/(secondInOneDay*float64(lastWeekDays))*100, 'f', 1, 64)
	}
	if thisWeekDays == 0 {
		outputData.ThisWeek = "0"
	} else {
		outputData.ThisWeek = strconv.FormatFloat(thisWeekTotal/float64(len(workplaceStateRecords))/(secondInOneDay*float64(thisWeekDays))*100, 'f', 1, 64)
	}

	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(outputData)
	logInfo("MAIN", "Parsing data ended in "+time.Since(start).String())
	return
}

func initiate(previousTime time.Time, record database.StateRecord, initial int) (time.Time, int) {
	if previousTime.Year() == 1 {
		previousTime = record.DateTimeStart
	}
	if initial == 0 {
		initial = record.StateID
	}
	return previousTime, initial
}

func downloadWorkplaceStateRecords(data LiveDataInput, db *gorm.DB, lastMonthStart time.Time) map[database.Workplace][]database.StateRecord {
	logInfo("MAIN", "Downloading data for "+data.Input)
	var workplacesData map[database.Workplace][]database.StateRecord
	workplacesData = make(map[database.Workplace][]database.StateRecord)
	switch data.Input {
	case "workplace":
		{
			var workplace database.Workplace
			var stateRecords []database.StateRecord
			db.Where("name = ?", data.Selection).Find(&workplace)
			db.Where("date_time_start > ?", lastMonthStart).Where("workplace_id = ?", workplace.ID).Order("date_time_start asc").Find(&stateRecords)
			workplacesData[workplace] = stateRecords
		}
	case "group":
		{
			var workplaceSection database.WorkplaceSection
			db.Where("name = ?", data.Selection).Find(&workplaceSection)
			var workplaces []database.Workplace
			var stateRecords []database.StateRecord
			db.Where("workplace_section_id = ?", workplaceSection.ID).Find(&workplaces)
			for _, workplace := range workplaces {
				db.Where("date_time_start > ?", lastMonthStart).Where("workplace_id = ?", workplace.ID).Order("date_time_start asc").Find(&stateRecords)
				workplacesData[workplace] = stateRecords
			}
			// INFO: if slow, check to use below
			//var workplaceByStateId []database.StateRecord
			//db.Debug().Select("MAX(date_time_start), workplace_id, state_id").Group("workplace_id").Group("state_id").Find(&workplaceByStateId)
			//for _, record := range workplaceByStateId {
			//	fmt.Println(record.WorkplaceID, record.StateID)
			//}
		}
	default:
		{
			var workplaces []database.Workplace
			var stateRecords []database.StateRecord
			db.Find(&workplaces)
			for _, workplace := range workplaces {
				db.Where("date_time_start > ?", lastMonthStart).Where("workplace_id = ?", workplace.ID).Order("date_time_start asc").Find(&stateRecords)
				workplacesData[workplace] = stateRecords
			}
			// INFO: if slow, check to use below
			//var workplaceByStateId []database.StateRecord
			//db.Debug().Select("MAX(date_time_start), workplace_id, state_id").Group("workplace_id").Group("state_id").Find(&workplaceByStateId)
			//for _, record := range workplaceByStateId {
			//	fmt.Println(record.WorkplaceID, record.StateID)
			//}
		}
	}
	return workplacesData
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
	var productionWorkplaces []WorkplaceData
	var downtimeWorkplaces []WorkplaceData
	var poweroffWorkplaces []WorkplaceData
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData LiveDataOutput
		responseData.Result = "nok: " + err.Error()
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Verifying user ended")
		return
	}
	var workplaces []database.Workplace
	if data.Input == "company" {
		db.Find(&workplaces)
	} else {
		var workplaceSection database.WorkplaceSection
		db.Where("name = ?", data.Selection).Find(&workplaceSection)
		db.Where("workplace_section_id = ?", workplaceSection.ID).Find(&workplaces)
	}
	for _, workplace := range workplaces {
		var stateRecord database.StateRecord
		db.Where("workplace_id = ?", workplace.ID).Last(&stateRecord)
		if stateRecord.StateID == 1 {
			productionWorkplaces = append(productionWorkplaces, WorkplaceData{WorkplaceName: workplace.Name, State: "production", StateDuration: time.Now().Sub(stateRecord.DateTimeStart).Round(1 * time.Second)})
		}
		if stateRecord.StateID == 2 {
			downtimeWorkplaces = append(downtimeWorkplaces, WorkplaceData{WorkplaceName: workplace.Name, State: "downtime", StateDuration: time.Now().Sub(stateRecord.DateTimeStart).Round(1 * time.Second)})
		}
		if stateRecord.StateID == 3 {
			poweroffWorkplaces = append(poweroffWorkplaces, WorkplaceData{WorkplaceName: workplace.Name, State: "poweroff", StateDuration: time.Now().Sub(stateRecord.DateTimeStart).Round(1 * time.Second)})
		}
	}

	sort.Slice(productionWorkplaces, func(i, j int) bool {
		return productionWorkplaces[i].StateDuration > productionWorkplaces[j].StateDuration
	})
	sort.Slice(downtimeWorkplaces, func(i, j int) bool {
		return downtimeWorkplaces[i].StateDuration > downtimeWorkplaces[j].StateDuration
	})
	sort.Slice(poweroffWorkplaces, func(i, j int) bool {
		return poweroffWorkplaces[i].StateDuration > poweroffWorkplaces[j].StateDuration
	})
	var productionWorkplacesString []WorkplaceDataString
	var downtimeWorkplacesString []WorkplaceDataString
	var poweroffWorkplacesString []WorkplaceDataString
	for _, workplace := range productionWorkplaces {
		productionWorkplacesString = append(productionWorkplacesString, WorkplaceDataString{
			WorkplaceName: workplace.WorkplaceName,
			State:         workplace.State,
			StateDuration: workplace.StateDuration.String(),
		})
	}
	for _, workplace := range downtimeWorkplaces {
		downtimeWorkplacesString = append(downtimeWorkplacesString, WorkplaceDataString{
			WorkplaceName: workplace.WorkplaceName,
			State:         workplace.State,
			StateDuration: workplace.StateDuration.String(),
		})
	}
	for _, workplace := range poweroffWorkplaces {
		poweroffWorkplacesString = append(poweroffWorkplacesString, WorkplaceDataString{
			WorkplaceName: workplace.WorkplaceName,
			State:         workplace.State,
			StateDuration: workplace.StateDuration.String(),
		})
	}

	var outputData LiveOverviewDataOutput
	outputData.Result = "ok"
	outputData.Production = productionWorkplacesString
	outputData.Downtime = downtimeWorkplacesString
	outputData.Poweroff = poweroffWorkplacesString
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
			var workplaces []database.Workplace
			db.Where("workplace_section_id = ?", group.ID).First(&workplaces)
			if len(workplaces) > 0 {
				selectionData = append(selectionData, group.Name)
			}
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
	logInfo("MAIN", "Order record start at "+orderRecord.DateTimeStart.String())
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
