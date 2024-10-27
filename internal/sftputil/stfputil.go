package sftputil

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/logger"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func GetConnection(ctx context.Context, connectionUrl string, checkKnownHosts bool) (*sftp.Client, func() error, error) {
	parsedUrl, err := url.Parse(connectionUrl)
	if err != nil {
		logger.Error(ctx, err, "failed to parse sftp url")
		return nil, nil, err
	}

	user := parsedUrl.User.Username()
	password, hasPassword := parsedUrl.User.Password()

	if !hasPassword {
		logger.Warn(ctx, "no password provided")
	}

	host := parsedUrl.Hostname()

	// default sft port
	port := 2022

	hostKeyCallback := ssh.InsecureIgnoreHostKey()

	if checkKnownHosts {
		hostKey, err := getHostKey(ctx, host)
		if err != nil {
			logger.Error(ctx, err, "failed to get host key")
			return nil, nil, err
		}
		hostKeyCallback = ssh.FixedHostKey(*hostKey)

	}

	var auths []ssh.AuthMethod
	// Try to use $SSH_AUTH_SOCK which contains the path of the unix file socket that the sshd agent uses
	// for communication with other processes.
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}
	// Use password authentication if provided
	if password != "" {
		auths = append(auths, ssh.Password(password))
	}

	// Initialize client configuration
	logger.Info(ctx, "connecting to ftp server")

	config := ssh.ClientConfig{
		User:            user,
		Auth:            auths,
		HostKeyCallback: hostKeyCallback,
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	// Connect to server
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connecto to [%s]: %v\n", addr, err)
		os.Exit(1)
	}

	sc, err := sftp.NewClient(conn)
	if err != nil {
		logger.Error(ctx, err, "failed to create sftp client")
		return nil, conn.Close, err
	}

	return sc, func() error {
		err := sc.Close()
		if err != nil {
			logger.Error(ctx, err, "failed to close sftp client")
			return err
		}

		err = conn.Close()
		if err != nil {
			logger.Error(ctx, err, "failed to close connection")
			return err
		}
		return nil
	}, nil
}

// Get host key from local know hosts file
func getHostKey(ctx context.Context, host string) (*ssh.PublicKey, error) {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		logger.Error(ctx, err, "failed to open known_hosts file")
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			hostKey, _, _, _, err := ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				logger.Error(ctx, err, "failed to parse authorized key")
				return nil, err
			}

			return &hostKey, nil
		}
	}

	return nil, fmt.Errorf("host key not found")
}
