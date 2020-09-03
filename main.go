package main

import (
	"github.com/chayupon/Calculator/internal/calculate"
	_ "github.com/lib/pq"
	//	"strconv"
)

func main() {
	var a calculate.App
	a.Initialize("db","5432","postgres", "tonkla727426", "calculator")
	a.Run(":8090")

}
