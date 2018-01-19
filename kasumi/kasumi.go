package kasumi

import (
	"flag"
	"io"
	"net"
)

func init() {
	flag.Parse()
}

/*
 * Global Command Line Options
 */

var Debug = flag.Bool("d", false, "print debug information to stderr")
var Realm = flag.String("realm", "192.168.0.103", "set the realm of server")
var Port = flag.Uint("p", 6667, "set the port number")

/*
 * use this interface in place of the type `net.Conn`
 * this is so we can pass fake connections through a
 * channel in order to represent members in a discord channel
 */
type Conn interface {
	io.ReadWriteCloser
	RemoteAddr() net.Addr
}
