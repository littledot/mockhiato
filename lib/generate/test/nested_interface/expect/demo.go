package example

import "net"

// Target contains a nested interface. Mockhiato must mock inerhited methods and record dependencies used by those methods.
type Target interface {
	net.Conn
}
