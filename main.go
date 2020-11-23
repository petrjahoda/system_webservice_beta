package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/petrjahoda/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

type ContentInput struct {
	Content string
}

type ContentOutput struct {
	Html string
}

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

const config = "user=postgres password=Zps05..... dbname=version3 host=localhost port=5432 sslmode=disable"

func main() {
	router := httprouter.New()
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/fonts/*filepath", http.Dir("fonts"))
	router.GET("/", system)
	router.POST("/get_content", getContent)
	router.POST("/verify_user", verifyUser)
	_ = http.ListenAndServe(":80", router)
}

func verifyUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	logInfo("MAIN", "Checking login")
	var data VerifyUserInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logError("MAIN", "Error parsing data: "+err.Error())
		var responseData VerifyUserOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}
	logInfo("MAIN", "Checking login started for "+data.UserEmail+":"+data.UserPassword+":"+data.UserSessionStorage)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Cannot connect to database")
		var responseData VerifyUserOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Parsing data from page ended")
		return
	}

	if len(data.UserSessionStorage) == 0 {
		logInfo("MAIN", "Session storage empty")
		if len(data.UserEmail) == 0 || len(data.UserPassword) == 0 {
			logError("MAIN", "Password or email empty")
			var responseData VerifyUserOutput
			responseData.Result = "nok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("MAIN", "Parsing data from page ended")
			return
		}
	} else {
		logInfo("MAIN", "Session storage not empty")
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
			logInfo("MAIN", "Checking login ended")
			return
		} else {
			logError("MAIN", "User does not match")
			var responseData VerifyUserOutput
			responseData.Result = "nok"
			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(responseData)
			logInfo("MAIN", "Parsing data from page ended")
			return
		}
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
		logInfo("MAIN", "Checking login ended")
		return
	} else {
		logInfo("MAIN", "User does not match")
		var responseData VerifyUserOutput
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(responseData)
		logInfo("MAIN", "Checking login ended")
		return
	}

}

func getContent(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
	logInfo("MAIN", "Html read from file")
	var responseData ContentOutput
	responseData.Html = string(file)
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(responseData)
	logInfo("MAIN", "Parsing data from page ended")
	return

}

func system(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "./html/system.html")
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
