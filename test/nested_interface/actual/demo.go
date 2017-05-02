package example

import "net"

// A contains a nested interface. Mockhiato must mock inerhited methods and record dependencies used by those methods.
type A interface {
	net.Conn
}
