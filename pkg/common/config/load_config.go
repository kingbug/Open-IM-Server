package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"github.com/spf13/viper"
)

func Load(configDirectory string, configFileName string, envPrefix string, config any) error {
	if runtimeenv.RuntimeEnvironment() == KUBERNETES {
		mountPath := os.Getenv(MountConfigFilePath)
		if mountPath == "" {
			return errs.ErrArgs.WrapMsg(MountConfigFilePath + " env is empty")
		}

		return loadConfig(filepath.Join(mountPath, configFileName), envPrefix, config)
	}

	return loadConfig(filepath.Join(configDirectory, configFileName), envPrefix, config)
}

func loadConfig(path string, envPrefix string, config any) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return errs.WrapMsg(err, "failed to read config file", "path", path, "envPrefix", envPrefix)
	}

	if err := v.Unmarshal(config, func(config *mapstructure.DecoderConfig) {
		config.TagName = StructTagName
	}); err != nil {
		return errs.WrapMsg(err, "failed to unmarshal config", "path", path, "envPrefix", envPrefix)
	}
	return nil
}
