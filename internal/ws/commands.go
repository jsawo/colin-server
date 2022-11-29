package ws

import (
	"fmt"
	"strings"

	"github.com/jsawo/colin-server/internal/model"
)

func handleClientCommands(remoteAddr, command string) {
	if strings.HasPrefix(command, "SUB ") {
		command = strings.TrimPrefix(command, "SUB ")
		topics := strings.Split(command, ",")
		handleClientSubscription(remoteAddr, topics)
	} else {
		fmt.Printf("Unrecognized client command (%v) %q \n", remoteAddr, command)
	}
}

func handleClientSubscription(remoteAddr string, topics []string) {
	fmt.Printf("Client %v requests subscription to %q \n", remoteAddr, topics)

	for topic, _ := range model.CollectorInstances {
		for _, topicRequest := range topics {
			if topicRequest == topic || topicRequest == "*" {
				model.SubscribeClient(remoteAddr, topic)
			}
		}
	}
}
