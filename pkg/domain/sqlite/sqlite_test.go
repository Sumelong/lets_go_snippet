package sqlite

import (
	"database/sql"
	"snippetbox/pkg/logger"
	"snippetbox/storing/store"
	"testing"
)

func TestSnippetModel_Insert(t *testing.T) {
	/*err := godotenv.Load() // Load variables from .env file
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	*/
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
				expires: "5",
			},
			want:    5,
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
