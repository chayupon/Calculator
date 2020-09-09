func (a *App) CountDetail(w http.ResponseWriter, r *http.Request) {
	sqlStr := `SELECT operate FROM history`
	rows, err := a.DB.Query(sqlStr)
	if err != nil {
		log.Println("Fail", err)
		fmt.Println("error")

		return
	}
	defer rows.Close()
	resulthistory := []result{}
	var operateall []string

	//fmt.Printf("%+v",u)
	for rows.Next() {
		//	var count int
		var operate string
		//	var errordescripe string

		if err := rows.Scan(&operate); err != nil {

			log.Println(err)
			respondWithJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		//	fmt.Println("operate :",op)
		//history.
		//log.Println("inputall :", s)
		operateall = append(operateall, operate)

		//resulthistory= append(resulthistory, re)
	}
	counthis := countHistory{
		Resultlist: resulthistory,
	}
	countadd :=0
	countdiff :=0
	countmulti :=0
	countdiv :=0
	for _,operate := range operateall{
		if operate =="+"{
			countadd++
		}else if operate =="-"{
			countdiff++
		}else if operate =="*"{
			countmulti++
		}else if operate =="/"{
			countdiv++
		}
	}
	fmt.Println("count1:",countadd,"count2:",countdiv,"count3:",countmulti,"count4:",countdiff)
	fmt.Println("operateall", operateall)
	output, _ := json.Marshal(&counthis)
	fmt.Println(string(output))
	respondWithJSON(w, http.StatusOK, counthis)

}



	re :=result{
			Operate: operate,
			Count: count,
			ErrorDescripe: errordescripe,

		}
		resulthistory= append(resulthistory, re)