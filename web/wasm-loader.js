// WasmLoader - handles loading and managing the Wasm module
class WasmLoader {
    constructor() {
        if (typeof Go === 'undefined') {
            throw new Error('Go Wasm runtime is not loaded');
        }

        this.go = new Go();
        this.module = null;
        this.memory = null;
        this.instance = null;
        this.resultBuffer = new ArrayBuffer(1024 * 1024); // 1MB buffer
    }

    async load(wasmPath) {
        try {
            console.log('Loading Wasm module from:', wasmPath);

            // Fetch the Wasm binary
            const response = await fetch(wasmPath);
            if (!response.ok) {
                throw new Error(`Failed to fetch Wasm module: ${response.statusText}`);
            }

            const buffer = await response.arrayBuffer();

            // Go's browser target needs the import object from wasm_exec.js.
            this.module = await WebAssembly.instantiate(buffer, this.go.importObject);

            this.instance = this.module.instance;

            // Start Go's runtime without awaiting its long-lived run promise.
            this.go.run(this.instance);
            await new Promise((resolve) => setTimeout(resolve, 0));

            this.memory = this.instance.exports.memory;

            if (!this.memory) {
                throw new Error('Wasm module does not export memory');
            }

            console.log('✓ Wasm module loaded successfully');
            return true;
        } catch (error) {
            console.error('Failed to load Wasm module:', error);
            throw error;
        }
    }

    call(functionName, ...args) {
        if (!this.instance) {
            throw new Error('Wasm module not loaded');
        }

        const func = this.instance.exports[functionName];
        if (!func) {
            throw new Error(`Function ${functionName} not found in Wasm module`);
        }

        try {
            const result = func(...args);
            return result;
        } catch (error) {
            console.error(`Error calling ${functionName}:`, error);
            throw error;
        }
    }

    readMemory(offset, length) {
        if (!this.memory) {
            throw new Error('Memory not available');
        }

        const buffer = new Uint8Array(this.memory.buffer, offset, length);
        return new TextDecoder().decode(buffer);
    }

    getMemory() {
        return new Uint8Array(this.memory.buffer);
    }

    // Helper functions for common operations
    createProject() {
        try {
            const result = this.call('CreateProject');
            return { success: true, result };
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    getTimelineState() {
        try {
            const result = this.call('GetTimelineState');
            return { success: true, result };
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    addClip(trackID, mediaRef, inTime, outTime, position) {
        try {
            const result = this.call('AddClip', trackID, mediaRef, inTime, outTime, position);
            return { success: true, result };
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    undo() {
        try {
            const result = this.call('Undo');
            return { success: true, result };
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    redo() {
        try {
            const result = this.call('Redo');
            return { success: true, result };
        } catch (error) {
            return { success: false, error: error.message };
        }
    }
}

// Export for use in app.js
window.WasmLoader = WasmLoader;
