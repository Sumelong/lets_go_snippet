package sqlite

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"snippetbox/pkg/domain/models"
	"snippetbox/pkg/logger"
	"strings"
	"sync"
)

type UserRepository struct {
	db *sql.DB
	mu sync.Mutex
	lg logger.ILogger
}

func NewUserRepository(db *sql.DB, lg *logger.ILogger) *UserRepository {
	return &UserRepository{
		db: db,
		lg: *lg,
	}
}

func (r *UserRepository) Create(user models.User) (uint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stmt := `INSERT INTO users (name, email, hashed_password, created, active)
			 VALUES(?, ?,?,CURRENT_TIMESTAMP, ?)`

	result, err := r.db.Exec(stmt, user.Name, user.Email, user.HashedPassword, user.Active)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, models.ErrDuplicateEmail
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.lg.Error(models.ErrEntityCreationFailed.Error(), user)
		return 0, err
	}

	return uint(id), nil
}

func (r *UserRepository) ReadAll() ([]*models.User, error) {

	stmt := `SELECT id, name, email, hashed_password, created, active FROM users
	 ORDER BY created DESC LIMIT 10`

	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {

		u := &models.User{}

		err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.Created, &u.Active)
		if err != nil {
			return nil, err
		}

		// Append it to the slice of users.
		users = append(users, u)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Authenticate(email, password string) (int, error) {
	// Retrieve the id and hashed password associated with the given email. If no
	// matching email exists, or the user is not active, we return the
	// ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := r.db.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

func (r *UserRepository) ReadOne(id int) (*models.User, error) {

	stmt := `SELECT id, name, email,  created, active FROM users
	WHERE id = ?`

	row := r.db.QueryRow(stmt, id)
	// Initialize a pointer to a new zeroed Snippet struct.
	u := &models.User{}

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If everything went OK then return the Snippet object.
	return u, nil
}

func (r *UserRepository) ReadBy(user *models.User) ([]*models.User, error) {
	return nil, nil
}

func (r *UserRepository) Update(user *models.User) (uint, error) {
	return 0, nil
}

func (r *UserRepository) Delete(id uint) (uint, error) {
	return 0, nil
}
