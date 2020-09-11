package calculate

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// "github.com/chayupon/Calculator/internal/count"
	"log"
	"net/http"
	"time"

	"github.com/chayupon/Calculator/internal/count"
	"github.com/chayupon/Calculator/internal/operate"
	"github.com/gorilla/mux"
)

//App input to db
type App struct {
	Router     *mux.Router
	DB         *sql.DB
	DateJoined time.Time
}

//Cal calculate
type Cal struct {
	Input1    float64 `json:"input1"`
	Input2    float64 `json:"input2"`
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
	//Count     int     `json:"count"`
}

//ResponseCount count
type ResponseCount struct {
	CountOperate []CountOperate `json:"result"`
}

//CountOperate count
type CountOperate struct {
	Operation string `json:"operation"`
	Count     int    `json:"count"`
}

//CountOperateError error
type CountOperateError struct {
	Errordescription string `json:"errordescription"`
}

//CountRequest operation
type CountRequest struct {
	Operation string `json:"operation"`
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

//Initialize connect db
func (a *App) Initialize(dbHost, dbPort, user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, user, password, dbname)
	fmt.Println("Database Connect")
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

//Run port
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/calculate", a.Calculate).Methods("POST")
	a.Router.HandleFunc("/calculate/request/operate", a.OperateRequest).Methods("POST")
	a.Router.HandleFunc("/calculate/detail", a.Detail).Methods("GET")
	a.Router.HandleFunc("/calculate/count/detail", a.CountOperateAll).Methods("GET")
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//CountOperateRequest count operate
func CountOperateRequest(operateall []string) (int, int, int, int) {
	countadd := 0
	countdiff := 0
	countmulti := 0
	countdiv := 0

	for _, operate := range operateall {
		if operate == "+" {
			countadd++
		} else if operate == "-" {
			countdiff++
		} else if operate == "*" {
			countmulti++
		} else if operate == "/" {
			countdiv++
		}
	}
	return countadd, countdiff, countmulti, countdiv
}

/*
func OutputError(outerror outputError, e error, cal Cal) {
	outerror.Errordescription = e.Error()
	fmt.Println("error :", outerror.Errordescription)
	s := fmt.Sprintf("%f %s %f = %f", cal.Input1, cal.Operation, cal.Input2, cal.Result)
	fmt.Println("ResultAll :", s)
	outerror.InputAll = s
}*/

//Calculate Insert
func (a *App) Calculate(w http.ResponseWriter, r *http.Request) {
	var cal Cal
	var out Output
	var outerror outputError
	sqlStr := `INSERT INTO history( time, input1, operate, input2, result, errordescripe) VALUES($1,$2,$3,$4,$5,$6)`
	e := json.NewDecoder(r.Body).Decode(&cal)
	if e != nil {
		//OutputError(outerror, e, cal)
		outerror.Errordescription = e.Error()
		fmt.Println("error :", outerror.Errordescription)
		s := fmt.Sprintf("%f %s %f = %f", cal.Input1, cal.Operation, cal.Input2, cal.Result)
		fmt.Println("ResultAll :", s)
		outerror.InputAll = s
		respondWithJSON(w, http.StatusBadRequest, outerror)
		return
	}
	result, err := operate.Add(cal.Input1, cal.Input2, cal.Operation)
	fmt.Println("Result :", result, err)
	currentime := time.Now()
	out.Time = currentime.Format(time.RFC3339)
	out.Result = result
	if err != nil {
		outerror.Errordescription = err.Error()
		fmt.Println("error :", outerror.Errordescription)
		s := fmt.Sprintf("%f %s %f = %f", cal.Input1, cal.Operation, cal.Input2, cal.Result)
		fmt.Println("ResultAll :", s)
		outerror.InputAll = s
	}
	_, err = a.DB.Exec(sqlStr, out.Time, cal.Input1, cal.Operation, cal.Input2, out.Result, outerror.Errordescription)
	if err != nil {
		fmt.Println(err)
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}
	if outerror.Errordescription != "" {
		respondWithJSON(w, http.StatusBadRequest, outerror)
		fmt.Println("Incorrect")
		return
	}
	respondWithJSON(w, http.StatusOK, out)
}

//Detail select db
func (a *App) Detail(w http.ResponseWriter, r *http.Request) {
	sqlStr := `SELECT  sequence, "time", input1, operate, input2, result, errordescripe FROM history`
	rows, err := a.DB.Query(sqlStr)
	if err != nil {
		log.Println("Fail", err)
		return
	}
	defer rows.Close()
	h := []history{}

	for rows.Next() {
		var sequence int
		var time string
		var input1 float64
		var operate string
		var input2 float64
		var result float64
		var errordescripe string
		if err := rows.Scan(&sequence, &time, &input1, &operate, &input2, &result, &errordescripe); err != nil {
			log.Println(err)
			respondWithJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		s := fmt.Sprintf("%f %s %f = %f", input1, operate, input2, result)
		//history.
		log.Println("inputall :", s)
		his := history{
			Sequence:      sequence,
			Time:          time,
			InputAll:      s,
			ErrorDescripe: errordescripe,
		}
		h = append(h, his)
	}
	output, _ := json.Marshal(&h)
	fmt.Println(string(output))
	respondWithJSON(w, http.StatusOK, h)
}

//CountOperateAll count history
func (a *App) CountOperateAll(w http.ResponseWriter, r *http.Request) {
	sqlStr := `SELECT  sequence,operate FROM history`
	rows, err := a.DB.Query(sqlStr)
	if err != nil {
		log.Println("Fail", err)
		fmt.Println("error")
		return
	}
	defer rows.Close()

	var operateall []string

	for rows.Next() {
		//  var count int
		var sequence int
		var operate string
		//  var errordescripe string
		if err := rows.Scan(&sequence, &operate); err != nil {
			log.Println(err)
			respondWithJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		operateall = append(operateall, operate)
	}
	countadd, countdiff, countmulti, countdiv := CountOperateRequest(operateall)
	countop := ResponseCount{
		CountOperate: []CountOperate{
			{
				Operation: "+",
				Count:     countadd,
			},
			{
				Operation: "-",
				Count:     countdiff,
			},
			{
				Operation: "*",
				Count:     countmulti,
			},
			{
				Operation: "/",
				Count:     countdiv,
			},
		},
	}
	fmt.Println("operateall", operateall)
	output, _ := json.Marshal(&countop)
	fmt.Println(string(output))
	respondWithJSON(w, http.StatusOK, countop)
}

// OperateRequest to count
func (a *App) OperateRequest(w http.ResponseWriter, r *http.Request) {
	var countOperateRQ CountRequest
	var out Output
	var countOperate CountOperate
	var counterror CountOperateError

	e := json.NewDecoder(r.Body).Decode(&countOperateRQ)
	if e != nil {
		counterror.Errordescription = e.Error()
		fmt.Println("error :", counterror.Errordescription)
		respondWithJSON(w, http.StatusBadRequest, counterror)
		return
	}
	count, err := count.CheckOperate(countOperateRQ.Operation)
	fmt.Println("operation :", count, err)
	currentime := time.Now()
	out.Time = currentime.Format(time.RFC3339)
	if err != nil {
		counterror.Errordescription = err.Error()
		fmt.Println("error :", counterror.Errordescription)
		respondWithJSON(w, http.StatusBadRequest, counterror)
		return
	}
	//select to show operate
	sqlStr := `SELECT sequence,operate FROM history`
	var sequence int
	rows, err := a.DB.Query(sqlStr)
	if err != nil {
		log.Println("Fail", err)
		fmt.Println("error")
		return
	}
	defer rows.Close()

	var operateall []string

	for rows.Next() {

		var operate string
		if err := rows.Scan(&sequence, &operate); err != nil {
			log.Println(err)
			respondWithJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		operateall = append(operateall, operate)
	}

	countadd, countdiff, countmulti, countdiv := CountOperateRequest(operateall)
	//show count operate
	if countOperateRQ.Operation == "+" {
		countOperate = CountOperate{
			Operation: "+",
			Count:     countadd,
		}
	} else if countOperateRQ.Operation == "-" {
		countOperate = CountOperate{
			Operation: "-",
			Count:     countdiff,
		}
	} else if countOperateRQ.Operation == "*" {
		countOperate = CountOperate{
			Operation: "*",
			Count:     countmulti,
		}

	} else if countOperateRQ.Operation == "/" {
		countOperate = CountOperate{
			Operation: "/",
			Count:     countdiv,
		}
	}

	fmt.Println("operateall", operateall)
	output, _ := json.Marshal(&countOperate)
	fmt.Println(string(output))
	respondWithJSON(w, http.StatusOK, countOperate)
}
