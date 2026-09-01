// Main application logic for Wasmcut Web UI
class WasmcutApp {
    constructor() {
        this.wasm = null;
        this.projectState = null;
        this.projectCreated = false;
        this.currentMode = 'editing';
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
        const bind = (id, event, handler) => {
            const element = document.getElementById(id);
            if (element) element.addEventListener(event, handler);
        };

        bind('newProjectBtn', 'click', () => this.createProject());
        bind('newProjectBtn2', 'click', () => this.createProject());
        bind('getStateBtn', 'click', () => this.refreshState());
        bind('undoBtn', 'click', () => this.undo());
        bind('redoBtn', 'click', () => this.redo());
        bind('addClipForm', 'submit', (event) => this.addClip(event));
        bind('editClipForm', 'submit', (event) => this.trimSelectedClip(event));
        bind('deleteClipBtn', 'click', () => this.deleteSelectedClip());

        // Mode selector buttons
        document.querySelectorAll('.mode-btn').forEach(btn => {
            btn.addEventListener('click', (e) => this.switchMode(e.target.dataset.mode));
        });
    }

    switchMode(mode) {
        // Update active button
        document.querySelectorAll('.mode-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        const activeButton = document.querySelector(`[data-mode="${mode}"]`);
        if (activeButton) activeButton.classList.add('active');
        
        // Handle mode switching (for future implementation)
        console.log('Switched to mode:', mode);
        this.currentMode = mode;
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
        if (this.projectCreated) {
            if (confirm('Are you sure you want to create a new project? Current project will be lost.')) {
                this.showStartScreen();
            }
        } else {
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

    addClip(event) {
        event.preventDefault();
        const mediaRef = document.getElementById('mediaRefInput').value.trim();
        const inTime = Number(document.getElementById('clipInInput').value);
        const outTime = Number(document.getElementById('clipOutInput').value);
        const position = Number(document.getElementById('clipPositionInput').value);
        if (!mediaRef || !Number.isFinite(inTime) || !Number.isFinite(outTime) || outTime <= inTime) {
            this.log('Enter a media reference and a valid In/Out range', 'error');
            return;
        }
        const result = this.wasm.addClip('track-1', mediaRef, inTime, outTime, position);
        if (result.success) {
            this.log('Clip added');
            this.refreshState();
        } else {
            this.log(`Error: ${result.error}`, 'error');
        }
    }

    selectClip(clip, trackID) {
        this.selectedClip = { ...clip, trackID };
        document.getElementById('clipEditor').hidden = false;
        document.getElementById('selectedClipName').textContent = `${clip.id} · ${clip.media_ref}`;
        document.getElementById('editInInput').value = clip.in;
        document.getElementById('editOutInput').value = clip.out;
        document.querySelectorAll('.clip.selected').forEach(element => element.classList.remove('selected'));
        const selectedElement = document.querySelector(`[data-clip-id="${clip.id}"]`);
        if (selectedElement) selectedElement.classList.add('selected');
    }

    trimSelectedClip(event) {
        event.preventDefault();
        if (!this.selectedClip) return;
        const newIn = Number(document.getElementById('editInInput').value);
        const newOut = Number(document.getElementById('editOutInput').value);
        if (!Number.isFinite(newIn) || !Number.isFinite(newOut) || newOut <= newIn) {
            this.log('Out must be greater than In', 'error');
            return;
        }
        const result = this.wasm.trimClip(this.selectedClip.id, this.selectedClip.trackID, this.selectedClip.in, this.selectedClip.out, newIn, newOut);
        if (result.success) {
            this.selectedClip.in = newIn;
            this.selectedClip.out = newOut;
            this.log('Clip trimmed');
            this.refreshState();
        } else this.log(`Error: ${result.error}`, 'error');
    }

    deleteSelectedClip() {
        if (!this.selectedClip) return;
        const result = this.wasm.deleteClip(this.selectedClip.id, this.selectedClip.trackID);
        if (result.success) {
            this.selectedClip = null;
            document.getElementById('clipEditor').hidden = true;
            this.log('Clip deleted');
            this.refreshState();
        } else this.log(`Error: ${result.error}`, 'error');
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

        const tracks = state.tracks || [];
        timelineArea.classList.toggle('has-content', tracks.some(track => track.clips && track.clips.length));
        timelineArea.innerHTML = tracks.map(track => `
            <div class="track">
                <span class="track-label">${track.id}</span>
                <div class="clips-container">
                    ${(track.clips || []).map(clip => `<button class="clip" data-clip-id="${clip.id}" type="button">${clip.id}: ${clip.media_ref} (${clip.in}s - ${clip.out}s)</button>`).join('') || '<span class="placeholder">Empty track</span>'}
                </div>
            </div>
        `).join('') || '<div class="placeholder">No tracks</div>';
        timelineArea.querySelectorAll('.clip').forEach(element => {
            element.addEventListener('click', () => {
                const track = tracks.find(item => item.clips.some(clip => clip.id === element.dataset.clipId));
                const clip = track && track.clips.find(item => item.id === element.dataset.clipId);
                if (clip) this.selectClip(clip, track.id);
            });
        });
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
