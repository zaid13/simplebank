package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/zaid13/simplebank/db/sqlc"
	"github.com/zaid13/simplebank/token"
	"net/http"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fromAccount,ok :=server.validAccount(ctx, req.FromAccountID, req.Currency)

	if !ok {
		return
	}
	authPayload:=ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if fromAccount.Owner != authPayload.Username {

		err:=errors.New("from account doent belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_,ok =server.validAccount(ctx, req.ToAccountID, req.Currency)

	if !ok {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)

}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account,bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {

		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account,false

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account,false
	}

	if account.Currency != currency {
		err := fmt.Errorf("Account  %d mismatch %s vs %s ", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account,false
	}

	return account, true

}

//type getAccountRequest struct {
//	ID int64 `uri:"id" binding:"required,min=1"`
//}
//
//func (server *Server) getAccount(ctx *gin.Context) {
//
//	var req getAccountRequest
//	if err := ctx.ShouldBindUri(&req); err != nil {
//		ctx.JSON(http.StatusExpectationFailed, errorResponse(err))
//		return
//	}
//	account, err := server.store.GetAccount(ctx, req.ID)
//	fmt.Println(req.ID)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			ctx.JSON(http.StatusExpectationFailed, errorResponse(err))
//
//		}
//		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//
//	ctx.JSON(http.StatusOK, account)
//}
//
//type listAccountRequest struct {
//	PageID   int32 `form:"page_id" binding:"required,min=1"`
//	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
//}
//
//func (server *Server) listAccount(ctx *gin.Context) {
//
//	var req listAccountRequest
//	if err := ctx.ShouldBindQuery(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	listAccountRequest := db.ListAccountParams{
//		req.PageSize,
//		(req.PageID - 1) * req.PageSize,
//	}
//
//	accounts, err := server.store.ListAccount(ctx, listAccountRequest)
//	if err != nil {
//
//		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//
//	ctx.JSON(http.StatusOK, accounts)
//}
