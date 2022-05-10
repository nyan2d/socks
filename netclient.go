package main

import (
	"errors"
	"io"
	"net"
)

type NetClient struct {
	connection net.Conn
}

func NewNetClient(c net.Conn) *NetClient {
	return &NetClient{
		connection: c,
	}
}

// basic level

func (c *NetClient) Close() error {
	return c.connection.Close()
}

func (c *NetClient) Read(b []byte) (n int, err error) {
	return c.connection.Read(b)
}

func (c *NetClient) Write(b []byte) (n int, err error) {
	return c.connection.Write(b)
}

// advanced level

func (c *NetClient) MustReadByte() byte {
	buf := make([]byte, 1)
	c.Read(buf)
	return buf[0]
}

func (c *NetClient) ReadBytes(len int) ([]byte, error) {
	buf := make([]byte, len)
	if _, err := io.ReadFull(c, buf); err != nil {
		return buf, errors.New("can't read bytes")
	}
	return buf, nil
}

func (c *NetClient) MustReadBytes(len int) []byte {
	b, _ := c.ReadBytes(len)
	return b
}

func (c *NetClient) WriteBytes(b ...byte) error {
	_, err := c.Write(b)
	return err
}
