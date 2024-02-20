package entities

type Tag struct {
	Value    string     `json:"value" gorm:"primaryKey"`
	Projects []*Project `gorm:"many2many:project_tags" json:"-"`
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
