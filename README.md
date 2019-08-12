# Tropos
Cloud-native development on Kubernetes utilizing Visual Studio Code's Remote Development features

## Build proxy container
```
$ docker build -t tropos-proxy -f docker/proxy/Dockerfile .
```

## Run proxy container
Note!, starting the proxy container will eventually be done by the Tropos CLI but for now it's done manually. E.g.
```
$ docker run -d --rm -v <src>:/workspace -v /home/olathe/.minikube:/home/olathe/.minikube -v $HOME/.kube/config:/kube/config -p 2022:22 tropos-proxy
```

## Tropos Pod

The Tropos pod will mount the workspace mounted in the proxy container by the equivalent of 
```
$ ssh username@server -i -R 10000:localmachine:22
$ sshfs -p 10000 localusername@127.0.0.1:/workspace /workspace 
```
in order to establish a synchronisation of the workspace between the development host and the pod (via the proxy container).

The pod will only allow key based login thus the Tropos executable must to the equivalent of
```
$ docker cp ~/.ssh/id_rsa.pub sshd:/root/.ssh/authorized_keys 
$ docker exec sshd chown root:root /root/.ssh/authorized_keys
```
to the pod container in order to be able to SSH in (taken from https://github.com/rastasheep/ubuntu-sshd). SSH by
```
$ ssh -i ~/.ssh/id_rsa root@localhost -p 2022
```
and possibly delete a cached key by `ssh-keygen -R [localhost]:2022`.

### Pod security policies
```
securityContext:
    privileged: true
    capabilities:
    add:
        - SYS_ADMIN
```
