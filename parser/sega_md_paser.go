package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
)

// GenerateSegaMDPlayerHTML 生成用于在Web上运行Sega Mega Drive游戏的HTML播放器
func GenerateSegaMDPlayerHTML(playerID, romPath, fileName, fileExt string) string {
	// 确保所有用户输入都被正确转义
	safePlayerID := template.HTMLEscapeString(playerID)
	safeRomPath := template.HTMLEscapeString(romPath)
	safeFileName := template.HTMLEscapeString(fileName)
	safeFileExt := template.HTMLEscapeString(fileExt)
	log.Println("Generating Sega Mega Drive player for ROM:", safeRomPath)
	log.Println("Generating Sega Mega Drive player for FileName:", safeFileName)
	// 生成唯一ID
	h := md5.Sum([]byte(playerID))
	uniqueID := fmt.Sprintf("smd_%x", h)[:12]

	return fmt.Sprintf(`<div class="sega-md-player-wrapper" id="%[1]s-wrapper">
  <style>
    #%[1]s-wrapper {
      margin: 20px auto;
      padding: 0;
      background: #000;
      border-radius: 16px;
      max-width: 800px;
      width: 100%%;
      font-family: 'Segoe UI', 'Arial', sans-serif;
      box-shadow: 0 20px 60px rgba(0, 0, 0, 0.8);
      overflow: hidden;
      border: 3px solid #1a1a1a;
    }

    /* Sega Header Style */
    #%[1]s-wrapper .sega-header {
      background: linear-gradient(135deg, #0050a0 0%%, #003070 50%%, #0050a0 100%%);
      padding: 15px 20px;
      border-bottom: 4px solid #ffd700;
      position: relative;
      overflow: hidden;
    }

    #%[1]s-wrapper .sega-header::before {
      content: '';
      position: absolute;
      top: 0;
      left: -100%%;
      width: 100%%;
      height: 100%%;
      background: linear-gradient(90deg, transparent, rgba(255,255,255,0.1), transparent);
      animation: shimmer 3s infinite;
    }

    @keyframes shimmer {
      0%% { left: -100%%; }
      100%% { left: 100%%; }
    }

    #%[1]s-wrapper .sega-logo-container {
      display: flex;
      align-items: center;
      justify-content: space-between;
      position: relative;
      z-index: 1;
    }

    #%[1]s-wrapper .sega-logo {
      font-size: 28px;
      font-weight: 900;
      color: #fff;
      text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
      letter-spacing: 2px;
      font-style: italic;
      display: flex;
      align-items: center;
      gap: 10px;
    }

    #%[1]s-wrapper .sega-logo .logo-text {
      background: linear-gradient(180deg, #fff 0%%, #e0e0e0 100%%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    #%[1]s-wrapper .mega-drive-badge {
      background: #ffd700;
      color: #003070;
      padding: 4px 12px;
      border-radius: 20px;
      font-size: 11px;
      font-weight: bold;
      text-transform: uppercase;
      letter-spacing: 1px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.3);
    }

    /* Game Info Bar */
    #%[1]s-wrapper .game-info-bar {
      background: linear-gradient(90deg, #1a1a1a 0%%, #2a2a2a 100%%);
      padding: 10px 20px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      border-bottom: 2px solid #333;
    }

    #%[1]s-wrapper .game-title {
      color: #ffd700;
      font-size: 16px;
      font-weight: 600;
      text-shadow: 1px 1px 2px rgba(0,0,0,0.8);
      display: flex;
      align-items: center;
      gap: 8px;
    }

    #%[1]s-wrapper .game-title::before {
      content: '🎮';
      font-size: 18px;
    }

    #%[1]s-wrapper .control-buttons {
      display: flex;
      gap: 8px;
    }

    #%[1]s-wrapper .ctrl-btn {
      background: linear-gradient(135deg, #333 0%%, #222 100%%);
      color: #fff;
      border: 1px solid #444;
      padding: 6px 12px;
      border-radius: 20px;
      cursor: pointer;
      font-size: 12px;
      font-weight: 500;
      transition: all 0.3s;
      text-decoration: none;
      display: inline-flex;
      align-items: center;
      gap: 5px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.3);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    #%[1]s-wrapper .ctrl-btn:hover {
      background: linear-gradient(135deg, #0050a0 0%%, #003070 100%%);
      border-color: #0080ff;
      transform: translateY(-2px);
      box-shadow: 0 4px 8px rgba(0,80,160,0.4);
    }

    #%[1]s-wrapper .ctrl-btn:active {
      transform: translateY(0);
    }

    #%[1]s-wrapper .ctrl-btn.primary {
      background: linear-gradient(135deg, #0050a0 0%%, #003070 100%%);
      border-color: #0080ff;
    }

    #%[1]s-wrapper .ctrl-btn.success {
      background: linear-gradient(135deg, #4CAF50 0%%, #388E3C 100%%);
      border-color: #66BB6A;
    }

    /* Emulator Container with CRT TV Style */
    #%[1]s-wrapper .emulator-outer {
      background: linear-gradient(135deg, #2a2a2a 0%%, #1a1a1a 100%%);
      padding: 20px;
      position: relative;
    }

    #%[1]s-wrapper .tv-frame {
      background: #111;
      border-radius: 12px;
      padding: 15px;
      box-shadow: 
        inset 0 0 20px rgba(0,0,0,0.8),
        0 0 30px rgba(0,80,160,0.2);
      position: relative;
    }

    #%[1]s-wrapper .screen-bezel {
      background: #000;
      border-radius: 8px;
      padding: 2px;
      box-shadow: inset 0 0 10px rgba(0,0,0,0.9);
    }

    #%[1]s-wrapper .emulator-container {
      background: #000;
      border-radius: 6px;
      position: relative;
      width: 100%%;
      max-width: 640px;
      margin: 0 auto;
      aspect-ratio: 4/3;
    }

    #%[1]s-wrapper #game {
      width: 100%%;
      height: 100%%;
      border-radius: 6px;
    }

    /* Power LED indicator */
    #%[1]s-wrapper .power-led {
      position: absolute;
      top: 10px;
      right: 10px;
      width: 8px;
      height: 8px;
      background: #00ff00;
      border-radius: 50%%;
      box-shadow: 0 0 10px #00ff00;
      animation: pulse 2s infinite;
    }

    @keyframes pulse {
      0%%, 100%% { opacity: 1; }
      50%% { opacity: 0.5; }
    }

    /* Loading Overlay */
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
      border-radius: 6px;
    }

    #%[1]s-wrapper .loading-overlay.hidden {
      display: none;
    }

    #%[1]s-wrapper .sega-loader {
      width: 80px;
      height: 80px;
      position: relative;
    }

    #%[1]s-wrapper .sega-loader::before {
      content: 'SEGA';
      position: absolute;
      top: 50%%;
      left: 50%%;
      transform: translate(-50%%, -50%%);
      font-size: 20px;
      font-weight: bold;
      color: #0080ff;
      animation: fadeInOut 1.5s infinite;
    }

    #%[1]s-wrapper .sega-loader::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      border: 3px solid rgba(0,128,255,0.2);
      border-top-color: #0080ff;
      border-radius: 50%%;
      animation: spin 1s linear infinite;
    }

    @keyframes fadeInOut {
      0%%, 100%% { opacity: 0.3; }
      50%% { opacity: 1; }
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    #%[1]s-wrapper .loading-text {
      color: #0080ff;
      margin-top: 100px;
      font-size: 14px;
      text-transform: uppercase;
      letter-spacing: 2px;
    }

    /* System Info Strip */
    #%[1]s-wrapper .system-info {
      background: #1a1a1a;
      padding: 8px 20px;
      display: flex;
      justify-content: center;
      gap: 30px;
      border-top: 2px solid #333;
      font-size: 11px;
      color: #888;
      text-transform: uppercase;
      letter-spacing: 1px;
    }

    #%[1]s-wrapper .info-item {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    #%[1]s-wrapper .info-item .label {
      color: #666;
    }

    #%[1]s-wrapper .info-item .value {
      color: #0080ff;
      font-weight: bold;
    }

    /* Controls Guide */
    #%[1]s-wrapper .controls-panel {
      background: linear-gradient(135deg, #1a1a1a 0%%, #0a0a0a 100%%);
      padding: 15px 20px;
      border-top: 2px solid #0050a0;
    }

    #%[1]s-wrapper .controls-header {
      color: #ffd700;
      font-size: 13px;
      font-weight: bold;
      margin-bottom: 10px;
      text-transform: uppercase;
      letter-spacing: 1px;
      display: flex;
      align-items: center;
      gap: 8px;
    }

    #%[1]s-wrapper .controls-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
      gap: 8px;
    }

    #%[1]s-wrapper .control-item {
      background: rgba(0,80,160,0.1);
      padding: 6px 10px;
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      border: 1px solid rgba(0,80,160,0.2);
      transition: all 0.2s;
    }

    #%[1]s-wrapper .control-item:hover {
      background: rgba(0,80,160,0.2);
      border-color: rgba(0,80,160,0.4);
    }

    #%[1]s-wrapper .control-action {
      color: #aaa;
      font-size: 11px;
    }

    #%[1]s-wrapper .control-key {
      background: #0050a0;
      color: #fff;
      padding: 2px 6px;
      border-radius: 4px;
      font-family: monospace;
      font-size: 11px;
      font-weight: bold;
      box-shadow: 0 1px 3px rgba(0,0,0,0.3);
    }

    /* Error Message */
    #%[1]s-wrapper .error-message {
      background: linear-gradient(135deg, rgba(255,0,0,0.2) 0%%, rgba(200,0,0,0.2) 100%%);
      border: 1px solid #ff0000;
      color: #ff6b6b;
      padding: 10px 15px;
      border-radius: 6px;
      margin: 10px 20px;
      display: none;
      text-align: center;
    }

    #%[1]s-wrapper .error-message.show {
      display: block;
    }

    /* Responsive Design */
    @media (max-width: 768px) {
      #%[1]s-wrapper {
        margin: 10px;
        border-radius: 12px;
      }

      #%[1]s-wrapper .sega-logo {
        font-size: 20px;
      }

      #%[1]s-wrapper .control-buttons {
        flex-wrap: wrap;
      }

      #%[1]s-wrapper .ctrl-btn {
        padding: 5px 10px;
        font-size: 11px;
      }

      #%[1]s-wrapper .game-info-bar {
        flex-direction: column;
        gap: 10px;
        align-items: flex-start;
      }

      #%[1]s-wrapper .controls-grid {
        grid-template-columns: 1fr;
      }

      #%[1]s-wrapper .system-info {
        flex-wrap: wrap;
        gap: 15px;
        justify-content: flex-start;
      }
    }

    /* Special Genesis/MD styling */
    #%[1]s-wrapper .genesis-style {
      background: repeating-linear-gradient(
        90deg,
        #0050a0,
        #0050a0 10px,
        #003070 10px,
        #003070 20px
      );
      height: 4px;
      margin: 0;
    }
  </style>
  
  <!-- Sega Header -->
  <div class="sega-header">
    <div class="sega-logo-container">
      <div class="sega-logo">
        <span class="logo-text">SEGA</span>
        <span class="mega-drive-badge">MEGA DRIVE</span>
      </div>
    </div>
  </div>
  
  <!-- Genesis Style Divider -->
  <div class="genesis-style"></div>
  
  <!-- Game Info Bar -->
  <div class="game-info-bar">
    <div class="game-title">%[3]s</div>
    <div class="control-buttons">
      <button class="ctrl-btn" id="%[4]s-status-btn">
        <span>⏳</span>
        <span>Loading</span>
      </button>
      <button class="ctrl-btn" id="%[4]s-fullscreen-btn" title="Fullscreen">
        <span>⛶</span>
      </button>
      <!-- <button class="ctrl-btn" id="%[4]s-save-btn" title="Save State">
        <span>💾</span>
      </button>
      <button class="ctrl-btn" id="%[4]s-load-btn" title="Load State">
        <span>📂</span>
      </button>
      <button class="ctrl-btn" id="%[4]s-reset-btn" title="Reset">
        <span>🔄</span>
      </button> -->
      <a href="%[2]s" download class="ctrl-btn" title="Download ROM">
        <span>⬇</span>
      </a>
    </div>
  </div>
  
  <div class="error-message" id="%[4]s-error"></div>
  
  <!-- Emulator Container with TV Style -->
  <div class="emulator-outer">
    <div class="tv-frame">
      <div class="power-led"></div>
      <div class="screen-bezel">
        <div class="emulator-container" id="%[4]s-container">
          <div class="loading-overlay" id="%[4]s-loading">
            <div class="sega-loader"></div>
            <div class="loading-text">LOADING...</div>
          </div>
          <div id="game"></div>
        </div>
      </div>
    </div>
  </div>
  
  <!-- System Info Strip -->
  <div class="system-info">
    <div class="info-item">
      <span class="label">System:</span>
      <span class="value">16-BIT</span>
    </div>
    <div class="info-item">
      <span class="label">Core:</span>
      <span class="value">Genesis Plus GX</span>
    </div>
    <div class="info-item">
      <span class="label">Format:</span>
      <span class="value">%[5]s</span>
    </div>
    <div class="info-item">
      <span class="label">Region:</span>
      <span class="value">NTSC/PAL</span>
    </div>
  </div>
  
  <!-- Controls Panel -->
  <div class="controls-panel">
    <div class="controls-header">
      <span>🎮</span>
      <span>Controller Mapping</span>
    </div>
    <div class="controls-grid">
      <div class="control-item">
        <span class="control-action">D-Pad</span>
        <span class="control-key">WSAD</span>
      </div>
      <div class="control-item">
        <span class="control-action">A</span>
        <span class="control-key">U</span>
      </div>
      <div class="control-item">
        <span class="control-action">B</span>
        <span class="control-key">I</span>
      </div>
      <div class="control-item">
        <span class="control-action">X</span>
        <span class="control-key">O</span>
      </div>
      <div class="control-item">
        <span class="control-action">Y</span>
        <span class="control-key">P</span>
      </div>
      <div class="control-item">
        <span class="control-action">Select</span>
        <span class="control-key">J</span>
      </div>
      <div class="control-item">
        <span class="control-action">Start</span>
        <span class="control-key">K</span>
      </div>
      <div class="control-item">
        <span class="control-action">Menu</span>
        <span class="control-key">Tab</span>
      </div>
    </div>
  </div>
  
  <!-- EmulatorJS Scripts -->
  <script>
    (function() {
      const ROM_URL = '%[2]s';
      const ROM_NAME = '%[3]s';
      const UNIQUE_ID = '%[4]s';
      const FILE_EXT = '%[5]s'.toLowerCase();
      
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
          const shortMessages = {
            'Loading...': 'Loading',
            'Running': 'Playing',
            'Error': 'Error',
            'State Saved': 'Saved',
            'State Loaded': 'Loaded'
          };
          const displayMessage = shortMessages[message] || message;
          statusBtn.innerHTML = '<span>' + (icons[type] || '⏳') + '</span><span>' + displayMessage + '</span>';
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
          setTimeout(() => {
            loadingEl.classList.add('hidden');
          }, 500);
        }
      }
      
      // Update loading text
      function updateLoadingText(text) {
        const loadingTextEl = document.querySelector('#' + UNIQUE_ID + '-loading .loading-text');
        if (loadingTextEl) {
          loadingTextEl.textContent = text.toUpperCase();
        }
      }
      
      // Initialize EmulatorJS
      function initializeEmulator() {
        try {
		// 设置键盘映射
		updateStatus('Setting Keyboard...', 'loading');
    	window.EJS_defaultControls = {
			0: {
				// WASD 方向
				4: { 'value': 'w', 'value2': 'DPAD_UP' },
				5: { 'value': 's', 'value2': 'DPAD_DOWN' },
				6: { 'value': 'a', 'value2': 'DPAD_LEFT' },
				7: { 'value': 'd', 'value2': 'DPAD_RIGHT' },
				
				// JKL; 动作键
				0: { 'value': 'i', 'value2': 'BUTTON_2' },    // B
				8: { 'value': 'u', 'value2': 'BUTTON_1' },    // A
				1: { 'value': 'o', 'value2': 'BUTTON_4' },    // Y
				9: { 'value': 'p', 'value2': 'BUTTON_3' },    // X
				
				// UI 系统键
				2: { 'value': 'j', 'value2': 'SELECT' },
				3: { 'value': 'k', 'value2': 'START' },
				
				// OP 肩键
				10: { 'value': '[', 'value2': 'LEFT_TOP_SHOULDER' },
				11: { 'value': ']', 'value2': 'RIGHT_TOP_SHOULDER' },
				
				// 功能键
				24: { 'value': '1' },
				25: { 'value': '3' },
				27: { 'value': 'space' }
			},
			1: {}, 2: {}, 3: {}
   	 	 };
          updateStatus('Loading', 'loading');
          updateLoadingText('Initializing');
          
          // EmulatorJS configuration
          window.EJS_player = '#game';
          window.EJS_gameUrl = ROM_URL;
          window.EJS_core = 'segaMD';
          window.EJS_gameName = ROM_NAME;
          window.EJS_language = 'zh-CN'; //language

          // Paths
          window.EJS_pathtodata = 'https://cdn.emulatorjs.org/4.2.3/data/';
          window.EJS_startOnLoaded = true;
          
          // Visual settings
          window.EJS_color = '#0050a0';
          
          // Features
          window.EJS_multitap = false;
          window.EJS_cheats = true;
          window.EJS_saveStateName = ROM_NAME.replace(/[^a-zA-Z0-9]/g, '_') + '.state';
          
          // Settings
          window.EJS_language = 'en-US';
          window.EJS_settingsMenu = true;
          
          // Callbacks
          window.EJS_onGameStart = function() {
            console.log('Sega MD Game started:', ROM_NAME);
            updateStatus('Playing', 'running');
            hideLoading();
            // Animate power LED
            const led = document.querySelector('#' + UNIQUE_ID + '-wrapper .power-led');
            if (led) {
              led.style.background = '#00ff00';
              led.style.boxShadow = '0 0 15px #00ff00';
            }
          };
          
          window.EJS_onLoadState = function() {
            updateStatus('Loaded', 'success');
            setTimeout(() => updateStatus('Playing', 'running'), 2000);
          };
          
          window.EJS_onSaveState = function() {
            updateStatus('Saved', 'success');
            setTimeout(() => updateStatus('Playing', 'running'), 2000);
          };
          
          // Setup button handlers
          setupButtonHandlers();
          
          // Load EmulatorJS
          updateLoadingText('Loading Core');
          const script = document.createElement('script');
          script.src = '/emulatorjs/loader.js';
          script.onerror = function() {
            showError('Failed to load EmulatorJS. Please check installation.');
            hideLoading();
          };
          script.onload = function() {
            console.log('EmulatorJS loaded');
            updateLoadingText('Starting Game');
            updateStatus('Starting', 'loading');
          };
          document.body.appendChild(script);
          
        } catch (error) {
          console.error('Failed to initialize:', error);
          showError('Initialization failed: ' + error.message);
          hideLoading();
        }
      }
      
      // Setup button handlers
      function setupButtonHandlers() {
        // Fullscreen
        const fullscreenBtn = document.getElementById(UNIQUE_ID + '-fullscreen-btn');
        if (fullscreenBtn) {
          fullscreenBtn.addEventListener('click', function() {
            const container = document.getElementById(UNIQUE_ID + '-container');
            if (container) {
              if (container.requestFullscreen) {
                container.requestFullscreen();
              } else if (container.webkitRequestFullscreen) {
                container.webkitRequestFullscreen();
              }
            }
          });
        }
        
        // Save state
        const saveBtn = document.getElementById(UNIQUE_ID + '-save-btn');
        if (saveBtn) {
          saveBtn.addEventListener('click', function() {
            if (window.EJS_emulator && window.EJS_emulator.saveState) {
              window.EJS_emulator.saveState();
            }
          });
        }
        
        // Load state
        const loadBtn = document.getElementById(UNIQUE_ID + '-load-btn');
        if (loadBtn) {
          loadBtn.addEventListener('click', function() {
            if (window.EJS_emulator && window.EJS_emulator.loadState) {
              window.EJS_emulator.loadState();
            }
          });
        }
        
        // Reset
        const resetBtn = document.getElementById(UNIQUE_ID + '-reset-btn');
        if (resetBtn) {
          resetBtn.addEventListener('click', function() {
            if (window.EJS_emulator && window.EJS_emulator.reset) {
              window.EJS_emulator.reset();
              updateStatus('Reset', 'success');
              setTimeout(() => updateStatus('Playing', 'running'), 1000);
            }
          });
        }
      }
      
      // Check ROM format
      function checkRomFormat() {
        const validFormats = ['md', 'smd', 'gen', 'bin', 'sms', 'gg', '32x', 'cue', 'iso'];
        return validFormats.includes(FILE_EXT);
      }
      
      // Initialize
      if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
          if (!checkRomFormat()) {
            console.warn('Unusual format:', FILE_EXT);
          }
          initializeEmulator();
        });
      } else {
        if (!checkRomFormat()) {
          console.warn('Unusual format:', FILE_EXT);
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
