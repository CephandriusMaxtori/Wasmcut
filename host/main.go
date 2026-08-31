package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	wasmPath := flag.String("wasm", "../dist/core.wasm", "Path to Wasm module")
	flag.Parse()

	// Load Wasm module
	fmt.Println("Loading Wasmcut core module...")
	module, err := NewWasmModule(*wasmPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load Wasm module: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Wasm module loaded successfully")

	// Run interactive CLI
	if err := module.RunCLI(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "CLI error: %v\n", err)
		os.Exit(1)
	}
}
