package sftp

import (
	"github.com/harrylee2015/monitor/model"
	"testing"
)

func TestSftpClient_Exists(t *testing.T) {
	hostInfo := &model.HostInfo{
		HostIp:   "127.0.0.1",
		SSHPort:  22,
		UserName: "harry",
		PassWd:   "123456",
	}
	client, err := NewSftpClient(hostInfo)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer client.client.Close()
	isExist := client.Exists("/tmp/monitor/gopsutil")
	t.Logf("isExist:%v", isExist)
}
func TestSftpClient_UploadFile(t *testing.T) {
	hostInfo := &model.HostInfo{
		HostIp:   "192.168.0.194",
		UserName: "ubuntu",
		PassWd:   "123456",
		SSHPort:  22,
	}
	client, err := NewSftpClient(hostInfo)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer client.client.Close()
	err = client.UploadFile("../../build/gopsutil", "/tmp/monitor/gopsutil")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSftpClient_DeleteFile(t *testing.T) {
	hostInfo := &model.HostInfo{
		HostIp:   "192.168.0.194",
		UserName: "ubuntu",
		PassWd:   "123456",
		SSHPort:  22,
	}
	client, err := NewSftpClient(hostInfo)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	defer client.client.Close()
	err = client.DeleteFile("/tmp/monitor/gopsutil")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func BenchmarkSftpClient_UploadFile(b *testing.B) {
	for i:=0;i<b.N;i++{
		hostInfo := &model.HostInfo{
			HostIp:   "192.168.0.194",
			UserName: "ubuntu",
			PassWd:   "123456",
			SSHPort:  22,
		}
		client, err := NewSftpClient(hostInfo)
		if err != nil {
			b.Error(err)
			b.Fail()
		}
		defer client.client.Close()
		err = client.UploadFile("../../build/gopsutil", "/tmp/monitor/gopsutil")
		if err != nil {
			b.Error(err)
			b.Fail()
		}
	}
}