package server

import (
	"net/http"
	"reflect"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"testing"
)

/*
lg := logger.NewLogger()
	db := store.NewStoreSqlite(lg)
	snippet := sqlite.NewSnippet(db, lg)

	// Create a mock request
	req := httptest.NewRequest(http.MethodPost, "/snippets/create", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json") // Set content type

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()
*/

func TestHandlers_CreateSnippet(t *testing.T) {
	type fields struct {
		lg       logger.Logger
		snippets models.ISnippet
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handlers{
				lg:       tt.fields.lg,
				snippets: tt.fields.snippets,
			}
			h.CreateSnippet(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_HealthChecker(t *testing.T) {
	type fields struct {
		lg       logger.Logger
		snippets models.ISnippet
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handlers{
				lg:       tt.fields.lg,
				snippets: tt.fields.snippets,
			}
			h.HealthChecker(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_Home(t *testing.T) {
	type fields struct {
		lg       logger.Logger
		snippets models.ISnippet
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handlers{
				lg:       tt.fields.lg,
				snippets: tt.fields.snippets,
			}
			h.Home(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_ServeHTTP(t *testing.T) {
	type fields struct {
		lg       logger.Logger
		snippets models.ISnippet
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handlers{
				lg:       tt.fields.lg,
				snippets: tt.fields.snippets,
			}
			h.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_ShowSnippet(t *testing.T) {
	type fields struct {
		lg       logger.Logger
		snippets models.ISnippet
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handlers{
				lg:       tt.fields.lg,
				snippets: tt.fields.snippets,
			}
			h.ShowSnippet(tt.args.w, tt.args.r)
		})
	}
}

func TestNewHandler(t *testing.T) {
	type args struct {
		snippet models.ISnippet
		lg      logger.Logger
	}
	tests := []struct {
		name string
		args args
		want *Handlers
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.snippet, tt.args.lg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
