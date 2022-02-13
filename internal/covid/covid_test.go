package covid

import (

	//"fmt"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	//	_ "github.com/lib/pq"

	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var a App

func Test_ProcessError(t *testing.T) {

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/inputdata", strings.NewReader(`{  "input1" :12, "input2"  :"lampang","input3" :"b"}`))

	dbTest, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbTest.Close()
	sqlStr := `INSERT INTO history `
	currentime := time.Now()
	timeTest := currentime.Format(time.RFC3339)
	//timeTest := currentime.Format("RFC1111119")
	a.DateJoined = currentime
	mock.ExpectExec(sqlStr).WithArgs(timeTest, 12, "lampang", "b", "Invalid Syntax").WillReturnResult(sqlmock.NewResult(1, 1))
	a.DB = dbTest
	a.Process(response, request)
	assert.Equal(t, 400, response.Code)

}

func Test_Pass(t *testing.T) {

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/inputdata", strings.NewReader(`{  "input1" :12.0000, "input2"  :"lampang","input3" :"Bangkok"}`))
	dbTest, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbTest.Close()
	sqlStr := `INSERT INTO history `
	currentime := time.Now()
	timeTest := currentime.Format(time.RFC3339)
	mock.ExpectExec(sqlStr).WithArgs(timeTest, float64(12), "lampang", "Bangkok").WillReturnResult(sqlmock.NewResult(1, 1))
	a.DB = dbTest

	a.Process(response, request)

	var showout Output

	json.NewDecoder(response.Body).Decode(&showout)

	fmt.Println("result is :", float64(showout.Result))
	//assert.Equal(t, 200, response.Code)
	assert.Equalf(t, 0, showout.Result, "result 4.0000")

}

func Test_CountAll_Add(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/inputdata/count/detail", strings.NewReader(`{"input3" :"Bangkok"},`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "age", "province", "subprovince"}).
		AddRow(1, float64(12), "Bangkok", "lampang")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.CountHistoryAll(response, request)

	expect := "{\"Province\":{\"Samut Sakhon\":0,\"Bangkok\":0},\"AgeGroup\":{\"0-30\":1,\"31-60\":0,\"61+\":0,\"N/A\":0}}"
	t.Log(response.Body.String())
	assert.Equal(t, expect, response.Body.String())
	assert.Equal(t, 200, response.Code)

}

func Test_CountAll_Error(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/inputdata/count/detail", strings.NewReader(`{"0-30" :1}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}).
		AddRow("a", "12")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.CountHistoryAll(response, request)

	expect := "\"sql: expected 2 destination arguments in Scan, not 4\""
	t.Log(response.Body.String())

	assert.Equal(t, expect, response.Body.String())
	assert.Equal(t, 400, response.Code)

}
