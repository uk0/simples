package parser

import (
	"fmt"
	"html/template"
)

func GenerateNESPlayerHTML(playerID, romPath, fileName, fileExt string) string {
	// Ensure all user inputs are properly escaped
	safePlayerID := template.HTMLEscapeString(playerID)
	safeRomPath := template.HTMLEscapeString(romPath)
	safeFileName := template.HTMLEscapeString(fileName)

	return fmt.Sprintf(`<div class="nes-player-container" id="%[1]s-container">
  <style>
    #%[1]s-container {
      margin: 20px 0;
      padding: 20px;
      background: #1a1a1a;
      border-radius: 8px;
      display: inline-block;
      max-width: 100%%;
      position: relative;
    }
    #%[1]s-container .nes-header {
      color: #fff;
      margin-bottom: 15px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-wrap: wrap;
      gap: 10px;
    }
    #%[1]s-container .nes-title {
      font-size: 16px;
      font-weight: bold;
      word-break: break-word;
    }
    #%[1]s-container .nes-controls {
      display: flex;
      gap: 10px;
      margin-bottom: 10px;
      flex-wrap: wrap;
    }
    #%[1]s-container button {
      background: #4CAF50;
      color: white;
      border: none;
      padding: 8px 16px;
      border-radius: 4px;
      cursor: pointer;
      font-size: 14px;
      transition: background 0.3s;
    }
    #%[1]s-container button:hover:not(:disabled) {
      background: #45a049;
    }
    #%[1]s-container button:disabled {
      background: #666;
      cursor: not-allowed;
      opacity: 0.6;
    }
    #%[1]s-container .nes-canvas-wrapper {
      position: relative;
      display: inline-block;
    }
    #%[1]s-container canvas {
      border: 2px solid #333;
      display: block;
      image-rendering: pixelated;
      image-rendering: -moz-crisp-edges;
      image-rendering: crisp-edges;
      width: 512px;
      height: 480px;
      max-width: 100%%;
      height: auto;
      background: #000;
    }
    #%[1]s-container .nes-loading {
      color: #fff;
      text-align: center;
      padding: 50px;
    }
    #%[1]s-container .nes-error {
      color: #ff4444;
      padding: 20px;
      background: rgba(255, 68, 68, 0.1);
      border-radius: 4px;
      margin-top: 10px;
    }
    #%[1]s-container .nes-warning {
      color: #ffa500;
      padding: 15px;
      background: rgba(255, 165, 0, 0.1);
      border-radius: 4px;
      margin-top: 10px;
      font-size: 13px;
    }
    #%[1]s-container .nes-debug {
      color: #4CAF50;
      padding: 15px;
      background: rgba(76, 175, 80, 0.1);
      border-radius: 4px;
      margin-top: 10px;
      font-family: monospace;
      font-size: 12px;
      white-space: pre-wrap;
      word-break: break-all;
      max-height: 300px;
      overflow-y: auto;
    }
    #%[1]s-container .nes-instructions {
      color: #aaa;
      font-size: 12px;
      margin-top: 10px;
      line-height: 1.5;
    }
    #%[1]s-container .nes-info {
      color: #888;
      font-size: 11px;
      margin-top: 5px;
    }
    #%[1]s-container .nes-alternative {
      margin-top: 15px;
      padding: 15px;
      background: rgba(255, 255, 255, 0.05);
      border-radius: 4px;
    }
    #%[1]s-container .nes-alternative h4 {
      color: #fff;
      margin: 0 0 10px 0;
      font-size: 14px;
    }
    @media (max-width: 600px) {
      #%[1]s-container canvas {
        width: 100%%;
      }
    }
  </style>
  
  <div class="nes-header">
    <div class="nes-title">🎮 %[3]s</div>
    <a href="%[2]s" download style="color: #4CAF50; text-decoration: none; font-size: 14px;">⬇ Download ROM</a>
  </div>
  
  <div id="%[1]s-loading" class="nes-loading">
    <div>Loading NES emulator...</div>
    <div style="font-size: 12px; color: #888; margin-top: 10px;">Analyzing ROM format...</div>
  </div>
  
  <div id="%[1]s-player" style="display: none;">
    <div class="nes-controls">
      <button id="%[1]s-play">▶ Play</button>
      <button id="%[1]s-pause" style="display: none;">⏸ Pause</button>
      <button id="%[1]s-reset">🔄 Reset</button>
      <button id="%[1]s-fullscreen">⛶ Fullscreen</button>
      <button id="%[1]s-mute">🔇 Mute</button>
      <button id="%[1]s-switch-emu" style="display: none;">🔄 Switch Emulator</button>
    </div>
    <div class="nes-canvas-wrapper">
      <canvas id="%[1]s-canvas" width="256" height="240"></canvas>
    </div>
    <div class="nes-instructions">
      <strong>Controls:</strong><br>
      Arrow Keys: D-Pad | Z: A Button | X: B Button<br>
      Enter: Start | Shift: Select
    </div>
    <div id="%[1]s-info" class="nes-info"></div>
  </div>
  
  <div id="%[1]s-error" class="nes-error" style="display: none;"></div>
  <div id="%[1]s-warning" class="nes-warning" style="display: none;"></div>
  <div id="%[1]s-debug" class="nes-debug" style="display: none;"></div>
  
  <div id="%[1]s-alternative" class="nes-alternative" style="display: none;">
    <h4>Alternative: Try Web-Based Emulator</h4>
    <div>
      <button onclick="window.open('https://jsnes.org/#' + encodeURIComponent(window.location.origin + '/' + '%[2]s'), '_blank')" 
              style="background: #2196F3; color: white; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer;">
        Open in JSNES.org
      </button>
      <button onclick="window.open('https://emulatorjs.org/nes?rom=' + encodeURIComponent(window.location.origin + '/' + '%[2]s'), '_blank')" 
              style="background: #FF9800; color: white; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer; margin-left: 10px;">
        Open in EmulatorJS
      </button>
    </div>
  </div>
  
  <script>
  (function() {
    'use strict';
    
    // Configuration
    const PLAYER_ID = '%[1]s';
    const ROM_PATH = '%[2]s';
    const FILE_EXT = '%[4]s';
    const DEBUG_MODE = true;
    const USE_FALLBACK = true; // Enable fallback to alternative emulator
    
    // State management
    const state = {
      emulatorType: 'jsnes', // 'jsnes' or 'nes-embed'
      jsnesLoaded: false,
      nesInstance: null,
      animationId: null,
      audioContext: null,
      scriptProcessor: null,
      isRunning: false,
      isMuted: false,
      audioBuffer: [],
      romData: null,
      romInfo: null
    };
    
    // Get DOM elements
    const elements = {
      container: document.getElementById(PLAYER_ID + '-container'),
      loading: document.getElementById(PLAYER_ID + '-loading'),
      player: document.getElementById(PLAYER_ID + '-player'),
      error: document.getElementById(PLAYER_ID + '-error'),
      warning: document.getElementById(PLAYER_ID + '-warning'),
      debug: document.getElementById(PLAYER_ID + '-debug'),
      alternative: document.getElementById(PLAYER_ID + '-alternative'),
      canvas: document.getElementById(PLAYER_ID + '-canvas'),
      playBtn: document.getElementById(PLAYER_ID + '-play'),
      pauseBtn: document.getElementById(PLAYER_ID + '-pause'),
      resetBtn: document.getElementById(PLAYER_ID + '-reset'),
      fullscreenBtn: document.getElementById(PLAYER_ID + '-fullscreen'),
      muteBtn: document.getElementById(PLAYER_ID + '-mute'),
      switchBtn: document.getElementById(PLAYER_ID + '-switch-emu'),
      info: document.getElementById(PLAYER_ID + '-info')
    };
    
    // Debug output
    function debug(message, data) {
      console.log('[NES Player ' + PLAYER_ID + '] ' + message, data || '');
      if (DEBUG_MODE && elements.debug) {
        elements.debug.style.display = 'block';
        const timestamp = new Date().toLocaleTimeString();
        elements.debug.textContent += '[' + timestamp + '] ' + message + (data ? ': ' + JSON.stringify(data) : '') + '\\n';
        elements.debug.scrollTop = elements.debug.scrollHeight;
      }
    }
    
    // Show warning
    function showWarning(message) {
      if (elements.warning) {
        elements.warning.innerHTML = message;
        elements.warning.style.display = 'block';
      }
    }
    
    // Error handling
    function showError(message, canRetry, debugInfo) {
      console.error('[NES Player ' + PLAYER_ID + ']', message);
      if (elements.error) {
        let errorHTML = message;
        if (canRetry) {
          errorHTML += '<br><button onclick="location.reload()" style="margin-top: 10px; background: #4CAF50; color: white; border: none; padding: 5px 10px; border-radius: 3px; cursor: pointer;">Retry</button>';
        }
        if (debugInfo) {
          errorHTML += '<br><br><details style="margin-top: 10px;"><summary style="cursor: pointer;">Technical Details</summary><pre style="margin-top: 10px; font-size: 11px; overflow-x: auto;">' + debugInfo + '</pre></details>';
        }
        elements.error.innerHTML = errorHTML;
        elements.error.style.display = 'block';
      }
      if (elements.loading) {
        elements.loading.style.display = 'none';
      }
      
      // Show alternative options
      if (USE_FALLBACK && elements.alternative) {
        elements.alternative.style.display = 'block';
      }
    }
    
    // Analyze ROM header in detail
    function analyzeROMHeader(data) {
      const uint8 = new Uint8Array(data);
      if (uint8.length < 16) return null;
      
      const header = {
        magic: String.fromCharCode(uint8[0], uint8[1], uint8[2], uint8[3]),
        prgRomSize: uint8[4] * 16384, // in bytes
        chrRomSize: uint8[5] * 8192,  // in bytes
        flags6: uint8[6],
        flags7: uint8[7],
        flags8: uint8[8],
        flags9: uint8[9],
        flags10: uint8[10],
        // Derived values
        mapper: ((uint8[7] & 0xF0) | ((uint8[6] & 0xF0) >> 4)),
        mirroring: (uint8[6] & 0x01) ? 'Vertical' : 'Horizontal',
        battery: (uint8[6] & 0x02) ? true : false,
        trainer: (uint8[6] & 0x04) ? true : false,
        fourScreen: (uint8[6] & 0x08) ? true : false,
        vsUnisystem: (uint8[7] & 0x01) ? true : false,
        playChoice10: (uint8[7] & 0x02) ? true : false,
        nes2Format: ((uint8[7] & 0x0C) === 0x08) ? true : false
      };
      
      // Calculate expected file size
      header.expectedSize = 16 + // header
                           (header.trainer ? 512 : 0) + // trainer
                           header.prgRomSize + // PRG ROM
                           header.chrRomSize;  // CHR ROM
      
      header.actualSize = uint8.length;
      header.sizeMatch = header.expectedSize === header.actualSize;
      
      return header;
    }
    
    // Validate ROM more thoroughly
    function validateROM(data) {
      const uint8 = new Uint8Array(data);
      const header = analyzeROMHeader(data);
      
      if (!header) {
        return { valid: false, reason: 'File too small to be a NES ROM' };
      }
      
      if (header.magic !== 'NES\\x1A') {
        return { valid: false, reason: 'Invalid NES header magic bytes' };
      }
      
      // Check for common issues
      if (header.prgRomSize === 0) {
        return { valid: false, reason: 'PRG ROM size is 0' };
      }
      
      if (!header.sizeMatch) {
        debug('Size mismatch - Expected: ' + header.expectedSize + ', Actual: ' + header.actualSize);
        // Some ROMs have extra data, so this might not be fatal
        if (header.actualSize < header.expectedSize) {
          return { valid: false, reason: 'ROM file is truncated (missing ' + (header.expectedSize - header.actualSize) + ' bytes)' };
        }
      }
      
      // Check if mapper is supported by JSNES
      const supportedMappers = [0, 1, 2, 3, 4, 5, 7, 9, 10, 11, 15, 19, 21, 22, 23, 24, 26, 32, 33, 34, 66, 68, 69, 71, 78, 87, 94, 140, 180];
      if (!supportedMappers.includes(header.mapper)) {
        showWarning('⚠️ This ROM uses Mapper ' + header.mapper + ' which may not be fully supported. Some games might not work correctly.');
      }
      
      // Store ROM info for display
      state.romInfo = header;
      
      return { valid: true, header: header };
    }
    
    // Process ROM data with validation
    function processROMData(data) {
      const uint8Data = new Uint8Array(data);
      
      debug('Analyzing ROM file...');
      const validation = validateROM(data);
      
      if (!validation.valid) {
        debug('ROM validation failed: ' + validation.reason);
        
        // Try to fix common issues
        if (validation.reason.includes('Invalid NES header')) {
          // Try to find NES header at offset
          for (let i = 1; i < Math.min(512, uint8Data.length - 16); i++) {
            if (uint8Data[i] === 0x4E && 
                uint8Data[i+1] === 0x45 && 
                uint8Data[i+2] === 0x53 && 
                uint8Data[i+3] === 0x1A) {
              debug('Found NES header at offset ' + i + ', attempting to fix...');
              const fixedData = uint8Data.slice(i);
              const revalidation = validateROM(fixedData);
              if (revalidation.valid) {
                debug('Successfully fixed ROM by removing ' + i + ' byte offset');
                return fixedData;
              }
            }
          }
        }
        
        const first32 = Array.from(uint8Data.slice(0, 32));
        const hexView = first32.map(b => ('0' + b.toString(16)).slice(-2).toUpperCase()).join(' ');
        
        showError('Invalid NES ROM: ' + validation.reason, false, 
                  'File size: ' + uint8Data.length + ' bytes\\n' +
                  'First 32 bytes (hex):\\n' + hexView);
        return null;
      }
      
      // Display ROM information
      if (validation.header && elements.info) {
        const h = validation.header;
        elements.info.innerHTML = 
          'Mapper: ' + h.mapper + ' | ' +
          'PRG: ' + (h.prgRomSize/1024) + 'KB | ' +
          'CHR: ' + (h.chrRomSize/1024) + 'KB | ' +
          'Mirror: ' + h.mirroring +
          (h.battery ? ' | Battery' : '') +
          (h.trainer ? ' | Trainer' : '');
      }
      
      debug('ROM validation passed', validation.header);
      return uint8Data;
    }
    
    // Try loading with nes-embed as fallback
    function tryNesEmbed() {
      debug('Attempting to use nes-embed library as fallback...');
      
      const script = document.createElement('script');
      script.src = 'https://unpkg.com/nes-embed@1.0.0/dist/nes-embed.min.js';
      script.onload = function() {
        debug('nes-embed loaded, initializing...');
        
        try {
          const nesEmbed = new window.NesEmbed(elements.canvas, {
            rom: state.romData,
            autoStart: false
          });
          
          state.nesInstance = nesEmbed;
          state.emulatorType = 'nes-embed';
          
          elements.loading.style.display = 'none';
          elements.player.style.display = 'block';
          
          setupAlternativeControls();
          
          showWarning('⚠️ Using alternative emulator (nes-embed). Some features may differ.');
        } catch (e) {
          debug('nes-embed failed: ' + e.message);
          showError('Both JSNES and nes-embed failed to load this ROM.', false, e.message);
        }
      };
      script.onerror = function() {
        debug('Failed to load nes-embed library');
        showError('Failed to load alternative emulator library.', true);
      };
      document.head.appendChild(script);
    }
    
    // Load JSNES library
    function loadJSNES() {
      if (window.jsnes && window.jsnes.NES) {
        state.jsnesLoaded = true;
        initializePlayer();
        return;
      }
      
      const existingScript = document.querySelector('script[data-jsnes="true"]');
      if (existingScript) {
        existingScript.addEventListener('load', function() {
          state.jsnesLoaded = true;
          initializePlayer();
        });
        return;
      }
      
      const script = document.createElement('script');
      script.src = 'https://unpkg.com/jsnes@1.2.1/dist/jsnes.min.js';
      script.setAttribute('data-jsnes', 'true');
      script.onload = function() {
        state.jsnesLoaded = true;
        initializePlayer();
      };
      script.onerror = function() {
        showError('Failed to load JSNES emulator library. Please check your internet connection.', true);
      };
      document.head.appendChild(script);
    }
    
    // Initialize the NES player
    function initializePlayer() {
      if (!window.jsnes || !window.jsnes.NES) {
        showError('JSNES library not available', true);
        return;
      }
      
      const ctx = elements.canvas.getContext('2d');
      if (!ctx) {
        showError('Failed to get canvas context', false);
        return;
      }
      
      // Fetch and process ROM
      debug('Fetching ROM from: ' + ROM_PATH);
      fetch(ROM_PATH)
        .then(response => {
          if (!response.ok) {
            throw new Error('HTTP ' + response.status);
          }
          const contentLength = response.headers.get('content-length');
          debug('ROM response received, size: ' + (contentLength || 'unknown') + ' bytes');
          return response.arrayBuffer();
        })
        .then(buffer => {
          debug('Processing ROM data...');
          const romData = processROMData(buffer);
          
          if (!romData) {
            return; // Error already shown
          }
          
          state.romData = romData; // Store for potential fallback
          
          // Try to load with JSNES
          try {
            const imageData = ctx.getImageData(0, 0, 256, 240);
            
            // Initialize audio
            try {
              state.audioContext = new (window.AudioContext || window.webkitAudioContext)();
              state.scriptProcessor = state.audioContext.createScriptProcessor(1024, 0, 2);
              state.scriptProcessor.connect(state.audioContext.destination);
              
              state.scriptProcessor.onaudioprocess = function(e) {
                if (state.isMuted) return;
                
                const left = e.outputBuffer.getChannelData(0);
                const right = e.outputBuffer.getChannelData(1);
                const len = left.length;
                
                for (let i = 0; i < len; i++) {
                  left[i] = state.audioBuffer.shift() || 0;
                  right[i] = state.audioBuffer.shift() || 0;
                }
                
                if (state.audioBuffer.length > 8192) {
                  state.audioBuffer = state.audioBuffer.slice(-8192);
                }
              };
            } catch (e) {
              debug('Audio initialization failed: ' + e.message);
            }
            
            // Create NES instance with error handling
            state.nesInstance = new jsnes.NES({
              onFrame: function(framebuffer) {
                let i = 0;
                for (let y = 0; y < 240; y++) {
                  for (let x = 0; x < 256; x++) {
                    const offset = y * 256 + x;
                    imageData.data[i] = framebuffer[offset] & 0xFF;
                    imageData.data[i + 1] = (framebuffer[offset] >> 8) & 0xFF;
                    imageData.data[i + 2] = (framebuffer[offset] >> 16) & 0xFF;
                    imageData.data[i + 3] = 0xFF;
                    i += 4;
                  }
                }
                ctx.putImageData(imageData, 0, 0);
              },
              onAudioSample: function(left, right) {
                if (!state.isMuted && state.audioBuffer.length < 8192) {
                  state.audioBuffer.push(left);
                  state.audioBuffer.push(right);
                }
              },
              onStatusUpdate: function(status) {
                debug('NES Status: ' + status);
              },
              preferredFrameRate: 60,
              sampleRate: state.audioContext ? state.audioContext.sampleRate : 44100
            });
            
            // Try to load ROM
            state.nesInstance.loadROM(romData);
            debug('ROM loaded successfully into JSNES!');
            
            elements.loading.style.display = 'none';
            elements.player.style.display = 'block';
            
            setupControls();
            
          } catch (e) {
            debug('JSNES failed to load ROM: ' + e.message);
            
            // Show detailed error
            let debugInfo = 'JSNES Error: ' + e.message + '\\n\\n';
            if (state.romInfo) {
              debugInfo += 'ROM Header Info:\\n';
              debugInfo += '  Mapper: ' + state.romInfo.mapper + '\\n';
              debugInfo += '  PRG ROM: ' + (state.romInfo.prgRomSize/1024) + 'KB\\n';
              debugInfo += '  CHR ROM: ' + (state.romInfo.chrRomSize/1024) + 'KB\\n';
              debugInfo += '  Mirroring: ' + state.romInfo.mirroring + '\\n';
              debugInfo += '  File size: ' + state.romInfo.actualSize + ' bytes\\n';
              debugInfo += '  Expected size: ' + state.romInfo.expectedSize + ' bytes\\n';
            }
            
            if (USE_FALLBACK && e.message.includes('Not a valid NES ROM')) {
              showWarning('JSNES could not load this ROM. Trying alternative emulator...');
              tryNesEmbed();
            } else {
              showError('Failed to load ROM: ' + e.message, false, debugInfo);
            }
          }
        })
        .catch(error => {
          debug('Failed to fetch ROM: ' + error.message);
          showError('Failed to fetch ROM file: ' + error.message, true);
        });
    }
    
    // Setup controls for JSNES
    function setupControls() {
      elements.playBtn.onclick = function() {
        if (state.audioContext && state.audioContext.state === 'suspended') {
          state.audioContext.resume();
        }
        startEmulation();
        this.style.display = 'none';
        elements.pauseBtn.style.display = 'inline-block';
      };
      
      elements.pauseBtn.onclick = function() {
        stopEmulation();
        this.style.display = 'none';
        elements.playBtn.style.display = 'inline-block';
      };
      
      elements.resetBtn.onclick = function() {
        if (state.nesInstance) {
          state.nesInstance.reset();
        }
      };
      
      elements.fullscreenBtn.onclick = function() {
        const canvas = elements.canvas;
        if (canvas.requestFullscreen) {
          canvas.requestFullscreen();
        } else if (canvas.webkitRequestFullscreen) {
          canvas.webkitRequestFullscreen();
        } else if (canvas.mozRequestFullScreen) {
          canvas.mozRequestFullScreen();
        } else if (canvas.msRequestFullscreen) {
          canvas.msRequestFullscreen();
        }
      };
      
      elements.muteBtn.onclick = function() {
        state.isMuted = !state.isMuted;
        this.textContent = state.isMuted ? '🔊 Unmute' : '🔇 Mute';
        if (state.isMuted) {
          state.audioBuffer = [];
        }
      };
      
      // Keyboard controls
      const keyMap = {
        88: 0, // X - A button
        90: 1, // Z - B button  
        16: 2, // Shift - Select
        13: 3, // Enter - Start
        38: 4, // Up
        40: 5, // Down
        37: 6, // Left
        39: 7  // Right
      };
      
      const keydownHandler = function(e) {
        if (!state.isRunning || !state.nesInstance) return;
        
        if (!elements.container.contains(document.activeElement) && 
            document.activeElement !== elements.canvas) {
          return;
        }
        
        if (keyMap[e.keyCode] !== undefined) {
          state.nesInstance.buttonDown(1, keyMap[e.keyCode]);
          e.preventDefault();
        }
      };
      
      const keyupHandler = function(e) {
        if (!state.isRunning || !state.nesInstance) return;
        
        if (!elements.container.contains(document.activeElement) && 
            document.activeElement !== elements.canvas) {
          return;
        }
        
        if (keyMap[e.keyCode] !== undefined) {
          state.nesInstance.buttonUp(1, keyMap[e.keyCode]);
          e.preventDefault();
        }
      };
      
      elements.canvas.tabIndex = 0;
      elements.canvas.style.outline = 'none';
      elements.canvas.addEventListener('click', function() {
        this.focus();
      });
      elements.canvas.addEventListener('keydown', keydownHandler);
      elements.canvas.addEventListener('keyup', keyupHandler);
    }
    
    // Setup controls for alternative emulator
    function setupAlternativeControls() {
      elements.playBtn.onclick = function() {
        if (state.nesInstance && state.nesInstance.start) {
          state.nesInstance.start();
        }
        state.isRunning = true;
        this.style.display = 'none';
        elements.pauseBtn.style.display = 'inline-block';
      };
      
      elements.pauseBtn.onclick = function() {
        if (state.nesInstance && state.nesInstance.stop) {
          state.nesInstance.stop();
        }
        state.isRunning = false;
        this.style.display = 'none';
        elements.playBtn.style.display = 'inline-block';
      };
      
      elements.resetBtn.onclick = function() {
        if (state.nesInstance && state.nesInstance.reset) {
          state.nesInstance.reset();
        }
      };
    }
    
    // Start emulation
    function startEmulation() {
      if (state.isRunning || !state.nesInstance) return;
      state.isRunning = true;
      
      elements.canvas.focus();
      
      const frameTime = 1000 / 60;
      let lastTime = performance.now();
      
      function frame() {
        if (!state.isRunning) return;
        
        const now = performance.now();
        const elapsed = now - lastTime;
        
        if (elapsed >= frameTime) {
          try {
            state.nesInstance.frame();
          } catch (e) {
            console.error('Frame error:', e);
            stopEmulation();
            showError('Emulation error: ' + e.message, false);
            return;
          }
          lastTime = now - (elapsed %% frameTime);
        }
        
        state.animationId = requestAnimationFrame(frame);
      }
      
      state.animationId = requestAnimationFrame(frame);
    }
    
    // Stop emulation
    function stopEmulation() {
      state.isRunning = false;
      if (state.animationId) {
        cancelAnimationFrame(state.animationId);
        state.animationId = null;
      }
    }
    
    // Cleanup
    window.addEventListener('beforeunload', function() {
      stopEmulation();
      if (state.audioContext) {
        state.audioContext.close();
      }
    });
    
    // Initialize
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', loadJSNES);
    } else {
      loadJSNES();
    }
  })();
  </script>
</div>`, safePlayerID, safeRomPath, safeFileName, fileExt)
}
