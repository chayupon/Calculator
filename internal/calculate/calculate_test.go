package calculate

import (
	"encoding/json"
	"fmt"
	//"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	//	_ "github.com/lib/pq"
	"log"

	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var a App


func Test_OutputError(t *testing.T) {
	
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/calculator", strings.NewReader(`{  "input1" :0.22, "input2"  :0,"operation" :"/"}`))

	dbTest, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbTest.Close()
	sqlStr := `INSERT INTO history `
	currentime := time.Now()
	timeTest := currentime.Format(time.RFC3339)
	
	a.DateJoined = currentime
	mock.ExpectExec(sqlStr).WithArgs(timeTest, 4.00, "/", 0.00, 0.00, "error_divide_Zero").WillReturnResult(sqlmock.NewResult(1, 1))
	a.DB = dbTest
	a.Calculate(response, request)

	var outerr outputError

	json.NewDecoder(response.Body).Decode(&outerr)
	
	assert.Equal(t, 400, response.Code)


}

func Test_OperateError(t *testing.T) {
	
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/calculator", strings.NewReader(`{  "input1" :4.00, "input2"  :2,"operation" :"b"}`))

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
	mock.ExpectExec(sqlStr).WithArgs(timeTest, 4.00, "b", 2.00, 0.00, "Invalid Operate").WillReturnResult(sqlmock.NewResult(1, 1))
	a.DB = dbTest
	a.Calculate(response, request)
	assert.Equal(t, 400, response.Code)

}

func Test_StructEror(t *testing.T) {

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/calculator", strings.NewReader(`{  "input1" :"2", "input2"  :2,"operation" :"+"}`))

	a.Calculate(response, request)

	var outerr outputError

	json.NewDecoder(response.Body).Decode(&outerr)
	
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "json: cannot unmarshal string into Go struct field Cal.input1 of type float64", outerr.Errordescription)
	assert.Equal(t, "0.000000 + 2.000000 = 0.000000", outerr.InputAll)

}
func Test_Pass(t *testing.T) {

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/calculate", strings.NewReader(`{  "input1" :2, "input2"  :2,"operation" :"+"}`))
	dbTest, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbTest.Close()
	sqlStr := `INSERT INTO history `
	currentime := time.Now()
	timeTest := currentime.Format(time.RFC3339)
	mock.ExpectExec(sqlStr).WithArgs(timeTest, 2.00, "+", 2.00, 4.00, "").WillReturnResult(sqlmock.NewResult(1, 1))
	a.DB = dbTest

	a.Calculate(response, request)

	var showout Output

	json.NewDecoder(response.Body).Decode(&showout)

	
	fmt.Println("result is :", showout.Result)
	assert.Equal(t, 200, response.Code)
	assert.Equalf(t, 4.00, showout.Result, "result 4")

}

func Test_SelectDb(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/detail", strings.NewReader(`{  "input1" :2, "input2"  :2,"operation" :"+"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "time", "input1", "operate", "input2", "result", " errordescripe"}).
		AddRow(1, "2020-08-31T17:02:10.076232+07:00", 22.000000, "+", 3123.000000, 0.000000, "").
		AddRow(2, "2020-08-31T17:04:32.483087+07:00", 22.000000, "+", 31.000000, 0.000000, "")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.Detail(response, request)

	
	expect := `[{"sequence": 1,"time": "2020-08-31T17:02:10.076232+07:00","input_all": "22.000000 + 3123.000000 = 0.000000","error_descripe": ""},{"sequence": 2,"time": "2020-08-31T17:04:32.483087+07:00","input_all": "22.000000 + 31.000000 = 0.000000","error_descripe": ""}]`
	t.Log(response.Body.String())
	assert.JSONEq(t, expect, response.Body.String())
	
	assert.Equal(t, 200, response.Code)
}

func Test_SelectErrorType(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/detail", strings.NewReader(`{  "input1" :2, "input2"  :2,"operation" :"+"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "time", "input1", "operate", "input2", "result", " errordescripe"}). //RowError(1,fmt.Errorf("error"))
																AddRow("a", "2020-08-31T17:02:10.076232+07:00", 22, "+", 3123.000000, 0.000000, "")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.Detail(response, request)

	expect := "\"sql: Scan error on column index 0, name \\\"sequence\\\": converting driver.Value type string (\\\"a\\\") to a int: invalid syntax\""
	t.Log(response.Body.String())

	assert.Equal(t, expect, response.Body.String())
	assert.Equal(t, 400, response.Code)

}

func Test_StructOperateEror(t *testing.T) {

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/calculate/OperateRequest", strings.NewReader(`{ operation" :"+"}`))

	a.OperateRequest(response, request)

	var counterror CountOperateError

	json.NewDecoder(response.Body).Decode(&counterror)

	fmt.Println("error :", counterror)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "invalid character 'o' looking for beginning of object key string", counterror.Errordescription)

}
func Test_InvalidOperate(t *testing.T) {

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/calculate/OperateRequest", strings.NewReader(`{"operation" :"b"}`))

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
	mock.ExpectExec(sqlStr).WithArgs(timeTest, 0.00, "b", 0.00, 0.00, "Invalid Operate").WillReturnResult(sqlmock.NewResult(1, 1))
	a.DB = dbTest
	a.OperateRequest(response, request)
	assert.Equal(t, 400, response.Code)

}

func Test_OperateRequest_Add(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/OperateRequest", strings.NewReader(`{  "operation" :"+"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}).
		AddRow(1, "+").
		AddRow(2, "+")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.OperateRequest(response, request)

	expect := `{"operation": "+","count": 2}`
	t.Log(response.Body.String())
	assert.JSONEq(t, expect, response.Body.String())

	assert.Equal(t, 200, response.Code)
}

