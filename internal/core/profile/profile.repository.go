package profile

type ProfileRepository interface {
	ListAll() (map[string]*Profile, error)
	FindBySlug(slug string) (*Profile, error)
}
