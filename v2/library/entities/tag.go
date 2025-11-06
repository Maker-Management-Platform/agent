package entities

type Tag struct {
	Value  string   `json:"value" gorm:"primaryKey"`
	Assets []*Asset `gorm:"many2many:asset_tags" json:"-"`
}

func StringToTag(s string) *Tag {
	return &Tag{Value: s}
}

func StringsToTags(ss []string) []*Tag {
	rtn := make([]*Tag, 0)
	for _, s := range ss {
		rtn = append(rtn, StringToTag(s))
	}
	return rtn
}
