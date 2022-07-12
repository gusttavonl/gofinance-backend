package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/GustavoNoronha0/gofinance-backend/db/sqlc"
	"github.com/GustavoNoronha0/gofinance-backend/util"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	UserID      int32     `json:"user_id" binding:"required"`
	CategoryID  int32     `json:"category_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Value       int32     `json:"value" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}
	var req createAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var categoryId = req.CategoryID
	var accountType = req.Type

	category, err := server.store.GetCategory(ctx, categoryId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
	}

	var categoryTypeIsDifferentOfAccountType = category.Type != accountType
	if categoryTypeIsDifferentOfAccountType {
		ctx.JSON(http.StatusBadRequest, "Account type is different of Category type")
	} else {
		arg := db.CreateAccountParams{
			UserID:      req.UserID,
			CategoryID:  categoryId,
			Title:       req.Title,
			Type:        accountType,
			Description: req.Description,
			Value:       req.Value,
			Date:        req.Date,
		}

		account, err := server.store.CreateAccount(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		ctx.JSON(http.StatusOK, account)
	}
}

type getAccountRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}
	var req getAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountGraphRequest struct {
	UserID int32  `uri:"user_id" binding:"required"`
	Type   string `uri:"type" binding:"required"`
}

func (server *Server) getAccountGraph(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}
	var req getAccountGraphRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.GetAccountsGraphParams{
		UserID: req.UserID,
		Type:   req.Type,
	}

	countGraph, err := server.store.GetAccountsGraph(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, countGraph)
}

type getAccountReportsRequest struct {
	UserID int32  `uri:"user_id" binding:"required"`
	Type   string `uri:"type" binding:"required"`
}

func (server *Server) getAccountReports(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}
	var req getAccountReportsRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.GetAccountsReportsParams{
		UserID: req.UserID,
		Type:   req.Type,
	}

	sumReports, err := server.store.GetAccountsReports(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sumReports)
}

type deleteAccountRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}
	var req deleteAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, true)
}

