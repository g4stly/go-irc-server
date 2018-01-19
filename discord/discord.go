package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/kasumi/kasumi"
)

func Log(v ...string) {
	if *debug == true {
		log.Printf("DISCORD: %s", fmt.Sprint(v))
	}
}

type Bot struct {
	c	chan kasumi.Conn
	id	string
	sess	*discordgo.Session
	quit chan int
}

func Initialize(token string) Bot, error {
	bot := &Lain{
		c:	make(chan kasumi.Conn),
		quit:	make(chan int)}
	bot.sess, err = discordgo.New(token)
	if err != nil { return Bot{}, err }
	Log("Discord session created, are we authorized?")

	/* register callbacks here */

	return bot
}
