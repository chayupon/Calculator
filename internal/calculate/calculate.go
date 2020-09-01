package calculate

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/chayupon/Calculator/internal/operate"

	"net/http"

	"log"
	"time"

	"github.com/gorilla/mux"
)

//App input to db
type App struct {
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

//Initialize connect db
func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	fmt.Println("hello")
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
	a.Router.HandleFunc("/calculate/detail", a.Detail).Methods("GET")

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//Calculate Insert
func (a *App) Calculate(w http.ResponseWriter, r *http.Request) {

	var cal Cal
	var out Output
	var outerror outputError
	//var history History
	sqlStr := `INSERT INTO history( "time", input1, operate, input2, result, errordescripe) VALUES($1,$2,$3,$4,$5,$6)`

	e := json.NewDecoder(r.Body).Decode(&cal)
	if e != nil {

		outerror.Errordescription = e.Error()
		//outerror.Errordescription = err
		fmt.Println(outerror.Errordescription)

		s := fmt.Sprintf("%f %s %f = %f", cal.Input1, cal.Operation, cal.Input2, cal.Result)
		fmt.Println("ResultAll :", s)
		outerror.InputAll = s
		respondWithJSON(w, http.StatusBadRequest, outerror)

		return

	}

	result, err := operate.Add(cal.Input1, cal.Input2, cal.Operation)
	fmt.Println("Result :", result, err)
	currentime := time.Now()
	out.Time = currentime.Format(time.RFC3339Nano)
	out.Result = result
	if err != nil {

		outerror.Errordescription = err.Error()
		fmt.Println("error :", outerror.Errordescription)
		s := fmt.Sprintf("%f %s %f = %f", cal.Input1, cal.Operation, cal.Input2, cal.Result)
		fmt.Println("ResultAll :", s)
		outerror.InputAll = s

	}
	//save insert db
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

//History in db
type history struct {
	Sequence      int    `json:"sequence"`
	Time          string `json:"time"`
	InputAll      string `json:"input_all"`
	ErrorDescripe string `json:"error_descripe"`
}

//Detail select db
func (a *App) Detail(w http.ResponseWriter, r *http.Request) {

	sqlStr := `SELECT  sequence, "time", input1, operate, input2, result, errordescripe FROM history`
	rows, err := a.DB.Query(sqlStr)
	if err != nil {
		log.Println("Fail")
		return
	}
	defer rows.Close()
	h := []history{}
	//fmt.Printf("%+v",u)
	for rows.Next() {
		//var username string
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
		//fmt.Println("inputall :", s)
		his := history{
			Sequence:      sequence,
			Time:          time,
			InputAll:      s,
			ErrorDescripe: errordescripe,
		}
		h = append(h, his)
	}
	if !rows.NextResultSet() {
		log.Println(rows.Err())
	}
	output, _ := json.Marshal(&h)
	fmt.Println(string(output))
	respondWithJSON(w, http.StatusOK, h)

}