type updateAccountRequest struct {
	ID          int32  `json:"id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       int32  `json:"value"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}
	var req updateAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.UpdateAccountParams{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Value:       req.Value,
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountsRequest struct {
	UserID      int32     `json:"user_id" binding:"required"`
	Type        string    `json:"type" binding:"required"`
	CategoryID  int32     `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func (server *Server) getAccounts(ctx *gin.Context) {
	errOnValiteToken := util.GetTokenInHeaderAndVerify(ctx)
	if errOnValiteToken != nil {
		return
	}
	var req getAccountsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	var accounts interface{}
	var paremetersHasUserIdAndType = req.UserID > 0 && len(req.Type) > 0

	filterAsByUserIdAndType := req.CategoryID == 0 && req.Date.IsZero() && len(req.Description) == 0 && len(req.Title) == 0 && paremetersHasUserIdAndType
	if filterAsByUserIdAndType {
		arg := db.GetAccountsByUserIdAndTypeParams{
			UserID: req.UserID,
			Type:   req.Type,
		}

		accountsByUserIdAndType, err := server.store.GetAccountsByUserIdAndType(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		accounts = accountsByUserIdAndType
	}

	filterAsByUserIdAndTypeAndCategoryId := req.CategoryID != 0 && req.Date.IsZero() && len(req.Description) == 0 && len(req.Title) == 0 && paremetersHasUserIdAndType
	if filterAsByUserIdAndTypeAndCategoryId {
		arg := db.GetAccountsByUserIdAndTypeAndCategoryIdParams{
			UserID:     req.UserID,
			Type:       req.Type,
			CategoryID: req.CategoryID,
		}

		accountsByUserIdAndTypeAndCategoryId, err := server.store.GetAccountsByUserIdAndTypeAndCategoryId(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		accounts = accountsByUserIdAndTypeAndCategoryId
	}

	filterAsByUserIdAndTypeAndCategoryIdAndTitle := req.CategoryID != 0 && req.Date.IsZero() && len(req.Description) == 0 && len(req.Title) > 0 && paremetersHasUserIdAndType
	if filterAsByUserIdAndTypeAndCategoryIdAndTitle {
		arg := db.GetAccountsByUserIdAndTypeAndCategoryIdAndTitleParams{
			UserID:     req.UserID,
			Type:       req.Type,
			CategoryID: req.CategoryID,
			Title:      req.Title,
		}

		accountsByUserIdAndTypeAndCategoryIdAndTitle, err := server.store.GetAccountsByUserIdAndTypeAndCategoryIdAndTitle(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		accounts = accountsByUserIdAndTypeAndCategoryIdAndTitle
	}

	filterAsByUserIdAndTypeAndCategoryIdAndTitleAndDescription := req.CategoryID != 0 && req.Date.IsZero() && len(req.Description) > 0 && len(req.Title) > 0 && paremetersHasUserIdAndType
	if filterAsByUserIdAndTypeAndCategoryIdAndTitleAndDescription {
		arg := db.GetAccountsByUserIdAndTypeAndCategoryIdAndTitleAndDescriptionParams{
			UserID:      req.UserID,
			Type:        req.Type,
			CategoryID:  req.CategoryID,
			Title:       req.Title,
			Description: req.Description,
		}

		accountsByUserIdAndTypeAndCategoryIdAndTitle, err := server.store.GetAccountsByUserIdAndTypeAndCategoryIdAndTitleAndDescription(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		accounts = accountsByUserIdAndTypeAndCategoryIdAndTitle
	}

	filterAsByUserIdAndTypeAndDate := req.CategoryID == 0 && !req.Date.IsZero() && len(req.Description) == 0 && len(req.Title) == 0 && paremetersHasUserIdAndType
	if filterAsByUserIdAndTypeAndDate {
		arg := db.GetAccountsByUserIdAndTypeAndDateParams{
			UserID: req.UserID,
			Type:   req.Type,
			Date:   req.Date,
		}

		accountsByUserIdAndTypeAndDate, err := server.store.GetAccountsByUserIdAndTypeAndDate(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		accounts = accountsByUserIdAndTypeAndDate
	}

	filterAsByUserIdAndTypeAndDescription := req.CategoryID == 0 && req.Date.IsZero() && len(req.Description) > 0 && len(req.Title) == 0 && paremetersHasUserIdAndType
	if filterAsByUserIdAndTypeAndDescription {
		arg := db.GetAccountsByUserIdAndTypeAndDescriptionParams{
			UserID:      req.UserID,
			Type:        req.Type,
			Description: req.Description,
		}

		accountsByUserIdAndTypeAndDescription, err := server.store.GetAccountsByUserIdAndTypeAndDescription(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		accounts = accountsByUserIdAndTypeAndDescription
	}

	filterAsByUserIdAndTypeAndTitle := req.CategoryID == 0 && req.Date.IsZero() && len(req.Description) == 0 && len(req.Title) > 0 && paremetersHasUserIdAndType
	if filterAsByUserIdAndTypeAndTitle {
		arg := db.GetAccountsByUserIdAndTypeAndTitleParams{
			UserID: req.UserID,
			Type:   req.Type,
			Title:  req.Title,
		}

		accountsByUserIdAndTypeAndTitle, err := server.store.GetAccountsByUserIdAndTypeAndTitle(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		accounts = accountsByUserIdAndTypeAndTitle
	}

	filterAsAllParameters := req.CategoryID > 0 && !req.Date.IsZero() && len(req.Description) > 0 && len(req.Title) > 0 && paremetersHasUserIdAndType
	if filterAsAllParameters {
		arg := db.GetAccountsParams{
			UserID:      req.UserID,
			Type:        req.Type,
			Title:       req.Title,
			CategoryID:  req.CategoryID,
			Description: req.Description,
			Date:        req.Date,
		}

		accountsFilterAsAllParameters, err := server.store.GetAccounts(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		accounts = accountsFilterAsAllParameters
	}

	ctx.JSON(http.StatusOK, accounts)
}
