package logger

import (
	"io"
	"net"
)

func NewUDPWriter(addr string) (w io.Writer, err error) {
	var remoteAddr *net.UDPAddr
	if remoteAddr, err = net.ResolveUDPAddr("udp", addr); err != nil {
		return nil, err
	}

	var conn net.Conn
	if conn, err = net.DialUDP("udp", nil, remoteAddr); err != nil {
		return nil, err
	}

	return conn, nil
}
