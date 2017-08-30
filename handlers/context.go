package handlers

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
)

type AppContext struct {
	echo.Context
	Session *mgo.Session
}

type AppHandler func(AppContext) error

func Error(err error) interface{} {
	return map[string]string{"message": err.Error()}
}
