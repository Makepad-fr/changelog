.PHONY: build
build:
	go build -o ./bin/changelog ./changelog/changelog.go

.PHONY: install
install:
	go install ./changelog/changelog.go
