package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/zaid13/simplebank/db/sqlc"
	"github.com/zaid13/simplebank/token"
	"net/http"
)

type createAccountRequest struct {

	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload:=ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {

			switch pqError.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context)  {

	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusExpectationFailed, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)
	fmt.Println(req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusExpectationFailed, errorResponse(err))

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload:=ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.Username!=account.Owner{
		err:=errors.New("account doesnt belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {

	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload:=ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	listAccountRequest := db.ListAccountParams{
		authPayload.Username,
		req.PageSize,
		(req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccount(ctx, listAccountRequest)
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
