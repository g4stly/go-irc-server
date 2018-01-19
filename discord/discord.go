package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/kasumi/kasumi"
	"log"
)

func Log(v ...string) {
	if *kasumi.Debug == true {
		log.Printf("DISCORD: %s", fmt.Sprint(v))
	}
}

type Bot struct {
	C    chan kasumi.Conn
	id   string
	sess *discordgo.Session
	quit chan int
}

func New(token string) (*Bot, error) {
	/* initialize bot / discord session */
	var err error
	bot := &Bot{
		C:    make(chan kasumi.Conn),
		quit: make(chan int)}
	bot.sess, err = discordgo.New(token)
	if err != nil {
		Log("discordgo.New(): ", err.Error())
		return &Bot{}, err
	}
	Log("discord.New(): Session created, are we authorized?")
	/* connect to discord */
	err = bot.sess.Open()
	if err != nil {
		Log("(*discordgo.Session).Open(): ", err.Error())
		return &Bot{}, err
	}
	Log("discord.New(): We are now allegedly connected.")

	/* register callbacks here */

	return bot, nil
}
