package db

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

const dbName = "sts"

func InitDB(config *viper.Viper) *mgo.Session {
	DialInfo := &mgo.DialInfo{
		Addrs:    []string{config.GetString("addrs")},
		Timeout:  config.GetDuration("timeout"),
		Database: config.GetString("database"),
		Username: config.GetString("username"),
		Password: config.GetString("password"),
	}

	session, err := mgo.DialWithInfo(DialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	session.SetMode(mgo.Monotonic, true)

	return session
}

func Drop(session *mgo.Session) error {
	return session.DB(dbName).DropDatabase()
}
