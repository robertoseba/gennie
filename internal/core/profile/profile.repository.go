package profile

type IProfileRepository interface {
	ListAll() (map[string]*Profile, error)
	FindBySlug(slug string) (*Profile, error)
}
