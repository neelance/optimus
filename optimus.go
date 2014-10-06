package optimus

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type Host struct {
	Name       string
	Connection Connection
	Properties map[string]interface{}
}

type Connection interface {
	Run(cmd string, out io.Writer) error
}

func (h *Host) Run(cmd string) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	bufR := bufio.NewReader(r)
	done := make(chan bool)
	go func() {
		for {
			line, err := bufR.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			fmt.Printf("[%s] %s", h.Name, line)
		}
		done <- true
	}()
	if err := h.Connection.Run(cmd, w); err != nil {
		return err
	}
	w.Close()

	<-done

	return nil
}

func (h *Host) DownloadFile(path string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := h.Connection.Run(fmt.Sprintf(`cat "%s"`, path), buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
