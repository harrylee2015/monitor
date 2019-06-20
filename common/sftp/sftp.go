package sftp

import (
	"fmt"
	"github.com/harrylee2015/monitor/model"
	"io"
	"io/ioutil"
	"os"
	"time"

	"path/filepath"

	log "github.com/inconshreveable/log15"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//构造sftpClient
func NewSftpClient(host *model.HostInfo) (*SftpClient, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	) // get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(host.PassWd))
	clientConfig = &ssh.ClientConfig{User: host.UserName,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	} // connet to ssh
	addr = fmt.Sprintf("%s:%d", host.HostIp, host.SSHPort)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		sshClient.Close()
		return nil, err
	}
	return &SftpClient{sftpClient, sshClient}, nil
}

type SftpClient struct {
	Sftpclient *sftp.Client
	SSHClient  *ssh.Client
}

func (ftp *SftpClient) UploadFile(localFile string, remoteFile string) error {
	srcFile, err := os.Open(localFile)
	if err != nil {
		log.Error("open local file error", "error", err)
		return err
	}
	defer srcFile.Close()

	dir := filepath.Dir(remoteFile)
	if !ftp.Exists(dir) {
		ftp.Sftpclient.MkdirAll(dir)
	}
	dstFile, err := ftp.Sftpclient.Create(remoteFile)
	if err != nil {
		log.Error("sftpClient.Create error :", "error", err, "file", remoteFile)
		return err
	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		log.Error("ReadAll error : ", "error", err, "file:", localFile)
		return err
	}
	_, err = dstFile.Write(ff)
	return err
}

func (ftp *SftpClient) DownloadFile(localFile string, remoteFile string) error {
	fs, err := ftp.Sftpclient.Open(remoteFile)
	if err != nil {
		log.Error("open remote file error", "error", err, "file", remoteFile)
		return err
	}
	defer fs.Close()
	info, _ := fs.Stat()
	local, _ := os.OpenFile(localFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	_, err = io.Copy(local, io.LimitReader(fs, info.Size()))
	return err
}

func (ftp *SftpClient) Exists(remoteFile string) bool {
	// check it's there
	_, err := ftp.Sftpclient.Lstat(remoteFile)
	if err != nil {
		log.Info(err.Error())
		return false
	}
	return true
}

func (ftp *SftpClient) DeleteFile(remoteFile string) error {
	return ftp.Sftpclient.Remove(remoteFile)
}

//用完一定要关闭链接，不然会造成资源泄漏
func (ftp *SftpClient) Close() {
	ftp.Sftpclient.Close()
	ftp.SSHClient.Close()
}
