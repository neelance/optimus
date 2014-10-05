package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/neelance/optimus"
	"github.com/neelance/optimus/user"
	"github.com/spf13/cobra"
)

// make example work on any machine
func dockerHost() string {
	var u, _ = url.Parse(os.Getenv("DOCKER_HOST"))
	if u.Host == "" {
		fmt.Println("Please install docker daemon or boot2docker and set DOCKER_HOST")
		os.Exit(0)
	}
	return strings.Split(u.Host, ":")[0]
}

// host inventory
var webserver1 = &optimus.Host{
	Name:     "webserver1",
	Addr:     dockerHost() + ":50001",
	User:     "root",
	Password: "root",
}
var webserver2 = &optimus.Host{
	Name:     "webserver2",
	Addr:     dockerHost() + ":50002",
	User:     "root",
	Password: "root",
}
var webserver3 = &optimus.Host{
	Name:     "webserver3",
	Addr:     dockerHost() + ":50003",
	User:     "root",
	Password: "root",
}

// groups of hosts
var webservers = optimus.NewGroup(webserver1, webserver2, webserver3)

// user inventory
var alice = &user.User{
	Name:     "alice",
	Password: "$1$QPjAgnGi$CJUA.BQihAq1DUKB4Or8R0",
}
var bob = &user.User{
	Name:     "bob",
	Password: "$1$ZeCDBmNs$aWrXQeF0PPQkwxZjD27FQ0",
}
var oscar = &user.User{
	Name:     "oscar",
	Password: "$1$Hega65rU$7.H8UTos70caDftgOidap1",
}

func main() {
	// this function is executed with each host's current state to determine the desired actions
	config := func(state *optimus.HostState) {
		// add alice to all hosts
		user.Present(state, alice, true)

		// bob has only access to webserver3
		user.Present(state, bob, state.Host == webserver3)

		// remove oscar from all hosts
		user.Present(state, oscar, false)
	}

	// create command line interface
	var rootCmd = &cobra.Command{Use: "example"}
	rootCmd.AddCommand(optimus.CommandRun(webservers))
	rootCmd.AddCommand(optimus.CommandUp(webservers, config))
	rootCmd.Execute()
}
