Optimus
=======

A framework for building your own configuration management tool

Motivation
----------
There are [Puppet](http://puppetlabs.com/), [Chef](https://www.getchef.com/), [SaltStack](http://www.saltstack.com/) and [Ansible](http://www.ansible.com/). These are good products, so why create another configuration management tool? Sure, all of them have minor downsides, but that's not enough for kind of reinventing the wheel. But there is one property that they all share: They are tools that you feed with configuration files and scripts. They are big machines with thousands of knobs and switches. Imagine a factory intended for building every model of cars that is available nowadays. How huge and complex would it be? Could you really build ALL models? How many robotic arms would be hanging useless, because they are not required for the current model? Wouldn't it be better to just pick the robotic arms that you need and combine them into a relatively small and effective factory? That's the idea behind Optimus. It is not a tool, it is a framework for building your own configuration management tool.

Status
------
Proof of concept.

Installation
------------
```
go get github.com/neelance/optimus
```

Running the example
-------------------
Install Docker and set `$DOCKER_HOST` accordingly, see https://docs.docker.com/installation/#installation.

Create three docker containers that run `sshd`, available at ports 50001, 50002 and 50003 of the docker host with `root` as username and password:
```
$GOPATH/src/github.com/neelance/optimus/example/create-docker-containers.sh
```
Those encapsulated containers will be our test hosts. Run this script again at any time to reset them to their initial state.

Now you can build and run the example tool:
```
go run $GOPATH/src/github.com/neelance/optimus/example/example.go
```

This will show the available commands. The `run` command will execute a shell command on each host in parallel, for example:
```
go run $GOPATH/src/github.com/neelance/optimus/example/example.go run hostname
```

More powerful is the `up` command, which will bring the hosts to the configuration state described inside of `example.go`:
```
go run $GOPATH/src/github.com/neelance/optimus/example/example.go up
```
On first run, it will create some users. Subsequent runs will do nothing, since the users already exist. You can modify `example.go` to describe other configurations and then apply them with the `up` command.

The `up` command in detail
--------------------------
The `up` command consists of two phases, the *analyze* and the *modify* phase.

In the *analyze* phase, the `optimus.Configurator` function gets executed multiple times in parallel for each host. It can fetch information about the current state of the host and apply actions (changes) to the `optimus.HostState`. However, these actions are not yet applied to the actual server, but queued for the *modify* phase. During the *analyze* phase, all communication with the host is required to be **read-only**. This way, any information about the host's state need to be fetched only once per *analyze* phase. The phase ends by listing all queued actions and asking for confirmation to proceed.

The *modify* phase runs the queued actions on the actual hosts in parallel. If an error occurs while executing an action, the affected host will be marked as *failed* and it will be ignored for the rest of the `up` command. The *modify* phase continues for the remaining hosts.

Another *analyze* phase begins after all actions have been executed. It examines the hosts (except the *failed* ones) and ideally it should detect that no further actions are required. If this is the case, the `up` command terminates. In some scenarios, however, multiple cycles are required to reach a stable state.
