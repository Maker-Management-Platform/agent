package config

import (
	"log/slog"
	"path/filepath"

	"github.com/eduardooliveira/stLib/v2/utils"
	"github.com/spf13/viper"
)

var Cfg *Config
var v *viper.Viper
var configFile string

func Load(dataFolder string) error {

	configFile = filepath.Join(dataFolder, "config.toml")
	if err := utils.CreateFileIfNotExist(configFile); err != nil {
		return err
	}

	v = viper.New()
	v.SetConfigFile(configFile)
	defaults()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Error("Config file not found")
		} else {
			return err
		}
	}

	v.Set("core.dataFolder", dataFolder)
	if err := v.Unmarshal(&Cfg); err != nil {
		return err
	}

	if err := v.WriteConfigAs(configFile); err != nil {
		return err
	}

	slog.Info("Config loaded", slog.Any("Config", Cfg))
	return nil
}

func defaults() {
	v.SetDefault("core.log.enableFile", false)
	v.SetDefault("core.log.path", "log.log")
	v.SetDefault("server.port", 8000)
	if utils.IsDocker {
		v.SetDefault("library.fileSystems", []map[string]string{
			{"name": "Main Library", "path": "/library", "kind": "local"},
		})
	} else {
		v.SetDefault("library.paths", []string{})
		v.SetDefault("library.fileSystems", []map[string]string{
			{"name": "Main Library", "path": "change_me", "kind": "local"},
		})
	}
	v.SetDefault("library.blacklist", []string{".git", ".svn", ".hg", ".bzr", ".DS_Store", ".project.stlib", ".thumb.png", ".render.png"})
	v.SetDefault("library.renderBundles", false)
	v.SetDefault("library.extractImages", false)
	v.SetDefault("library.assetTypes", map[string]AssetType{
		"model": {
			Name:       "model",
			Label:      "Models",
			Extensions: []string{".stl", ".3mf"},
			Order:      0,
		},
		"image": {
			Name:       "image",
			Label:      "Images",
			Extensions: []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp"},
			Order:      1,
		},
		"slice": {
			Name:       "slice",
			Label:      "Slices",
			Extensions: []string{".gcode"},
			Order:      2,
		},
		"source": {
			Name:       "source",
			Label:      "Sources",
			Extensions: []string{".stp", ".step", ".ste", ".fbx", ".f3d", ".f3z", ".iam", ".ipt", ".scad"},
			Order:      99,
		},
	})
	v.SetDefault("library.ignoreDotFiles", true)
	v.SetDefault("render.maxWorkers", 4)
	v.SetDefault("render.modelColor", "#167DF0")
	v.SetDefault("render.backgroundColor", "#FFFFFF")
	v.SetDefault("integrations.thingiverse.token", "")
}
