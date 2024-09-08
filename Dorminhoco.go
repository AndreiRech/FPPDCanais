package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NJ = 5
const M = 4 

type carta string

var ch [NJ]chan carta

func escolherCarta(mao []carta) int {
	count := make(map[carta]int)
	for _, value := range mao {
		count[value]++
	}

	minCarta := mao[0]
	for _, value := range mao {
		if count[value] < count[minCarta] {
			minCarta = value
		}
	}

	for i, value := range mao {
		if value == minCarta {
			return i
		}
	}
	return -1
}

func bateu(mao []carta) bool {
	for i := 0; i < len(mao); i++ {
		count := 0
		for j := 0; j < len(mao); j++ {
			if mao[i] == mao[j] || mao[j] == "@" {
				count++
			}
		}
		if count >= 4 {
			return true
		}
	}
	return false
}

func jogador(id int, in chan carta, out chan carta, cartasIniciais []carta, acabou, batida chan int) {
	mao := cartasIniciais
	batido := false

	for {
		if batido {
			cartaRecebida := <-in
			//fmt.Printf("Jogador batido %d recebeu a carta: %s\n", id, cartaRecebida)
			out <- cartaRecebida
		} else {
			select {
			case <-batida:
				batida <- id
				acabou <- id
				batido = true
				fmt.Printf("Jogador %d bateu e saiu do jogo\n", id)
			default:
				cartaRecebida := <-in
				//fmt.Printf("Jogador %d recebeu a carta: %s\n", id, cartaRecebida)
				mao = append(mao, cartaRecebida)
				if bateu(mao) {
					batida <- id
					acabou <- id
					batido = true
					fmt.Printf("Jogador %d bateu e saiu do jogo\n", id)
				}
				cartaParaSair := escolherCarta(mao)
				//fmt.Printf("Jogador %d enviou a carta: %s\n", id, mao[cartaParaSair])
				out <- mao[cartaParaSair]
				mao = append(mao[:cartaParaSair], mao[cartaParaSair+1:]...)
			}
		}
	}
}

func embaralharBaralho(baralho []carta) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(baralho), func(i, j int) {
		baralho[i], baralho[j] = baralho[j], baralho[i]
	})
}

func main() {
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta)
	}

	baralho := []carta{"A", "A", "A", "A",
		"B", "B", "B", "B", "C", "C", "C", "C",
		"D", "D", "D", "D", "E", "E", "E", "E", "@"}

	embaralharBaralho(baralho[:])

	batida := make(chan int)
	acabou := make(chan int, 5)

	for i := 0; i < NJ; i++ {
		cartasEscolhidas := baralho[i*M : (i+1)*M]
		go jogador(i, ch[i], ch[(i+1)%NJ], cartasEscolhidas, batida, acabou)
	}

	// Inicia o jogo enviando a primeira carta
	ch[0] <- baralho[NJ*M]

	// Continuar até restar um jogador
	for i := 0; i < NJ-1; i++{
		fmt.Printf("Jogador %d terminou e está liberado\n", <-acabou)
	}

	fmt.Print("Acabou")
}