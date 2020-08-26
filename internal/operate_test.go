package operate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_zero_add_zero(t *testing.T) {

	//cal(input1,input2,condition)

	result, err := cal(0, 0, "+")

	assert.Equal(t, 0, result, "result 0")
	assert.Equal(t, "", err)
}

func Test_nine_add_nine(t *testing.T) {

	//cal(input1,input2,condition)

	result, err := cal(9, 9, "+")

	assert.Equal(t, 18, result, "result 18")
	assert.Equal(t, "", err)
}

func Test_zero_diff_zero(t *testing.T) {

	//cal(input1,input2,condition)

	result, err := cal(0, 0, "-")

	assert.Equal(t, 0, result, "result 0")
	assert.Equal(t, "", err)
}

func Test_nine_diff_nine(t *testing.T) {

	//cal(input1,input2,condition)

	result, err := cal(9, 9, "-")

	assert.Equal(t, 0, result, "result 18")
	assert.Equal(t, "", err)
}

func Test_four_diff_five(t *testing.T) {

	//cal(input1,input2,condition)

	result, err := cal(4, 5, "-")

	assert.Equal(t, -1, result, "result -1")
	assert.Equal(t, "", err)
}

func Test_four_mul_five(t *testing.T) {

	//cal(input1,input2,condition)

	result, err := cal(4, 5, "*")

	assert.Equal(t, 20, result, "result 20")
	assert.Equal(t, "", err)
}

func Test_twenth_mul_five(t *testing.T) {

	//cal(input1,input2,condition)

	result, err := cal(20, 5, "/")

	assert.Equal(t, 4, result, "result 4")
	assert.Equal(t, "", err)
}

func Test_twenth_mul_zero(t *testing.T) {

	//cal(input1,input2,condition)

	result, err := cal(20, 0, "/")

	assert.Equal(t, 0, result, "result 0")
	assert.Equal(t, "error", err)
}