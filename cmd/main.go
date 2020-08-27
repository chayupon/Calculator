package main

import (
	"fmt"
	"time"

	"github.com/chayupon/Calculator/internal/service/operate"
	_ "github.com/lib/pq"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//Cal calculate
type Cal struct {
	Input1    float64 `json:"input1"`
	Input2    float64 `json:"input2"`
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
}

//Output Value
type Output struct {
	Result float64 `json:"result"`
	Time   string  `json:"time"`
}

//var cals []Cal
func calculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cal Cal
	var out Output

	json.NewDecoder(r.Body).Decode(&cal)
	//fmt.Println(err)
	fmt.Println("Input1", cal.Input1)
	fmt.Println("Input2", cal.Input2)
	result, err := operate.Add(cal.Input1, cal.Input2, cal.Operation)
	fmt.Println("Result :", result, err)

	out.Result = result
	currentTime := time.Now()
	//fmt.Println("RFC3339Nano: ", currentTime.Format(time.RFC3339Nano))
	out.Time = currentTime.Format(time.RFC3339Nano)

	json.NewEncoder(w).Encode(&out)

}

func main() {
	//fmt.Println("TTTTTT")
	router := mux.NewRouter()
	router.HandleFunc("/calculator", calculate).Methods("POST")

	fmt.Println(http.ListenAndServe(":8081", router))

	/*
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
	*/
}
