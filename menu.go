package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
)

type ContentInput struct {
	Content string
}

type ContentOutput struct {
	Html        string
	MenuLocales []MenuLocale
}

func getMenu(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData ContentOutput
		responseData.Html = ""
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Verifying user ended")
		return
	}

	var databaseMenuLocales []database.Locale
	db.Where("name like ?", data.Content+"-%").Find(&databaseMenuLocales)
	var menuLocales []MenuLocale
	for _, locale := range databaseMenuLocales {
		var menuLocale MenuLocale
		menuLocale.Name = locale.Name
		menuLocale.CsCZ = locale.CsCZ
		menuLocale.DeDE = locale.DeDE
		menuLocale.EnUS = locale.EnUS
		menuLocale.EsES = locale.EsES
		menuLocale.FrFR = locale.FrFR
		menuLocale.ItIT = locale.ItIT
		menuLocale.PlPL = locale.PlPL
		menuLocale.PtPT = locale.PtPT
		menuLocale.SkSK = locale.SkSK
		menuLocale.RuRU = locale.RuRU
		menuLocales = append(menuLocales, menuLocale)
	}
	logInfo("MAIN", "Html read from file")
	var responseData ContentOutput
	responseData.Html = string(file)
	responseData.MenuLocales = menuLocales
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("MAIN", "Parsing data from page ended")
	return
}
