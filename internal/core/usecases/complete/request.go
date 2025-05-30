package complete

type ctxKey string

var ctxRequestKey ctxKey = "request"

type Request struct {
	Question       string
	ProfileSlug    string
	ModelSlug      string
	IsFollowUp     bool
	AppendFilename string
}
