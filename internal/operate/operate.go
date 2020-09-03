package operate

import "fmt"

//Add input
func Add(input1 float64, input2 float64, cons string) (float64, error) {

	if cons == "+" {

		result := input1 + input2

		return result, nil
	} else if cons == "*" {

		result := input1 * input2

		return result, nil
	} else if cons == "/" {
		if input2 != 0 {
			result := input1 / input2
			return result, nil
		}
		return 0, fmt.Errorf("error_divide_Zero")
	} else if cons == "-" {
		result := input1 - input2
		return result, nil
	}

	return 0, fmt.Errorf("Invalid Operate")

}
