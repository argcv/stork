// Configuration loader
package config

import (
	"github.com/spf13/viper"
)

func CheckKeyIsExists(key string) bool {
	return viper.Get(key) != nil
}

func GetStringOrDefault(key, or string) string {
	if CheckKeyIsExists(key) {
		return viper.GetString(key)
	} else {
		return or
	}
}

func GetStringSliceOrDefault(key string, or []string) []string {
	if CheckKeyIsExists(key) {
		return viper.GetStringSlice(key)
	} else {
		return or
	}
}

func GetIntOrDefault(key string, or int) int {
	if CheckKeyIsExists(key) {
		return viper.GetInt(key)
	} else {
		return or
	}
}

func GetInt64OrDefault(key string, or int64) int64 {
	if CheckKeyIsExists(key) {
		return viper.GetInt64(key)
	} else {
		return or
	}
}

func GetFloat32OrDefault(key string, or float32) float32 {
	if CheckKeyIsExists(key) {
		return float32(viper.GetFloat64(key))
	} else {
		return or
	}
}

func GetFloat64OrDefault(key string, or float64) float64 {
	if CheckKeyIsExists(key) {
		return viper.GetFloat64(key)
	} else {
		return or
	}
}

func GetBoolOrDefault(key string, or bool) bool {
	if CheckKeyIsExists(key) {
		return viper.GetBool(key)
	} else {
		return or
	}
}

func SetConfig(key string, value interface{}) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}
