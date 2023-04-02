package assets

import (
	"embed"
)

// content is our static web server content.
//
//go:embed frontend/* swagger/*
var (
	Content embed.FS
)
