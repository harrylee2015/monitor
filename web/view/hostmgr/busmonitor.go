package hostmgr

import (
	"github.com/harrylee2015/monitor/common/db"
	"github.com/harrylee2015/monitor/model"
	"github.com/harrylee2015/monitor/types"
	. "github.com/harrylee2015/monitor/web/view/webutil"
	"github.com/kataras/iris"
)

func ListMonitors(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	items := types.GetDB().QueryMonitor(id)
	ctx.JSON(items)
}

func GetBusWaringCount(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	count := types.GetDB().QueryBusWarningCount(id)
	ctx.JSON(count)
}

func GetBusWaringByGroupId(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	items := types.GetDB().QueryWarningByGroupId(id)
	ctx.JSON(items)
}

func DeletewarningById(ctx iris.Context) {
	var warning model.Warning
	if err := ctx.ReadJSON(&warning); err != nil {
		ClientErr(ctx, err)
		return
	}
	types.GetDB().UpdateData(warning)
	ServOK(ctx)
}

func DeletewarningByList(ctx iris.Context) {
	var list []model.Warning
	if err := ctx.ReadJSON(&list); err != nil {
		ClientErr(ctx, err)
		return
	}
	for _, warning := range list {
		types.GetDB().UpdateData(&warning)
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
	items := types.GetDB().QueryLastBalance(id)
	ctx.JSON(items)
}

func GetBalanceListByTime(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	//TODO:
	items := types.GetDB().QueryBalance(id)
	ctx.JSON(items)
}
func GetHistoryWarning(ctx iris.Context) {
	var page model.Page
	if err := ctx.ReadJSON(&page); err != nil {
		ClientErr(ctx, err)
		return
	}
	items := types.GetDB().QueryHistoryWarning(&page)
	ctx.JSON(items)
}

func GetBlockHash(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	items := types.GetDB().QueryWarningByGroupIdAndType(id,db.HASH_WARING)
	if len(items)==0{
		hash :=&model.Hash{
			IsConsistent:true,
		}
		ctx.JSON(hash)
		return
	}
	hash :=&model.Hash{
		IsConsistent:false,
		Values:items,
	}
	ctx.JSON(hash)
}