package handler

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/holive/doc/templates"

	"github.com/go-chi/chi"
	"github.com/holive/doc/app/docApi"
	"github.com/pkg/errors"
)

func (h *Handler) CreateDoc(w http.ResponseWriter, r *http.Request) {
	doc, err := h.getDocFromRequest(r)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	folderPath := path.Join(docApi.FilesFolder, doc.Squad, doc.Projeto, doc.Versao)

	err = h.receiveFile(r, folderPath)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.Services.DocApi.Create(r.Context(), folderPath, docApi.FileName, doc)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	doc.Doc = nil
	respondWithJSON(w, http.StatusOK, doc)
}

func (h *Handler) GetDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	doc, err := h.getDocFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := h.Services.DocApi.Find(r.Context(), doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	htmlData := templates.DocHtml{
		DocUrl: string(f.Doc),
	}

	tmpl, err := template.ParseFiles(path.Join(templates.TemplateDirectory, "doc.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, htmlData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) ListBySquad(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	squad := chi.URLParam(r, "squad")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	result, err := h.Services.DocApi.FindBySquad(r.Context(), squad, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	htmlData := h.searchResultToTemplate(result)
	newHtmlData, err := h.getAllSquadsToTemplate(r, htmlData, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(path.Join(templates.TemplateDirectory, "home.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, newHtmlData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) SearchByProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	projeto := chi.URLParam(r, "projeto")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	result, err := h.Services.DocApi.SearchProject(r.Context(), projeto, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	htmlData := h.searchResultToTemplate(result)
	newHtmlData, err := h.getAllSquadsToTemplate(r, htmlData, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(path.Join(templates.TemplateDirectory, "search.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, newHtmlData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetAllDocs(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	result, err := h.Services.DocApi.FindAll(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	htmlData := h.searchResultToTemplate(result)
	newHtmlData, err := h.getAllSquadsToTemplate(r, htmlData, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(path.Join(templates.TemplateDirectory, "home.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, newHtmlData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) DeleteDoc(w http.ResponseWriter, r *http.Request) {
	doc, err := h.getDocFromRequest(r)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.Services.DocApi.Delete(r.Context(), doc); err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{})
}

func (h *Handler) getDocFromRequest(r *http.Request) (*docApi.DocApi, error) {
	squad := chi.URLParam(r, "squad")
	projeto := chi.URLParam(r, "projeto")
	versao := chi.URLParam(r, "versao")

	if squad == "" || projeto == "" || versao == "" {
		return nil, errors.New("missing url param")
	}

	if squad == "squad" {
		return nil, errors.New("cannot use \"squad\" as squad name")
	}

	return &docApi.DocApi{
		Squad:   squad,
		Projeto: projeto,
		Versao:  versao,
		Doc:     nil,
	}, nil
}

func (h *Handler) receiveFile(r *http.Request, folderPath string) error {
	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		return errors.Wrap(err, "could not parse multipart form")
	}

	src, _, err := r.FormFile("fileupload")
	if err != nil {
		return errors.Wrap(err, "could not get src from request")
	}
	defer src.Close()

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "could not create the folderPath")
	}

	dst, err := os.Create(path.Join(folderPath, docApi.FileName))
	if err != nil {
		return errors.Wrap(err, "could not create the src")
	}

	if _, err = io.Copy(dst, src); err != nil {
		return errors.Wrap(err, "could not copy src to dst")
	}

	dst.Sync()

	return nil
}

func (h *Handler) searchResultToTemplate(result *docApi.SearchResult) templates.HomeHtml {
	var docUrls []templates.DocHtml

	for _, doc := range result.Docs {
		filePath := path.Join(doc.Squad, doc.Projeto, doc.Versao)

		docUrls = append(docUrls, templates.DocHtml{
			Squad:   doc.Squad,
			Projeto: doc.Projeto,
			Versao:  doc.Versao,
			DocUrl:  filePath,
		})
	}

	return templates.HomeHtml{
		Docs:    docUrls,
		Results: result.Result,
	}
}

func (h *Handler) getAllSquadsToTemplate(r *http.Request,
	homeHtml templates.HomeHtml,
	limit string,
	offset string) (templates.HomeHtml, error) {

	// limit squad list on select filter
	result, err := h.Services.DocApi.FindAll(r.Context(), "30", "0")
	if err != nil {
		return templates.HomeHtml{}, errors.Wrap(err, "error at getAllSquads")
	}

	return templates.HomeHtml{
		Docs:    homeHtml.Docs,
		Squads:  h.reduceSquads(result),
		Results: homeHtml.Results,
	}, nil
}

func (h *Handler) reduceSquads(docs *docApi.SearchResult) []string {
	m := make(map[string]interface{})
	for _, doc := range docs.Docs {
		m[doc.Squad] = nil
	}

	var squads []string
	for key, _ := range m {
		squads = append(squads, key)
	}

	return squads
}
