package api

import (
	"github.com/zaid13/simplebank/token"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	testCases:=[]struct{
		name string
		setupAuth func( t *testing.T, request http.Request,tokenmaker token.Maker  )
		checkResponse func( t *testing.T, request httptest.ResponseRecorder  )

	}{}
	for  i :=range  testCases{
		tc:=testCases[i]
		t.Run(tc.name, func(t *testing.T) {


		})


	}

}
