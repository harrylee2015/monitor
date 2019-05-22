package web

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"

	"github.com/harrylee2015/monitor/web/view/hostmgr"
	"github.com/harrylee2015/monitor/web/view/taskmgr"
)

type WebServer struct {
	// 监听信息
	addr string

	// 服务实例
	app *iris.Application
}

func NewWebServer(addr string) *WebServer {
	server := &WebServer{addr, iris.New()}
	return server
}

func (server *WebServer) Start() {
	server.init()
	//crontab task
	go taskmgr.CronTask()

	server.app.Logger().SetLevel("debug")
	server.app.Run(iris.Addr(server.addr))
}

func (s *WebServer) init() {

	// 注册访问路由
	s.route("GET", "/", index)

	s.routeHostmgr()
	s.routeHostGroupmgr()
	s.routeAddressmgr()
	s.routeMonitormgr()
}

// 添加主机组管理路由信息
func (s *WebServer) routeHostGroupmgr() {
	// 分组列表
	s.route("POST", "/hostmgr/group", hostmgr.ListHostGroups)
	// 查看单个分组
	//s.route("GET", "/hostmgr/group/{id:int}", hostmgr.GetHostGroup)
	// 创建分组
	s.route("POST", "/hostmgr/group", hostmgr.AddHostGroup)
	// 修改分组
	s.route("POST", "/hostmgr/group", hostmgr.UpdateHostGroup)
	// 删除分组
	s.route("DELETE", "/hostmgr/group/{groupId:int}", hostmgr.DeleteHostGroup)
}

func (s *WebServer) routeHostmgr() {
	// 主机列表
	s.route("POST", "/hostmgr", hostmgr.ListHosts)
	// 查看单个主机
	//s.route("GET", "/hostmgr/{id:int}", hostmgr.GetHost)
	// 创建主机
	s.route("POST", "/hostmgr", hostmgr.AddHost)
	// 修改主机
	s.route("POST", "/hostmgr", hostmgr.UpdateHost)
	// 删除单个主机
	s.route("DELETE", "/hostmgr/{hostId:int}", hostmgr.DeleteHost)
}

func (s *WebServer) routeAddressmgr() {
	// 主机列表
	s.route("POST", "/addressmgr", hostmgr.ListAddress)

	// 创建主机
	s.route("POST", "/addressmgr", hostmgr.AddAddress)
	// 修改主机
	s.route("POST", "/addressmgr", hostmgr.UpdateAddress)
	// 删除单个主机
	s.route("DELETE", "/addressmgr/{groupId:int}", hostmgr.DeleteAddress)
}

func (s *WebServer) routeMonitormgr() {
	// 主机列表
	s.route("GET", "/monitormgr/{groupId:int}", hostmgr.ListMonitors)
	s.route("GET", "/monitormgr/buswarningcount/{groupId:int}", hostmgr.GetBusWaringCount)
	s.route("GET", "/monitormgr/reswarningcount/{groupId:int}", hostmgr.GetResWaringCount)

	s.route("GET", "/monitormgr/buswarning/{groupId:int}", hostmgr.GetBusWaringByGroupId)
	s.route("GET", "/monitormgr/reswarning/{hostId:int}", hostmgr.GetResWaringByHostId)

	s.route("POST", "/monitormgr/warning/remove", hostmgr.DeletewarningById)
	s.route("POST", "/monitormgr/warning/batchremove", hostmgr.DeletewarningByList)
	s.route("POST", "/monitormgr/warning/history", hostmgr.GetHistoryWarning)

	s.route("GET", "/monitormgr/balance/{groupId:int}", hostmgr.GetPaymentAddressBalance)
	s.route("GET", "/monitormgr/balancelist/{groupId:int}", hostmgr.GetBalanceListByTime)
	s.route("GET", "/monitormgr/resourcelist/{hostId:int}", hostmgr.ListResource)
}

// 添加路由信息
func (server *WebServer) route(method string, uri string, handlers ...context.Handler) {
	switch method {
	case "GET":
		server.app.Get(uri, handlers...)
	case "POST":
		server.app.Post(uri, handlers...)
	case "PUT":
		server.app.Put(uri, handlers...)
	case "DELETE":
		server.app.Delete(uri, handlers...)
	default:
		fmt.Println(fmt.Sprintf("no support for %v:%v", method, uri))
	}
}

func profile(ctx iris.Context) {
	// finally, render the template.
	ctx.View("user/profile.html")
}

func index(ctx iris.Context) {
	ctx.JSON(iris.Map{"message": "Hello World"})
	//ctx.View("index.html")
}
