package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/zaid13/simplebank/db/sqlc"
	"github.com/zaid13/simplebank/db/util"
	"net/http"
	"time"
)

type createUserRequest struct {
	Username     string `json:"username" binding:"required,alphanum"`
	HashPassword string `json:"password" binding:"required,min=6"`
	FullName     string `json:"full_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username string `json:"username"`

	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {

	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hash, err := util.HashPassword(req.HashPassword)
	arg := db.CreateUsersParams{
		req.FullName,
		hash,
		req.FullName,
		req.Email,
		time.Now(),
	}
	user, err := server.store.CreateUsers(ctx, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {

			switch pqError.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)

}

type logininUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}
type logininUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req logininUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUsers(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

	}
	err = util.CheckPassword(req.Password, user.HashPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := logininUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
