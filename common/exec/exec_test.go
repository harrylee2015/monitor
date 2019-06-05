package exec
import (
	"github.com/harrylee2015/monitor/model"
	"testing"
)
func TestExec_CollectResource(t *testing.T) {
	host :=&model.HostInfo{
		HostIp:"127.0.0.1",
		UserName:"harry",
		PassWd:"123456",
		SSHPort:22,
	}
	err :=Exec_Scp("../../build/gopsutil","/tmp/monitor",host)
	if err !=nil {
		t.Error(err)
		t.Fail()
	}
	res,err :=Exec_CollectResource(host)
	if err !=nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(res)
}
func TestExec_Scp(t *testing.T) {
	host :=&model.HostInfo{
		HostIp:"192.168.0.194",
		UserName:"ubuntu",
		PassWd:"123456",
		SSHPort:22,
	}
	err :=Exec_Scp("../../build/gopsutil","/tmp/monitor",host)
	if err !=nil {
		t.Error(err)
		t.Fail()
	}
}
