package main

import (
	"fmt"
	"log"

	"github.com/samuelngs/kubernetes-bot/bot"
	"github.com/samuelngs/kubernetes-bot/client"
)

func main() {

	bot, err := bots.Slack(
		bots.EnvUsername(),
		bots.EnvToken(),
		bots.EnvChannel(),
	)

	cli, err := client.New(
		client.EnvHost(),
		client.EnvUsername(),
		client.EnvPassword(),
		client.EnvInsecure(),
		client.EnvInterval(),
	)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case e := <-cli.Watch():
				switch e.Type() {
				case client.EventNewPod:
					table := bots.SlackTable(
						"good",
						bots.Field("Namespace", e.Namespace()),
						bots.Field("Message", fmt.Sprintf("Created pod: %s", e.GenerateName())),
						bots.Field("Node", e.NodeName()),
					)
					bot.Emit(table)
				case client.EventDelPod:
					table := bots.SlackTable(
						"warning",
						bots.Field("Namespace", e.Namespace()),
						bots.Field("Message", fmt.Sprintf("Deleted pod: %s", e.GenerateName())),
						bots.Field("Node", e.NodeName()),
					)
					bot.Emit(table)
				case client.EventUpdatePod:
					table := bots.SlackTable(
						"warning",
						bots.Field("Namespace", e.Namespace()),
						bots.Field("Message", fmt.Sprintf("Update pod: %s", e.GenerateName())),
						bots.Field("Node", e.NodeName()),
					)
					bot.Emit(table)
				case client.EventNewNode:
					table := bots.SlackTable(
						"good",
						bots.Field("Namespace", e.Namespace()),
						bots.Field("Message", fmt.Sprintf("Created cluster node: %s", e.GenerateName())),
						bots.Field("Node", e.NodeName()),
					)
					bot.Emit(table)
				case client.EventDelNode:
					table := bots.SlackTable(
						"danger",
						bots.Field("Namespace", e.Namespace()),
						bots.Field("Message", fmt.Sprintf("Deleted cluster node: %s", e.GenerateName())),
						bots.Field("Node", e.NodeName()),
					)
					bot.Emit(table)
				case client.EventUpdateNode:
					table := bots.SlackTable(
						"warning",
						bots.Field("Namespace", e.Namespace()),
						bots.Field("Message", fmt.Sprintf("Update cluster node: %s", e.GenerateName())),
						bots.Field("Node", e.NodeName()),
					)
					bot.Emit(table)
				}
				log.Printf("server event: %s", e)
			case s := <-bot.Receive():
				log.Printf("receive message: %s", s)
			}
		}
	}()

	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}

}
