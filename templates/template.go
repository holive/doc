package templates

const TemplateDirectory = "templates"

type HomeHtml struct {
	Docs    []DocHtml
	Squads  []string
	Results struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit"`
		Total  int64 `json:"total"`
	} `json:"_result"`
}

type DocHtml struct {
	Squad     string
	Projeto   string
	Versao    string
	Descricao string
	DocUrl    string
}
