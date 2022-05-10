package main

import (
	"net"
	"strconv"
)

type SocksAddress struct {
	AddressType byte
	IPAddress   net.IP
	DomainName  string
	Port        int
}

func TCPAddrToSocksAddress(addr net.TCPAddr) SocksAddress {
	x := SocksAddress{
		AddressType: socksAddressTypeIPv4,
		IPAddress:   addr.IP,
		Port:        addr.Port,
	}
	if len(addr.IP) > 4 {
		x.AddressType = socksAddressTypeIPv6
	}
	return x
}

func ReadSocksAddress(client *NetClient) SocksAddress {
	res := SocksAddress{}

	// read address
	res.AddressType = client.MustReadByte()
	switch res.AddressType {
	case socksAddressTypeIPv4:
		res.IPAddress = net.IP(client.MustReadBytes(4))
	case socksAddressTypeIPv6:
		res.IPAddress = net.IP(client.MustReadBytes(16))
	case socksAddressTypeDomain:
		domainLength := client.MustReadByte()
		res.DomainName = string(client.MustReadBytes(int(domainLength)))
	}

	// read port
	port := client.MustReadBytes(2)
	res.Port = int(int(port[0])<<8 | int(port[1]))

	return res
}

func (sa *SocksAddress) Bytes() []byte {
	switch sa.AddressType {
	case socksAddressTypeIPv4, socksAddressTypeIPv6:
		return append(
			append([]byte{sa.AddressType},
				sa.IPAddress...,
			),
			byte(sa.Port>>8),
			byte(sa.Port&255),
		)
	case socksAddressTypeDomain:
		return append(
			append([]byte{sa.AddressType},
				[]byte(sa.DomainName)...,
			),
			byte(sa.Port>>8),
			byte(sa.Port&255),
		)
	default:
		return []byte{
			socksAddressTypeIPv4,
			0x00, 0x00, 0x00, 0x00, // ip
			0x00, 0x00, // port
		}
	}
}

func (sa *SocksAddress) Address(dns *DNS) string {
	ip := sa.IPAddress
	if sa.AddressType == socksAddressTypeDomain {
		ip = dns.Resolve(sa.DomainName)
	}
	address := net.JoinHostPort(ip.String(), strconv.Itoa(sa.Port))
	return address
}
