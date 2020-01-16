// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
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
	stdout       io.Reader
	stdin        io.Writer
	stderr       io.Reader
	closeHandler func() error
	closed       bool
}

type Config struct {
	Username   string
	Host       string
	Port       uint
	Password   string // 用密码连接
	PrivateKey string // 用私钥连接，如果设置了私钥，优先使用私钥连接
	Width      int    // pty width
	Height     int    // pty height
}

func (t *Terminal) SetCloseHandler(h func() error) {
	t.closeHandler = h
}

// 终端是否已断臂
func (t *Terminal) IsClosed() bool {
	return t.closed
}

func (t *Terminal) Close() (err error) {
	if t.IsClosed() {
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

	termType := os.Getenv("TERM")

	if termType == "" {
		termType = "vt100"
	}

	if err = t.session.RequestPty(termType, t.config.Height, t.config.Width, ssh.TerminalModes{}); err != nil {
		return err
	}

	t.session.Stdin = stdin
	t.session.Stderr = stderr
	t.session.Stdout = stdout

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
	var authMethods []ssh.AuthMethod

	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		BannerCallback:  ssh.BannerDisplayStderr(),
	}

	if config.PrivateKey != "" {
		if pk, err := ssh.ParsePrivateKey([]byte(config.PrivateKey)); err != nil {
			return nil, err
		} else {
			authMethods = append(authMethods, ssh.PublicKeys(pk))
		}
	} else {
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	sshConfig.Auth = authMethods

	addr := net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port))

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
