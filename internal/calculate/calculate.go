package calculate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chayupon/Calculator/internal/operate"

	"net/http"

	"github.com/gorilla/mux"
)

//Query input to db
type Query struct {
	Router *mux.Router
	DB     *sql.DB
}

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

//History in db
type history struct {
	Sequence      int    `json:"sequence"`
	Time          string `json:"time"`
	InputAll      string `json:"input_all"`
	ErrorDescripe string `json:"error_descripe"`
}

//Calculate select
/*
func (q Query) createResult(w http.ResponseWriter, r *http.Request) {
	var sequence int
	var time string
	var inputAll string
	var errorDescripe string

	sqlStr := `SELECT "Sequence","Time","InputAll","ErrorDescripe" FROM history`
	rows, err := q.DB.Query(sqlStr)
	if err != nil {
		log.Println("Fail")
		return 
	}
	defer rows.Close()
	h := []history{}
	//fmt.Printf("%+v",u)
	for rows.Next() {
		//var username string
		if err := rows.Scan(&sequence, &time, &inputAll, &errorDescripe); err != nil {
			log.Println(err)
		}
		his := history{
			Sequence	 :  sequence,
			Time		 :  time,
			InputAll	 :  inputAll,
			ErrorDescripe:  errorDescripe,
		}
		h = append(h, his)
	}
	if !rows.NextResultSet() {
		log.Println(rows.Err())
	}
	output, _ := json.Marshal(&h)
	fmt.Println(string(output))
	json.NewEncoder(w).Encode(&h)
	
}*/

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
		//w.WriteHeader()
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
