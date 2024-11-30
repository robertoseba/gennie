package profile

import "errors"

var ErrNoProfileSlug = errors.New("profile not found. Try using refresh command if you're sure the profile exists")
var ErrLoadingToml = errors.New("error loading toml file")
var ErrProfileDirNotExist = errors.New("error reading profile directory")
