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
	GOOS=js GOARCH=wasm go build \
		-o dist/core.wasm \
		-ldflags="-s -w" \
		./wasm
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" dist/wasm_exec.js

test-wasm:
	cd wasm && go test -v ./...

clean:
	rm -rf dist/

.DEFAULT_GOAL := help
