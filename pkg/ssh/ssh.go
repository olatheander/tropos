package ssh

//Basically a clone of the https://gist.github.com/codref/473351a24a3ef90162cf10857fac0ff3

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// Open an SSH tunnel (ssh -R)
func NewSSHTunnel(user string,
	publicKeyPath string,
	serverEndpoint *Endpoint,
	localEndpoint *Endpoint,
	remoteEndpoint *Endpoint) error {
	sshClientConfig := &ssh.ClientConfig{
		// SSH connection username
		User: user,
		Auth: []ssh.AuthMethod{
			// put here your private key path
			publicKeyFile(publicKeyPath),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH remote server using serverEndpoint
	serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshClientConfig)
	if err != nil {
		panic(err)
	}

	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, err := net.Dial("tcp", localEndpoint.String())
		if err != nil {
			panic(err)
		}

		client, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		handleClient(client, local)
	}

	return nil
}

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			panic(err)
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			panic(err)
		}
		chDone <- true
	}()

	<-chDone
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(key)
}
