package optimus

import (
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"code.google.com/p/go.crypto/ssh"
)

type SSH struct {
	Addr     string
	User     string
	Password string
	client   *ssh.Client
}

func (c *SSH) connect() error {
	if c.client != nil {
		return nil
	}

	key, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh/id_rsa"))
	if err != nil {
		return err
	}

	if c.User == "" {
		u, err := user.Current()
		if err != nil {
			return err
		}
		c.User = u.Username
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
			ssh.Password(c.Password),
		},
	}
	c.client, err = ssh.Dial("tcp", c.Addr, config)
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
