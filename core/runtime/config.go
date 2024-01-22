package runtime

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

type Config struct {
	Core struct {
		Log struct {
			EnableFile bool   ` json:"enable_file"`
			Path       string ` json:"path"`
		} ` json:"log"`
	} ` json:"core"`
	Server struct {
		Port int ` json:"port"`
	} ` json:"server"`
	Library struct {
		Path           string   ` json:"path"`
		Blacklist      []string ` json:"blacklist"`
		IgnoreDotFiles bool     ` json:"ignore_dot_files"`
	} ` json:"library"`
	Render struct {
		MaxWorkers      int    ` json:"max_workers"`
		ModelColor      string ` json:"model_color"`
		BackgroundColor string ` json:"background_color"`
	} ` json:"render"`
	Integrations struct {
		Thingiverse struct {
			Token string ` json:"token"`
		} ` json:"thingiverse"`
	} ` json:"integrations"`
}

var Cfg *Config

func init() {
	viper.BindEnv("DATA_FOLDER")
	if viper.GetString("DATA_FOLDER") == "" {
		log.Panic("data folder not defined")
	}

	bindEnv()

	viper.SetDefault("server.port", viper.GetInt("PORT"))
	viper.SetDefault("server.hostname", "localhost")
	viper.SetDefault("library.path", viper.GetString("LIBRARY_PATH"))
	viper.SetDefault("library.blacklist", []string{})
	viper.SetDefault("library.ignore_dot_files", true)
	viper.SetDefault("render.max_workers", 5)
	viper.SetDefault("render.model_color", viper.GetString("MODEL_RENDER_COLOR"))
	viper.SetDefault("render.background_color", viper.GetString("MODEL_BACKGROUND_COLOR"))
	viper.SetDefault("integrations.thingiverse.token", viper.GetString("THINGIVERSE_TOKEN"))
	viper.SetDefault("core.log.path", viper.GetString("LOG_PATH"))
	viper.SetDefault("core.log.enable_file", false)

	viper.SetConfigName("config")
	viper.AddConfigPath(viper.GetString("DATA_FOLDER"))
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("error config file: %w", err)
	}
	cfg := &Config{}
	viper.Unmarshal(cfg)

	l, _ := json.Marshal(cfg)
	log.Println(string(l))
	viper.WriteConfigAs("wee.toml")

	cfg.Library.Blacklist = append(cfg.Library.Blacklist, ".project.stlib")

	if cfg.Library.Path == "" {
		log.Panic("library path is empty")
	}
	if cfg.Server.Port == 0 {
		log.Panic("server port is empty")
	}

	Cfg = cfg
}

func bindEnv() {
	viper.BindEnv("PORT")
	viper.BindEnv("LIBRARY_PATH")
	viper.BindEnv("MAX_RENDER_WORKERS")
	viper.BindEnv("MODEL_RENDER_COLOR")
	viper.BindEnv("MODEL_BACKGROUND_COLOR")
	viper.BindEnv("LOG_PATH")
	viper.BindEnv("THINGIVERSE_TOKEN")
}

func GetDataFolder() string {
	return viper.GetString("DATA_FOLDER")
}

func SaveConfig(cfg *Config) error {
	f, err := os.OpenFile(filepath.Join(GetDataFolder(), "config.toml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	if err := toml.NewEncoder(f).Encode(cfg); err != nil {
		log.Println(err)
	}
	if err := f.Close(); err != nil {
		log.Println(err)
	}
	Cfg = cfg
	return err
}
