package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type VerifyUserInput struct {
	UserEmail          string
	UserPassword       string
	UserSessionStorage string
}

type VerifyUserOutput struct {
	Result      string
	Data        string
	Locale      string
	MenuLocales []MenuLocale
}
type MenuLocale struct {
	Name string
	CsCZ string
	DeDE string
	EnUS string
	EsES string
	FrFR string
	ItIT string
	PlPL string
	PtPT string
	SkSK string
	RuRU string
}

func verifyUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Parsing data")
	var data VerifyUserInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData VerifyUserOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data ended")
		return
	}
	logInfo("MAIN", "Verifying user started for "+data.UserEmail+":"+data.UserPassword+":"+data.UserSessionStorage)
	if len(data.UserSessionStorage) != 0 {
		logInfo("MAIN", "Session storage not empty")
		db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		if err != nil {
			logError("MAIN", "Cannot connect to database")
			var responseData VerifyUserOutput
			responseData.Result = "nok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("MAIN", "Verifying user ended")
			return
		}
		var user database.User
		db.Where("password = ?", data.UserSessionStorage).Find(&user)
		logInfo("MAIN", user.Password)
		menuLocales := downloadMenuLocales(db)
		if user.Password == data.UserSessionStorage {
			logInfo("MAIN", "User matches")
			var responseData VerifyUserOutput
			responseData.Result = "ok"
			responseData.Data = user.Password
			responseData.Locale = user.Locale
			responseData.MenuLocales = menuLocales
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("MAIN", "Verifying user ended")
			return
		} else {
			logError("MAIN", "User does not match")
			var responseData VerifyUserOutput
			responseData.Result = "nok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("MAIN", "Verifying user ended")
			return
		}

	} else {
		logInfo("MAIN", "Session storage empty")
		if len(data.UserEmail) == 0 || len(data.UserPassword) == 0 {
			logError("MAIN", "Password/email empty")
			var responseData VerifyUserOutput
			responseData.Result = "nok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("MAIN", "Verifying user ended")
			return
		} else {
			logInfo("MAIN", "Password/email not empty")
			db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
			if err != nil {
				logError("MAIN", "Cannot connect to database")
				var responseData VerifyUserOutput
				responseData.Result = "nok"
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				logInfo("MAIN", "Verifying user ended")
				return
			}
			var user database.User
			db.Where("email = ?", data.UserEmail).Find(&user)
			logInfo("MAIN", user.Password)
			userMatchesPassword := comparePasswords(user.Password, []byte(data.UserPassword))
			menuLocales := downloadMenuLocales(db)
			if userMatchesPassword {
				logInfo("MAIN", "User matches")
				var responseData VerifyUserOutput
				responseData.Result = "ok"
				responseData.Data = user.Password
				responseData.Locale = user.Locale
				responseData.MenuLocales = menuLocales
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				logInfo("MAIN", "Verifying user ended")
				return
			} else {
				logInfo("MAIN", "User does not match")
				var responseData VerifyUserOutput
				responseData.Result = "nok"
				writer.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(writer).Encode(responseData)
				logInfo("MAIN", "Verifying user ended")
				return
			}
		}
	}
}

func downloadMenuLocales(db *gorm.DB) []MenuLocale {
	var databaseMenuLocales []database.Locale
	db.Where("name like ?", "navbar-%").Find(&databaseMenuLocales)
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
	return menuLocales
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
