package optimus

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func Run(target HostOrGroup, cmd string) {
	EachHostParallel(target, func(h *Host) {
		err := h.Run(cmd)
		if err != nil {
			fmt.Printf("[%s] %s\n", h.Name, err)
		}
	})
}

func RunCommand(target HostOrGroup) *cobra.Command {
	return &cobra.Command{
		Use:   "run [shell command]",
		Short: "Run a shell command on hosts",
		Long:  `Run a shell command on hosts`,
		Run: func(cmd *cobra.Command, args []string) {
			Run(target, strings.Join(args, " "))
		},
	}
}

type Configurator func(state *HostState)

func Up(target HostOrGroup, config Configurator) {
	failed := make(map[*Host]bool)
	for {
		fmt.Println(yellow("Analyzing configurations..."))
		states := make(map[*Host]*HostState)
		done := true
		EachHostParallel(target, func(h *Host) {
			if failed[h] {
				return
			}
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf(red("[%s] Error: %s\n"), h.Name, err)
					failed[h] = true
				}
			}()

			s := &HostState{Host: h, Modules: make(map[Module]HostStateModule)}
			config(s)

			states[h] = s // add after sucessful config
			if len(s.actions) != 0 {
				done = false
			}
		})

		if done {
			switch len(failed) {
			case 0:
				fmt.Println(green("No pending actions."))
			case 1:
				fmt.Printf(red("No pending actions, but 1 host was ignored because of errors.\n"))
			default:
				fmt.Printf(red("No pending actions, but %d hosts were ignored because of errors.\n"), len(failed))
			}
			break
		}

		fmt.Println("Actions to be done:")
		for _, s := range states {
			for _, a := range s.actions {
				fmt.Printf("[%s] %s\n", s.Host.Name, a.Description())
			}
		}
		fmt.Print("Proceed? (Y/n) ")
		var answer string
		fmt.Scanln(&answer)
		if answer != "" && answer != "y" && answer != "Y" {
			break
		}

		EachHostParallel(target, func(h *Host) {
			if failed[h] {
				return
			}
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf(red("[%s] Error: %s\n"), h.Name, err)
					failed[h] = true
				}
			}()
			for _, a := range states[h].actions {
				a.Run()
			}
		})

		switch len(failed) {
		case 0:
			fmt.Println(green("All changes sucessfully applied."))
		case 1:
			fmt.Printf(red("Errors occurred for 1 host. Ignoring this host.\n"))
		default:
			fmt.Printf(red("Errors occurred for %d hosts. Ignoring these hosts.\n"), len(failed))
		}
	}
}

func UpCommand(target HostOrGroup, config Configurator) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Analyze the hosts and apply configuration",
		Long:  `Analyze the hosts and apply configuration`,
		Run: func(cmd *cobra.Command, args []string) {
			Up(target, config)
		},
	}
}

func EachHostParallel(target HostOrGroup, f func(h *Host)) {
	hosts := target.Hosts()
	wait := make(chan bool, len(hosts))
	for _, h := range hosts {
		go func(h *Host) {
			f(h)
			wait <- true
		}(h)
	}
	for _ = range hosts {
		<-wait
	}
}

func green(s string) string {
	return "\x1B[32m" + s + "\x1B[0m"
}

func yellow(s string) string {
	return "\x1B[33m" + s + "\x1B[0m"
}

func red(s string) string {
	return "\x1B[31m" + s + "\x1B[0m"
}
