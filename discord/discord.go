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
	id   string
	sess *discordgo.Session
	quit chan int
}

func (bot *Bot) ready(s *discordgo.Session, r *discordgo.Ready) {
	Log("Bot: Got ready event!")
	bot.id = r.User.ID
	kasumi.AddGuilds(r.Guilds)
}

func (bot *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	kasumi.CreateMessage(m.Author, m.ChannelID, m.Content)
}

func New(token string) (*Bot, error) {
	/* initialize bot / discord session */
	var err error
	Log("Using token: ", token)
	bot := &Bot{
		quit: make(chan int)}
	bot.sess, err = discordgo.New(token)
	if err != nil {
		Log("discordgo.New(): ", err.Error())
		return &Bot{}, err
	}
	Log("discord.New(): Session created, are we authorized?")
	/* register callbacks here */
	bot.sess.AddHandler(bot.ready)
	bot.sess.AddHandler(bot.messageCreate)
	/* connect to discord */
	err = bot.sess.Open()
	if err != nil {
		Log("(*discordgo.Session).Open(): ", err.Error())
		return &Bot{}, err
	}
	Log("discord.New(): We are now allegedly connected.")

	return bot, nil
}
