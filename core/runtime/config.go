package runtime

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

type Config struct {
	Core struct {
		Log struct {
			EnableFile bool   `json:"enable_file" mapstructure:"enable_file"`
			Path       string `json:"path" mapstructure:"path"`
		} `json:"log" mapstructure:"log"`
	} `json:"core" mapstructure:"core"`
	Server struct {
		Port int `json:"port" mapstructure:"port"`
	} `json:"server" mapstructure:"server"`
	Library struct {
		Path           string   `json:"path" mapstructure:"path"`
		Blacklist      []string `json:"blacklist" mapstructure:"blacklist"`
		IgnoreDotFiles bool     `json:"ignore_dot_files" mapstructure:"ignore_dot_files"`
	} `json:"library" mapstructure:"library"`
	Render struct {
		MaxWorkers      int    `json:"max_workers" mapstructure:"max_workers"`
		ModelColor      string `json:"model_color" mapstructure:"model_color"`
		BackgroundColor string `json:"background_color" mapstructure:"background_color"`
	} `json:"render" mapstructure:"render"`
	Integrations struct {
		Thingiverse struct {
			Token string `json:"token" mapstructure:"token"`
		} `json:"thingiverse" mapstructure:"thingiverse"`
	} `json:"integrations" mapstructure:"integrations"`
}

var Cfg *Config

var dataPath = "/data"

func init() {
	viper.BindEnv("DATA_PATH")
	if viper.GetString("DATA_PATH") != "" {
		dataPath = viper.GetString("DATA_PATH")
	}
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		err := os.MkdirAll(dataPath, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
	}

	bindEnv()

	if v := viper.GetInt("PORT"); v == 0 {
		viper.SetDefault("server.port", 8000)
	} else {
		viper.SetDefault("server.port", v)
	}
	if v := viper.GetString("LIBRARY_PATH"); v == "" {
		viper.SetDefault("library.path", "/library")
	} else {
		viper.SetDefault("library.path", v)
	}
	if v := viper.GetString("MODEL_RENDER_COLOR"); v == "" {
		viper.SetDefault("render.model_color", "#167DF0")
	} else {
		viper.SetDefault("render.model_color", v)
	}
	if v := viper.GetString("MODEL_BACKGROUND_COLOR"); v == "" {
		viper.SetDefault("render.background_color", "#FFFFFF")
	} else {
		viper.SetDefault("render.background_color", v)
	}

	viper.SetDefault("library.blacklist", []string{})
	viper.SetDefault("library.ignore_dot_files", true)
	viper.SetDefault("render.max_workers", 5)
	viper.SetDefault("core.log.enable_file", false)

	viper.SetDefault("server.hostname", "localhost")

	viper.SetConfigName("config")
	viper.AddConfigPath(dataPath)
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	cfg := &Config{}
	viper.Unmarshal(cfg)

	cfg.Library.Blacklist = append(cfg.Library.Blacklist, ".project.stlib", ".thumb.png", ".render.png")

	if _, err := os.Stat(path.Join(dataPath, "config.toml")); os.IsNotExist(err) {
		log.Println("config.toml not found, creating...")
		SaveConfig(cfg)
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

func GetDataPath() string {
	return dataPath
}

func SaveConfig(cfg *Config) error {
	f, err := os.OpenFile(filepath.Join(GetDataPath(), "config.toml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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
