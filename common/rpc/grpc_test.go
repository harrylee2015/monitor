package rpc

import (
	"testing"
)

var remoteClient = "101.37.227.226:8802,39.97.20.242:8802,47.107.15.126:8802,jiedian2.bityuan.com,cloud.bityuan.com"

func TestQueryBlockHashByGrpc(t *testing.T) {
	reply, err := QueryBlockHashByGrpc(remoteClient, 10000)
	if err != nil {
		t.Error(err)
	}
	t.Log(reply.GetHash())
}

func TestQueryLastHeaderByGrpc(t *testing.T) {
	reply, err := QueryLastHeaderByGrpc(remoteClient)
	if err != nil {
		t.Error(err)
	}
	t.Log(reply.GetHeight())
}

func TestQueryIsSyncByGrpc(t *testing.T) {
	reply, err := QueryIsSyncByGrpc(remoteClient)
	if err != nil {
		t.Error(err)
	}
	t.Log(reply.GetIsOk())
	t.Log(reply.GetMsg())
}
