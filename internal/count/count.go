package count

import "fmt"

//CheckOperate count
func CheckOperate(cons string) (int, error) {
	count := 0
	if cons == "+" {

		return count, nil
	} else if cons == "*" {

		return count, nil
	} else if cons == "/" {

		return count, nil
	} else if cons == "-" {

		return count, nil
	}

	return 0, fmt.Errorf("Invalid Operate")

}
