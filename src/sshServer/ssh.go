package ssh

import (
	"io"
	"net"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

var config *ssh.ServerConfig

func RemoteServer() {
	// listen on port 2222
	listener, err := net.Listen("tcp", " :2222")
	if err != nil {
		panic("Failed to listen on 2222 ( " + err.Error() + " )")
	}

	// accept connection
	nConn, err := listener.Accept()
	if err != nil {
		panic("Failed to accept incoming connection ( " + err.Error() + " )")
	}
	// handshake
	_, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		panic("Failed to handshake ( " + err.Error() + " )")
	}
	// discard all global out-of-band Requests
	go ssh.DiscardRequests(reqs)
	// accept all channels
	for newChannel := range chans {

		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			panic("Could not accept channel ( " + err.Error() + " )")
		}
		// allocate a terminal for this channel
		term := NewTerminal(channel)
		// start the shell
		go term.Shell()
		// start stdin -> channel
		go io.Copy(channel, os.Stdin)
		// start channel -> stdout
		go io.Copy(os.Stdout, channel)
		// start out-of-band Requests
		go func() {
			for req := range requests {
				if req.Type == "exec" {
					// execute a command
					cmd := string(req.Payload[4:])
					term.Execute(cmd)
				}
			}
		}()
	}

}

func NewTerminal(channel ssh.Channel) *Terminal {
	term := &Terminal{
		channel: channel,
		in:      channel,
		out:     channel,
	}
	return term
}

type Terminal struct {
	channel ssh.Channel
	in      io.Reader
	out     io.Writer
}

func (t *Terminal) Shell() {
	// start a shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	cmd := exec.Command(shell)
	cmd.Stdin = t.in
	cmd.Stdout = t.out
	cmd.Stderr = t.out
	cmd.Run()
}

func (t *Terminal) Execute(cmd string) {
	// execute a command
	c := exec.Command(cmd)
	c.Stdin = t.in
	c.Stdout = t.out
	c.Stderr = t.out
	c.Run()
}

func (t *Terminal) Close() {
	t.channel.Close()
}

func main() {
	RemoteServer()
}
