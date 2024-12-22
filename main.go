package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    string
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Println("Seja bem vindo(a) ao QuizGo!")
	fmt.Println("Escreva o seu nome:")
	reader := bufio.NewReader(os.Stdin)

	name, err := reader.ReadString('\n')
	if err != nil {
		panic("Erro ao ler o nome")
	}

	g.Name = name

	fmt.Printf("Vamos ao jogo %s", g.Name)
}

func (g *GameState) ProccessCSV() {
	file, err := os.Open("./quiz.csv")
	if err != nil {
		panic("Erro ao ler arquivo csv")
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic("Erro ao ler csv")
	}

	for i, record := range records {
		fmt.Println(i, record)

		if i > 0 {
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  toInt(record[5]),
			}

			g.Questions = append(g.Questions, question)
		}
	}
}

func toInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}

	return i
}

func main() {
	game := &GameState{}
	go game.ProccessCSV()
	game.Init()
	fmt.Println(game)
}
