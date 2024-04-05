package sqlite

import (
	"errors"
	"io"
	"log"
	"reflect"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"snippetbox/storing/migration"
	"testing"
	"time"
)

func TestUserRepository_ReadOne(t *testing.T) {

	// Skip the test if the `-short` flag is provided when running the test.
	if testing.Short() {
		t.Skip("Sqlite_UserRepository: skipping integration test")
	}

	var lg logger.ILogger = logger.StdLogger{
		ErrLog:  log.New(io.Discard, "", 0),
		InfoLog: log.New(io.Discard, "", 0),
	}
	tests := []struct {
		name    string
		userID  int
		want    *models.User
		wantErr error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want: &models.User{
				ID:      1,
				Name:    "Alice Jones",
				Email:   "alice@example.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active:  true,
			},
			wantErr: nil,
		},
		{
			name:    "Zero ID",
			userID:  0,
			want:    nil,
			wantErr: models.ErrNoRecord,
		},
		{
			name:    "Non-existent ID",
			userID:  2,
			want:    nil,
			wantErr: models.ErrNoRecord,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Initialize a connection pool to our test database, and defer a
			// call to the teardown function, so it is always run immediately
			// before this sub-test returns.
			db, teardown := migration.NewTestSqliteDB(t)
			defer teardown()
			// Create a new instance of the UserModel.
			r := NewUserRepository(db, &lg)

			got, err := r.ReadOne(tt.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ReadOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}
