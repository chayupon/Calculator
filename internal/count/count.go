package count

import "fmt"

//Count operate
func Count(cons string) (int, error) {
	count := 0
	if cons == "+" {

		count++
		return count, nil
	} else if cons == "*" {

		count++
		return count, nil
	} else if cons == "/" {

		count++
		return count, nil
	} else if cons == "-" {

		count++
		return count, nil
	}

	return 0, fmt.Errorf("Invalid Operate")

}
