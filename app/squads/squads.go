package squads

type Squad struct {
	Name string `json:"name" bson:"name"`
	Key  string `json:"key" bson:"key"`
}
