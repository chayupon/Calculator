package main

import (
	"github.com/chayupon/Calculator/internal/covid"
	_ "github.com/lib/pq"
	//	"strconv"
)

func main() {
	var a covid.App
	a.Initialize("db", "5432", "postgres", "tonkla727426", "calculator")
	a.Run(":8090")

}
