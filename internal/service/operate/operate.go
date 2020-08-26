package operate

//Add input
func Add(input1 int, input2 int, cons string) (int, string) {

	if cons == "+" {

		result := input1 + input2

		return result, ""
	} else if cons == "*" {

		result := input1 * input2

		return result, ""
	} else if cons == "/" {
		if input2 != 0 {
			result := input1 / input2
			return result, ""
		}
		return 0, "error"
	}

	result := input1 - input2

	return result, ""
	
}