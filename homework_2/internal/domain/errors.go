package domain

import "errors"

var (
	ErrorProductNotFound   = errors.New("Информацию можно получить только по существующему товару")
	ErrorProductInvalidSKU = errors.New("Идентификатор товара должен быть натуральным числом (больше нуля)")
	ErrorInternal          = errors.New("Проблемы из-за неисправностей в системе")
)
