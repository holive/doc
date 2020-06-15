package docApi

const (
	FilesFolder = "static"
	FileName    = "document.yaml"
	SquadKey    = "X-Squad-Key"
)

var FileTypes = []string{".css", ".js", "json", ".yaml", ".yml"}

type DocApi struct {
	Squad   string `json:"squad" bson:"squad"`
	Projeto string `json:"projeto" bson:"projeto"`
	Versao  string `json:"versao"bson:"versao"`
	Doc     []byte `json:"doc,omitempty"bson:"doc"`
}

type SearchResult struct {
	Docs   []DocApi `json:"docs"`
	Result struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit"`
		Total  int64 `json:"total"`
	} `json:"_result"`
}
