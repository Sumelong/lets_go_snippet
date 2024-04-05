package ports

import "snippetbox/pkg/domain/models"

/*
type StoreReader[T any] interface {
	ReadAll() ([]*T, error)
	ReadOne(uint) (*T, error)
	ReadBy(*T) ([]*T, error)
}

type StoreWriter[T any] interface {
	Create(*T) (uint, error)
	Update(*T) (uint, error)
	Delete(uint) (uint, error)
}

type Repository[T any] interface {
	StoreWriter[T]
	StoreReader[T]
}

*/

type Repository[T any] interface {
	Create(T) (uint, error)
	ReadAll() ([]*T, error)
	ReadOne(int) (*T, error)
	ReadBy(*T) ([]*T, error)
	Update(*T) (uint, error)
	Delete(uint) (uint, error)
}

type IUserRepository interface {
	Repository[models.User]
	Authenticate(email, password string) (int, error)
}

type ISnippetRepository interface {
	Repository[models.Snippet]
}
