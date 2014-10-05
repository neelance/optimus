package user

import (
	"bytes"
	"fmt"

	"github.com/neelance/optimus"
)

type User struct {
	Name     string
	Password string
}

var module = &struct{}{}

type stateModule struct {
	users []string
}

func (s *stateModule) userIndex(name string) int {
	for i, u := range s.users {
		if u == name {
			return i
		}
	}
	return -1
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

func Present(state *optimus.HostState, user *User, shouldBePresent bool) {
	s := ownState(state)
	index := s.userIndex(user.Name)
	if (index != -1) == shouldBePresent {
		return
	}

	if !shouldBePresent {
		s.users[index] = s.users[len(s.users)-1]
		s.users = s.users[:len(s.users)-1]
		state.AddAction(optimus.SimpleAction(fmt.Sprintf(`Remove user "%s"`, user.Name), func() {
			state.Host.Run(fmt.Sprintf(`userdel "%s"`, user.Name))
		}))
		return
	}

	s.users = append(s.users, user.Name)
	state.AddAction(optimus.SimpleAction(fmt.Sprintf(`Add user "%s"`, user.Name), func() {
		state.Host.Run(fmt.Sprintf(`useradd --create-home --password "%s" "%s"`, user.Password, user.Name))
	}))
}
