package store

import (
	"database/sql"
	"reflect"
	"snippetbox/pkg/logger"
	"testing"
)

func TestNewStorePostgres(t *testing.T) {

	lg := logger.NewLogger()

	type args struct {
		lg logger.Logger
	}
	tests := []struct {
		name string
		args args
		want *sql.DB
	}{
		{
			name: "correct pg",
			args: args{lg: lg},
			want: NewStorePostgres(lg),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStorePostgres(tt.args.lg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStorePostgres() = %v, want %v", got, tt.want)
			}
		})
	}
}
