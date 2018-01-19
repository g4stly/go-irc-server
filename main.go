package main

import (
	"flag"
	"log"
	"io/ioutil"
	"github.com/g4stly/kasumi/irc"
	"github.com/g4stly/kasumi/kasumi"
	"os"
)

/*
 * DISCORD BOT 	: MODEL
 * KASUMI	: CONTROLLER
 * IRC SERVER	: VIEW
 */

 /* TODO: flesh out discord.go, then flesh out kasumi.go */

func main() {
	flag.Parse()
	/* read discord token */
	t, err := ioutil.ReadFile("discord_token")
	if err != nil {
		log.Printf("Error reading discord token: %s", err)
		os.Exit(-1)
	}
	log.Printf("Using token %s\n", strings.TrimSpace(string(t)))
	/* initialize discord bot */
	bot, err := discord.Initialize(strings.TrimSpace(string(t)))
	if err != nil {
		log.Printf("Failed to initialize discord bot: ", err.Error())
		os.Exit(return -1)
	}
	go bot.Connect()
	/* start up irc server */
	c := make(chan kasumi.Conn)
	irc.Server(c)

	os.Exit(0)
}
