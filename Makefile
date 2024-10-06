LATEST_TAG := $(shell git describe --tags --abbrev=0 | sed 's/^v//')

run:
	@go run .
build:
	@go build -ldflags="-s -w" -o bin/ .

test:
	@go test -v ./...

release:
	@if [ -z "$(VERSION)" ]; then \
		echo "VERSION variable is required"; \
		exit 1; \
	fi
	@if [ "$(VERSION)" = "$(LATEST_TAG)" ]; then \
		echo "Version has already been released"; \
		exit 1; \
	fi

	git tag -a v$(VERSION) -m "Release version $(VERSION)"
	echo "Building version $(VERSION)" 
	@echo "$(VERSION) - Released at: $$(date +"%Y-%m-%d")" > version.txt
	@go build -ldflags="-s -w" -o bin/ .