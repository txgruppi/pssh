package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

var auth goph.Auth

func init() {
	var err error
	auth, err = goph.UseAgent()
	if err != nil {
		panic(err)
	}
}

func sshUploadRunAndCleanup(host Host, filepath string, useSudo, dryRun, noKnownHosts bool) (stdout, stderr []byte, err error) {
	client, err := goph.NewConn(&goph.Config{
		Auth: auth,
		Addr: host.Addr,
		Port: 22,
		User: host.User,
		Callback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			if noKnownHosts {
				return nil
			}
			found, err := goph.CheckKnownHost(hostname, remote, key, "")
			if found && err != nil {
				return fmt.Errorf("error checking known host %s: %w", hostname, err)
			}
			if found && err == nil {
				return fmt.Errorf("error checking known host %s: %w", hostname, err)
			}
			slog.Warn("adding new host to known_hosts", "host", hostname, "addr", remote.String())
			err = goph.AddKnownHost(hostname, remote, key, "")
			if err != nil {
				return fmt.Errorf("error adding known host %s: %w", hostname, err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, nil, err
	}
	defer client.Close()

	if dryRun {
		slog.Info("dry-run: upload", "host", host.Name, "file", filepath)
	} else {
		err = upload(client, filepath, filepath)
		if err != nil {
			return nil, nil, err
		}
	}

	defer func() {
		if dryRun {
			slog.Info("dry-run: cleanup", "host", host.Name, "file", filepath)
		} else {
			_, err = client.Run("rm -f " + filepath)
			if err != nil {
				slog.Error("cleanup failed", "host", host.Name, "file", filepath, "error", err)
			}
		}
	}()

	var cmd *goph.Cmd
	if useSudo {
		if dryRun {
			slog.Info("dry-run: sudo run", "host", host.Name, "file", filepath)
		} else {
			if host.SudoPassword == "" {
				cmd, err = client.Command("sudo", filepath)
				if err != nil {
					return nil, nil, err
				}
			} else {
				cmd, err = client.Command("sudo", "--stdin", "--prompt=", filepath)
				if err != nil {
					return nil, nil, err
				}
				pw := host.SudoPassword
				if !strings.HasSuffix(pw, "\n") {
					pw += "\n"
				}
				cmd.Stdin = strings.NewReader(pw)
			}
		}
	} else {
		if dryRun {
			slog.Info("dry-run: run", "host", host.Name, "file", filepath)
		} else {
			cmd, err = client.Command(filepath)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	var outbuf, errbuf bytes.Buffer
	if !dryRun {
		cmd.Stdout = &outbuf
		cmd.Stderr = &errbuf
		err = cmd.Run()
		if err != nil {
			return outbuf.Bytes(), errbuf.Bytes(), err
		}
	}

	return outbuf.Bytes(), errbuf.Bytes(), nil
}

func upload(client *goph.Client, local, remote string) error {
	localFile, err := os.Open(local)
	if err != nil {
		return err
	}
	defer localFile.Close()

	info, err := localFile.Stat()
	if err != nil {
		return err
	}

	ftp, err := client.NewSftp()
	if err != nil {
		return err
	}
	defer ftp.Close()

	remoteFile, err := ftp.Create(remote)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return err
	}

	err = remoteFile.Chmod(info.Mode())
	if err != nil {
		return err
	}

	return nil
}
