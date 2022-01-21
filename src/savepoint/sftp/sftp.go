package sftp

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Savepoint(filePath string, identity map[string]string) error {
	host := identity["host"]
	port := identity["port"]
	user := identity["username"]
	password := identity["password"]
	remotePath := identity["path"]

	connection, err := createConnection(host, port, user, password)
	if err != nil {
		return errors.New("failed to connect sftp server: " + err.Error())
	}

	err = connection.put(filePath, remotePath + "/" + strings.Replace(filePath, "./tmp/", "", 1))
	if err != nil {
		return errors.New("failed to send a backup to sftp: " + err.Error())
	}

	return nil
}

type sftpClient struct {
	host		string
	user		string
	password	string
	port		string
	*sftp.Client
}

func (sc *sftpClient) connect() (err error) {
	config := &ssh.ClientConfig{
		User:            sc.user,
		Auth:            []ssh.AuthMethod{ssh.Password(sc.password)},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%s", sc.host, sc.port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	sc.Client = client

	return nil
}

func (sc *sftpClient) put(localFile, remoteFile string) (err error) {
	srcFile, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer srcFile.Close()

	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		_ = sc.Mkdir(path)
	}

	dstFile, err := sc.Create(remoteFile)
	if err != nil {
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return
}

func createConnection(host string, port string, user string, password string) (client *sftpClient, err error) {
	switch {
	case `` == strings.TrimSpace(host),
		`` == strings.TrimSpace(user),
		`` == strings.TrimSpace(password),
		`` == strings.TrimSpace(port):
		return nil, errors.New("invalid parameters")
	}

	client = &sftpClient{
		host:     host,
		user:     user,
		password: password,
		port:     port,
	}

	if err = client.connect(); nil != err {
		return nil, err
	}
	return client, nil
}