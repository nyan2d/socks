package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

type UserAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct {
	users           []UserAccount
	AllowNoAuthMode bool
}

func NewAuth(allowguest bool) *Auth {
	return &Auth{
		users:           make([]UserAccount, 0),
		AllowNoAuthMode: allowguest,
	}
}

func (a *Auth) ReadFromJsonFile(jsonFile string) error {
	var accs []UserAccount

	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &accs)
	if err != nil {
		return err
	}

	a.users = accs
	return nil
}

func (a *Auth) Contains(username string, password string) bool {
	if len(a.users) < 1 {
		return false
	}
	for _, user := range a.users {
		if strings.EqualFold(username, user.Username) {
			if password == user.Password {
				return true
			}
		}
	}
	return false
}

func (a *Auth) AuthenticateClient(client *NetClient) error {
	// read auth methods
	authMethods, err := readAuthMethods(client)
	if err != nil {
		return err
	}

	// pick auth method
	if isAuthMethodAvailable(authMethods, authModeUsernamePassword) {
		client.WriteBytes(socksProtocolVersion, authModeUsernamePassword)
		return a.authenticateWithUsernamePassword(client)
	} else if a.AllowNoAuthMode && isAuthMethodAvailable(authMethods, authModeNoauth) {
		client.WriteBytes(socksProtocolVersion, authModeNoauth)
		return nil
	} else {
		client.WriteBytes(socksProtocolVersion, authModeNoAuthMethods)
		return errors.New("no auth methods avalible")
	}
}

func (a *Auth) authenticateWithUsernamePassword(client *NetClient) error {
	// read auth proto version
	clientauthprotoversion := client.MustReadByte()
	if clientauthprotoversion != authProtocolVersion {
		return errors.New("wrong auth protocol version")
	}

	// read username
	usernameLength := int(client.MustReadByte())
	username, err := client.ReadBytes(usernameLength)
	if err != nil {
		return err
	}

	// read password
	passwordLength := int(client.MustReadByte())
	password, err := client.ReadBytes(passwordLength)
	if err != nil {
		return err
	}

	// trying to authenticate the client
	isLoginPasswordValid := a.Contains(string(username), string(password))
	if isLoginPasswordValid {
		return client.WriteBytes(authProtocolVersion, authStatusSuccess)
	}
	return client.WriteBytes(authProtocolVersion, authStatusFailure)
}

func readAuthMethods(client *NetClient) ([]byte, error) {
	methodsCount := client.MustReadByte()
	if methodsCount == 0 {
		return []byte{}, errors.New("no auth methods avalible")
	}

	methods, err := client.ReadBytes(int(methodsCount))
	if err != nil {
		return methods, errors.New("can't read auth methods")
	}
	return methods, nil
}

func isAuthMethodAvailable(methods []byte, method byte) bool {
	if len(methods) == 0 {
		return false
	}
	for _, v := range methods {
		if v == method {
			return true
		}
	}
	return false
}
