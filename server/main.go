package main

import (
	"C2-D2/server/controllers"
	"C2-D2/server/database"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// TO-DO: Read this in from a config/env file
const (
	DB_USER     = "postgres"
	DB_PASSWORD = "docker"
	DB_NAME     = "postgres"
	DB_PORT     = "5432"
	SERVER_PORT = "8080"
)

func main() {
	database.Initialize(DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	database.Migrate()
	Router := mux.NewRouter()
	controllers.InitializeRoutes(Router)
	logrus.Println(fmt.Sprintf("Starting Server on port %s", SERVER_PORT))
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", SERVER_PORT), Router))
}
