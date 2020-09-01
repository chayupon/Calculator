package main

import (
	_ "github.com/lib/pq"
	"github.com/chayupon/Calculator/internal/calculate"
	//	"strconv"
)

func main() {
	var a calculate.App
	a.Initialize("postgres", "tonkla727426", "calculator")
	a.Run(":8090")
}
