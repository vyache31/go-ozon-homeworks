package domain

import "errors"

var (
	ErrorProductNotFound   = errors.New("Информацию можно получить только по существующему товару")
	ErrorProductInvalidSKU = errors.New("Идентификатор товара должен быть натуральным числом (больше нуля)")
	ErrorProductInternal   = errors.New("Проблемы из-за неисправностей в системе")

	ErrorUserInvalidID        = errors.New("Идентификатор пользователя должен быть натуральным числом (больше нуля)")
	ErrorPositionInvalidCount = errors.New("Количество должно быть натуральным числом (больше нуля)")
	ErrorCartInvalidSKU       = errors.New("SKU должен быть натуральным числом (больше нуля)")
	ErrorCartNotFoundSKU      = errors.New("SKU должен существовать в сервисе product-service")
	ErrorCartInternal         = errors.New("Проблемы из-за неисправностей в системе")
	ErrorCartNoExistsOrEmpty  = errors.New("Корзины ещё нет или она пуста")
)
