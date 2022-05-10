package main

import (
	"errors"
	"io"
	"net"
)

type SocksRequest struct {
	Command            byte
	DestinationAddress SocksAddress
}

func ReadSocksRequest(client *NetClient) (*SocksRequest, error) {
	res := &SocksRequest{}

	// check protocol version
	if client.MustReadByte() != socksProtocolVersion {
		return res, errors.New("wrong socks protocol version")
	}

	//read command
	res.Command = client.MustReadByte()

	// read reserved field (should be 0x00)
	client.MustReadByte()

	// read destination address
	res.DestinationAddress = ReadSocksAddress(client)

	return res, nil
}

func (req *SocksRequest) Execute(client *NetClient, dns *DNS) {
	switch req.Command {
	case socksRequestCommandConnect:
		req.executeConnect(client, dns)
	// case socksRequestCommandAssociate:
	// case socksRequestCommandBind:
	default:
		NewResponse(socksReplyCommandNotSupported, req.DestinationAddress).Send(client)
	}
}

func (req *SocksRequest) executeConnect(client *NetClient, dns *DNS) {
	// try to connect to the target host
	connection, err := net.Dial("tcp", req.DestinationAddress.Address(dns))
	if err != nil {
		NewResponse(socksReplyCommandHostUnreachable, req.DestinationAddress)
		return
	}
	defer connection.Close()

	// get local tcp address
	localAddress := TCPAddrToSocksAddress(*connection.LocalAddr().(*net.TCPAddr))

	// send OK status
	NewResponse(socksReplyCommandSucceeded, localAddress).Send(client)

	// tunnel traffic
	canclose := make(chan int, 2)
	go tunnel(client, connection, canclose)
	go tunnel(connection, client, canclose)
	closed := 0
	for closed < 2 {
		<-canclose
		closed++
	}
}

func tunnel(from io.Reader, to io.Writer, canclose chan int) {
	io.Copy(to, from)
	canclose <- 0
}
