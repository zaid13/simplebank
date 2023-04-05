package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/zaid13/simplebank/db/util"
)

var ValidCurrency validator.Func = func(flLevel validator.FieldLevel) bool {

	if currency, ok := flLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
		///check if currency is OK
	}
	return false

}
