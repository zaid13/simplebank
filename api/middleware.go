package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zaid13/simplebank/token"
	"net/http"
	"strings"
)

const (
	autherizationHeaderKey = "authorization"
	autherizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(token token.Maker) gin.HandlerFunc {
	return func(context *gin.Context) {
		autherizationHeader := context.GetHeader(autherizationHeaderKey)
		//no authorization header passed


		if len(autherizationHeader) == 0 {
			err := errors.New("Authorization key is not provided")
			context.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}
		fields:=strings.Fields(autherizationHeader)
		//invaid token is passed

		if len(fields) <2{
			err:=errors.New("invaid authorization format")
			context.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return

		}

		authorizationTypes:=strings.ToLower(fields[0])


		if authorizationTypes!=autherizationTypeBearer{
			err := fmt.Errorf("Unotherised Authorization %s",authorizationTypes)
			context.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		accessToken:=fields[1]
		payload,err:=token.VerifyToken(accessToken)
		fmt.Println(accessToken )

		if err != nil {
			err := errors.New("Authorization key is not valid")
			context.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}


		context.Set(authorizationPayloadKey,payload)
		context.Next()

	}
}
