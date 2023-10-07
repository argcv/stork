package config

import (
	"time"
)

/**
 *  Database is the default database name used when the Session.DB method
 *  is called with an empty name, and is also used during the initial
 *  authentication if Source is unset.
 *
 *  Username and Password inform the credentials for the initial authentication
 *  done on the database defined by the Source field. See Session.Login.
 *
 */
type MongoAuth struct {
	Source    string
	Username  string
	Password  string
	Mechanism string
}

type MongoConfig struct {
	Addrs     [] string
	Timeout   time.Duration
	DefaultDb string
	Auth      *MongoAuth
}

func LoadMongoConfig(profile string) (cfg *MongoConfig) {
	ckb := NewKeyBuilder(KeyMongoBase).WithProfile(profile)
	var defaultAddrs []string
	cfg = &MongoConfig{
		Addrs:     ckb.GetStringSliceOrDefault(KeyMongoAddrs, defaultAddrs),
		Timeout:   time.Duration(ckb.GetInt64OrDefault(KeyMongoTimeoutSec, 0)) * time.Second,
		DefaultDb: ckb.GetStringOrDefault(KeyMongoDatabase, ""),
		Auth:      GetDBMongoAuth(ckb),
	}

	return
}

func GetDBMongoAuth(ckb *KeyBuilder) *MongoAuth {
	if ckb.GetBoolOrDefault(KeyMongoPerformAuth, false) {
		auth := ckb.Clone().WithClass(KeyMongoAuthBase)
		source := auth.GetStringOrDefault(KeyMongoAuthDatabase, "")
		user := auth.GetStringOrDefault(KeyMongoAuthUser, "admin")
		pass := auth.GetStringOrDefault(KeyMongoAuthPass, "")
		mech := auth.GetStringOrDefault(KeyMongoAuthMechanism, "")

		return &MongoAuth{
			Source:    source,
			Username:  user,
			Password:  pass,
			Mechanism: mech,
		}
	}
	return nil
}
