package repository

type Repository struct {
	*UserRepo
}

func NewRepository() *Repository {
	return &Repository{
		UserRepo: NewUserRepo(),
	}
}
