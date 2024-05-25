package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "tutorial.sqlc.dev/app/db/sqlc"
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
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	// if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
	// 	return
	// }

	// if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
	// 	return
	// }

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountID,
		ToAccountId:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, transfer)

}

// func (server *Server) validAccount(ctx *gin.Context, accounID int64, currency string) bool {
// 	acc, err := server.store.GetAccount(ctx, accounID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusBadRequest, errorMessage(err))
// 			return false
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
// 		return false
// 	}
// 	if currency != acc.Currency {
// 		err := fmt.Errorf("account [%d] currency mismatch: %s cs %s ", accounID, acc.Currency, currency)
// 		ctx.JSON(http.StatusBadRequest, errorMessage(err))
// 		return false
// 	}
// 	return true
// }
