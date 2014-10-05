#!/bin/bash

docker rm -f exampleServer1 exampleServer2 exampleServer3 2>&1 > /dev/null

docker run --name exampleServer1 --detach --publish 50001:22 rastasheep/ubuntu-sshd
docker run --name exampleServer2 --detach --publish 50002:22 rastasheep/ubuntu-sshd
docker run --name exampleServer3 --detach --publish 50003:22 rastasheep/ubuntu-sshd