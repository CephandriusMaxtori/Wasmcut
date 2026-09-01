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

            // Log available exports for debugging
            const exportedNames = Object.keys(this.instance.exports);
            console.log('Wasm exports:', exportedNames.join(', '));

            // Start Go's runtime without awaiting its long-lived run promise.
            this.go.run(this.instance);
            await new Promise((resolve) => setTimeout(resolve, 0));

            // Try to get memory - it may be named differently or managed by Go runtime
            this.memory = this.instance.exports.memory;
            
            if (!this.memory && this.go.mem) {
                // Fallback: use Go runtime's memory buffer if available
                this.memory = this.go.mem;
            }

            // Log memory status for debugging
            if (this.memory) {
                console.log('✓ Memory available, size:', this.memory.buffer.byteLength, 'bytes');
            } else {
                console.warn('⚠ Memory not directly available - some features may be limited');
            }

            // Check for expected exported functions
            const expectedFuncs = ['CreateProject', 'GetTimelineState', 'AddClip', 'Undo', 'Redo'];
            const missingFuncs = expectedFuncs.filter(fn => !this.instance.exports[fn]);
            if (missingFuncs.length > 0) {
                console.warn('⚠ Missing exported functions:', missingFuncs.join(', '));
            }

            console.log('✓ Wasm module loaded successfully');
            return true;
        } catch (error) {
            console.error('Failed to load Wasm module:', error);
            throw error;
        }
    }

    call(functionName, ...args) {
        // In Go's js/wasm target, functions are registered on the global scope
        const func = typeof window !== 'undefined' 
            ? window[functionName] 
            : globalThis[functionName];

        if (!func || typeof func !== 'function') {
            throw new Error(`Function ${functionName} not found`);
        }

        try {
            const result = func(...args);
            
            // Parse JSON response if it's a string (from Go)
            if (typeof result === 'string') {
                try {
                    return JSON.parse(result);
                } catch (e) {
                    return { success: true, data: result };
                }
            }
            
            return result;
        } catch (error) {
            console.error(`Error calling ${functionName}:`, error);
            throw error;
        }
    }

    readMemory(offset, length) {
        if (!this.memory) {
            console.warn('Memory not available');
            return '';
        }

        try {
            const buffer = new Uint8Array(this.memory.buffer, offset, length);
            return new TextDecoder().decode(buffer);
        } catch (error) {
            console.error('Error reading memory:', error);
            return '';
        }
    }

    getMemory() {
        if (!this.memory) {
            return new Uint8Array(0);
        }
        return new Uint8Array(this.memory.buffer);
    }

    // Helper functions for common operations
    createProject() {
        try {
            const result = this.call('createProject');
            return result;
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    getTimelineState() {
        try {
            const result = this.call('getTimelineState');
            return result;
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    addClip(trackID, mediaRef, inTime, outTime, position) {
        try {
            const result = this.call('addClip', trackID, mediaRef, inTime, outTime, position);
            return result;
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    trimClip(clipID, trackID, oldIn, oldOut, newIn, newOut) {
        try {
            return this.call('trimClip', clipID, trackID, oldIn, oldOut, newIn, newOut);
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    deleteClip(clipID, trackID) {
        try {
            return this.call('deleteClip', clipID, trackID);
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    undo() {
        try {
            const result = this.call('undo');
            return result;
        } catch (error) {
            return { success: false, error: error.message };
        }
    }

    redo() {
        try {
            const result = this.call('redo');
            return result;
        } catch (error) {
            return { success: false, error: error.message };
        }
    }
}

// Export for use in app.js
window.WasmLoader = WasmLoader;
