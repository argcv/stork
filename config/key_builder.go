package config

import (
	"github.com/spf13/viper"
	"strings"
)

type KeyBuilder struct {
	Base    string
	Profile string
	Class   []string
	Src     *viper.Viper
}

func (ckb *KeyBuilder) Join(elems ...string) string {
	var neElems []string
	// remove empty elements
	for _, e := range elems {
		if len(e) > 0 {
			neElems = append(neElems, e)
		}
	}
	return strings.Join(neElems, ".")
}

func (ckb *KeyBuilder) Current() string {
	return ckb.GetKey("")
}

func (ckb *KeyBuilder) WithProfile(p string) *KeyBuilder {
	ckb.Profile = p
	return ckb
}

func (ckb *KeyBuilder) WithClass(c string) *KeyBuilder {
	ckb.Class = append(ckb.Class, c)
	return ckb
}

func (ckb *KeyBuilder) Clone() *KeyBuilder {
	class := make([]string, len(ckb.Class))
	copy(class, ckb.Class)
	newKb := &KeyBuilder{
		Base:    ckb.Base,
		Profile: ckb.Profile,
		Class:   class,
		Src:     ckb.Src,
	}
	return newKb
}

func (ckb *KeyBuilder) SubViper(key string) (*viper.Viper) {
	return ckb.Src.Sub(ckb.GetKey(key))
}

func (ckb *KeyBuilder) GetKey(key string) string {
	keys := []string{}
	keys = append(keys, ckb.Base, ckb.Profile)
	keys = append(keys, ckb.Class...)
	keys = append(keys, key)
	return ckb.Join(keys...)
}

func (ckb *KeyBuilder) GetStringOrDefault(key string, or string) string {
	return GetStringOrDefault(ckb.GetKey(key), or)
}

func (ckb *KeyBuilder) CheckKeyIsExists(key string) bool {
	return CheckKeyIsExists(ckb.GetKey(key))
}

func (ckb *KeyBuilder) GetIntOrDefault(key string, or int) int {
	return GetIntOrDefault(ckb.GetKey(key), or)
}

func (ckb *KeyBuilder) GetInt64OrDefault(key string, or int64) int64 {
	return GetInt64OrDefault(ckb.GetKey(key), or)
}

func (ckb *KeyBuilder) GetFloat32OrDefault(key string, or float32) float32 {
	return GetFloat32OrDefault(ckb.GetKey(key), or)
}

func (ckb *KeyBuilder) GetFloat64OrDefault(key string, or float64) float64 {
	return GetFloat64OrDefault(ckb.GetKey(key), or)
}

func (ckb *KeyBuilder) GetBoolOrDefault(key string, or bool) bool {
	return GetBoolOrDefault(ckb.GetKey(key), or)
}

func (ckb *KeyBuilder) GetStringSliceOrDefault(key string, or []string) []string {
	return GetStringSliceOrDefault(ckb.GetKey(key), or)
}

func NewKeyBuilder(base string, rest ...string) (ckb *KeyBuilder) {
	ckb = &KeyBuilder{
		Base:    base,
		Profile: KeyDefaultProfile,
		Class:   []string{},
		Src:     viper.GetViper(),
	}
	if len(rest) > 0 {
		ckb.Base = ckb.Join(base, ckb.Join(rest...))
	}
	return
}
