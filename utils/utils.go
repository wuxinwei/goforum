package utils

import (
	"errors"

	"github.com/fpay/gopress"
	"github.com/wuxinwei/goforum/services"
)

// GetDB get a db service reference
func GetDB(app *gopress.App) (*services.DbService, error) {
	dbRaw := app.Services.Get("db")
	if dbRaw == nil {
		return nil, errors.New("no db service")
	}
	db, ok := dbRaw.(*services.DbService)
	if !ok {
		return nil, errors.New("db service is invalid format")
	}
	return db, nil
}

// GetCache get a cache service reference
func GetCache(app *gopress.App) (*services.DbService, error) {
	cacheRaw := app.Services.Get("cache")
	if cacheRaw == nil {
		return nil, errors.New("no cache service")
	}
	cache, ok := cacheRaw.(*services.DbService)
	if !ok {
		return nil, errors.New("cache service is invalid format")
	}
	return cache, nil
}
