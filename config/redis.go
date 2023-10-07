package config

type RedisConfig = struct {
	Host string
	Port int
	Pass string
	Db   int
}

/*
Example Config File:

	redis:
		default: # profile
			host: 127.0.0.1
			port: 6379
			pass: "" # empty for not encrypted
			db: 0 # db number
 */
func LoadRedisConfig(profile string) *RedisConfig {
	ckb := NewKeyBuilder(KeyRedisBase).WithProfile(profile)
	return &RedisConfig{
		Host: ckb.GetStringOrDefault(KeyRedisHost, "127.0.0.1"),
		Port: ckb.GetIntOrDefault(KeyRedisPort, 6379),
		Pass: ckb.GetStringOrDefault(KeyRedisPass, ""),
		Db:   ckb.GetIntOrDefault(KeyRedisDb, 0),
	}
}
