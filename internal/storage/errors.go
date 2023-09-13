package storage

import "errors"

// Основные ошибки проекта, используемые для работы с хранилищем.
var (
	ErrNotFound          = errors.New("information not found")
	ErrSubZero           = errors.New("amount now is sub-zero")
	ErrTransactionFailed = errors.New("transaction failed")
)
