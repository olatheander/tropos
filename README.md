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
$ ssh username@server -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null -R 10000:localmachine:22
$ sshfs -p 10000 localusername@127.0.0.1:/workspace /workspace 
```
in order to establish a synchronisation of the workspace between the development host and the pod (via the proxy container), see e.g.
* https://superuser.com/questions/616182/how-to-mount-local-directory-to-remote-like-sshfs
* https://askubuntu.com/questions/1090715/fuse-bad-mount-point-mnt-transport-endpoint-is-not-connected

The pod will only allow key based login thus the Tropos executable must to the equivalent of
```
$ docker cp ~/.ssh/id_rsa.pub sshd:/root/.ssh/authorized_keys 
$ docker exec sshd chown root:root /root/.ssh/authorized_keys
```
to the pod container in order to be able to SSH in (taken from https://github.com/rastasheep/ubuntu-sshd). SSH by
```
$ ssh -i ~/.ssh/id_rsa -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null root@localhost -p 2022
```
and possibly delete a cached key by `ssh-keygen -R [localhost]:2022`.

## Tropos base image
Have the Tropos container generate a new pair of keys on the fly like:
```
ssh-keygen -t rsa
```
and then copy the key to the pod running in K8s for password-less SSH (improvement is to generate a new key-pair in the pod too and copy public back to Tropos container):
```
# kubectl cp ~/.ssh/id_rsa.pub tropos-58d96c958d-d4799:/root/.ssh/authorized_keys
# kubectl exec tropos-58d96c958d-d4799 -- chown root:root /root/.ssh/authorized_keys
```
Set up port forwarding from the Tropos pod to the Tropos container:
```
# kubectl port-forward pod/tropos-58d96c958d-d4799 2022:22 &
```
and then connect with SSH enabling remote forwarding over SSH and run `sshfs` mounting `/workspace` in the Tropos container:
```
# ssh -i ~/.ssh/id_rsa root@localhost -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null -p 2022 -R 10000:localhost:22
# mkdir /workspace
# sshfs -o IdentityFile=/root/id_rsa -p 10000 -C root@127.0.0.1:/workspace /workspace
```
Here `-p 10000` is the port number and `-C` enable compression. 

### Pod security policies
```
securityContext:
    privileged: true
    capabilities:
    add:
        - SYS_ADMIN
```
