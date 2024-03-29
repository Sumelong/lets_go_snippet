package sqlite

import (
	"database/sql"
	"reflect"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"snippetbox/storing/store"
	"testing"
)

func TestNewSnippet(t *testing.T) {

	lg := logger.NewLogger()
	db := store.NewStoreFactory(store.StorageInstanceSqlite, lg)

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
			name: "correct db",
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

func TestSnippetModel_Insert(t *testing.T) {

	lg := logger.NewLogger()
	db := store.NewStoreFactory(store.StorageInstanceSqlite, lg)

	type fields struct {
		DB *sql.DB
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
		{
			name:   "valid test",
			fields: fields{DB: db},
			args: args{
				title:   "test-title",
				content: "test-content",
				expires: "6",
			},
			want:    9,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SnippetModel{
				DB: tt.fields.DB,
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

func TestSnippetModel_Get(t *testing.T) {

	lg := logger.NewLogger()
	db := store.NewStoreSqlite(lg)

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
			name: "valid get",
			fields: fields{
				DB: db,
				lg: lg,
			},
			args: args{id: 1},
			want: &models.Snippet{
				ID: 1,
			},
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
			if !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("Get() got = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

func TestSnippetModel_Latest(t *testing.T) {
	lg := logger.NewLogger()
	db := store.NewStoreSqlite(lg)

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
		{
			name: "valid rows",
			fields: fields{
				DB: db,
				lg: lg,
			},
			want:    make([]*models.Snippet, 0),
			wantErr: false,
		},
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
