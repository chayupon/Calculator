package main

import (
	"fmt"

	"github.com/chayupon/Calculator/internal/service/operate"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
)
//Cal calculate
type Cal struct{
	Input1 int `json:"input1"` 
	Input2 int `json:"input2"`
	Result float32 `json:"result"`  
}

func calculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cal Cal
	_ = json.NewDecoder(r.Body).Decode(&cal)
	cal.Result=1000
	json.NewEncoder(w).Encode(&cal)
	
  }

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/calculator", calculate).Methods("POST")
	http.ListenAndServe(":8080", router)


  getCons:
	var cons string
	fmt.Print("Please select an operation: +, -, *, / : ")
	fmt.Scanln(&cons)

	var input1 int
	fmt.Print("Please input the first number: ")
	fmt.Scanln(&input1)

	var input2 int
	fmt.Print("Please input the second number: ")
	fmt.Scanln(&input2)

	switch cons {
	case "+":
		fmt.Print("Result: ")
		fmt.Println(operate.Add(input1, input2, cons))
	case "-":
		fmt.Print("Result: ")
		fmt.Println(operate.Add(input1, input2, cons))

	case "*":
		fmt.Print("Result: ")
		fmt.Println(operate.Add(input1, input2, cons))

	case "/":
		fmt.Print("Result: ")
		fmt.Println(operate.Add(input1, input2, cons))

	default:
		fmt.Println("Invalid operation selected!")
		goto getCons
	}
}


	// POST /calculator ,body >> {"input1":0,"input2":0,"operator":""}

	// input number > 0-9 , other error

	// +,-,*,/
	// result

	//case 1 0+0=0

