.PHONY: help build-wasm test-wasm clean dist

help:
	@echo "Wasmcut build targets:"
	@echo "  make build-wasm    - Build Wasm core module to wasm32-wasip1"
	@echo "  make test-wasm     - Run tests on Wasm module"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make dist          - Create dist directory"

dist:
	mkdir -p dist

build-wasm: dist
	GOOS=wasip1 GOARCH=wasm go build \
		-o dist/core.wasm \
		-ldflags="-s -w" \
		./wasm

test-wasm:
	cd wasm && go test -v ./...

clean:
	rm -rf dist/

.DEFAULT_GOAL := help
