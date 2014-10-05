package optimus

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func CommandRun(target HostOrGroup) *cobra.Command {
	return &cobra.Command{
		Use:   "run [shell command]",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			EachHostParallel(target, func(h *Host) {
				err := h.Run(strings.Join(args, " "))
				if err != nil {
					fmt.Printf("[%s] %s\n", h.Name, err)
				}
			})
		},
	}
}

type Configurator func(state *HostState)

func CommandUp(target HostOrGroup, config Configurator) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			for {
				fmt.Println(yellow("Analyzing configurations..."))
				states := make(map[*Host]*HostState)
				done := true
				EachHostParallel(target, func(h *Host) {
					s := &HostState{Host: h, Modules: make(map[Module]HostStateModule)}
					config(s)
					states[h] = s
					if len(s.actions) != 0 {
						done = false
					}
				})

				if done {
					fmt.Println(green("No actions required."))
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
					for _, a := range states[h].actions {
						a.Run()
					}
				})
				fmt.Println(green("All changes sucessfully applied."))
			}
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
