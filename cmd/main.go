package main

import (
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
	if err != nil {
		log.Fatal(err)
	}

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
				table := bots.SlackTable(
					e.Level(),
					bots.Field("Name", e.Name()),
					bots.Field("Namespace", e.Namespace()),
					bots.Field("Message", e.Message()),
					bots.Field("Reason", e.Reason()),
					bots.Field("Kind", e.Kind()),
					bots.Field("Component", e.Component()),
				)
				bot.Emit(table)
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
