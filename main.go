package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// questions,answers
type QA struct {
	csvFile   string
	questions []string
	answers   []string
}

func shufle(arr [][]string) {
	for i := range arr {
		source := rand.NewSource(time.Now().UnixNano())
		randNum := rand.New(source)
		newInd := randNum.Intn(len(arr) - 1)
		arr[i], arr[newInd] = arr[newInd], arr[i]
	}
}

// read csv file and get questions and answers
func (qa *QA) getQA(shufleQN bool) ([]string, []string) {
	fl, err := os.Open(qa.csvFile)
	if err != nil {
		fmt.Println("Err: ", err)
	}
	reader := csv.NewReader(fl)
	records, _ := reader.ReadAll()
	if shufleQN {
		shufle(records)
	}
	// get question[0] and answer[1]
	for _, record := range records {
		qa.questions = append(qa.questions, record[0])
		qa.answers = append(qa.answers, record[1])
	}
	return qa.questions, qa.answers
}

type Quiz struct {
	total_quizes    float64
	num_of_corrects float64
	input           *bufio.Reader
}

func (q Quiz) readInput() string {
	s, _, _ := q.input.ReadLine()
	return string(s)
}
func (q Quiz) calcScore() float64 {
	return (q.num_of_corrects / q.total_quizes) * 100
}

func main() {
	csv_file := flag.String("csv", "D:\\BBserver\\Codes\\QuizGame\\quiezes.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 5, "the time limit for the quiz in seconds")
	shufleQN := flag.Bool("shuffle", true, "Shuffle the questions")

	flag.Parse()

	qn_ans := QA{
		csvFile:   *csv_file,
		questions: []string{},
		answers:   []string{},
	}
	// if you want shuffle questions
	qn, an := qn_ans.getQA(*shufleQN)

	quiz := Quiz{
		total_quizes:    float64(len(qn_ans.questions)),
		num_of_corrects: 0,
		input:           bufio.NewReader(os.Stdin),
	}
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	for ind, q := range qn {
		fmt.Printf("#Quetion[%d]: %s = ", ind+1, q)
		ansChan := make(chan string)
		go func() {
			ansChan <- quiz.readInput()
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nScore:%2.2f%%\n", quiz.calcScore())
			return
		case user := <-ansChan:
			answer := strings.ToLower(strings.TrimSpace(an[ind]))
			if answer == strings.ToLower(user) {
				fmt.Println("Correct!")
				quiz.num_of_corrects += 1
			} else {
				fmt.Printf("Not correct! Answer is %s\n", answer)
			}
		}
	}
}
