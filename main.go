package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/chayupon/Calculator/internal/calculate"
	//	"strconv"
)

func main() {
	//fmt.Println("TTTTTT")
	router := mux.NewRouter()
	router.HandleFunc("/calculator", calculate.Calculate).Methods("POST")

	fmt.Println(http.ListenAndServe(":8082", router))

}
