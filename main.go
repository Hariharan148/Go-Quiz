package main

import (
	"fmt"
	"flag"
	"os"
	"time"
	"encoding/csv"

)

type Problem struct{
	q string
	a string
}

func problemPuller(file string)([]Problem, error){

	if file, err := os.Open(file); err ==nil {

		csvR := csv.NewReader(file)
		if lines, err := csvR.ReadAll(); err == nil{
			return problemParser(lines), nil
		} else{
			fmt.Println("Problem while reading the file: ", err)
			return nil, err
		} 
	} else{
		fmt.Println("Problem while opening the file: ", err)
		return nil, err
	}


}


func problemParser(lines [][]string)[]Problem{
	problems := make([]Problem, len(lines))

	for i := 0; i < len(lines); i++ {
		problems[i] = Problem{
			q: lines[i][0],
			a: lines[i][1],
		}
	}
	return problems

}

func main(){

	//1. Get the file name
	fileName := flag.String("f", "quiz.csv", "path of the file")

	//2. Set the timer
	timer := flag.Int("t", 30, "Time to end the program")

	flag.Parse()
	//3. Call problemPuller
	problems, err := problemPuller(*fileName)

	//4. HandleError
	if err != nil {
		exit(err)
	}

	tObj := time.NewTimer(time.Duration(*timer) * time.Second)

	correctAns := 0

	ansC := make(chan string)
	//5. loop through the problems 
	problemLoop:

		for i, problem := range problems{
			fmt.Printf("Problem %d: %s = ", i+1, problem.q)
			var answer string

			go func(){
				fmt.Scan(&answer)
				ansC <-answer
			}()
			select{
				case <-tObj.C:
					break problemLoop
				case ans := <-ansC:
					if ans == problem.a{
						correctAns++
					}
					if i+1 == len(problems){
						close(ansC)
					}
			}
		}
	//6. get the answer and check the answer + check weather the timer has end concurrently using go routines
	
	//7. Print the result

	fmt.Printf("Your result is %d out of %d", correctAns, len(problems))

}


func exit(err error){
	fmt.Printf("Error: ", err.Error())
	os.Exit(1)
}