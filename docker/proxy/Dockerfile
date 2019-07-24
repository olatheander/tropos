FROM golang:1.12-alpine AS builder

WORKDIR /src

COPY . .

RUN apk add --no-cache git

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
#RUN go get -d -v ./...

# Dirty hack for https://github.com/moby/moby/issues/36664 (TODO: remove after improving dependency managment)
#RUN rm -r $GOPATH/src/github.com/docker/docker/vendor/github.com/docker/go-connections/nat

# Build the package
RUN go build
RUN echo $PWD

FROM tropos-base

LABEL maintainer="Ola Theander <ola.theander@myola.se>"

# Copy tropos executable
COPY --from=builder /src/tropos /usr/local/bin

# Install kubectl
#RUN apt-get update && apt-get install -y apt-transport-https curl gnupg
#RUN curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
#RUN echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list
#RUN  apt-get update
#RUN apt-get install -y kubectl
#
#RUN apt-get clean && \
#    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

VOLUME [ "/workspace" ]
