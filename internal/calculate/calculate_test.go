package calculate

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func CheckEror(t *testing.T) {

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/calculator", strings.NewReader(`{  "input1" :0.22, "input2"  :0,"operation" :"/"}`))

	Calculate(response, request)

	var outerr outputError

	json.NewDecoder(response.Body).Decode(&outerr)
	// test
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "error", outerr.Errordescription)

}

func CheckStructEror(t *testing.T) {

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/calculator", strings.NewReader(`{  "input1" :"2", "input2"  :2,"operation" :"+"}`))

	Calculate(response, request)

	var outerr outputError

	json.NewDecoder(response.Body).Decode(&outerr)
	// test
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "json: cannot unmarshal string into Go struct field Cal.input1 of type float64", outerr.Errordescription)
	assert.Equal(t, "0.000000 + 2.000000 = 0.000000", outerr.InputAll)

}

func CheckPass(t *testing.T) {

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/calculator", strings.NewReader(`{  "input1" :2, "input2"  :2,"operation" :"+"}`))

	Calculate(response, request)

	var showout Output

	json.NewDecoder(response.Body).Decode(&showout)

	// test
	assert.Equal(t, 200, response.Code)
	assert.Equalf(t, 4.00, showout.Result, "result 4")

}

