package kasumi

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net"
	"strings"
)

func init() { flag.Parse() }

func Log(v ...string) {
	if *Debug == true {
		log.Printf("*KASUMI: %s", fmt.Sprint(v))
	}
}

/*
 * use this interface in place of the type `net.Conn`
 * this is so we can pass fake connections through a
 * channel in order to represent members in a discord channel
 */
type Conn interface {
	io.ReadWriteCloser
	RemoteAddr() net.Addr
}

type FauxAddr struct {
	net  string
	addr string
}

type FauxConn struct {
	addr *FauxAddr
	IN   chan string
}

type Client struct {
	con  *FauxConn
	user *discordgo.User
}

type Channel struct {
	name string
	id   string
}

var clients map[string]*Client = make(map[string]*Client)
var channels map[string]*Channel = make(map[string]*Channel)
var ConnSpawner chan Conn = make(chan Conn)

/*
 * Global Command Line Options
 */

var Debug = flag.Bool("d", false, "print debug information to stderr")
var Realm = flag.String("realm", "192.168.0.103", "set the realm of server")
var Port = flag.Uint("p", 6667, "set the port number")

/*
 * FauxAddr Methods
 */

func newFauxAddr() *FauxAddr {
	return &FauxAddr{net: "tcp", addr: "69.4.20.69"}
}

func (a *FauxAddr) Network() string {
	return a.net
}

func (a *FauxAddr) String() string {
	return a.addr

}

/*
 * FauxConn Methods
 */

func newFauxConn() *FauxConn {
	return &FauxConn{addr: newFauxAddr(), IN: make(chan string, 3)}
}

func (con *FauxConn) Read(p []byte) (n int, err error) {
	str := <-con.IN
	max := len(p)
	for n := 0; n < len(str); n++ {
		if n > max {
			return n, errors.New("Read(): buffer is full")
		}
		p[n] = str[n]
	}
	return n, nil
}

func (con *FauxConn) Write(p []byte) (n int, err error) {
	// used to satisfy interface
	return len(p), nil
}

func (con *FauxConn) Close() error {
	// used to satisfy interface
	return nil
}

func (con *FauxConn) RemoteAddr() net.Addr {
	return con.addr
}

/*
 * Client Methods...
 */

func newClient(usr *discordgo.User) *Client {
	localClient := &Client{con: newFauxConn(), user: usr}
	ConnSpawner <- localClient.con
	localClient.nickMsg(usr.Username)
	localClient.userMsg(usr.Username)
	return localClient
}

func (c *Client) joinIRCChannel(guild discordgo.Guild, channel *discordgo.Channel) {
	dest, ok := channels[channel.ID]
	if !ok {
		channels[channel.ID] = newChannel(channel.Name, channel.ID)
		dest = channels[channel.ID]
	}
	guildName := strings.Split(guild.Name, " ")
	firstWordOfName := guildName[0]
	if firstWordOfName == "The" {
		firstWordOfName = guildName[1]
	}
	c.joinMsg(fmt.Sprintf("#" + firstWordOfName + dest.name))
}

func (c *Client) nickMsg(nick string) {
	c.sendMessage(fmt.Sprintf("NICK %s", nick))
}

func (c *Client) userMsg(real_name string) {
	c.sendMessage(fmt.Sprintf("USER fake 0 * :%s", real_name))
}

func (c *Client) joinMsg(channel string) {
	c.sendMessage(fmt.Sprintf("JOIN %s", channel))
}

func (c *Client) say(id string, msg string) {
	c.sendMessage(fmt.Sprintf("PRIVMSG %s %s", id, msg))
}

func (c *Client) sendMessage(msg string) {
	Log("Client.sendMessage(): sending:", msg)
	c.con.IN <- msg + "\r\n"
}

/*
 * Channel Methods
 */

func newChannel(n string, i string) *Channel {
	localChannel := &Channel{name: n, id: i}
	return localChannel
}

/*
 * Exported Package Functions
 */

func AddGuilds(guilds []*discordgo.Guild) {
	for _, g := range guilds {
		Log("AddGuilds: Adding guild ", g.Name)
		AddGuild(*g)
	}
}

func AddGuild(guild discordgo.Guild) {
	for _, m := range guild.Members {
		client, ok := clients[m.User.ID]
		if !ok {
			clients[m.User.ID] = newClient(m.User)
			client = clients[m.User.ID]
		}
		for _, c := range guild.Channels {
			client.joinIRCChannel(guild, c)
		}
	}
}

func CreateMessage(user *discordgo.User, id string, msg string) {
	clientID := user.Username + user.Discriminator
	client, ok := clients[clientID]
	if !ok {
		Log("CreateMessage(): Failed to retrieve client:", clientID)
		return
	}
	client.say(id, msg)
}
