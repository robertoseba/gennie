LATEST_TAG := $(shell git describe --tags --abbrev=0 | sed 's/^v//')

run:
	@go run .
build:
	@go build -ldflags="-s -w" -o bin/ .

test:
	@go test -v ./... | colorize.sh

lint:
	@golangci-lint run ./...

test-cover:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

current-tag:
	@echo "$(LATEST_TAG)"

release:
	@if ! make test; then \
		exit 1; \
	fi

	@if ! make lint; then \
		exit 1; \
	fi

	@if [ -z "$(VERSION)" ]; then \
		echo "VERSION variable is required"; \
		exit 1; \
	fi
	@if [ "$(VERSION)" = "$(LATEST_TAG)" ]; then \
		echo "Version has already been released"; \
		exit 1; \
	fi

	@echo "$(VERSION) - Released at: $$(date +"%Y-%m-%d")" > version.txt
	git add version.txt
	git commit -m "Release version $(VERSION)"
	git tag -a v$(VERSION) -m "Release version $(VERSION)"
	echo "Building version $(VERSION)" 
	@go build -ldflags="-s -w" -o bin/ .
	@read -p "Want to push it to remote? [y/N] " answer; \
         if [ "$$answer" != "y" ]; then \
             exit 1; \
         fi
	git push origin --follow-tags 
