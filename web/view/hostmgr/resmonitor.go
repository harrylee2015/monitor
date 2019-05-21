package hostmgr

import (
	"github.com/harrylee2015/monitor/model"
	. "github.com/harrylee2015/monitor/web/view/webutil"
	"github.com/kataras/iris"
)

func ListResource(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("hostId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	//TODO
	items := model.GetDB().QueryResourceInfo(id, 100)
	ctx.JSON(items)
}

func GetResWaringCount(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	count := model.GetDB().QueryResWarningCount(id)
	ctx.JSON(count)
}

func GetResWaringByHostId(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("hostId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	items := model.GetDB().QueryWarningByHostId(id)
	ctx.JSON(items)
}
