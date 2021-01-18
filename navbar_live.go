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
	Production []WorkplaceData
	Downtime   []WorkplaceData
	Poweroff   []WorkplaceData
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
	for datetime, totalDuration := range productionRate {
		averageDayDurationInSeconds := totalDuration / float64(len(workplaceStateRecords))
		var oneCalendarData CalendarData
		oneCalendarData.Date = datetime
		oneCalendarData.ProductionValue = strconv.FormatFloat(averageDayDurationInSeconds*100/secondInOneDay, 'f', 1, 64)
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
			format := record.DateTimeStart.Format("2006-01-02")
			if record.StateID != 1 {
				if previousRecord.DateTimeStart.YearDay() == record.DateTimeStart.YearDay() {
					productionRate[format] = productionRate[format] + record.DateTimeStart.Sub(previousRecord.DateTimeStart).Seconds()
				} else {
					if !previousRecord.DateTimeStart.IsZero() {
						//fmt.Println("We have production in another day, beginning at: " + previousRecord.DateTimeStart.String())
						// you have to calculate that previous date should be more days before
						// add proper language translations into database
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
	yesterday := time.Now().AddDate(0, 0, -1)
	lastMonthStart := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	thisMonthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	todayStart := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	year, week := time.Now().ISOWeek()
	thisWeekStart := isoweek.StartTime(year, week, time.UTC)
	lastWeekStart := thisWeekStart.AddDate(0, 0, -7)
	workplaceStateRecords := downloadWorkplaceStateRecords(data, db, lastMonthStart)
	logInfo("MAIN", "Found total of "+strconv.Itoa(len(workplaceStateRecords))+" state records")
	var lastMonthTotalTimeAll time.Duration
	var lastMonthProductivityTimeAll time.Duration
	var thisMonthTotalTimeAll time.Duration
	var thisMonthProductivityTimeAll time.Duration
	var lastWeekTotalTimeAll time.Duration
	var lastWeekProductivityTimeAll time.Duration
	var thisWeekTotalTimeAll time.Duration
	var thisWeekProductivityTimeAll time.Duration
	var yesterdayTotalTimeAll time.Duration
	var yesterdayProductivityTimeAll time.Duration
	var todayTotalTimeAll time.Duration
	var todayProductivityTimeAll time.Duration
	for _, records := range workplaceStateRecords {
		var lastMonthTotalTime time.Duration
		var lastMonthProductivityTime time.Duration
		var thisMonthTotalTime time.Duration
		var thisMonthProductivityTime time.Duration
		var lastWeekTotalTime time.Duration
		var lastWeekProductivityTime time.Duration
		var thisWeekTotalTime time.Duration
		var thisWeekProductivityTime time.Duration
		var yesterdayTotalTime time.Duration
		var yesterdayProductivityTime time.Duration
		var todayTotalTime time.Duration
		var todayProductivityTime time.Duration
		var previousTime time.Time
		var initial int
		var lastMonthCompleted bool
		var lastWeekCompleted bool
		var yesterdayCompleted bool
		for _, record := range records {
			previousTime, initial = initiate(previousTime, record, initial)
			if record.DateTimeStart.Before(thisMonthStart) {
				if lastMonthTotalTime.Seconds() == 0 {
					lastMonthTotalTime += record.DateTimeStart.Sub(lastMonthStart)
					if initial == 1 {
						lastMonthProductivityTime += record.DateTimeStart.Sub(lastMonthStart)
					}
				} else {
					lastMonthTotalTime += record.DateTimeStart.Sub(previousTime)
				}
				if initial == 1 {
					lastMonthProductivityTime += record.DateTimeStart.Sub(previousTime)
				}
			} else {
				if !lastMonthCompleted {
					lastMonthTotalTime += thisMonthStart.Sub(previousTime)
					lastMonthCompleted = true
				}
				if thisMonthTotalTime.Seconds() == 0 {
					thisMonthTotalTime += record.DateTimeStart.Sub(thisMonthStart)
					if initial == 1 {
						thisMonthProductivityTime += record.DateTimeStart.Sub(thisMonthStart)
					}
				} else {
					thisMonthTotalTime += record.DateTimeStart.Sub(previousTime)
				}
				if initial == 1 {
					thisMonthProductivityTime += record.DateTimeStart.Sub(previousTime)
				}
			}
			if record.DateTimeStart.After(lastWeekStart) && record.DateTimeStart.Before(thisWeekStart) {
				if lastWeekTotalTime.Seconds() == 0 {
					lastWeekTotalTime += record.DateTimeStart.Sub(lastWeekStart)
					if initial == 1 {
						lastWeekProductivityTime += record.DateTimeStart.Sub(lastWeekStart)
					}
				} else {
					lastWeekTotalTime += record.DateTimeStart.Sub(previousTime)
				}
				if initial == 1 {
					lastWeekProductivityTime += record.DateTimeStart.Sub(previousTime)
				}
			} else if record.DateTimeStart.After(thisWeekStart) {
				if !lastWeekCompleted {
					lastWeekTotalTime += thisWeekStart.Sub(previousTime)
					lastWeekCompleted = true
				}
				if thisWeekTotalTime.Seconds() == 0 {
					thisWeekTotalTime += record.DateTimeStart.Sub(thisWeekStart)
					if initial == 1 {
						thisWeekProductivityTime += record.DateTimeStart.Sub(thisWeekStart)
					}
				} else {
					thisWeekTotalTime += record.DateTimeStart.Sub(previousTime)
				}
				if initial == 1 {
					thisWeekProductivityTime += record.DateTimeStart.Sub(previousTime)
				}
			}
			if record.DateTimeStart.After(yesterdayStart) && record.DateTimeStart.Before(todayStart) {
				if yesterdayTotalTime.Seconds() == 0 {
					yesterdayTotalTime += record.DateTimeStart.Sub(yesterdayStart)
					if initial == 1 {
						yesterdayProductivityTime += record.DateTimeStart.Sub(yesterdayStart)
					}
				} else {
					yesterdayTotalTime += record.DateTimeStart.Sub(previousTime)
				}
				if initial == 1 {
					yesterdayProductivityTime += record.DateTimeStart.Sub(previousTime)
				}
			} else if record.DateTimeStart.After(todayStart) {
				if !yesterdayCompleted {
					yesterdayTotalTime += todayStart.Sub(previousTime)
					yesterdayCompleted = true
				}
				if todayTotalTime.Seconds() == 0 {
					todayTotalTime += record.DateTimeStart.Sub(todayStart)
					if initial == 1 {
						todayProductivityTime += record.DateTimeStart.Sub(todayStart)
					}
				} else {
					todayTotalTime += record.DateTimeStart.Sub(previousTime)
				}
				if initial == 1 {
					todayProductivityTime += record.DateTimeStart.Sub(previousTime)
				}
			}
			previousTime = record.DateTimeStart
			initial = record.StateID

		}
		thisMonthTotalTime += time.Now().Sub(previousTime)
		thisWeekTotalTime += time.Now().Sub(previousTime)
		todayTotalTime += time.Now().Sub(previousTime)
		if initial == 1 {
			thisMonthProductivityTime += time.Now().Sub(previousTime)
			thisWeekProductivityTime += time.Now().Sub(previousTime)
			todayProductivityTime += time.Now().Sub(previousTime)
		}
		lastMonthTotalTimeAll = lastMonthTotalTimeAll + lastMonthTotalTime
		lastMonthProductivityTimeAll = lastMonthProductivityTimeAll + lastMonthProductivityTime
		thisMonthTotalTimeAll = thisMonthTotalTimeAll + thisMonthTotalTime
		thisMonthProductivityTimeAll = thisMonthProductivityTimeAll + thisMonthProductivityTime
		lastWeekTotalTimeAll = lastWeekTotalTimeAll + lastWeekTotalTime
		lastWeekProductivityTimeAll = lastWeekProductivityTimeAll + lastWeekProductivityTime
		thisWeekTotalTimeAll = thisWeekTotalTimeAll + thisWeekTotalTime
		thisWeekProductivityTimeAll = thisWeekProductivityTimeAll + thisWeekProductivityTime
		thisWeekProductivityTimeAll = thisWeekProductivityTimeAll + thisWeekProductivityTime
		yesterdayProductivityTimeAll = yesterdayProductivityTimeAll + yesterdayProductivityTime
		todayTotalTimeAll = yesterdayProductivityTimeAll + yesterdayProductivityTime
		todayProductivityTimeAll = todayProductivityTimeAll + todayProductivityTime
	}
	// TODO: pocita to tady nejake nesmysly pro viceropracovist
	logInfo("MAIN", "Last month, total time: "+lastMonthTotalTimeAll.String()+",  productivity time: "+lastMonthProductivityTimeAll.String())
	logInfo("MAIN", "This month, total time: "+thisMonthTotalTimeAll.String()+",  productivity time: "+thisMonthProductivityTimeAll.String())
	logInfo("MAIN", "Last week, total time: "+lastWeekTotalTimeAll.String()+",  productivity time: "+lastWeekProductivityTimeAll.String())
	logInfo("MAIN", "This week, total time: "+thisWeekTotalTimeAll.String()+",  productivity time: "+thisWeekProductivityTimeAll.String())
	logInfo("MAIN", "Yesterday, total time: "+yesterdayTotalTimeAll.String()+",  productivity time: "+yesterdayProductivityTimeAll.String())
	logInfo("MAIN", "Today, total time: "+todayTotalTimeAll.String()+",  productivity time: "+todayProductivityTimeAll.String())

	var outputData LiveDataOutput
	outputData.Result = "ok"
	outputData.Today = strconv.FormatFloat(todayProductivityTimeAll.Seconds()*100/todayTotalTimeAll.Seconds(), 'f', 1, 64)
	outputData.Yesterday = strconv.FormatFloat(yesterdayProductivityTimeAll.Seconds()*100/yesterdayTotalTimeAll.Seconds(), 'f', 1, 64)
	outputData.LastWeek = strconv.FormatFloat(lastWeekProductivityTimeAll.Seconds()*100/lastWeekTotalTimeAll.Seconds(), 'f', 1, 64)
	outputData.ThisWeek = strconv.FormatFloat(thisWeekProductivityTimeAll.Seconds()*100/thisWeekTotalTimeAll.Seconds(), 'f', 1, 64)
	outputData.LastMonth = strconv.FormatFloat(lastMonthProductivityTimeAll.Seconds()*100/lastMonthTotalTimeAll.Seconds(), 'f', 1, 64)
	outputData.ThisMonth = strconv.FormatFloat(thisMonthProductivityTimeAll.Seconds()*100/thisMonthTotalTimeAll.Seconds(), 'f', 1, 64)
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
			productionWorkplaces = append(productionWorkplaces, WorkplaceData{WorkplaceName: workplace.Name, State: "production", StateDuration: time.Now().Sub(stateRecord.DateTimeStart).Round(1 * time.Second).String()})
		}
		if stateRecord.StateID == 2 {
			downtimeWorkplaces = append(downtimeWorkplaces, WorkplaceData{WorkplaceName: workplace.Name, State: "downtime", StateDuration: time.Now().Sub(stateRecord.DateTimeStart).Round(1 * time.Second).String()})
		}
		if stateRecord.StateID == 3 {
			poweroffWorkplaces = append(poweroffWorkplaces, WorkplaceData{WorkplaceName: workplace.Name, State: "poweroff", StateDuration: time.Now().Sub(stateRecord.DateTimeStart).Round(1 * time.Second).String()})
		}
	}

	sort.Slice(productionWorkplaces, func(i, j int) bool {
		return productionWorkplaces[i].WorkplaceName < productionWorkplaces[j].WorkplaceName
	})
	sort.Slice(downtimeWorkplaces, func(i, j int) bool {
		return downtimeWorkplaces[i].WorkplaceName < downtimeWorkplaces[j].WorkplaceName
	})
	sort.Slice(poweroffWorkplaces, func(i, j int) bool {
		return poweroffWorkplaces[i].WorkplaceName < poweroffWorkplaces[j].WorkplaceName
	})
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
