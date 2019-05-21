package hostmgr

import (
	"github.com/harrylee2015/monitor/model"
	. "github.com/harrylee2015/monitor/web/view/webutil"
	"github.com/kataras/iris"
)

func ListHosts(ctx iris.Context) {
	var page model.Page
	if err := ctx.ReadJSON(&page); err != nil {
		ClientErr(ctx, err)
		return
	}
	items := model.GetDB().QueryHostInfoByPageNum(&page)
	ctx.JSON(items)
}

func AddHost(ctx iris.Context) {
	var host model.HostInfo
	if err := ctx.ReadJSON(&host); err != nil {
		ClientErr(ctx, err)
		return
	}
	model.GetDB().InsertData(&host)
	ServOK(ctx)
}

func UpdateHost(ctx iris.Context) {
	var host model.HostInfo
	if err := ctx.ReadJSON(&host); err != nil {
		ClientErr(ctx, err)
		return
	}
	model.GetDB().UpdateData(&host)
	ServOK(ctx)
}

func DeleteHost(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("hostId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	model.GetDB().DeleteDataByHostId(id)
	ServOK(ctx)
}
