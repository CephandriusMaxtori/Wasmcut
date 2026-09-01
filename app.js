// Main application logic for Wasmcut Web UI
class WasmcutApp {
    constructor() {
        this.wasm = null;
        this.projectState = null;
        this.projectCreated = false;
        this.init();
    }

    async init() {
        console.log('Initializing Wasmcut Web UI...');
        
        try {
            // Load Wasm module
            this.wasm = new WasmLoader();
            
            // Try to load Wasm module from local directory or dist directory
            const wasmPaths = [
                './core.wasm',
                '../dist/core.wasm',
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
            this.showStartScreen();
            this.log('Wasm module ready');

        } catch (error) {
            this.updateWasmStatus(false);
            this.log(`Error initializing: ${error.message}`, 'error');
            console.error('Initialization error:', error);
        }
    }

    setupEventListeners() {
        document.getElementById('newProjectBtn').addEventListener('click', () => this.createProject());
        document.getElementById('newProjectBtn2').addEventListener('click', () => this.returnToStart());
        document.getElementById('getStateBtn').addEventListener('click', () => this.refreshState());
        document.getElementById('undoBtn').addEventListener('click', () => this.undo());
        document.getElementById('redoBtn').addEventListener('click', () => this.redo());
    }

    showStartScreen() {
        this.projectCreated = false;
        document.getElementById('startScreen').style.display = 'flex';
        document.getElementById('editorSection').style.display = 'none';
    }

    showEditor() {
        this.projectCreated = true;
        document.getElementById('startScreen').style.display = 'none';
        document.getElementById('editorSection').style.display = 'block';
    }

    returnToStart() {
        if (confirm('Are you sure you want to create a new project? Current project will be lost.')) {
            this.showStartScreen();
        }
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
                this.showEditor();
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
                this.projectState = result;
                this.displayState(result.data);
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

    displayState(response) {
        const stateDisplay = document.getElementById('stateDisplay');
        const timelineArea = document.getElementById('timelineArea');

        if (!response) {
            stateDisplay.textContent = 'No state available';
            timelineArea.innerHTML = '<div class="placeholder">No project data</div>';
            return;
        }

        // Handle response format from Go functions
        let stateData = response;
        if (response.data) {
            stateData = response.data;
        }

        // Display state as JSON (for development)
        try {
            stateDisplay.textContent = JSON.stringify(stateData, null, 2);
        } catch (error) {
            stateDisplay.textContent = String(stateData);
        }

        // Render timeline visualization
        this.renderTimeline(stateData);
    }

    renderTimeline(state) {
        const timelineArea = document.getElementById('timelineArea');

        if (!state) {
            timelineArea.innerHTML = '<div class="placeholder">No timeline data</div>';
            return;
        }

        // For now, show a placeholder with the state data
        // In a real implementation, we'd parse tracks and clips
        timelineArea.innerHTML = `
            <div class="placeholder">
                <p>Timeline active with ${state.tracks ? state.tracks.length : 0} track(s)</p>
                <p><em>Clip editing coming soon...</em></p>
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
