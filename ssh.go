package optimus

import (
	"io"
	"io/ioutil"
	"os/user"

	"code.google.com/p/go.crypto/ssh"
)

type SSH struct {
	Addr       string
	User       string
	Password   string
	PrivateKey string
	client     *ssh.Client
}

func (c *SSH) connect() error {
	if c.client != nil {
		return nil
	}

	var config ssh.ClientConfig

	switch c.User {
	case "":
		u, err := user.Current()
		if err != nil {
			return err
		}
		config.User = u.Username
	default:
		config.User = c.User
	}

	if c.Password != "" {
		config.Auth = append(config.Auth, ssh.Password(c.Password))
	}

	if c.PrivateKey != "" {
		key, err := ioutil.ReadFile(c.PrivateKey)
		if err != nil {
			return err
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return err
		}

		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	var err error
	c.client, err = ssh.Dial("tcp", c.Addr, &config)
	if err != nil {
		return err
	}

	return nil
}

func (c *SSH) Run(cmd string, out io.Writer) error {
	if err := c.connect(); err != nil {
		return err
	}

	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = out
	session.Stderr = out

	return session.Run(cmd)
}
