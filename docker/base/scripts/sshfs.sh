#!/bin/bash

keyPath=$1
userAndHost=$2:root@127.0.0.1
port=$3:10000

sshfs -o IdentityFile=$keyPath -p $port -C -F /dev/null -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null $userAndHost:/workspace /workspace