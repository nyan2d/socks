package main

const (
	socksProtocolVersion byte = 0x05
)

const (
	authModeNoauth           byte = 0x00
	authModeGSSAPI           byte = 0x01
	authModeUsernamePassword byte = 0x02
	authModeNoAuthMethods    byte = 0xFF

	authStatusSuccess byte = 0x00
	authStatusFailure byte = 0xFF

	authProtocolVersion byte = 0x01
)

const (
	socksRequestCommandConnect   byte = 0x01
	socksRequestCommandBind      byte = 0x02
	socksRequestCommandAssociate byte = 0x03
)

const (
	socksAddressTypeIPv4   byte = 0x01
	socksAddressTypeDomain byte = 0x03
	socksAddressTypeIPv6   byte = 0x04
)

const (
	socksReplyCommandSucceeded       byte = 0x00
	socksReplyCommandHostUnreachable byte = 0x04
	socksReplyCommandNotSupported    byte = 0x07
)
