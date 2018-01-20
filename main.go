package main

import (
	"os"
	"log"
	"strings"
	"io/ioutil"
	"github.com/g4stly/kasumi/irc"
	"github.com/g4stly/kasumi/discord"
	"github.com/g4stly/kasumi/kasumi"
)

/*
 * DISCORD BOT 	: MODEL
 * KASUMI	: CONTROLLER
 * IRC SERVER	: VIEW
 */


func main() {
	/* read discord token */
	t, err := ioutil.ReadFile("discord_token")
	if err != nil {
		log.Printf("Error reading discord token: %s", err)
		os.Exit(-1)
	}
	/* initialize discord bot */
	_, err = discord.New(strings.TrimSpace(string(t)))
	if err != nil {
		log.Printf("Failed to initialize discord bot: ", err.Error())
		os.Exit(-1)
	}
	/* start up irc server */
	irc.Server(kasumi.ConnSpawner)

	os.Exit(0)
}
