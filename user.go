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
	Result string
	Data   string
	Locale string
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
		if user.Password == data.UserSessionStorage {
			logInfo("MAIN", "User matches")
			var responseData VerifyUserOutput
			responseData.Result = "ok"
			responseData.Data = user.Password
			responseData.Locale = user.Locale
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
			logError("MAIN", "Password/email not empty")
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
			if userMatchesPassword {
				logInfo("MAIN", "User matches")
				var responseData VerifyUserOutput
				responseData.Result = "ok"
				responseData.Data = user.Password
				responseData.Locale = user.Locale
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

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
