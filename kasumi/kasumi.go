package kasumi

import (
	"fmt"
	"flag"
	"io"
	"log"
	"errors"
	"net"
	"github.com/bwmarrin/discordgo"
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
	net	string
	addr	string
}

type FauxConn struct {
	addr	*FauxAddr
	IN	chan string
}

type Client struct {
	con	*FauxConn
	user	*discordgo.User
}

var clients map[string]*Client = make(map[string]*Client)
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
	str	:= <-con.IN
	max	:= len(p)
	for n := 0; n < len(str); n++ {
		if n > max { return n, errors.New("Read(): buffer is full")}
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
	c.sendMessage(fmt.Sprintf("PRIVMSG #test %s", msg))
}

func (c *Client) sendMessage(msg string) {
	Log("Client.sendMessage(): sending:", msg)
	c.con.IN <- msg+"\r\n"
}

/*
 * Exported Package Functions
 */

func AddGuilds(guilds []*discordgo.Guild) {
	for _, g := range guilds {
		AddGuild(*g)
	}
}

func AddGuild(guild discordgo.Guild) {
	/* create a faux connection for each member */
	Log("AddGuild: Registering users")
	for _, m := range guild.Members {
		clientID := m.User.Username+m.User.Discriminator
		if client, ok := clients[clientID]; ok {
			client.joinMsg("#test")
			continue
		}
		clients[clientID] = newClient(m.User)
		clients[clientID].joinMsg("#test")
	}
}

func CreateMessage(user *discordgo.User, id string, msg string) {
	clientID := user.Username+user.Discriminator
	client, ok := clients[clientID]
	if !ok {
		Log("CreateMessage(): Failed to retrieve client:", clientID)
		return
	}
	client.say(id, msg)
}

