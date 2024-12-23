package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Println("--------------------------------")
	fmt.Println("| Seja bem vindo(a) ao QuizGo! |")
	fmt.Println("--------------------------------")
	fmt.Print("Digite o seu nome: ")
	reader := bufio.NewReader(os.Stdin)

	name, err := reader.ReadString('\n')
	if err != nil {
		panic("Erro ao ler o nome")
	}

	g.Name = name

	fmt.Printf("Vamos ao jogo %s\n!", g.Name)
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

		if i > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}

			g.Questions = append(g.Questions, question)
		}
	}
}

func (g *GameState) Run() {
	// Exibir a pergunta para o usuário
	for index, question := range g.Questions {
		fmt.Printf("\033[34mQuestão %d: %s\033[0m\n", index+1, question.Text)

		// Iterar sobre as opções que existem no game state
		// e exibir no terminal
		for pos, option := range question.Options {
			fmt.Printf("[%d] %s\n", pos+1, option)
		}

		fmt.Println("\nDigite uma alternativa")

		// Coletar a entrada do usuário
		// Validar o caractere que foi inserido
		// Se for inválido o usuário deve inserir novamente
		var answer int
		var err error

		for {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')

			answer, err = toInt(strings.TrimSpace(read))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			break
		}

		if answer == question.Answer {
			fmt.Println("Parabéns! Resposta correta.")
			g.Points += 10
			continue
		}
		fmt.Println("Resposta incorreta!")
		break
	}
}

func toInt(str string) (int, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.New("insira um número válido")
	}

	return i, nil
}

func main() {
	game := &GameState{}
	go game.ProccessCSV()
	game.Init()
	game.Run()

	fmt.Printf("Fim de jogo! Você fez %d pontos\n", game.Points)
}
