package templates

const TemplateDirectory = "templates"

type DocHtml struct {
	Squad   string
	Projeto string
	Versao  string
	DocUrl  string
}

type HomeHtml struct {
	Docs   []DocHtml
	Squads []string
}
