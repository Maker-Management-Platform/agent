package entities

type AssetType struct {
	Name       string   `json:"name" mapstructure:"name"`
	Label      string   `json:"label" mapstructure:"label"`
	Extensions []string `json:"extensions" mapstructure:"extensions"`
	Order      int      `json:"order" mapstructure:"order"`
}
