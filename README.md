Optimus
=======

Build your own configuration management tool

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
