package calculate

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chayupon/Calculator/internal/operate"

	"net/http"
	//	"strconv"
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

//outputError show error
type outputError struct {
	Errordescription string `json:"errordescription"`
	InputAll         string `json:"inputall"`
}

//Calculate cal
func Calculate(w http.ResponseWriter, r *http.Request) {

	var cal Cal
	var out Output
	var outerror outputError

	e := json.NewDecoder(r.Body).Decode(&cal)
	if e != nil {

		outerror.Errordescription = e.Error()
		//outerror.Errordescription = err
		fmt.Println(outerror.Errordescription)

		s := fmt.Sprintf("%f %s %f = %f", cal.Input1, cal.Operation, cal.Input2, cal.Result)
		fmt.Println("ResultAll :", s)
		outerror.InputAll = s
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&outerror)
		return
		//response
	}
	fmt.Println("Input2", cal.Input2)
	fmt.Println("Input1", cal.Input1)
	result, err := operate.Add(cal.Input1, cal.Input2, cal.Operation)
	fmt.Println("Result :", result, err)

	if err != nil {
		outerror.Errordescription = err.Error()
		fmt.Println(outerror.Errordescription)
		//เครื่องหมาย กับ status
		s := fmt.Sprintf("%f %s %f = %f", cal.Input1, cal.Operation, cal.Input2, cal.Result)
		fmt.Println("ResultAll :", s)
		outerror.InputAll = s
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&outerror)
		return
	} //response

	out.Result = result
	currentTime := time.Now()
	out.Time = currentTime.Format(time.RFC3339Nano)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&out)
	//response
} //response
