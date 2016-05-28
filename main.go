package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/heppu/jun/client"
	"github.com/heppu/pwp-bot/api"
	"github.com/heppu/pwp-bot/wiki"
)

func main() {
	delay := 500 * time.Microsecond
	c := client.New(
		"irc.nlnog.net:6667",
		"pwp-bot",
		[]string{"#otit.code.pwp"},
		nil,
		&delay,
	)

	api := api.NewApiClient("http://127.0.0.1:8000")
	wiki.NewWikiBot(c, api, "#otit.code.pwp")
	err := c.Connect()
	if err != nil {
		log.Fatal(err)
	}

	go func(c *client.Client) {
		err := <-c.Error
		log.Println(err)
		c.Disconnect()
		os.Exit(0)
	}(c)

	// Graceful shutdown for Ctrl+C
	go func(c *client.Client) {
		kill := make(chan os.Signal, 1)
		signal.Notify(kill, os.Interrupt)
		<-kill

		c.Disconnect()
		os.Exit(0)
	}(c)

	<-c.Quit
}
