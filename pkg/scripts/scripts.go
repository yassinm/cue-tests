package scripts

import "embed"

//go:embed pipeline/*
var StaticFs embed.FS
