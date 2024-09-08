// PROBLEMA:
//          como proximo passo, implemente que, durante a inundação de ida, a rota vai sendo
//          gravada.  A rota é a sequencia de nodos por onde a mensagem passa.
//          Pode ser uma pilha de inteiros.  Cada nodo antes de repassar, empilha seu id.
//          Desta forma, a resposta pode ser enviada somente pela rota de retorno.
//          Ou seja, a mensagem trafega pela rota reversa.  Basta que cada nodo intermediario
//          desempilhe o identificador do proximo e repasse a mensagem para este.
// ATENCAO
//          o codigo abaixo NAO apresenta a solucao.  é uma cópia do anterior.
//          VOCE DEVE DESENVOLVER ESTA SOLUCAO.
//          este desenvolvimento pode ser com seu grupo.
//          deverá ser entregue em data marcada.

package main

import (
	"fmt"
)

const N = 10 
const channelBufferSize = 5

type Topology [N][N]int

type Message struct {
	id       int   
	source   int     
	receiver int      
	data     string  
	route    []int  // Adicionamos a rota da mensagem
}

type inputChan [N]chan Message

type nodeStruct struct {
	id               int 
	topo             Topology
	inCh             inputChan
	received         map[int]Message 
	receivedMessages []Message      
}

func (n *nodeStruct) broadcast(m Message) {
	for j := 0; j < N; j++ {
		if n.topo[n.id][j] == 1 {
			n.inCh[j] <- m
		}
	}
}

func (n *nodeStruct) nodo() {
	fmt.Println(n.id, " ativo! ")
	for {
		m := <-n.inCh[n.id]
		if m.receiver == n.id {
			n.receivedMessages = append(n.receivedMessages, m)
			// Se a mensagem é de ida, envia a resposta pela rota reversa
			if m.id > 0 {
				fmt.Println("                                   ", n.id, " recebe de ", m.source, " msg ", m.id, "  ", m.data)

				// Cria a mensagem de resposta com a rota reversa
				response := Message{
					id:       -m.id,
					source:   n.id,
					receiver: m.source,
					data:     "resp to msg",
					route:    reverseRoute(m.route, n.id), // Chama a função que inverte a rota
				}
				go n.broadcast(response)
			} else {
				fmt.Println("                                                                      ", n.id, " recebe de ", m.source, " msg ", m.id, "  ", m.data)
			}
		} else {
			_, achou := n.received[m.id]
			if !achou {
				// Adiciona o ID do nodo à rota
				m.route = append(m.route, n.id)
				n.received[m.id] = m
				go n.broadcast(m)
			}
		}
	}
}

// Lógica de reversão da rota
func reverseRoute(route []int, currentID int) []int {
	reversed := []int{}
	for i := len(route) - 1; i >= 0; i-- {
		reversed = append(reversed, route[i])
	}
	// Adiciona o próprio ID como útlima posição
	reversed = append(reversed, currentID)
	return reversed
}

func carga(nodoInicial int, inCh chan Message) {
	for i := 1; i < 10; i++ {
		inCh <- Message{(nodoInicial * 1000) + i, nodoInicial, i, "req", []int{}}
	}
}

func main() {
	var topo Topology
	topo = [N][N]int{
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 1, 0, 0, 0, 1},
		{0, 0, 0, 0, 1, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
		{0, 0, 0, 0, 0, 1, 0, 0, 1, 0},
	}

	var inCh inputChan
	for i := 0; i < N; i++ {
		inCh[i] = make(chan Message, channelBufferSize)
	}

	for id := 0; id < N; id++ {
		n := nodeStruct{id, topo, inCh, make(map[int]Message), []Message{}}
		go n.nodo()
	}

	go carga(0, inCh[0])
	carga(5, inCh[5])

	for{}
}