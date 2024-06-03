package web

import "embed"

//go:embed certs/*
var Certs embed.FS
