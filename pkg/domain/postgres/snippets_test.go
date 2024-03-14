package postgres

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"reflect"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"snippetbox/storing/store"
	"testing"
)

func TestNewSnippet(t *testing.T) {
	err := godotenv.Load() // Load variables from .env file
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	lg := logger.NewLogger()
	db := store.NewStoreFactory(store.StorageInstancePostgres, lg)

	type args struct {
		db *sql.DB
		lg logger.ILogger
	}
	tests := []struct {
		name string
		args args
		want *SnippetModel
	}{
		{
			name: "valid test",
			args: args{
				db: db,
				lg: lg,
			},
			want: NewSnippet(db, lg),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSnippet(tt.args.db, tt.args.lg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSnippet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnippetModel_Get(t *testing.T) {
	_ = godotenv.Load()
	lg := logger.NewLogger()
	db := store.NewStoreFactory(store.StorageInstancePostgres, lg)

	type fields struct {
		DB *sql.DB
		lg logger.ILogger
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Snippet
		wantErr bool
	}{
		{
			name: "Get",
			fields: fields{
				DB: db,
				lg: lg,
			},
			args:    args{id: 3},
			want:    &models.Snippet{ID: 3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SnippetModel{
				DB: tt.fields.DB,
				lg: tt.fields.lg,
			}
			got, err := m.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnippetModel_Insert(t *testing.T) {
	type fields struct {
		DB *sql.DB
		lg logger.ILogger
	}
	type args struct {
		title   string
		content string
		expires string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SnippetModel{
				DB: tt.fields.DB,
				lg: tt.fields.lg,
			}
			got, err := m.Insert(tt.args.title, tt.args.content, tt.args.expires)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Insert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnippetModel_Latest(t *testing.T) {
	type fields struct {
		DB *sql.DB
		lg logger.ILogger
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*models.Snippet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SnippetModel{
				DB: tt.fields.DB,
				lg: tt.fields.lg,
			}
			got, err := m.Latest()
			if (err != nil) != tt.wantErr {
				t.Errorf("Latest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Latest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
