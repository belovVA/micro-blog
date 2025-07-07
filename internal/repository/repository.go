package repository

type Repository struct {
	*UserRepo
	*PostRepo
}

func NewRepository() *Repository {
	return &Repository{
		UserRepo: NewUserRepo(),
		PostRepo: NewPostRepo(),
	}
}
