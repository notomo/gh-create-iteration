package gqltest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
)

type Transport struct {
	Handler http.HandlerFunc
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	t.Handler(rec, req)
	return rec.Result(), nil
}

type Handler struct {
	QueryName  string
	MutateName string

	Handle http.HandlerFunc
}

type Option struct {
	Handlers []Handler
}

func New(t *testing.T, opts ...func(*Option)) (*api.GraphQLClient, error) {
	option := Option{}
	for _, opt := range opts {
		opt(&option)
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Restore the body so registered handlers can read it again (e.g. to capture variables).
		r.Body = io.NopCloser(bytes.NewReader(raw))

		var body struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, h := range option.Handlers {
			if h.QueryName != "" && !strings.HasPrefix(body.Query, "query "+h.QueryName) {
				continue
			}
			if h.MutateName != "" && !strings.HasPrefix(body.Query, "mutation "+h.MutateName) {
				continue
			}
			h.Handle(w, r)
			return
		}
		http.Error(w, "no matched handler for query: "+body.Query, http.StatusNotFound)
	}
	return api.NewGraphQLClient(api.ClientOptions{
		Transport: &Transport{Handler: handler},
	})
}

func WithQueryOK(queryName string, body string) func(*Option) {
	return WithQuery(queryName, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	})
}

func WithMutateOK(mutateName string, body string) func(*Option) {
	return WithMutate(mutateName, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	})
}

// WithQuery registers a full handler for a named query, letting tests inspect the
// request (e.g. captured variables) before responding.
func WithQuery(queryName string, handle http.HandlerFunc) func(*Option) {
	return func(o *Option) {
		o.Handlers = append(o.Handlers, Handler{QueryName: queryName, Handle: handle})
	}
}

// WithMutate registers a full handler for a named mutation, letting tests inspect the
// request (e.g. captured variables) before responding.
func WithMutate(mutateName string, handle http.HandlerFunc) func(*Option) {
	return func(o *Option) {
		o.Handlers = append(o.Handlers, Handler{MutateName: mutateName, Handle: handle})
	}
}
