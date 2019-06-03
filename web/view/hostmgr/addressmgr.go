package hostmgr

import (
	"github.com/harrylee2015/monitor/model"
	"github.com/harrylee2015/monitor/types"
	. "github.com/harrylee2015/monitor/web/view/webutil"
	"github.com/kataras/iris"
)

func ListAddress(ctx iris.Context) {
	var page model.Page
	if err := ctx.ReadJSON(&page); err != nil {
		ClientErr(ctx, err)
		return
	}
	items := types.GetDB().QueryPaymentAddressByPageNum(&page)
	ctx.JSON(items)
}

func AddAddress(ctx iris.Context) {
	var addr model.PaymentAddress
	if err := ctx.ReadJSON(&addr); err != nil {
		ClientErr(ctx, err)
		return
	}
	types.GetDB().InsertData(&addr)
	ServOK(ctx)
}

func UpdateAddress(ctx iris.Context) {
	var addr model.PaymentAddress
	if err := ctx.ReadJSON(&addr); err != nil {
		ClientErr(ctx, err)
		return
	}
	types.GetDB().UpdateData(&addr)
	ServOK(ctx)
}

func DeleteAddress(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	types.GetDB().DeleteAddressByGroupId(id)
	ServOK(ctx)
}
