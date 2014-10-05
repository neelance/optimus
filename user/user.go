package user

import (
	"bytes"
	"fmt"

	"github.com/neelance/optimus"
)

var module = &struct{}{}

type stateModule struct {
	users []string
}

func (s *stateModule) userExists(name string) bool {
	for _, u := range s.users {
		if u == name {
			return true
		}
	}
	return false
}

func ownState(state *optimus.HostState) *stateModule {
	s, ok := state.Modules[module].(*stateModule)
	if !ok {
		s = &stateModule{}
		state.Modules[module] = s

		passwd, err := state.Host.DownloadFile("/etc/passwd")
		if err != nil {
			panic(err)
		}

		for _, line := range bytes.Split(passwd, []byte{'\n'}) {
			parts := bytes.Split(line, []byte{':'})
			if len(parts) < 2 {
				continue
			}
			s.users = append(s.users, string(parts[0]))
		}
	}
	return s
}

func Add(state *optimus.HostState, name string) {
	s := ownState(state)
	if s.userExists(name) {
		return
	}

	s.users = append(s.users, name)
	state.AddAction(optimus.SimpleAction(fmt.Sprintf(`Add user "%s"`, name), func() {
		state.Host.Run(fmt.Sprintf(`adduser --disabled-password --gecos "" %s`, name))
	}))
}

func Remove(state *optimus.HostState, name string) {
	s := ownState(state)
	if !s.userExists(name) {
		return
	}

	s.users = append(s.users, name)
	state.AddAction(optimus.SimpleAction(fmt.Sprintf(`Remove user "%s"`, name), func() {
		state.Host.Run(fmt.Sprintf(`deluser %s`, name))
	}))
}
