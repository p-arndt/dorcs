package httpx

import (
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"dorcs-v2/internal/site"
)

// Handlers provides net/http endpoints backed by a site.Site.
//
// Routes (expected):
//   - GET /                -> index
//   - GET /doc/<key>       -> document renderer (extensionless; <key> matches site.Doc.Key)
type Handlers struct {
	Site      *site.Site
	SiteTitle string

	// DocBasePath is the URL path prefix under which docs are served.
	// Must start and end with "/" (default "/doc/").
	DocBasePath string

	// When false, draft documents are hidden from the index.
	IncludeDrafts bool

	templates *template.Template
}

// IndexDoc is the shape used by the index template.
type IndexDoc struct {
	Title       string
	Key         string
	URLPath     string
	RelPath     string
	Description string
	Date        string
	Tags        []string
	Draft       bool
	UpdatedAt   string
}

// IndexViewModel is passed to the index template.
type IndexViewModel struct {
	SiteTitle string
	DocsDir   string
	Docs      []IndexDoc
	Generated string
}

// DocViewModel is passed to the doc template.
type DocViewModel struct {
	SiteTitle string
	Doc       *site.Doc
	HTML      template.HTML
	TOCHTML   template.HTML
	// CanonicalURL can be optionally set by the caller; left blank by default.
	CanonicalURL string
	// Convenience for older templates; doc page can derive title from Doc.
	DocTitleFallback string
	// LastModified is a formatted timestamp for optional display.
	LastModified string
}

// New creates a Handlers instance with defaults.
// You must call LoadTemplates(...) before serving requests.
func New(s *site.Site) *Handlers {
	return &Handlers{
		Site:          s,
		SiteTitle:     "Markdown Docs",
		DocBasePath:   "/doc/",
		IncludeDrafts: false,
	}
}

// LoadTemplates parses templates from an embedded FS.
//
// It expects template files to define at least:
//   - "index" or "layout"+"content blocks" (see your templates)
//   - "doc" or "doc.html" depending on how you structure them
//
// To keep this flexible, the handlers will ExecuteTemplate using "index" and "doc" first,
// and fall back to "index.html"/"doc.html" if those are the defined names.
func (h *Handlers) LoadTemplates(fsys embed.FS, glob string) error {
	t, err := template.New("root").Funcs(template.FuncMap{
		"formatDate": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02")
		},
	}).ParseFS(fsys, glob)
	if err != nil {
		return err
	}
	h.templates = t
	return nil
}

// Register registers the handlers on the provided ServeMux.
func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.Index)
	mux.HandleFunc("GET "+h.DocBasePath+"{rest...}", h.Doc)
	// Also handle exact base path (e.g. /doc/) so we can redirect to index or 404 consistently.
	mux.HandleFunc("GET "+strings.TrimSuffix(h.DocBasePath, "/"), h.Doc)
}

// Index renders the index page.
func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	if h.Site == nil {
		http.Error(w, "site not configured", http.StatusInternalServerError)
		return
	}
	if h.templates == nil {
		http.Error(w, "templates not configured", http.StatusInternalServerError)
		return
	}

	docs := h.Site.ListDocs(h.IncludeDrafts)
	items := make([]IndexDoc, 0, len(docs))
	for _, d := range docs {
		items = append(items, IndexDoc{
			Title:       d.Title,
			Key:         d.Key,
			URLPath:     h.DocBasePath + escapePathSegments(d.Key),
			RelPath:     d.RelPath,
			Description: d.Description,
			Date:        formatDate(d.Date),
			Tags:        append([]string(nil), d.Tags...),
			Draft:       d.Draft,
			UpdatedAt:   d.UpdatedAt.Format(time.RFC3339),
		})
	}

	vm := IndexViewModel{
		SiteTitle: h.SiteTitle,
		DocsDir:   h.Site.RootDir,
		Docs:      items,
		Generated: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Try common names; keep compatibility with either "index" or "index.html".
	if err := h.templates.ExecuteTemplate(w, "index", vm); err == nil {
		return
	}
	if err := h.templates.ExecuteTemplate(w, "index.html", vm); err == nil {
		return
	}

	// If you get here, you likely don't have matching template names.
	http.Error(w, "index template not found or failed to execute", http.StatusInternalServerError)
}

// Doc renders a document page.
// It expects extensionless URL paths like /doc/guide/getting-started.
func (h *Handlers) Doc(w http.ResponseWriter, r *http.Request) {
	if h.Site == nil {
		http.Error(w, "site not configured", http.StatusInternalServerError)
		return
	}
	if h.templates == nil {
		http.Error(w, "templates not configured", http.StatusInternalServerError)
		return
	}

	key, ok := h.extractKey(r.URL)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if key == "" {
		// Someone hit /doc or /doc/; redirect to index.
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	rendered, err := h.Site.RenderDoc(key)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to render document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Optionally hide drafts even if addressable (your call). Here: allow direct access, only hide from index.
	// If you want to block direct access, uncomment below:
	// if rendered.Doc.Draft && !h.IncludeDrafts {
	// 	http.NotFound(w, r)
	// 	return
	// }

	vm := DocViewModel{
		SiteTitle:        h.SiteTitle,
		Doc:              rendered.Doc,
		HTML:             rendered.HTML,
		TOCHTML:          rendered.TocHTML,
		DocTitleFallback: fallbackTitleFromKey(rendered.Doc.Key),
		LastModified:     rendered.Doc.UpdatedAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Try common names; keep compatibility with either "doc" or "doc.html".
	if err := h.templates.ExecuteTemplate(w, "doc", vm); err == nil {
		return
	}
	if err := h.templates.ExecuteTemplate(w, "doc.html", vm); err == nil {
		return
	}

	http.Error(w, "doc template not found or failed to execute", http.StatusInternalServerError)
}

func (h *Handlers) extractKey(u *url.URL) (string, bool) {
	p := u.Path

	// Allow /doc (no trailing slash) as well.
	baseNoSlash := strings.TrimSuffix(h.DocBasePath, "/")
	if p == baseNoSlash || p == h.DocBasePath {
		return "", true
	}

	if !strings.HasPrefix(p, h.DocBasePath) {
		// If we were mounted under /doc/ but got /docsomething, it's not ours.
		return "", false
	}

	rest := strings.TrimPrefix(p, h.DocBasePath)
	rest = strings.Trim(rest, "/")
	rest, _ = url.PathUnescape(rest)

	// Normalize and reject traversal/oddities.
	rest = cleanKey(rest)
	if rest == "" {
		return "", true
	}
	return rest, true
}

func cleanKey(k string) string {
	k = strings.ReplaceAll(k, "\\", "/")
	k = path.Clean("/" + k)
	k = strings.TrimPrefix(k, "/")
	k = strings.Trim(k, "/")

	if k == "" {
		return ""
	}
	if strings.Contains(k, "..") || strings.HasPrefix(k, ".") {
		return ""
	}
	return k
}

func escapePathSegments(pth string) string {
	pth = strings.ReplaceAll(pth, "\\", "/")
	parts := strings.Split(pth, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/")
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func fallbackTitleFromKey(key string) string {
	key = strings.Trim(key, "/")
	if key == "" {
		return "Document"
	}
	base := path.Base(key)
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	return strings.TrimSpace(base)
}
