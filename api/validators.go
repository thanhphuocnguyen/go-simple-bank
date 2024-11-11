package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/thanhphuocnguyen/go-simple-bank/utils"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, err := fieldLevel.Field().Interface().(string); err {
		return utils.IsSupportedCurrency(currency)
	}
	return false
}
