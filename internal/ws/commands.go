package ws

import (
	"fmt"
	"github.com/jsawo/colin-server/internal/cache"
	"strings"
	"time"

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

	var topicsSubscribed []string
	for topic := range model.CollectorInstances {
		for _, topicRequest := range topics {
			if topicRequest == topic || topicRequest == "*" {
				if model.SubscribeClient(remoteAddr, topic) {
					topicsSubscribed = append(topicsSubscribed, topic)
				}
			}
		}
	}

	for _, topic := range topicsSubscribed {
		SendMessageToRecipients(Message{
			Topic:      "_info",
			Payload:    "OK - subscribed to " + topic,
			Timestamp:  time.Now(),
			Recipients: []string{remoteAddr},
		})

		SendMessageToRecipients(Message{
			Topic:      topic,
			Payload:    cache.GetLatestValue(topic),
			Timestamp:  time.Now(),
			Recipients: []string{remoteAddr},
		})
	}
}
