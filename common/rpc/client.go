package rpc

import (
	"fmt"

	"bytes"

	"github.com/33cn/chain33/common"
	log "github.com/inconshreveable/log15"

	"io/ioutil"
	"net/http"

	"github.com/33cn/chain33/types"
)

var (
	chain33_jrpcUrl = "http://localhost:8801"
	clog            = log.New("module", "common")
)

func SetJrpcURL(url string) {
	chain33_jrpcUrl = url
}
func GetJrpcURL() string {
	return chain33_jrpcUrl
}

// TODO: SetPostRunCb()
type RpcCtx struct {
	Addr   string
	Method string
	Params interface{}
	Res    interface{}

	cb Callback
}

type Callback func(res interface{}) (interface{}, error)

func NewRpcCtx(laddr, methed string, params, res interface{}) *RpcCtx {
	return &RpcCtx{
		Addr:   laddr,
		Method: methed,
		Params: params,
		Res:    res,
	}
}

func (c *RpcCtx) SetResultCb(cb Callback) {
	c.cb = cb
}

//func (c *RpcCtx) ReplyData() ([]byte, error) {
//	rpc, err := NewJSONClient(c.Addr)
//	if err != nil {
//		clog.Error(fmt.Sprintf("NewJsonClient have err:%v", err.Error()))
//		return nil, err
//	}
//
//	err = rpc.Call(c.Method, c.Params, c.Res)
//	if err != nil {
//		clog.Error(fmt.Sprintf("rpc.Call have err:%v", err.Error()))
//		return nil, err
//	}
//
//	// maybe format rpc result
//	var result interface{}
//	if c.cb != nil {
//		result, err = c.cb(c.Res)
//		if err != nil {
//			clog.Error(fmt.Sprintf("rpc.CallBack have err:%v", err.Error()))
//			return nil, err
//		}
//	} else {
//		result = c.Res
//	}
//
//	data, err := json.MarshalIndent(result, "", "    ")
//	if err != nil {
//		clog.Error(fmt.Sprintf("MarshalIndent have err:%v", err.Error()))
//		return nil, err
//	}
//	return data, nil
//}

func (c *RpcCtx) ReplyInterface() (interface{}, error) {
	rpc, err := NewJSONClient(c.Addr)
	if err != nil {
		clog.Error(fmt.Sprintf("NewJsonClient have err:%v", err.Error()))
		return nil, err
	}

	err = rpc.Call(c.Method, c.Params, c.Res)
	if err != nil {
		clog.Error(fmt.Sprintf("rpc.Call have err:%v", err.Error()))
		return nil, err
	}

	// maybe format rpc result
	var result interface{}
	if c.cb != nil {
		result, err = c.cb(c.Res)
		if err != nil {
			clog.Error(fmt.Sprintf("rpc.CallBack have err:%v", err.Error()))
			return nil, err
		}
	} else {
		result = c.Res
	}
	return result, nil
}

//func (c *RpcCtx) RunWithoutMarshal() (string, error) {
//	var res string
//	rpc, err := NewJSONClient(c.Addr)
//	if err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		return "", err
//	}
//
//	err = rpc.Call(c.Method, c.Params, &res)
//	if err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		return "", err
//	}
//	return res, nil
//}

func SendTransaction(tx *types.Transaction) (string, error) {
	poststr := fmt.Sprintf(`{"jsonrpc":"2.0","id":2,"method":"Chain33.SendTransaction","params":[{"data":"%v"}]}`,
		common.ToHex(types.Encode(tx)))

	resp, err := http.Post(GetJrpcURL(), "application/json", bytes.NewBufferString(poststr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(b), err
}

func QueryLastHeader(jrpcURl string) (*Header, error) {
	var res Header
	ctx := NewRpcCtx(jrpcURl, Chain33_GetLastHeader, nil, &res)
	_, err := ctx.ReplyInterface()
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func QueryBlockHash(jrpcURl string, height int64) (*ReplyHash, error) {
	var res ReplyHash
	params := types.ReqInt{
		Height: height,
	}
	ctx := NewRpcCtx(jrpcURl, Chain33_GetBlockHash, params, &res)
	_, err := ctx.ReplyInterface()
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func QueryIsSync(jrpcURl string) (bool, error) {
	var res bool
	ctx := NewRpcCtx(jrpcURl, Para_IsSync, nil, &res)
	_, err := ctx.ReplyInterface()
	if err != nil {
		return false, err
	}
	return res, nil
}

func QueryMainNetIsSync(jrpcURl string) (bool, error) {
	var res bool
	ctx := NewRpcCtx(jrpcURl, Chain33_IsSync, nil, &res)
	_, err := ctx.ReplyInterface()
	if err != nil {
		return false, err
	}
	return res, nil
}
func QueryBalance(jrpcURl string, addrList []string) ([]*Account, error) {
	var res []*Account

	var addrs []string
	addrs = append(addrs, addrList...)
	params := types.ReqBalance{
		Addresses: addrs,
		Execer:    "coins",
	}
	ctx := NewRpcCtx(jrpcURl, Chain33_GetBalance, params, &res)
	_, err := ctx.ReplyInterface()
	if err != nil {
		return nil, err
	}
	return res, nil
}
