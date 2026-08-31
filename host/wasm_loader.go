package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"unsafe"

	wasm "github.com/bytecodealliance/wasmtime-go/v14"
)

// WasmModule wraps a loaded Wasm module with its runtime
type WasmModule struct {
	Store     *wasm.Store
	Module    *wasm.Module
	Instance  *wasm.Instance
	Memory    *wasm.Memory
	ResultBuf []byte
}

// NewWasmModule loads and initializes the Wasm core module
func NewWasmModule(wasmPath string) (*WasmModule, error) {
	// Create engine and store
	engine := wasm.NewEngine()
	store := wasm.NewStore(engine)

	// Read Wasm binary
	data, err := os.ReadFile(wasmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read wasm file: %w", err)
	}

	// Compile module
	module, err := wasm.NewModule(store.Engine, data)
	if err != nil {
		return nil, fmt.Errorf("failed to compile wasm module: %w", err)
	}

	// Create WASI context for wasm32-wasip1 target
	wasiConfig := wasm.NewWasiConfig()
	wasiConfig.InheritStdout()
	wasiConfig.InheritStderr()

	// Create linker with WASI support
	linker := wasm.NewLinker(store.Engine)
	if err := linker.DefineWasi(wasiConfig); err != nil {
		return nil, fmt.Errorf("failed to define WASI: %w", err)
	}

	// Instantiate module with WASI
	instance, err := linker.Instantiate(store, module)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wasm module: %w", err)
	}

	// Get memory export
	memory := instance.GetExport(store, "memory").Memory()
	if memory == nil {
		return nil, fmt.Errorf("wasm module has no memory export")
	}

	return &WasmModule{
		Store:     store,
		Module:    module,
		Instance:  instance,
		Memory:    memory,
		ResultBuf: make([]byte, 1024*1024), // 1MB result buffer
	}, nil
}

// callFunction invokes an exported function and reads the result
func (w *WasmModule) callFunction(funcName string, args ...interface{}) ([]byte, error) {
	fn := w.Instance.GetExport(w.Store, funcName).Func()
	if fn == nil {
		return nil, fmt.Errorf("function %s not found in wasm module", funcName)
	}

	// Call function
	callResult, err := fn.Call(w.Store, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to call %s: %w", funcName, err)
	}

	// Get result size from return value
	size := int32(0)
	if callResult != nil {
		// Result might be int64 or int32 depending on Wasm target
		switch v := callResult.(type) {
		case int32:
			size = v
		case int64:
			size = int32(v)
		}
	}

	if size <= 0 {
		return nil, fmt.Errorf("function %s returned error or empty result", funcName)
	}

	// Read result from memory
	memData := w.Memory.Data(w.Store)
	memBytes := unsafe.Slice((*byte)(memData), size)
	resultBytes := make([]byte, size)
	copy(resultBytes, memBytes)
	return resultBytes, nil
}

// GetTimelineState queries the current project state from the Wasm module
func (w *WasmModule) GetTimelineState() (map[string]interface{}, error) {
	data, err := w.callFunction("GetTimelineState")
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse timeline state: %w", err)
	}

	return result, nil
}

// CreateProject initializes a new project in the Wasm module
func (w *WasmModule) CreateProject() (map[string]interface{}, error) {
	data, err := w.callFunction("CreateProject")
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse create project response: %w", err)
	}

	return result, nil
}

// Interactive CLI for testing
func (w *WasmModule) RunCLI(input io.Reader) error {
	scanner := bufio.NewScanner(input)

	fmt.Println("\n=== Wasmcut Host - M2 Proof ===")
	fmt.Println("Commands: state, create, help, exit")
	fmt.Print("> ")

	for scanner.Scan() {
		cmd := scanner.Text()

		switch cmd {
		case "state":
			state, err := w.GetTimelineState()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				data, _ := json.MarshalIndent(state, "", "  ")
				fmt.Printf("%s\n", string(data))
			}

		case "create":
			result, err := w.CreateProject()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				data, _ := json.MarshalIndent(result, "", "  ")
				fmt.Printf("%s\n", string(data))
			}

		case "help":
			fmt.Println(`
Available commands:
  state   - Show current timeline state
  create  - Create new project
  help    - Show this message
  exit    - Exit the program
`)

		case "exit":
			fmt.Println("Goodbye!")
			return nil

		default:
			fmt.Println("Unknown command. Type 'help' for options.")
		}

		fmt.Print("> ")
	}

	return scanner.Err()
}
