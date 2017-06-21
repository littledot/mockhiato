package lib

import "bytes"

// Bufferw is a wrapper around bytes.Buffer.
type Bufferw struct {
	buf bytes.Buffer
}

// WriteString wraps bytes.Buffer.WriteString()
func (r *Bufferw) WriteString(s string) int {
	n, err := r.buf.WriteString(s)
	if err != nil {
		panic(err)
	}
	return n
}

// Bytes wraps bytes.Buffer.Bytse()
func (r *Bufferw) Bytes() []byte {
	return r.buf.Bytes()
}
