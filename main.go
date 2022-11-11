package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jsawo/colin-server/ws"
)

func main() {
	go StartServer()

	fmt.Println("boo")
	go generateNoise()

	<-done
}

func generateNoise() {
	rand.Seed(time.Now().UnixNano())
	for {
		n := rand.Intn(10)
		time.Sleep(time.Duration(n) * time.Second)
		ws.WriteMessage("dummy", fmt.Sprintf("Message delayed %v seconds\n", n))
	}
}
