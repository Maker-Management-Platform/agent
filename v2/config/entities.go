package config

type Config struct {
	Core struct {
		Log struct {
			EnableFile bool   `json:"enableFile" mapstructure:"enableFile"`
			Path       string `json:"path" mapstructure:"path"`
		} `json:"log" mapstructure:"log"`
		DataFolder string `json:"dataFolder" mapstructure:"dataFolder"`
	} `json:"core" mapstructure:"core"`
	Server struct {
		Port int `json:"port" mapstructure:"port"`
	} `json:"server" mapstructure:"server"`
	Library struct {
		Paths          []string   `json:"paths" mapstructure:"paths"`
		Blacklist      []string   `json:"blacklist" mapstructure:"blacklist"`
		IgnoreDotFiles bool       `json:"ignoreDotFiles" mapstructure:"ignoreDotFiles"`
		AssetTypes     AssetTypes `json:"assetTypes" mapstructure:"assetTypes"`
	} `json:"library" mapstructure:"library"`
	Render struct {
		MaxWorkers      int    `json:"maxWorkers" mapstructure:"maxWorkers"`
		ModelColor      string `json:"modelColor" mapstructure:"modelColor"`
		BackgroundColor string `json:"backgroundColor" mapstructure:"backgroundColor"`
	} `json:"render" mapstructure:"render"`
	Integrations struct {
		Thingiverse struct {
			Token string `json:"token" mapstructure:"token"`
		} `json:"thingiverse" mapstructure:"thingiverse"`
	} `json:"integrations" mapstructure:"integrations"`
}
type AssetTypes map[string]AssetType

func (a AssetTypes) ByExtension(ext string) AssetType {
	for _, assetType := range a {
		for _, extension := range assetType.Extensions {
			if extension == ext {
				return assetType
			}
		}
	}
	return AssetType{
		Name:       "other",
		Label:      "Other",
		Extensions: []string{ext},
		Order:      100,
	}
}

type AssetType struct {
	Name       string   `json:"name" mapstructure:"name"`
	Label      string   `json:"label" mapstructure:"label"`
	Extensions []string `json:"extensions" mapstructure:"extensions"`
	Order      int      `json:"order" mapstructure:"order"`
}
