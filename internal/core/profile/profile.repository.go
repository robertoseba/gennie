package profile

type ProfileRepositoryInterface interface {
	ListAll() (map[string]*Profile, error)
	FindBySlug(slug string) (*Profile, error)
}
