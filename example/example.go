package main

import (
	"github.com/neelance/optimus"
	"github.com/neelance/optimus/user"
	"github.com/spf13/cobra"
)

func main() {
	server1 := &optimus.Host{
		Name:     "server1",
		Addr:     "192.168.59.103:49155",
		User:     "root",
		Password: "root",
	}
	server2 := &optimus.Host{
		Name:     "server2",
		Addr:     "192.168.59.103:49156",
		User:     "root",
		Password: "root",
	}
	server3 := &optimus.Host{
		Name:     "server3",
		Addr:     "192.168.59.103:49157",
		User:     "root",
		Password: "root",
	}

	servers := optimus.NewGroup(server1, server2, server3)

	config := func(state *optimus.HostState) {
		user.Add(state, "richard")
		user.Add(state, "chris")
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(optimus.CommandRun(servers))
	rootCmd.AddCommand(optimus.CommandUp(servers, config))
	rootCmd.Execute()
}
