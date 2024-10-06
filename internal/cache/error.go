package cache

import "errors"

var ErrNoStoreFile = errors.New("No cache file found")

var ErrNoProfileSlug = errors.New("Profile not found. Try using refresh command if you're sure the profile exists.")
