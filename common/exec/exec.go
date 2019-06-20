package exec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ThomasRooney/gexpect"
	"github.com/harrylee2015/monitor/conf"
	"net"
	"os"
	"strings"
	"time"

	"github.com/harrylee2015/monitor/model"
	log "github.com/inconshreveable/log15"
	"golang.org/x/crypto/ssh"
)

const (
	REMOTEPATH = "/tmp/monitor"
)

type ScpInfo struct {
	UserName      string
	PassWord      string
	HostIp        string
	Port          int
	LocalFilePath string
	RemoteDir     string
}

type CmdInfo struct {
	UserName string
	PassWord string
	HostIp   string
	Port     int
	Cmd      string
}
type ExecClient struct {
	SSHClient *ssh.Client

	SSHSession *ssh.Session
}

func (client *ExecClient) Close() {
	client.SSHSession.Close()
	client.SSHClient.Close()
}
func NewExecClient(user, password, host string, port int) (*ExecClient, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		//需要验证服务端，不做验证返回nil就可以
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create session
	if session, err = client.NewSession(); err != nil {
		client.Close()
		return nil, err
	}
	return &ExecClient{SSHClient: client, SSHSession: session}, nil
}

func RemoteExec(cmdInfo *CmdInfo) error {
	//A Session only accepts one call to Run, Start or Shell.
	client, err := NewExecClient(cmdInfo.UserName, cmdInfo.PassWord, cmdInfo.HostIp, cmdInfo.Port)
	if err != nil {
		return err
	}
	defer client.Close()
	client.SSHSession.Stdout = os.Stdout
	client.SSHSession.Stderr = os.Stderr
	err = client.SSHSession.Run(cmdInfo.Cmd)
	return err
}

//上传信息采集脚本
func UploadScriptsFile(host *model.HostInfo) error {
	return Exec_Scp(GetScriptsPath(), REMOTEPATH, host)
}

//TODO 这个方法有bug,goexpect bug 待修复,暂时用sftp替代
func Exec_Scp(localFilePath, remotePath string, host *model.HostInfo) error {
	cmdInfo := &CmdInfo{
		UserName: host.UserName,
		PassWord: host.PassWd,
		HostIp:   host.HostIp,
		Port:     int(host.SSHPort),
		Cmd:      fmt.Sprintf("mkdir -p %v", remotePath),
	}
	RemoteExec(cmdInfo)

	cmd := fmt.Sprintf("scp -P %v %s  %s@%s:%s", host.SSHPort, localFilePath, host.UserName, host.HostIp, remotePath)
	pwd := host.PassWd

	child, err := gexpect.Spawn(cmd)
	if err != nil {
		log.Error("spawn", err)
		return err
	}
	//如果是第一次使用scp的话，会出现（yes/no)选项
	//if err := child.ExpectTimeout("(yes/no)", 100*time.Microsecond); err != nil {
	//	log.Warn("Expect timieout error ", err)
	//}

	if err := child.SendLine(pwd); err != nil {
		log.Error("SendLine password error ", err)
		return err
	}

	if err := child.Wait(); err != nil {
		log.Error("Wait error: ", "", err.Error())
		return err
	}

	log.Info("Success")
	return nil
}
func Exec_CollectResource(host *model.HostInfo) (*model.ResourceInfo, error) {
	client, err := NewExecClient(host.UserName, host.PassWd, host.HostIp, int(host.SSHPort))
	if err != nil {
		log.Error("Exec_CollectResource", "sshconnect err:", err.Error())
		return nil, err
	}
	defer client.Close()

	// 创建输入管道，以支持顺序执行多个命令
	pipe, err := client.SSHSession.StdinPipe()
	if err != nil {
		log.Error("Exec_CollectResource", "StdinPipe err:", err.Error())
		return nil, err
	}
	defer pipe.Close()
	//定义一个buffer 字节数组用于接收ssh会话输出内容
	bout := bytes.NewBuffer(nil)
	berr := bytes.NewBuffer(nil)
	//session.Stdout = os.Stdout
	//session.Stderr = os.Stdout
	client.SSHSession.Stdout = bout
	client.SSHSession.Stderr = berr
	// 启动远程执行shell
	err = client.SSHSession.Shell()
	if err != nil {
		log.Error("Exec_CollectResource", "Shell err:", err.Error())
		return nil, err
	}
	commands := genCollectScript()
	log.Debug("Exec_CollectResource","cmds:", commands)
	// 执行指定脚本内容
	for _, cmd := range commands {
		_, err := fmt.Fprintf(pipe, "%s\n", cmd)
		if err != nil {
			log.Error("execute command error:", err, "command", cmd)
			return nil, err
		}
	}
	err = client.SSHSession.Wait()
	if err != nil {
		log.Error("Exec_CollectResource", "Wait err:", err.Error())
		return nil, err
	}

	str := bout.String()
	strs := strings.Split(str, "data=============:")
	var resource model.ResourceInfo
	err = json.Unmarshal([]byte(strs[1]), &resource)
	if err != nil {
		return nil, err
	}
	log.Info("Exec_CollectResource", "resource:", resource)
	return &resource, nil
}

// 生产采集数据命令
func genCollectScript() []string {
	commands := []string{
		fmt.Sprintf("cd %s", REMOTEPATH),
		"chmod +x gopsutil",
		"./gopsutil",
		"exit",
	}
	return commands
}
func GetScriptsPath() string {
	return fmt.Sprintf("%s/gopsutil", conf.CurrDir)
}
func GetRemoteScriptsPath() string {
	return fmt.Sprintf("%s/gopsutil", REMOTEPATH)
}
