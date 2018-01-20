package kasumi

import (
	"flag"
	"io"
	"error"
	"net"
	_ "github.com/bwmarrin/discordgo"
)

func init() { flag.Parse() }

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
	con	Conn
	user	*discordgo.User
}

var clients := make(map[string]Client)
var ConnSpawner := make(chan Conn)

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

func (a *FauxAddr) Network() {
	return a.net
}

func (a *FauxAddr) String() {
	return a.addr

}

/*
 * FauxConn Methods
 */

func newFauxConn() *FauxConn {
	return &{addr: newFauxAddr(), IN: make(chan string, 3)}
}

func (con *FauxConn) Read(p []byte) (n int, err error) {
	str	:= <-con.IN
	max	:= len(p)
	for n := 0; n < len(str); n++ {
		if n > max { return n, Error.New("Read(): buffer is full")}
		p[i] = str[i]
	}
	return n, nil
}

func (con *FauxConn) Write(p []byte) (n int, err error) {
	// used to satisfy interface
	return len(p), nil
}

func (con *FauxConn) Close(p []byte) error {
	// used to satisfy interface
	return nil
}

func (con *FauxConn) RemoteAddr() net.Addr {
	return *con.addr
}

/*
 * Client Methods...
 */

func newClient(usr *discordgo.User) *Client {
	localClient := &Client{con: newFauxConn(), user: usr}
	ConnSpawner <- localClient.con
	// here is where we do the nick/user commands to the irc server
	return localClient
}

/*
 * Private Package Functions
 */

func joinIRCChannel(c Client, channel string) {
	// here is where we craft a message to the irc server
	// then we send it into c.con.IN and hope for the best
}

/*
 * Exported Package Functions
 */

func AddGuild(guild discordgo.Guild) {
	/* create a faux connection for each member */
	for _, m := range g.Members {
		clients[m.User.Username+Discriminator] = newClient(m.User)
	}
	/* have each member join appropriate rooms */
	for _, c := range g.Channels {
		for _, r := range c.Recipients {
			client := clients[r.Username+r.Discriminator]
			/* just for now, have them join #test */
			joinIRCChannel(client, "test")
		}
	}
}

func CreateMessage(user discordgo.User, id string, msg string) {
}

