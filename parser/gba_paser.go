package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
)

// GenerateGBAPlayerHTML 生成用于在Web上运行GBA游戏的HTML播放器（使用EmulatorJS）
func GenerateGBAPlayerHTML(playerID, romPath, fileName, fileExt string) string {
	// 确保所有用户输入都被正确转义
	safePlayerID := template.HTMLEscapeString(playerID)
	safeRomPath := template.HTMLEscapeString(romPath)
	safeFileName := template.HTMLEscapeString(fileName)
	safeFileExt := template.HTMLEscapeString(fileExt)

	// 生成唯一ID
	h := md5.Sum([]byte(playerID))
	uniqueID := fmt.Sprintf("gba_%x", h)[:12]

	return fmt.Sprintf(`<div class="gba-player-wrapper" id="%[1]s-wrapper">
  <style>
    #%[1]s-wrapper {
      margin: 20px auto;
      padding: 0;
      background: linear-gradient(135deg, #4B0082 0%%, #1a0033 100%%);
      border-radius: 16px;
      max-width: 720px;
      width: 100%%;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      box-shadow: 0 15px 40px rgba(75, 0, 130, 0.6);
      overflow: hidden;
      border: 2px solid #6B46C1;
    }

    /* GBA Style Header */
    #%[1]s-wrapper .gba-header {
      background: linear-gradient(135deg, #6B46C1 0%%, #4B0082 100%%);
      padding: 15px 20px;
      border-bottom: 3px solid #FFD700;
      position: relative;
      overflow: hidden;
    }

    #%[1]s-wrapper .gba-header::before {
      content: '';
      position: absolute;
      top: 0;
      //left: -100%%;
      //width: 200%%;
      //height: 100%%;
      //background: repeating-linear-gradient(
      //  90deg,
      //  transparent,
      //  transparent 10px,
      //  rgba(255, 255, 255, 0.05) 10px,
      //  rgba(255, 255, 255, 0.05) 20px
      //);
      animation: slide 20s linear infinite;
    }

    @keyframes slide {
      0%% { transform: translateX(0); }
      100%% { transform: translateX(50%%); }
    }

    #%[1]s-wrapper .player-header-top {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 20px;
      position: relative;
      z-index: 1;
    }
    
    #%[1]s-wrapper .player-title {
      font-size: 20px;
      font-weight: 600;
      color: #fff;
      display: flex;
      align-items: center;
      gap: 10px;
      text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
    }

    #%[1]s-wrapper .gba-logo {
      background: linear-gradient(90deg, #FFD700 0%%, #FFA500 50%%, #FFD700 100%%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      font-weight: bold;
      font-style: italic;
      letter-spacing: 2px;
    }

    #%[1]s-wrapper .advance-badge {
      background: linear-gradient(135deg, #FFD700 0%%, #FFA500 100%%);
      color: #4B0082;
      padding: 2px 8px;
      border-radius: 12px;
      font-size: 10px;
      font-weight: bold;
      text-transform: uppercase;
      letter-spacing: 1px;
    }
    
    #%[1]s-wrapper .player-controls {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }
    
    #%[1]s-wrapper .player-info {
      background: rgba(0, 0, 0, 0.3);
      padding: 8px 20px;
      display: flex;
      gap: 20px;
      font-size: 13px;
      color: rgba(255, 255, 255, 0.9);
      border-bottom: 1px solid rgba(255, 215, 0, 0.3);
    }

    #%[1]s-wrapper .player-info-item {
      display: flex;
      align-items: center;
      gap: 5px;
    }

    #%[1]s-wrapper .player-info-item span:first-child {
      color: #FFD700;
    }
    
    #%[1]s-wrapper .player-tip {
      font-size: 13px;
      color: rgba(255, 255, 255, 0.9);
      background: linear-gradient(135deg, rgba(107, 70, 193, 0.3) 0%%, rgba(75, 0, 130, 0.3) 100%%);
      padding: 10px 15px;
      border-left: 3px solid #FFD700;
      display: flex;
      align-items: center;
      gap: 8px;
      margin: 15px 20px;
      border-radius: 6px;
    }
    
    #%[1]s-wrapper .player-tip .tip-icon {
      color: #FFD700;
      font-size: 16px;
    }
    
    #%[1]s-wrapper .ctrl-btn {
      background: linear-gradient(135deg, rgba(255, 255, 255, 0.2) 0%%, rgba(255, 255, 255, 0.1) 100%%);
      color: #fff;
      border: 1px solid rgba(255, 215, 0, 0.4);
      padding: 6px 12px;
      border-radius: 20px;
      cursor: pointer;
      font-size: 13px;
      transition: all 0.2s;
      text-decoration: none;
      display: inline-flex;
      align-items: center;
      gap: 5px;
      text-transform: uppercase;
      font-weight: 500;
      letter-spacing: 0.5px;
    }
    
    #%[1]s-wrapper .ctrl-btn:hover {
      background: linear-gradient(135deg, rgba(255, 215, 0, 0.3) 0%%, rgba(255, 165, 0, 0.3) 100%%);
      transform: translateY(-1px);
      border-color: #FFD700;
    }
    
    #%[1]s-wrapper .ctrl-btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    
    #%[1]s-wrapper .ctrl-btn.success {
      background: linear-gradient(135deg, rgba(76, 175, 80, 0.3) 0%%, rgba(56, 142, 60, 0.3) 100%%);
      border-color: rgba(76, 175, 80, 0.6);
    }

    #%[1]s-wrapper .ctrl-btn.primary {
      background: linear-gradient(135deg, rgba(107, 70, 193, 0.4) 0%%, rgba(75, 0, 130, 0.4) 100%%);
      border-color: rgba(107, 70, 193, 0.6);
    }
    
    /* EmulatorJS Container */
    #%[1]s-wrapper .emulator-container {
      background: #000;
      border-radius: 8px;
      padding: 0;
      display: block;
      position: relative;
      max-width: 640px;
      margin: 20px auto;
      aspect-ratio: 3/2;
      box-shadow: 
        0 10px 30px rgba(0, 0, 0, 0.5),
        inset 0 0 20px rgba(107, 70, 193, 0.2);
    }
    
    #%[1]s-wrapper #game {
      width: 100%%;
      height: 100%%;
      max-width: 100%%;
      border-radius: 8px;
    }
    
    /* Loading overlay */
    #%[1]s-wrapper .loading-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.95);
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      z-index: 100;
      border-radius: 8px;
    }
    
    #%[1]s-wrapper .loading-overlay.hidden {
      display: none;
    }
    
    #%[1]s-wrapper .loading-spinner {
      width: 60px;
      height: 60px;
      border: 4px solid rgba(107, 70, 193, 0.3);
      border-top-color: #FFD700;
      border-radius: 50%%;
      animation: spin 1s linear infinite;
    }
    
    @keyframes spin {
      to { transform: rotate(360deg); }
    }
    
    #%[1]s-wrapper .loading-text {
      color: #FFD700;
      margin-top: 20px;
      font-size: 14px;
      text-transform: uppercase;
      letter-spacing: 2px;
    }

    /* Control Instructions */
    #%[1]s-wrapper .controls-info {
      margin: 20px;
      padding: 15px;
      background: linear-gradient(135deg, rgba(107, 70, 193, 0.2) 0%%, rgba(75, 0, 130, 0.2) 100%%);
      border-radius: 8px;
      border: 1px solid rgba(255, 215, 0, 0.3);
      color: rgba(255, 255, 255, 0.9);
      font-size: 13px;
    }

    #%[1]s-wrapper .controls-info h4 {
      margin: 0 0 10px 0;
      color: #FFD700;
      font-size: 14px;
      text-transform: uppercase;
      letter-spacing: 1px;
    }

    #%[1]s-wrapper .controls-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
      gap: 8px;
      margin-top: 10px;
    }

    #%[1]s-wrapper .control-item {
      display: flex;
      align-items: center;
      gap: 8px;
      background: rgba(0, 0, 0, 0.3);
      padding: 5px 10px;
      border-radius: 6px;
    }

    #%[1]s-wrapper .key {
      background: linear-gradient(135deg, #6B46C1 0%%, #4B0082 100%%);
      border: 1px solid #FFD700;
      padding: 2px 6px;
      border-radius: 4px;
      font-family: monospace;
      font-size: 11px;
      color: #fff;
      font-weight: bold;
      min-width: 40px;
      text-align: center;
    }
    
    /* Error message */
    #%[1]s-wrapper .error-message {
      background: linear-gradient(135deg, rgba(255, 0, 0, 0.2) 0%%, rgba(139, 0, 0, 0.2) 100%%);
      border: 1px solid rgba(255, 0, 0, 0.4);
      color: #ff6b6b;
      padding: 10px 15px;
      border-radius: 6px;
      margin: 10px 20px;
      display: none;
    }
    
    #%[1]s-wrapper .error-message.show {
      display: block;
    }

    /* Game options panel */
    #%[1]s-wrapper .game-options {
      display: flex;
      gap: 15px;
      margin: 15px 20px;
      padding: 12px;
      background: linear-gradient(135deg, rgba(0, 0, 0, 0.4) 0%%, rgba(0, 0, 0, 0.3) 100%%);
      border-radius: 8px;
      flex-wrap: wrap;
      align-items: center;
      border: 1px solid rgba(255, 215, 0, 0.2);
    }

    #%[1]s-wrapper .option-group {
      display: flex;
      align-items: center;
      gap: 8px;
      color: rgba(255, 255, 255, 0.9);
      font-size: 13px;
    }

    #%[1]s-wrapper .option-group label {
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 5px;
    }

    #%[1]s-wrapper .option-group input[type="checkbox"] {
      cursor: pointer;
      accent-color: #FFD700;
    }

    #%[1]s-wrapper .option-group input[type="range"] {
      width: 100px;
      accent-color: #FFD700;
    }
    
    @media (max-width: 768px) {
      #%[1]s-wrapper {
        margin: 10px;
        border-radius: 12px;
      }
      
      #%[1]s-wrapper .emulator-container {
        margin: 10px;
        max-width: calc(100%% - 20px);
      }
      
      #%[1]s-wrapper .player-header-top {
        flex-direction: column;
        align-items: flex-start;
      }
      
      #%[1]s-wrapper .controls-grid {
        grid-template-columns: 1fr;
      }

      #%[1]s-wrapper .player-controls {
        width: 100%%;
        justify-content: flex-start;
      }

      #%[1]s-wrapper .player-info {
        flex-wrap: wrap;
        gap: 10px;
      }
    }
  </style>
  
  <div class="gba-header">
    <div class="player-header-top">
      <div class="player-title">
        <span class="gba-logo">GAME BOY</span>
        <span class="advance-badge">ADVANCE</span>
        <span style="color: #FFD700;">•</span>
        <span>%[3]s</span>
      </div>
      <div class="player-controls">
        <button class="ctrl-btn" id="%[4]s-status-btn">
          <span>⏳</span>
          <span>Loading...</span>
        </button>
        <button class="ctrl-btn" id="%[4]s-fullscreen-btn">
          <span>⛶</span>
          <span>Full</span>
        </button>
        <a href="%[2]s" download class="ctrl-btn">
          <span>⬇</span>
          <span>ROM</span>
        </a>
      </div>
    </div>
  </div>
  
  <div class="player-info">
    <div class="player-info-item">
      <span>📦</span>
      <span>Format: %[5]s</span>
    </div>
    <div class="player-info-item">
      <span>🎮</span>
      <span>System: Game Boy Advance</span>
    </div>
    <div class="player-info-item">
      <span>🏷️</span>
      <span>Core: mGBA</span>
    </div>
  </div>

  <div class="player-tip">
    <span class="tip-icon">💡</span>
    <span>提示：按 Tab 键打开模拟器菜单。支持手柄控制和键盘映射自定义。</span>
  </div>
  
  <div class="error-message" id="%[4]s-error"></div>
  
  <!-- EmulatorJS Container -->
  <div class="emulator-container" id="%[4]s-container">
    <div class="loading-overlay" id="%[4]s-loading">
      <div class="loading-spinner"></div>
      <div class="loading-text">Initializing GBA...</div>
    </div>
    <div id="game"></div>
  </div>

  <!-- Game Options -->
  <div class="game-options">
    <div class="option-group">
      <label>
        <input type="checkbox" id="%[4]s-turbo" />
        <span>⚡ Turbo Mode</span>
      </label>
    </div>
    <div class="option-group">
      <label>
        <span>🔊 Volume:</span>
        <input type="range" id="%[4]s-volume" min="0" max="100" value="50" />
      </label>
    </div>
    <div class="option-group">
      <label>
        <input type="checkbox" id="%[4]s-lcd" />
        <span>📱 LCD Filter</span>
      </label>
    </div>
  </div>

  <!-- Controls Information -->
  <div class="controls-info">
    <h4>🎮 Keyboard Controls</h4>
    <div class="controls-grid">
      <div class="control-item">
        <span class="key">WASD</span> D-Pad
      </div>
      <div class="control-item">
        <span class="key">U</span> A Button
      </div>
      <div class="control-item">
        <span class="key">I</span> B Button
      </div>
      <div class="control-item">
        <span class="key">[</span> L Button
      </div>
      <div class="control-item">
        <span class="key">]</span> R Button
      </div>
      <div class="control-item">
        <span class="key">K</span> Start
      </div>
      <div class="control-item">
        <span class="key">J</span> Select
      </div>
      <div class="control-item">
        <span class="key">Tab</span> Menu
      </div>
      <div class="control-item">
        <span class="key">1/3</span> Save/Load
      </div>
    </div>
    <p style="margin-top: 10px; color: rgba(255,255,255,0.7); font-size: 12px;">
      💡 支持 USB/蓝牙游戏手柄。GBA 游戏支持存档和即时存档功能。
    </p>
  </div>
  
  <!-- EmulatorJS Scripts -->
  <script>
    (function() {
      const ROM_URL = '%[2]s';
      const ROM_NAME = '%[3]s';
      const UNIQUE_ID = '%[4]s';
      const FILE_EXT = '%[5]s'.toLowerCase();
      
      let emulatorReady = false;
      let volumeControl = null;
      let turboMode = false;
      let lcdFilter = false;
      
      // Status update function
      function updateStatus(message, type = 'loading') {
        const statusBtn = document.getElementById(UNIQUE_ID + '-status-btn');
        if (statusBtn) {
          const icons = {
            loading: '⏳',
            success: '✅',
            error: '❌',
            running: '🎮',
            paused: '⏸️'
          };
          statusBtn.innerHTML = '<span>' + (icons[type] || '⏳') + '</span><span>' + message + '</span>';
          statusBtn.className = 'ctrl-btn';
          if (type === 'success' || type === 'running') {
            statusBtn.classList.add('success');
          } else if (type === 'loading') {
            statusBtn.classList.add('primary');
          }
        }
      }
      
      // Show error message
      function showError(message) {
        const errorEl = document.getElementById(UNIQUE_ID + '-error');
        if (errorEl) {
          errorEl.textContent = '⚠️ ' + message;
          errorEl.classList.add('show');
        }
        updateStatus('Error', 'error');
      }
      
      // Hide loading overlay
      function hideLoading() {
        const loadingEl = document.getElementById(UNIQUE_ID + '-loading');
        if (loadingEl) {
          loadingEl.classList.add('hidden');
        }
      }
      
      // Update loading text
      function updateLoadingText(text) {
        const loadingTextEl = document.querySelector('#' + UNIQUE_ID + '-loading .loading-text');
        if (loadingTextEl) {
          loadingTextEl.textContent = text;
        }
      }
      
      // Initialize EmulatorJS
      function initializeEmulator() {
        try {
          updateStatus('Initializing...', 'loading');
	
          // 设置键盘映射
          updateStatus('Setting Controls...', 'loading');
          window.EJS_defaultControls = {
            0: {
              // WASD 方向
              4: { 'value': 'w', 'value2': 'DPAD_UP' },
              5: { 'value': 's', 'value2': 'DPAD_DOWN' },
              6: { 'value': 'a', 'value2': 'DPAD_LEFT' },
              7: { 'value': 'd', 'value2': 'DPAD_RIGHT' },
              
              // UI 动作键 (GBA布局)
              8: { 'value': 'u', 'value2': 'BUTTON_1' },    // A
              0: { 'value': 'i', 'value2': 'BUTTON_2' },    // B
              
              // JK 系统键
              2: { 'value': 'j', 'value2': 'SELECT' },
              3: { 'value': 'k', 'value2': 'START' },
              
              // [] 肩键
              10: { 'value': '[', 'value2': 'LEFT_TOP_SHOULDER' },   // L
              11: { 'value': ']', 'value2': 'RIGHT_TOP_SHOULDER' },  // R
              
              // 功能键
              24: { 'value': '1' },      // Quick Save
              25: { 'value': '3' },      // Quick Load  
              27: { 'value': 'space' }   // Fast Forward
            },
            1: {}, 2: {}, 3: {}
          };

          updateLoadingText('Setting up GBA emulator...');
          
          // EmulatorJS configuration
          window.EJS_player = '#game';
          window.EJS_gameUrl = ROM_URL;
          window.EJS_core = 'gba'; // GBA core
          window.EJS_gameName = ROM_NAME;
          
          // Optional configurations
          window.EJS_pathtodata = 'https://cdn.emulatorjs.org/stable/data/';
          window.EJS_startOnLoaded = true;
          window.EJS_biosUrl = ''; // GBA BIOS (optional, but recommended)
          window.EJS_DEBUG_XX = false;
          
          // Color and theme
          window.EJS_color = '#6B46C1';
          
          // Features
          window.EJS_cheats = true;
          window.EJS_saveStateName = ROM_NAME + '.state';
          
          // Language
          window.EJS_language = 'en-US';
          
          // Settings menu
          window.EJS_settingsMenu = true;
          
          // Screen settings (GBA native resolution)
          window.EJS_aspectRatio = '3:2';
          
          // Virtual gamepad for mobile
          window.EJS_VirtualGamepadSettings = {
            scale: 1.5,
            opacity: 0.8,
            enabled: true
          };
          
          // Callbacks
          window.EJS_onGameStart = function() {
            console.log('GBA game started:', ROM_NAME);
			// try check if game loaded successfully after 0.5 seconds
			setTimeout(() => {
				crt_shader_init()
			}, 1000);
            updateStatus('Running', 'running');
            hideLoading();
            emulatorReady = true;
            
            // Initialize custom controls
            initializeCustomControls();
          };
          
          window.EJS_onLoadState = function() {
            console.log('State loaded');
            updateStatus('State Loaded', 'success');
            setTimeout(() => updateStatus('Running', 'running'), 2000);
          };
          
          window.EJS_onSaveState = function() {
            console.log('State saved');
            updateStatus('State Saved', 'success');
            setTimeout(() => updateStatus('Running', 'running'), 2000);
          };
          
          // Setup button handlers
          setupButtonHandlers();
          
          // Setup option handlers
          setupOptionHandlers();
          
          // Load EmulatorJS loader script
          updateLoadingText('Loading EmulatorJS core...');
          const script = document.createElement('script');
          script.src = '/emulatorjs/loader.js';
          script.onerror = function() {
            showError('Failed to load EmulatorJS. Please check if EmulatorJS is properly installed.');
            hideLoading();
          };
          script.onload = function() {
            console.log('EmulatorJS loader loaded successfully');
            updateLoadingText('Starting GBA emulator...');
            updateStatus('Starting...', 'loading');
          };
          document.body.appendChild(script);
          
        } catch (error) {
          console.error('Failed to initialize emulator:', error);
          showError('Failed to initialize emulator: ' + error.message);
          hideLoading();
        }
      }
      
      // Initialize custom controls after emulator starts
      function initializeCustomControls() {
        // Get volume control reference
        volumeControl = document.getElementById(UNIQUE_ID + '-volume');
        if (volumeControl && window.EJS_emulator) {
          volumeControl.addEventListener('input', function() {
            const volume = this.value / 100;
            if (window.EJS_emulator.setVolume) {
              window.EJS_emulator.setVolume(volume);
            }
          });
        }
      }
      
      // Setup button handlers
      function setupButtonHandlers() {
        // Fullscreen button
        const fullscreenBtn = document.getElementById(UNIQUE_ID + '-fullscreen-btn');
        if (fullscreenBtn) {
          fullscreenBtn.addEventListener('click', function() {
            const container = document.getElementById(UNIQUE_ID + '-container');
            if (container) {
              if (container.requestFullscreen) {
                container.requestFullscreen();
              } else if (container.webkitRequestFullscreen) {
                container.webkitRequestFullscreen();
              } else if (container.msRequestFullscreen) {
                container.msRequestFullscreen();
              }
            }
          });
        }
      }
      
      // Setup game option handlers
      function setupOptionHandlers() {
        // Turbo mode toggle
        const turboCheckbox = document.getElementById(UNIQUE_ID + '-turbo');
        if (turboCheckbox) {
          turboCheckbox.addEventListener('change', function() {
            turboMode = this.checked;
            if (window.EJS_emulator && window.EJS_emulator.setSpeed) {
              window.EJS_emulator.setSpeed(turboMode ? 2.0 : 1.0);
            }
          });
        }
        
        // LCD filter toggle
        const lcdCheckbox = document.getElementById(UNIQUE_ID + '-lcd');
        if (lcdCheckbox) {
          lcdCheckbox.addEventListener('change', function() {
            lcdFilter = this.checked;
            if (window.EJS_emulator && window.EJS_emulator.setShader) {
              window.EJS_emulator.setShader(lcdFilter ? 'lcd' : 'none');
            }
          });
        }
      }
      
      // Check ROM format
      function checkRomFormat() {
        const validFormats = ['gba', 'agb', 'bin', 'mb', 'gbc', 'gb'];
        if (!validFormats.includes(FILE_EXT)) {
          console.warn('Unusual file format for GBA:', FILE_EXT);
          return true; // Still try to load
        }
        return true;
      }
      
      // Initialize when DOM is ready
      if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
          if (!checkRomFormat()) {
            showError('Warning: Unusual file format. Expected: .gba, .agb');
          }
          initializeEmulator();
        });
      } else {
        if (!checkRomFormat()) {
          showError('Warning: Unusual file format. Expected: .gba, .agb');
        }
        initializeEmulator();
      }
    })();
  </script>
</div>`,
		safePlayerID, // %[1]s - Player ID for wrapper
		safeRomPath,  // %[2]s - ROM file URL path
		safeFileName, // %[3]s - Display filename with extension
		uniqueID,     // %[4]s - Unique ID for elements
		safeFileExt,  // %[5]s - File extension
	)
}
