# Convenience targets for Stride Mechanica. Everything here can also be run directly, which is the way to go
# on Windows without make: see the commands below, the main README and docs/cad-export.md.

PYTHON ?= python

.PHONY: build test run web cad-setup cad-vendor cad-check

build:
	go build

test:
	go test ./...

run:
	go run .

# Build the browser version into web/ and serve it on http://localhost:8080
web:
	go run ./cmd/web

# Install CadQuery in cad/.venv (needs Python 3.9 - 3.12)
cad-setup:
	$(PYTHON) cad/tools.py setup

# Download the vendor STEP files that can be downloaded automatically, and explain the others
cad-vendor:
	$(PYTHON) cad/tools.py vendor

# Show what is installed, downloaded and measured
cad-check:
	$(PYTHON) cad/tools.py check
