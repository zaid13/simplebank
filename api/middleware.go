package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zaid13/simplebank/token"
	"net/http"
)

const (
	autherizationHeaderKey = "authorization"
)

func authMiddleware(token token.Maker) gin.HandlerFunc {
	return func(context *gin.Context) {
		autherizationHeader := context.GetHeader(autherizationHeaderKey)
		if len(autherizationHeader) == 0 {
			err := errors.New("Authorization key is not provided")
			context.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}
		//feilds:
	}
}
