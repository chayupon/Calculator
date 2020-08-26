package main
import (
	"fmt"
	"github.com/chayupon/Calculator/internal/service/operate"
	_ "github.com/lib/pq"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	
	getCons:
        var cons string
	fmt.Print("Please select an operation: +, -, *, / : ")
        fmt.Scanln(&cons)

	var input1 int
	fmt.Print("Please input the first number: ")
	fmt.Scanln(&input1)

	var input2 int
	fmt.Print("Please input the second number: ")
	fmt.Scanln(&input2)
	
switch cons {
	case "+":
		fmt.Print("Result: ")
		fmt.Println(operate.Add(input1,input2,cons))
	case "-":
		fmt.Print("Result: ")
		fmt.Println(operate.Add(input1,input2,cons))

	case "*":
		fmt.Print("Result: ")
		fmt.Println(operate.Add(input1,input2,cons))

	case "/":
		fmt.Print("Result: ")
		fmt.Println(operate.Add(input1,input2,cons))

	default:
		fmt.Println("Invalid operation selected. Please try again!")
		goto getCons
	}
	
	r := gin.Default()
	r.GET("/calculator", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Input1": input1,
			"Input2": input2,
			"Cons": cons,
		})
	})
	r.Run()
	
}


	// POST /calculator ,body >> {"input1":0,"input2":0,"operator":""}

	// input number > 0-9 , other error

	// +,-,*,/
	// result

	//case 1 0+0=0
