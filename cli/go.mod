module github.com/Makepad-fr/changelog/cli

go 1.23.4

replace github.com/Makepad-fr/changelog/parser => ../parser

replace github.com/Makepad-fr/changelog/core => ../core

require (
	github.com/Makepad-fr/changelog/core v0.0.0-00010101000000-000000000000
	github.com/Makepad-fr/changelog/parser v0.0.0-00010101000000-000000000000
	github.com/Makepad-fr/semver/semver v0.0.0-20240510163019-28f8831d8e0f
)
