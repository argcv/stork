package config

import (
	"github.com/pkg/errors"
)

type SMTPConfigOption struct {
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
}

type SMTPConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Sender string `mapstructure:"sender"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	// SkipTLS
	Option SMTPConfigOption `mapstructure:"option"`
}

func (cfg *SMTPConfig) verify() error {
	if cfg == nil {
		return errors.New("cfg is nil")
	}
	if len(cfg.Host) == 0 {
		return errors.New("host is not defined")
	}
	return nil
}

/*
Example Config File:

	mail: # base-0
		smtp: # base-1
			default: # profile
				host: 127.0.0.1
				port: 25
				sender: Alice Ez
				user: "aez@example.com"
				pass: "your-password" # empty for no authorization
				option:
					insecure_skip_verify: true
			nightly: # profile
				host: 127.0.0.1
				port: 25
				sender: Alice Ez
				user: "aez@example.com"
				pass: "your-password" # empty for no authorization
				option:
					insecure_skip_verify: true

 */
func LoadSMTPConfig(profile string) (cfg *SMTPConfig, err error) {
	ckb := NewKeyBuilder(KeyMailBase, KeyMailSMTPBase).WithProfile(profile).SubViper("")

	cfg = &SMTPConfig{}
	if err = ckb.Unmarshal(cfg); err == nil {
		return cfg, cfg.verify()
	}
	return
}
