package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

const version = "2020.4.2.20"
const serviceName = "System WebService"
const serviceDescription = "System web interface"
const config = "user=postgres password=Zps05..... dbname=version3 host=localhost port=5432 sslmode=disable"

type program struct{}

func main() {
	logInfo("MAIN", serviceName+" ["+version+"] starting...")
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
}

func (p *program) Start(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started")
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] stopped")
	return nil
}

func (p *program) run() {
	updateProgramVersion()
	router := httprouter.New()
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/fonts/*filepath", http.Dir("fonts"))
	router.GET("/", system)
	router.POST("/get_menu", getMenu)
	router.POST("/get_content", getContent)
	router.POST("/verify_user", verifyUser)
	err := http.ListenAndServe(":80", router)
	if err != nil {
		logError("MAIN", "Problem starting service: "+err.Error())
		os.Exit(-1)
	}
	logInfo("MAIN", serviceName+" ["+version+"] running")
}

func system(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "./html/system.html")
}

func updateProgramVersion() {
	logInfo("MAIN", "Writing program version into settings")
	timer := time.Now()
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	var existingSettings database.Setting
	db.Where("name=?", serviceName).Find(&existingSettings)
	existingSettings.Name = serviceName
	existingSettings.Value = version
	db.Save(&existingSettings)
	logInfo("MAIN", "Program version written into settings in "+time.Since(timer).String())
}
