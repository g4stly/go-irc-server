package kasumi

import (
	"io"
	"net"
)

/*
 * use this interface in place of the type `net.Conn`
 * this is so we can pass fake connections through a
 * channel in order to represent members in a discord channel
 */
type Conn interface {
	io.ReadWriteCloser
	RemoteAddr() net.Addr
}
