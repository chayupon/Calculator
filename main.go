package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"net/http"

	"github.com/chayupon/Calculator/internal/calculate"
	"github.com/gorilla/mux"
	//	"strconv"
)


func main() {
	//fmt.Println("TTTTTT")
	const (
		dbHost     = "localhost"
		dbName     = "calculator"
		dbUser     = "postgres"
		dbPassword = "tonkla727426"
		dbPort     = 5432
	)
	dbCal := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbCal)
	
	if err != nil {
		log.Fatal("Connect Fail")
	}
	log.Println("Connect")
	defer db.Close()
	
	router := mux.NewRouter()
	router.HandleFunc("/calculator", calculate.Calculate).Methods("POST")
	http.ListenAndServe(":8082", router)

}
