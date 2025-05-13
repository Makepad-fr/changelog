module github.com/Makepad-fr/changelog/parser

go 1.23.4

require (
	github.com/Makepad-fr/changelog/core v0.0.0-00010101000000-000000000000
	github.com/Makepad-fr/semver/semver v0.0.0-20240510163019-28f8831d8e0f
	golang.org/x/text v0.25.0
)

replace github.com/Makepad-fr/changelog/core => ../core
