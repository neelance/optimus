package optimus

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"code.google.com/p/go.crypto/ssh"
)

type Host struct {
	Name       string
	Addr       string
	User       string
	Password   string
	Properties map[string]interface{}
	client     *ssh.Client
}

func (h *Host) connect() error {
	if h.client != nil {
		return nil
	}

	key, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh/id_rsa"))
	if err != nil {
		return err
	}

	if h.User == "" {
		u, err := user.Current()
		if err != nil {
			return err
		}
		h.User = u.Username
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: h.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
			ssh.Password(h.Password),
		},
	}
	h.client, err = ssh.Dial("tcp", h.Addr, config)
	if err != nil {
		return err
	}

	return nil
}

func (h *Host) Run(cmd string) error {
	if err := h.connect(); err != nil {
		return err
	}

	session, err := h.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	h.forwardOutput(stdout, os.Stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	h.forwardOutput(stderr, os.Stderr)

	if err := session.Run(cmd); err != nil {
		return err
	}

	return nil
}

func (h *Host) DownloadFile(path string) ([]byte, error) {
	if err := h.connect(); err != nil {
		return nil, err
	}

	session, err := h.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	buf := bytes.NewBuffer(nil)
	session.Stdout = buf
	if err := session.Run(fmt.Sprintf(`cat "%s"`, path)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *Host) forwardOutput(r io.Reader, w io.Writer) {
	out := bufio.NewReader(r)
	go func() {
		for {
			line, err := out.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			fmt.Printf("[%s] %s", h.Name, line)
		}
	}()
}

type Group struct {
	Children   []HostOrGroup
	Properties map[string]interface{}
}

func NewGroup(hosts ...HostOrGroup) *Group {
	return &Group{Children: hosts}
}

type HostOrGroup interface {
	Hosts() []*Host
}

func (h *Host) Hosts() []*Host {
	return []*Host{h}
}

func (g *Group) Hosts() []*Host {
	var hosts []*Host
	for _, c := range g.Children {
		hosts = append(hosts, c.Hosts()...)
	}
	return hosts
}

type Module interface{}

type HostState struct {
	Host    *Host
	Modules map[Module]HostStateModule
	actions []Action
}

func (s *HostState) AddAction(action Action) {
	s.actions = append(s.actions, action)
}

type HostStateModule interface {
}

type Action interface {
	Description() string
	Run()
}

type simpleAction struct {
	desc string
	fun  func()
}

func (a *simpleAction) Description() string {
	return a.desc
}

func (a *simpleAction) Run() {
	a.fun()
}

func SimpleAction(desc string, fun func()) Action {
	return &simpleAction{desc, fun}
}
