package main

import (
	"fmt"
	"net"
)

type SocksServer struct {
	dns  *DNS
	auth *Auth
}

func NewSocksServer(dns *DNS, auth *Auth) *SocksServer {
	return &SocksServer{
		dns:  dns,
		auth: auth,
	}
}

func (s *SocksServer) Listen(address string) error {
	netListener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	for {
		connection, _ := netListener.Accept()
		go s.holdConnection(NewNetClient(connection))
	}
}

func (s *SocksServer) holdConnection(client *NetClient) error {
	defer client.Close()

	// checking socks protocol version
	protoversion := client.MustReadByte()
	if protoversion != socksProtocolVersion {
		return fmt.Errorf("wrong proto version: %b", protoversion)
	}

	// authenticate the client
	err := s.auth.AuthenticateClient(client)
	if err != nil {
		return err
	}

	// read request
	request, err := ReadSocksRequest(client)
	if err != nil {
		return err
	}

	// execute request
	request.Execute(client, s.dns)

	return nil
}
