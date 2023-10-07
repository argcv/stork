package config

import (
	"fmt"
	"github.com/argcv/stork/cntr"
	"github.com/argcv/stork/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Option struct {
	Project                  string
	Path                     string
	FileMustExists           bool
	DefaultPath              string
	DefaultConfigureName     string
	ConfigSearchPath         string
	ConfigSearchPaths        []string
	ConfigFallbackSearchPath string
}

func (c *Option) With(rhs *Option) *Option {
	if rhs.Project != "" {
		c.Project = rhs.Project
	}
	if rhs.Path != "" {
		c.Path = rhs.Path
	}
	if rhs.FileMustExists {
		c.FileMustExists = true
	}
	if rhs.DefaultPath != "" {
		c.DefaultPath = rhs.DefaultPath
	}
	if rhs.DefaultConfigureName != "" {
		c.DefaultConfigureName = rhs.DefaultConfigureName
	}
	if len(rhs.ConfigSearchPaths) > 0 {
		c.ConfigSearchPaths = append(c.ConfigSearchPaths, rhs.ConfigSearchPaths...)
	}
	if len(rhs.ConfigSearchPath) > 0 {
		c.ConfigSearchPaths = append(c.ConfigSearchPaths, rhs.ConfigSearchPath)
	}
	if len(rhs.ConfigFallbackSearchPath) > 0 {
		c.ConfigFallbackSearchPath = rhs.ConfigFallbackSearchPath
	}
	return c
}

func (c *Option) GetDefaultPath() string {
	if c.DefaultPath != "" {
		return c.DefaultPath
	} else {
		return c.Path
	}
}

// this function will search and load configurations
func LoadConfig(options ...Option) (err error) {
	option := Option{}

	for _, opt := range options {
		option.With(&opt)
	}

	project := option.Project

	if project == "" {
		return errors.New("Required parameter missing: project")
	}

	viper.SetConfigName(project)
	viper.SetEnvPrefix(project)

	if option.Path != "" {
		viper.SetConfigFile(option.Path)
	} else {
		cfgPaths := option.ConfigSearchPaths
		cfgPaths = append(cfgPaths, option.ConfigSearchPath)

		cfgPaths = append(cfgPaths, ".")  // current folder
		cfgPaths = append(cfgPaths, "..") // parent folder
		cfgPaths = append(cfgPaths, "$HOME/")
		cfgPaths = append(cfgPaths, fmt.Sprintf("$HOME/.%s/", project))
		cfgPaths = append(cfgPaths, "/etc/")
		cfgPaths = append(cfgPaths, fmt.Sprintf("/etc/%s/", project))

		cfgPaths = cntr.DistinctStrings(cfgPaths...)

		for _, cpath := range cfgPaths {
			if len(cpath) > 0 {
				viper.AddConfigPath(cpath)
			}
		}

		if len(option.ConfigFallbackSearchPath) > 0 {
			viper.AddConfigPath(option.ConfigFallbackSearchPath)
		}

		if conf := os.Getenv(fmt.Sprintf("%s_CFG", strings.ToUpper(project))); conf != "" {
			viper.SetConfigFile(conf)
		}
	}

	readAndTestConfig := func() (string, error) {
		err = viper.ReadInConfig()
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && err != nil {
			return "", err
		}
		if conf := viper.ConfigFileUsed(); conf != "" {
			wd, _ := os.Getwd()
			if rel, _ := filepath.Rel(wd, conf); rel != "" && strings.Count(rel, "..") < 3 {
				conf = rel
			}
			log.Infof("using config file: %s", conf)
			return conf, nil
		} else {
			return "", nil
		}
	}

	if conf, e := readAndTestConfig(); conf != "" {
		return nil
	} else if e != nil {
		return e
	} else if option.FileMustExists {
		defaultConfigPath := option.GetDefaultPath()
		defaultConfigName := option.DefaultConfigureName
		if defaultConfigName == "" {
			defaultConfigName = fmt.Sprintf("%s.yml", project)
		}
		if defaultConfigPath == "" {
			msg := "No configure file: default path Not Assigned"
			log.Warnf(msg)
			return errors.New(msg)
		}
		if e := os.MkdirAll(defaultConfigPath, 0700); e != nil {
			return e
		}
		defaultConfigPathFile := path.Join(defaultConfigPath, defaultConfigName)
		log.Infof("configure file NOT found, using default file: %v", defaultConfigPathFile)
		if e := viper.WriteConfigAs(defaultConfigPathFile); e != nil {
			return e
		}
		viper.SetConfigFile(defaultConfigPathFile)

		if conf, e := readAndTestConfig(); conf != "" {
			return nil
		} else if e != nil {
			return e
		} else {
			msg := "No configure file"
			log.Warnf(msg)
			return errors.New(msg)
		}
	} else {
		return nil
	}
}
