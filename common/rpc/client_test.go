package rpc

import (
	"fmt"
	"testing"
)

func getJrpc(ip string, port int64) string {
	return fmt.Sprintf("http://%s:%v", ip, port)
}

func TestQueryLastHeader(t *testing.T) {
	head, err := QueryLastHeader(getJrpc("192.168.0.194", 9901))
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(head)
	t.Logf("lashBlockHash:%v",head.Hash)
}

func TestQueryBalance(t *testing.T) {
	accounts, err := QueryBalance(getJrpc("192.168.0.194", 8801), []string{"1MNwSHugmoBcGxqfU8Gt4excUg68DA7hYE"})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(accounts)
	t.Logf("balance:%v",accounts[0].Balance)
}

func TestQueryIsSync(t *testing.T) {
    isSync,err:=QueryIsSync(getJrpc("192.168.0.194", 9901))
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(isSync)
}

func TestQueryBlockHash(t *testing.T) {
    reply,err:= QueryBlockHash(getJrpc("192.168.0.194", 9901),100)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Logf("height %d block, blockhash:%s",100,reply.Hash)
}
