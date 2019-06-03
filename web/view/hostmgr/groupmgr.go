package hostmgr

import (
	"github.com/harrylee2015/monitor/model"
	"github.com/harrylee2015/monitor/types"
	. "github.com/harrylee2015/monitor/web/view/webutil"
	"github.com/kataras/iris"
)

func ListHostGroups(ctx iris.Context) {
	var page model.Page
	if err := ctx.ReadJSON(&page); err != nil {
		ClientErr(ctx, err)
		return
	}
	items := types.GetDB().QueryGroupInfoByPageNum(&page)
	ctx.JSON(items)
}

func AddHostGroup(ctx iris.Context) {
	var group model.GroupInfo
	if err := ctx.ReadJSON(&group); err != nil {
		ClientErr(ctx, err)
		return
	}
	types.GetDB().InsertData(&group)
	ServOK(ctx)
}

func UpdateHostGroup(ctx iris.Context) {
	var group model.GroupInfo
	if err := ctx.ReadJSON(&group); err != nil {
		ClientErr(ctx, err)
		return
	}
	types.GetDB().UpdateData(&group)
	ServOK(ctx)
}

func DeleteHostGroup(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("groupId")
	if err != nil {
		ClientErr(ctx, err)
		return
	}
	types.GetDB().DeleteDataByGroupId(id)
	ServOK(ctx)
}
