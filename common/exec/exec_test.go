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
	err :=Init(host)
	if err !=nil {
		t.Error(err)
	}
	res,err :=Exec_CollectResource(host)
	if err !=nil {
		t.Error(err)
	}
	t.Log(res)
}