func Test_OperateRequest_Diff(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/OperateRequest", strings.NewReader(`{  "operation" :"-"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}).
		AddRow(1, "-").
		AddRow(2, "-")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.OperateRequest(response, request)

	expect := `{"operation": "-","count": 2}`
	t.Log(response.Body.String())
	assert.JSONEq(t, expect, response.Body.String())
	assert.Equal(t, 200, response.Code)
}

func Test_OperateRequest_Multi(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/OperateRequest", strings.NewReader(`{  "operation" :"*"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}).
		AddRow(1, "*")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.OperateRequest(response, request)

	
	expect := `{"operation": "*","count": 1}`
	t.Log(response.Body.String())
	assert.JSONEq(t, expect, response.Body.String())
	
	assert.Equal(t, 200, response.Code)
}

func Test_OperateRequest_Div(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/OperateRequest", strings.NewReader(`{  "operation" :"/"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}).
		AddRow(1, "/")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.OperateRequest(response, request)

	
	expect := `{"operation": "/","count": 1}`
	t.Log(response.Body.String())
	assert.JSONEq(t, expect, response.Body.String())
	
	assert.Equal(t, 200, response.Code)
}

func Test_OperateRequest_Error(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/OperateRequest", strings.NewReader(`{  "operation" :"/"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}).
		AddRow("a", "/")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.OperateRequest(response, request)

	
	expect := "\"sql: Scan error on column index 0, name \\\"sequence\\\": converting driver.Value type string (\\\"a\\\") to a int: invalid syntax\""
	t.Log(response.Body.String())
	assert.JSONEq(t, expect, response.Body.String())
	
	assert.Equal(t, 400, response.Code)
}

func Test_CountOperateAll_Add(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/count/detail", strings.NewReader(`{"operation" :"+"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}).
		AddRow(1, "+")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.CountOperateAll(response, request)

	expect := "{\"result\":[{\"operation\":\"+\",\"count\":1},{\"operation\":\"-\",\"count\":0},{\"operation\":\"*\",\"count\":0},{\"operation\":\"/\",\"count\":0}]}"
	t.Log(response.Body.String())
	assert.Equal(t, expect, response.Body.String())
	assert.Equal(t, 200, response.Code)

}

func Test_CountOperateAll_Diff(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/count/detail", strings.NewReader(`{"operation" :"-"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}).
		AddRow(1, "-")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.CountOperateAll(response, request)

	expect := "{\"result\":[{\"operation\":\"+\",\"count\":0},{\"operation\":\"-\",\"count\":1},{\"operation\":\"*\",\"count\":0},{\"operation\":\"/\",\"count\":0}]}"
	t.Log(response.Body.String())

	assert.Equal(t, expect, response.Body.String())
	assert.Equal(t, 200, response.Code)

}

func Test_CountOperateAll_Multi(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/count/detail", strings.NewReader(`{"operation" :"*"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}).
		AddRow(1, "*")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.CountOperateAll(response, request)

	expect := "{\"result\":[{\"operation\":\"+\",\"count\":0},{\"operation\":\"-\",\"count\":0},{\"operation\":\"*\",\"count\":1},{\"operation\":\"/\",\"count\":0}]}"
	t.Log(response.Body.String())

	assert.Equal(t, expect, response.Body.String())
	assert.Equal(t, 200, response.Code)

}

func Test_CountOperateAll_Div(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/count/detail", strings.NewReader(`{"operation" :"/"}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}). //RowError(1,fmt.Errorf("error"))
									AddRow(1, "/")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.CountOperateAll(response, request)

	expect := "{\"result\":[{\"operation\":\"+\",\"count\":0},{\"operation\":\"-\",\"count\":0},{\"operation\":\"*\",\"count\":0},{\"operation\":\"/\",\"count\":1}]}"
	t.Log(response.Body.String())

	assert.Equal(t, expect, response.Body.String())
	assert.Equal(t, 200, response.Code)

}

func Test_CountOperateAll_Error(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/calculate/count/detail", strings.NewReader(`{"operation" :1}`))

	dbTest, mock, err := sqlmock.New()
	a.DB = dbTest
	if err != nil {
		log.Print(err)
	}
	defer dbTest.Close()

	rows := sqlmock.NewRows([]string{"sequence", "operate"}).
		AddRow("a", "+")

	mock.ExpectQuery("^SELECT (.+) FROM history").
		WillReturnRows(rows)
	a.CountOperateAll(response, request)

	expect := "\"sql: Scan error on column index 0, name \\\"sequence\\\": converting driver.Value type string (\\\"a\\\") to a int: invalid syntax\""
	t.Log(response.Body.String())

	assert.Equal(t, expect, response.Body.String())
	assert.Equal(t, 400, response.Code)

}
