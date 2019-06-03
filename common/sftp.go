package common

import (
	"io"
	"io/ioutil"
	"os"

	"path/filepath"

	"github.com/harrylee2015/monitor/model"
	log "github.com/inconshreveable/log15"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SftpClient struct {
	*model.HostInfo
	conn   *ssh.Client
	client *sftp.Client
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
		ftp.client.MkdirAll(dir)
	}
	dstFile, err := ftp.client.Create(remoteFile)
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
	fs, err := ftp.client.Open(remoteFile)
	if err != nil {
		log.Error("open remote file error", "error", err, "file", remoteFile)
		return err
	}
	info, _ := fs.Stat()
	local, _ := os.OpenFile(localFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	_, err = io.Copy(local, io.LimitReader(fs, info.Size()))
	return err
}

func (ftp *SftpClient) Exists(remoteFile string) bool {
	// check it's there
	_, err := ftp.client.Lstat(remoteFile)
	if err != nil {
		log.Info(err.Error())
		return false
	}
	return true
}

func (ftp *SftpClient) DeleteFile(remoteFile string) error {
	return ftp.client.Remove(remoteFile)
}
