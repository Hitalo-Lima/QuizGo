package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	quizConhecimentosGerais string = "./files/quiz_conhecimentos_gerais.csv"
	quizHistoria            string = "./files/quiz_historia.csv"
	quizIngles              string = "./files/quiz_ingles.csv"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Theme     string
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

	validOption := false

	for !validOption {
		fmt.Println("Escolha uma opção de tema para as perguntas")
		fmt.Print("[1] Conhecimentos Gerais [2] História [3] Inglês: ")

		option, err := reader.ReadString('\n')
		if err != nil {
			panic("Erro ao ler a opção escolhida")
		}

		switch strings.TrimSpace(option) {
		case "1":
			g.Theme = quizConhecimentosGerais
			validOption = true
		case "2":
			g.Theme = quizHistoria
			validOption = true
		case "3":
			g.Theme = quizIngles
			validOption = true
		default:
			fmt.Println("Escolha uma opção válida!")
		}
	}

	fmt.Printf("Vamos ao jogo %s!\n", strings.TrimSpace(name))
}

func (g *GameState) ProccessCSV() {
	file, err := os.Open(g.Theme)
	if err != nil {
		panic("Erro ao ler arquivo csv")
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("erro: ", err)
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
		fmt.Printf("\033[34mQuestão %d: %s\033[0m\n\n", index+1, question.Text)

		// Iterar sobre as opções e exibir no terminal
		for pos, option := range question.Options {
			fmt.Printf("[%d] %s\n", pos+1, option)
		}

		fmt.Print("\nDigite uma alternativa: ")

		timeout := make(chan bool, 1)
		answerChan := make(chan int, 1)

		// Go routines para o timer e para entrada do usuário
		go func() {
			timer := time.NewTimer(time.Second * 15)
			<-timer.C
			timeout <- true
		}()

		go func() {
			reader := bufio.NewReader(os.Stdin)
			for {
				read, _ := reader.ReadString('\n')
				answer, err := toInt(strings.TrimSpace(read))

				// Se for uma resposta válida, sai do loop e envia para o canal
				if err == nil && answer >= 1 && answer <= len(question.Options) {
					answerChan <- answer
					return
				}

				// Se a resposta for inválida, avisa o jogador
				fmt.Println("Erro: Insira uma opção válida!")
			}
		}()

		// Espera até que um dos dois canais receba uma resposta: entrada do usuário ou timeout
		select {
		case answer := <-answerChan:
			if answer == question.Answer {
				fmt.Printf("\033[32mParabéns! Resposta correta.\033[0m\n")
				g.Points += 10
			} else {
				fmt.Printf("\033[31mResposta incorreta!\033[0m\n")
			}
		case <-timeout:
			fmt.Println("\033[31mTempo esgotado! Resposta não recebida.\033[0m")
			return
		}
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
	game.Init()
	game.ProccessCSV()
	game.Run()

	fmt.Printf("Fim de jogo! Você fez %d de %d pontos!", game.Points, len(game.Questions)*10)
}
