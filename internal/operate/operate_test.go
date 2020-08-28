package operate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_zero_add_zero(t *testing.T) {

	//Add(input1,input2,condition)

	result, err := Add(0, 0, "+")

	assert.Equal(t, 0.00, result, "result 0.00")
	assert.NoError(t, err)
}

func Test_nine_add_nine(t *testing.T) {

	//Add(input1,input2,condition)

	result, err := Add(9, 9, "+")

	assert.Equal(t, 18.00, result, "result 18.00")
	assert.NoError(t, err)
}

func Test_zero_diff_zero(t *testing.T) {

	//Add(input1,input2,condition)

	result, err := Add(0, 0, "-")

	assert.Equal(t, 0.00, result, "result 0.00")
	assert.NoError(t, err)
}

func Test_nine_diff_nine(t *testing.T) {

	//Add(input1,input2,condition)

	result, err := Add(9, 9, "-")

	assert.Equal(t, 0.00, result, "result 0.00")
	assert.NoError(t, err)
}

func Test_four_diff_five(t *testing.T) {

	//Add(input1,input2,condition)

	result, err := Add(4, 5, "-")

	assert.Equal(t, -1.00, result, "result -1")
	assert.NoError(t, err)
}

func Test_four_mul_five(t *testing.T) {

	//Add(input1,input2,condition)

	result, err := Add(4, 5, "*")

	assert.Equal(t, 20.00, result, "result 20")
	assert.NoError(t, err)
}

func Test_twenth_mul_five(t *testing.T) {

	//Add(input1,input2,condition)

	result, err := Add(20, 5, "/")

	assert.Equalf(t, 4.00, result, "result 4")
	assert.NoError(t, err)
}

func Test_twenth_mul_zero(t *testing.T) {

	//Add(input1,input2,condition)

	result, err := Add(20, 0, "/")

	assert.Equal(t, 0.00, result, "result 0")
	assert.EqualError(t, err, "error")
}
