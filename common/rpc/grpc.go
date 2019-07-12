package rpc

import (
	"context"

	"github.com/33cn/chain33/rpc/grpcclient"
	"github.com/33cn/chain33/types"
	"google.golang.org/grpc"
)

// func QueryRandNumByGrpc(remoteGrpcClient string, param *types.ReqRandHash) (*types.ReplyHash, error) {

// 	// remoteGrpcClient := "127.0.0.1:8004,127.0.0.1:8003,127.0.0.1"
// 	conn, err := grpc.Dial(grpcclient.NewMultipleURL(remoteGrpcClient), grpc.WithInsecure())
// 	if err != nil {
// 		return nil, err
// 	}
// 	grpcClient := types.NewChain33Client(conn)
// 	param := &types.ReqRandHash{
// 		ExecName: "ticket",
// 		BlockNum: 1000,
// 		Hash:     []byte("hello"),
// 	}
// 	reply, err := grpcClient.QueryRandNum(context.Background(), param)
// 	return reply, err
// }

func QueryBlockHashByGrpc(remoteGrpcClient string, height int64) (*types.ReplyHash, error) {
	param := &types.ReqInt{
		Height: height,
	}
	conn, err := grpc.Dial(grpcclient.NewMultipleURL(remoteGrpcClient), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	grpcClient := types.NewChain33Client(conn)
	reply, err := grpcClient.GetBlockHash(context.Background(), param)
	return reply, err
}

func QueryIsSyncByGrpc(remoteGrpcClient string) (*types.Reply, error) {
	param := &types.ReqNil{}
	conn, err := grpc.Dial(grpcclient.NewMultipleURL(remoteGrpcClient), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	grpcClient := types.NewChain33Client(conn)
	reply, err := grpcClient.IsSync(context.Background(), param)
	return reply, err
}

func QueryLastHeaderByGrpc(remoteGrpcClient string) (*types.Header, error) {
	param := &types.ReqNil{}
	conn, err := grpc.Dial(grpcclient.NewMultipleURL(remoteGrpcClient), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	grpcClient := types.NewChain33Client(conn)
	reply, err := grpcClient.GetLastHeader(context.Background(), param)
	return reply, err
}
