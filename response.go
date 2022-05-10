package main

import "io"

type SocksResponse struct {
	Status  byte
	Address SocksAddress
}

func NewResponse(status byte, address SocksAddress) *SocksResponse {
	return &SocksResponse{
		Status:  status,
		Address: address,
	}
}

func (sr *SocksResponse) Send(writer io.Writer) error {
	packet := []byte{
		socksProtocolVersion, // protocol version
		sr.Status,            // resp status
		0x00,                 // reserved
	}
	packet = append(packet, sr.Address.Bytes()...)
	_, err := writer.Write(packet)
	return err
}
