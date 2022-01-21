package sftp

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Savepoint(filePath string, identity map[string]string) error {
	host := identity["host"]
	port, _ := strconv.ParseInt(identity["port"], 10, 0)
	user := identity["username"]
	password := identity["password"]
	remotePath := identity["path"]

	connection, err := createConnection(host, int(port), user, password)
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
	port		int
	*sftp.Client
}

func (sc *sftpClient) connect() (err error) {
	config := &ssh.ClientConfig{
		User:            sc.user,
		Auth:            []ssh.AuthMethod{ssh.Password(sc.password)},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", sc.host, sc.port)
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

	// Make remote directories recursion
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

func createConnection(host string, port int, user string, password string) (client *sftpClient, err error) {
	switch {
	case `` == strings.TrimSpace(host),
		`` == strings.TrimSpace(user),
		`` == strings.TrimSpace(password),
		0 >= port || port > 65535:
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