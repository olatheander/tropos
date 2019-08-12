FROM golang:1.12-alpine AS builder

WORKDIR /src

COPY . .

RUN apk add --no-cache git

# Build the package
RUN go build

FROM tropos-base

LABEL maintainer="Ola Theander <ola.theander@myola.se>"

ENV KUBECONFIG=/kube/config

# Copy tropos executable
COPY --from=builder /src/tropos /usr/local/bin

# Install kubectl TODO: just for development, to be removed from image.
RUN apt-get update && apt-get install -y apt-transport-https curl gnupg
RUN curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
RUN echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list
RUN  apt-get update
RUN apt-get install -y kubectl
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Volume for source code
VOLUME [ "/workspace" ]
# Volume for Kubernetes context file
VOLUME [ "/kube" ]
