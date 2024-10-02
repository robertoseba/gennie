package profile

type ErrNoFolder struct {
	Path string
}

func (e ErrNoFolder) Error() string {
	return "No profiles folder in " + e.Path
}

type ErrNoProfiles struct {
	Path string
}

func (e ErrNoProfiles) Error() string {
	return "No profiles found in " + e.Path
}
