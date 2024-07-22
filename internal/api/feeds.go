package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) Feeds(ctx *gin.Context) {
	limitStr := ctx.Param("limit")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	if limit == 0 || limit < 0 {
		limit = 10
	}

	feeds, err := api.repository.Feeds.Feeds(limit)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.Header("Access-Control-Allow-Origin", "*")

	ctx.JSON(http.StatusOK, feeds)
}
