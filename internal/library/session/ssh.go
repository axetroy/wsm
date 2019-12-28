package session

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
)

type Terminal struct {
	client       *ssh.Client
	session      *ssh.Session
	config       Config
	exitMsg      string
	stdout       io.Reader
	stdin        io.Writer
	stderr       io.Reader
	closeHandler func() error
	closed       bool
}

type Config struct {
	Username string
	Host     string
	Port     uint
	Password string
	Width    int // pty width
	Height   int // pty height
}

func (t *Terminal) SetCloseHandler(h func() error) {
	t.closeHandler = h
}

func (t *Terminal) Close() (err error) {
	if t.closed {
		return nil
	}
	defer func() {
		if t.closeHandler != nil {
			err = t.closeHandler()
		}
		t.closed = true
	}()

	if err = t.session.Close(); err != nil {
		return
	}

	if err = t.client.Close(); err != nil {
		return
	}

	return
}

func (t *Terminal) Connect(stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	var err error
	defer func() {
		if t.exitMsg == "" {
			// _, _ = fmt.Fprintln(stdout, "the connection was closed on the remote side on ", time.Now().Format(time.RFC822))
		} else {
			_, _ = fmt.Fprintln(stdout, t.exitMsg)
		}
	}()

	termType := os.Getenv("TERM")

	if termType == "" {
		termType = "vt100"
	}

	if err = t.session.RequestPty(termType, t.config.Height, t.config.Width, ssh.TerminalModes{}); err != nil {
		return err
	}

	t.session.Stdin = stdin
	t.session.Stderr = stderr
	t.session.Stdout = stderr

	if err = t.session.Shell(); err != nil {
		return err
	}

	quit := make(chan int)
	go func() {
		_ = t.session.Wait()
		_ = t.Close()
		quit <- 1
	}()

	return nil
}

func NewTerminal(config Config) (*Terminal, error) {
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		BannerCallback:  ssh.BannerDisplayStderr(),
	}

	addr := net.JoinHostPort(config.Host, strconv.Itoa(int(config.Port)))

	client, err := ssh.Dial("tcp", addr, sshConfig)

	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()

	if err != nil {
		return nil, err
	}

	s := Terminal{
		client:  client,
		config:  config,
		session: session,
	}

	return &s, nil
}

// 测试服务器是否可用
func Test(config Config) bool {
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		BannerCallback:  ssh.BannerDisplayStderr(),
	}

	addr := net.JoinHostPort(config.Host, strconv.Itoa(int(config.Port)))

	client, err := ssh.Dial("tcp", addr, sshConfig)

	if err != nil {
		return false
	}

	defer client.Close()

	return true
}
