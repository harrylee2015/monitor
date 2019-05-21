package webutil

import (
	"github.com/kataras/iris"
	"net/http"
)

func ServErr(ctx iris.Context, err error) {
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.JSON(iris.Map{"error": err.Error()})
}

func ClientErr(ctx iris.Context, err error) {
	ctx.StatusCode(http.StatusBadRequest)
	ctx.JSON(iris.Map{"error": err.Error()})
}

func ServOK(ctx iris.Context) {
	ctx.StatusCode(http.StatusOK)
}
