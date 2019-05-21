package hostmgr

import (
	"github.com/harrylee2015/monitor/model"
	. "github.com/harrylee2015/monitor/web/view/webutil"
	"github.com/kataras/iris"
)

func ListMonitors(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	items := model.GetDB().QueryMonitor(id)
	ctx.JSON(items)
}

func GetBusWaringCount(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	count := model.GetDB().QueryBusWarningCount(id)
	ctx.JSON(count)
}

func GetBusWaringByGroupId(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	items := model.GetDB().QueryWarningByGroupId(id)
	ctx.JSON(items)
}

func DeletewarningById(ctx iris.Context) {
	var warning model.Warning
	if err := ctx.ReadJSON(&warning); err != nil {
		ClientErr(ctx, err)
		return
	}
	model.GetDB().UpdateData(warning)
	ServOK(ctx)
}

func DeletewarningByList(ctx iris.Context) {
	var list []model.Warning
	if err := ctx.ReadJSON(&list); err != nil {
		ClientErr(ctx, err)
		return
	}
	for _, warning := range list {
		model.GetDB().UpdateData(&warning)
	}
	ServOK(ctx)
}

func GetPaymentAddressBalance(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	//TODO:
	count := model.GetDB().QueryBusWarningCount(id)
	ctx.JSON(count)
}

func GetBalanceListByTime(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	//TODO:
	items := model.GetDB().QueryBalance(id)
	ctx.JSON(items)
}
func GetHistoryWarning(ctx iris.Context) {
	var page model.Page
	if err := ctx.ReadJSON(&page); err != nil {
		ClientErr(ctx, err)
		return
	}
	items := model.GetDB().QueryHistoryWarning(&page)
	ctx.JSON(items)
}
