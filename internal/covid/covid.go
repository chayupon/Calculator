package covid

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//App input to db
type App struct {
	Router     *mux.Router
	DB         *sql.DB
	DateJoined time.Time
}

//Cal calculate
type Request struct {
	Input1 float64 `json:"age"`
	Input2 string  `json:"subprovince"`
	Input3 string  `json:"province"`
	//Count     int     `json:"count"`
}

//ResponseCount count
type ResponseBody struct {
	Province province `json:"Province"`
	AgeGroup resp     `json:"AgeGroup"`
}

type resp struct {
	Age1       int `json:"0-30"`
	Age2       int `json:"31-60"`
	Age3       int `json:"61+"`
	AnotherAge int `json:"N/A"`
}
type province struct {
	Name1 int `json:"Samut Sakhon"`
	Name2 int `json:"Bangkok"`
}

//Output Value
type Output struct {
	Result int    `json:"result"`
	Time   string `json:"time"`
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
	a.Router.HandleFunc("/inputdata", a.Process).Methods("POST")
	a.Router.HandleFunc("/inputdata/count/detail", a.CountHistoryAll).Methods("GET")
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//check data to count
func CountRequest(ageAll []float64, subprovinceAll []string, provinceAll []string) (int, int, int, int, int, int) {
	countage := 0
	countage2 := 0
	countage3 := 0
	countage4 := 0
	countsub := 0
	countpro := 0

	//check sub province
	for _, data := range subprovinceAll {
		if data == "Samut Sakhon" {
			countsub++
		}
	}

	//check provice
	for _, data2 := range provinceAll {
		if data2 == "Bangkok" {
			countpro++
		}
	}

	for _, data3 := range ageAll {
		if data3 <= 30 && data3 > 0 {
			countage++
		} else if data3 <= 60 && data3 >= 31 {
			countage2++
		} else if data3 >= 61 {
			countage3++
		} else {
			countage4++
		}
	}
	return countage, countage2, countage3, countage4, countsub, countpro
}

//Insert data
func (a *App) Process(w http.ResponseWriter, r *http.Request) {
	var rq Request
	var out Output
	// var outerror outputError
	sqlStr := `INSERT INTO history(time,age,subprovince,province) VALUES($1,$2,$3,$4)`
	if sqlStr == "" {
		fmt.Println("ERROR")
	}
	e := json.NewDecoder(r.Body).Decode(&rq)
	if e != nil {
		s := fmt.Sprintf("%f %s %s", rq.Input1, rq.Input2, rq.Input3)
		fmt.Println("ResultAll :", s)
		respondWithJSON(w, http.StatusBadRequest, s)
		return
	}

	currentime := time.Now()
	out.Time = currentime.Format(time.RFC3339)

	_, err := a.DB.Exec(sqlStr, out.Time, rq.Input1, rq.Input2, rq.Input3)
	if err != nil {
		fmt.Println(err)
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, http.StatusOK, out.Time)
}

//count history
func (a *App) CountHistoryAll(w http.ResponseWriter, r *http.Request) {
	sqlStr := `SELECT  sequence,age,subprovince,province FROM history`
	rows, err := a.DB.Query(sqlStr)
	if err != nil {
		log.Println("Fail", err)
		fmt.Println("error")
		return
	}
	defer rows.Close()

	var operateall []float64
	var subprovinceAll []string
	var provinceAll []string

	for rows.Next() {
		//  var count int
		var sequence int
		var age float64
		var subprovince string
		var province string
		//var age int
		if err := rows.Scan(&sequence, &age, &subprovince, &province); err != nil {
			log.Println(err)
			respondWithJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		operateall = append(operateall, age)
		subprovinceAll = append(subprovinceAll, subprovince)
		provinceAll = append(provinceAll, province)
	}
	countage, countage2, countage3, countage4, countsub, countpro := CountRequest(operateall, subprovinceAll, provinceAll)
	countop := ResponseBody{
		Province: province{
			Name1: countsub,
			Name2: countpro,
		},
		AgeGroup: resp{

			Age1: countage,

			Age2: countage2,

			Age3: countage3,

			AnotherAge: countage4,
		},
	}
	fmt.Println("operateall", operateall)
	output, _ := json.Marshal(&countop)
	fmt.Println(string(output))
	respondWithJSON(w, http.StatusOK, countop)
}
