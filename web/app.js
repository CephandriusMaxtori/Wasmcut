// Main application logic for Wasmcut Web UI
class WasmcutApp {
    constructor() {
        this.wasm = null;
        this.projectState = null;
        this.init();
    }

    async init() {
        console.log('Initializing Wasmcut Web UI...');
        
        try {
            // Load Wasm module
            this.wasm = new WasmLoader();
            
            // Try to load Wasm module from dist directory or web directory
            const wasmPaths = [
                '../dist/core.wasm',
                './core.wasm',
                'https://cdn.jsdelivr.net/npm/wasmcut@latest/dist/core.wasm'
            ];

            let loaded = false;
            for (const path of wasmPaths) {
                try {
                    console.log(`Trying to load Wasm from: ${path}`);
                    await this.wasm.load(path);
                    loaded = true;
                    break;
                } catch (error) {
                    console.warn(`Failed to load from ${path}:`, error.message);
                }
            }

            if (!loaded) {
                throw new Error('Could not load Wasm module from any source');
            }

            this.updateWasmStatus(true);
            this.setupEventListeners();
            this.log('Wasm module ready');

        } catch (error) {
            this.updateWasmStatus(false);
            this.log(`Error initializing: ${error.message}`, 'error');
            console.error('Initialization error:', error);
        }
    }

    setupEventListeners() {
        document.getElementById('createProjectBtn').addEventListener('click', () => this.createProject());
        document.getElementById('getStateBtn').addEventListener('click', () => this.refreshState());
        document.getElementById('undoBtn').addEventListener('click', () => this.undo());
        document.getElementById('redoBtn').addEventListener('click', () => this.redo());
    }

    async createProject() {
        this.log('Creating new project...');
        try {
            if (!this.wasm || !this.wasm.instance) {
                throw new Error('Wasm module not initialized');
            }

            const result = this.wasm.createProject();
            if (result.success) {
                this.log('Project created successfully');
                this.refreshState();
            } else {
                this.log(`Error: ${result.error}`, 'error');
            }
        } catch (error) {
            this.log(`Error creating project: ${error.message}`, 'error');
        }
    }

    async refreshState() {
        this.log('Fetching timeline state...');
        try {
            if (!this.wasm || !this.wasm.instance) {
                throw new Error('Wasm module not initialized');
            }

            const result = this.wasm.getTimelineState();
            if (result.success) {
                // In a real implementation, we'd parse the JSON from Wasm memory
                // For now, show the raw result
                this.projectState = result.result;
                this.displayState(result.result);
                this.log('Timeline state updated');
            } else {
                this.log(`Error: ${result.error}`, 'error');
            }
        } catch (error) {
            this.log(`Error fetching state: ${error.message}`, 'error');
        }
    }

    async undo() {
        this.log('Undoing last action...');
        try {
            const result = this.wasm.undo();
            if (result.success) {
                this.refreshState();
                this.log('Action undone');
            } else {
                this.log(`Error: ${result.error}`, 'error');
            }
        } catch (error) {
            this.log(`Error during undo: ${error.message}`, 'error');
        }
    }

    async redo() {
        this.log('Redoing last action...');
        try {
            const result = this.wasm.redo();
            if (result.success) {
                this.refreshState();
                this.log('Action redone');
            } else {
                this.log(`Error: ${result.error}`, 'error');
            }
        } catch (error) {
            this.log(`Error during redo: ${error.message}`, 'error');
        }
    }

    displayState(state) {
        const stateDisplay = document.getElementById('stateDisplay');
        const timelineArea = document.getElementById('timelineArea');

        if (!state) {
            stateDisplay.textContent = 'No state available';
            timelineArea.innerHTML = '<div class="placeholder">Load or create a project to see timeline</div>';
            return;
        }

        // Display raw state as JSON (for development)
        try {
            stateDisplay.textContent = JSON.stringify(state, null, 2);
        } catch (error) {
            stateDisplay.textContent = String(state);
        }

        // Render timeline visualization
        this.renderTimeline(state);
    }

    renderTimeline(state) {
        const timelineArea = document.getElementById('timelineArea');

        if (!state || typeof state !== 'number') {
            timelineArea.innerHTML = '<div class="placeholder">Invalid state format</div>';
            return;
        }

        // For now, show a placeholder
        // In a real implementation, we'd parse the actual timeline data
        timelineArea.innerHTML = `
            <div class="placeholder">
                <p>Wasm module returned: ${state}</p>
                <p>Next: Parse timeline JSON from Wasm memory</p>
            </div>
        `;
    }

    updateWasmStatus(ready) {
        const statusEl = document.getElementById('wasmStatus');
        if (ready) {
            statusEl.textContent = 'Wasm: Ready';
            statusEl.classList.add('ready');
        } else {
            statusEl.textContent = 'Wasm: Error';
            statusEl.classList.remove('ready');
        }
    }

    log(message, type = 'info') {
        const statusText = document.getElementById('statusText');
        statusText.textContent = message;
        console.log(`[${type.toUpperCase()}] ${message}`);
    }
}

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    window.app = new WasmcutApp();
});
