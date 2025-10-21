package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
)

// GeneratePlayStationPlayerHTML 生成用于在Web上运行PlayStation游戏的HTML播放器（使用EmulatorJS）
func GeneratePlayStationPlayerHTML(playerID, romPath, fileName, fileExt string) string {
	// 确保所有用户输入都被正确转义
	safePlayerID := template.HTMLEscapeString(playerID)
	safeRomPath := template.HTMLEscapeString(romPath)
	safeFileName := template.HTMLEscapeString(fileName)
	safeFileExt := template.HTMLEscapeString(fileExt)

	// 生成唯一ID
	h := md5.Sum([]byte(playerID))
	uniqueID := fmt.Sprintf("ps_%x", h)[:12]

	return fmt.Sprintf(`<div class="ps-player-wrapper" id="%[1]s-wrapper">
  <style>
    @import url('https://fonts.googleapis.com/css2?family=Bebas+Neue&display=swap');
    
    #%[1]s-wrapper {
      margin: 20px auto;
      padding: 0;
      background: linear-gradient(180deg, #8B8B8B 0%%, #4A4A4A 50%%, #2C2C2C 100%%);
      border-radius: 8px;
      max-width: 800px;
      width: 100%%;
      font-family: 'Bebas Neue', 'Arial', sans-serif;
      box-shadow: 
        0 20px 50px rgba(0, 0, 0, 0.8),
        inset 0 1px 0 rgba(255, 255, 255, 0.3),
        inset 0 -1px 0 rgba(0, 0, 0, 0.5);
      overflow: hidden;
      border: 2px solid #1a1a1a;
      position: relative;
    }

    /* PlayStation Classic Header */
    #%[1]s-wrapper .ps-header {
      background: linear-gradient(135deg, #000000 0%%, #1a1a1a 100%%);
      padding: 12px 20px;
      border-bottom: 4px solid transparent;
      border-image: linear-gradient(90deg, #FF0000 0%%, #00FF00 33%%, #0080FF 66%%, #FFD700 100%%);
      border-image-slice: 1;
      position: relative;
      overflow: hidden;
    }

    /* PlayStation Logo Animation */
    #%[1]s-wrapper .ps-header::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      width: 100%%;
      height: 100%%;
      background: linear-gradient(90deg, 
        transparent 0%%, 
        rgba(255,255,255,0.05) 50%%, 
        transparent 100%%);
      transform: translateX(-100%%);
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
      font-size: 24px;
      font-weight: 400;
      color: #fff;
      display: flex;
      align-items: center;
      gap: 15px;
      letter-spacing: 2px;
    }

    /* PlayStation Logo Style */
    #%[1]s-wrapper .ps-logo {
      display: flex;
      align-items: center;
      gap: 3px;
    }

    #%[1]s-wrapper .ps-logo .ps-p {
      color: #FF0000;
      font-size: 28px;
      font-weight: bold;
      font-style: italic;
      transform: skewX(-10deg);
    }

    #%[1]s-wrapper .ps-logo .ps-s {
      color: #FFD700;
      font-size: 28px;
      font-weight: bold;
      font-style: italic;
      transform: skewX(-10deg);
      margin-left: -5px;
    }

    #%[1]s-wrapper .ps-badge {
      background: #000;
      color: #fff;
      padding: 4px 10px;
      border-radius: 4px;
      font-size: 10px;
      text-transform: uppercase;
      letter-spacing: 2px;
      border: 1px solid #333;
    }
    
    #%[1]s-wrapper .player-controls {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }
    
    /* PlayStation System Info */
    #%[1]s-wrapper .player-info {
      background: linear-gradient(90deg, #1a1a1a 0%%, #2c2c2c 100%%);
      padding: 10px 20px;
      display: flex;
      gap: 25px;
      font-size: 12px;
      color: #999;
      border-bottom: 1px solid #000;
      text-transform: uppercase;
      letter-spacing: 1px;
    }

    #%[1]s-wrapper .player-info-item {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    #%[1]s-wrapper .player-info-item .label {
      color: #666;
      font-size: 10px;
    }

    #%[1]s-wrapper .player-info-item .value {
      color: #fff;
      font-weight: bold;
    }
    
    /* PlayStation Tip */
    #%[1]s-wrapper .player-tip {
      font-size: 12px;
      color: #fff;
      background: linear-gradient(135deg, rgba(0, 128, 255, 0.2) 0%%, rgba(0, 0, 0, 0.4) 100%%);
      padding: 10px 15px;
      border-left: 3px solid #0080FF;
      display: flex;
      align-items: center;
      gap: 10px;
      margin: 15px 20px;
      border-radius: 4px;
      text-transform: uppercase;
      letter-spacing: 1px;
    }
    
    #%[1]s-wrapper .player-tip .tip-icon {
      color: #0080FF;
      font-size: 16px;
    }
    
    /* PlayStation Style Buttons */
    #%[1]s-wrapper .ctrl-btn {
      background: linear-gradient(135deg, #333 0%%, #1a1a1a 100%%);
      color: #fff;
      border: 1px solid #555;
      padding: 6px 14px;
      border-radius: 20px;
      cursor: pointer;
      font-size: 11px;
      transition: all 0.2s;
      text-decoration: none;
      display: inline-flex;
      align-items: center;
      gap: 6px;
      text-transform: uppercase;
      font-weight: 400;
      letter-spacing: 1px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.5);
    }
    
    #%[1]s-wrapper .ctrl-btn:hover {
      background: linear-gradient(135deg, #0080FF 0%%, #0050AA 100%%);
      transform: translateY(-1px);
      box-shadow: 0 4px 8px rgba(0,128,255,0.4);
    }
    
    #%[1]s-wrapper .ctrl-btn:active {
      transform: translateY(0);
    }
    
    #%[1]s-wrapper .ctrl-btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    
    #%[1]s-wrapper .ctrl-btn.success {
      background: linear-gradient(135deg, #00AA00 0%%, #006600 100%%);
      border-color: #00FF00;
    }

    #%[1]s-wrapper .ctrl-btn.primary {
      background: linear-gradient(135deg, #0080FF 0%%, #0050AA 100%%);
      border-color: #0080FF;
    }
    
    /* EmulatorJS Container - CRT TV Style */
    #%[1]s-wrapper .emulator-outer {
      background: linear-gradient(135deg, #2a2a2a 0%%, #1a1a1a 100%%);
      padding: 20px;
      position: relative;
    }

    #%[1]s-wrapper .emulator-container {
      background: #000;
      border-radius: 8px;
      padding: 0;
      display: block;
      position: relative;
      max-width: 720px;
      margin: 0 auto;
      aspect-ratio: 4/3;
      box-shadow: 
        0 15px 40px rgba(0, 0, 0, 0.7),
        inset 0 0 30px rgba(0, 0, 0, 0.8),
        inset 0 0 60px rgba(0, 128, 255, 0.1);
      border: 3px solid #222;
    }
    
    #%[1]s-wrapper #game {
      width: 100%%;
      height: 100%%;
      max-width: 100%%;
      border-radius: 6px;
    }

    /* Power LED */
    #%[1]s-wrapper .power-indicator {
      position: absolute;
      top: 10px;
      right: 10px;
      display: flex;
      align-items: center;
      gap: 8px;
      z-index: 10;
    }

    #%[1]s-wrapper .power-led {
      width: 8px;
      height: 8px;
      border-radius: 50%%;
      background: #00FF00;
      box-shadow: 0 0 10px #00FF00;
      animation: pulseLed 2s infinite;
    }

    @keyframes pulseLed {
      0%%, 100%% { opacity: 1; }
      50%% { opacity: 0.6; }
    }

    #%[1]s-wrapper .power-text {
      color: #666;
      font-size: 10px;
      text-transform: uppercase;
      letter-spacing: 1px;
    }
    
    /* Loading overlay */
    #%[1]s-wrapper .loading-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: #000;
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

    /* PlayStation Loading Animation */
    #%[1]s-wrapper .ps-loading {
      position: relative;
      width: 120px;
      height: 120px;
    }

    #%[1]s-wrapper .ps-symbols {
      position: absolute;
      width: 100%%;
      height: 100%%;
      animation: rotate 2s linear infinite;
    }

    #%[1]s-wrapper .ps-symbol {
      position: absolute;
      font-size: 24px;
      font-weight: bold;
    }

    #%[1]s-wrapper .ps-symbol.triangle {
      top: 0;
      left: 50%%;
      transform: translateX(-50%%);
      color: #00FF00;
    }

    #%[1]s-wrapper .ps-symbol.circle {
      top: 50%%;
      right: 0;
      transform: translateY(-50%%);
      color: #FF0000;
    }

    #%[1]s-wrapper .ps-symbol.x {
      bottom: 0;
      left: 50%%;
      transform: translateX(-50%%);
      color: #0080FF;
    }

    #%[1]s-wrapper .ps-symbol.square {
      top: 50%%;
      left: 0;
      transform: translateY(-50%%);
      color: #FF00FF;
    }

    @keyframes rotate {
      from { transform: rotate(0deg); }
      to { transform: rotate(360deg); }
    }
    
    #%[1]s-wrapper .loading-text {
      color: #fff;
      margin-top: 140px;
      font-size: 14px;
      text-transform: uppercase;
      letter-spacing: 3px;
      opacity: 0.8;
    }

    /* Control Instructions */
    #%[1]s-wrapper .controls-info {
      margin: 20px;
      padding: 15px;
      background: linear-gradient(135deg, rgba(0, 0, 0, 0.6) 0%%, rgba(0, 0, 0, 0.4) 100%%);
      border-radius: 8px;
      border: 1px solid #333;
      color: #ccc;
      font-size: 12px;
    }

    #%[1]s-wrapper .controls-info h4 {
      margin: 0 0 12px 0;
      color: #fff;
      font-size: 14px;
      text-transform: uppercase;
      letter-spacing: 2px;
      display: flex;
      align-items: center;
      gap: 8px;
    }

    #%[1]s-wrapper .controls-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
      gap: 8px;
      margin-top: 12px;
    }

    #%[1]s-wrapper .control-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      background: rgba(255, 255, 255, 0.05);
      padding: 6px 10px;
      border-radius: 4px;
      border: 1px solid rgba(255, 255, 255, 0.1);
    }

    #%[1]s-wrapper .control-item .action {
      color: #999;
      font-size: 11px;
    }

    #%[1]s-wrapper .key {
      background: linear-gradient(135deg, #333 0%%, #1a1a1a 100%%);
      border: 1px solid #555;
      padding: 3px 8px;
      border-radius: 4px;
      font-family: monospace;
      font-size: 10px;
      color: #fff;
      font-weight: bold;
      text-transform: uppercase;
    }

    /* PS Button Colors */
    #%[1]s-wrapper .ps-triangle { color: #00FF00; }
    #%[1]s-wrapper .ps-circle { color: #FF0000; }
    #%[1]s-wrapper .ps-x { color: #0080FF; }
    #%[1]s-wrapper .ps-square { color: #FF00FF; }
    
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
      background: linear-gradient(135deg, rgba(0, 0, 0, 0.5) 0%%, rgba(0, 0, 0, 0.3) 100%%);
      border-radius: 6px;
      flex-wrap: wrap;
      align-items: center;
      border: 1px solid #333;
    }

    #%[1]s-wrapper .option-group {
      display: flex;
      align-items: center;
      gap: 8px;
      color: #ccc;
      font-size: 12px;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    #%[1]s-wrapper .option-group label {
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 5px;
    }

    #%[1]s-wrapper .option-group input[type="checkbox"] {
      cursor: pointer;
      accent-color: #0080FF;
    }

    #%[1]s-wrapper .option-group input[type="range"] {
      width: 100px;
      accent-color: #0080FF;
    }
    
    @media (max-width: 768px) {
      #%[1]s-wrapper {
        margin: 10px;
        border-radius: 6px;
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
  
  <div class="ps-header">
    <div class="player-header-top">
      <div class="player-title">
        <div class="ps-logo">
          <span class="ps-p">P</span>
          <span class="ps-s">S</span>
        </div>
        <span class="ps-badge">PlayStation</span>
        <span style="color: #666;">|</span>
        <span style="font-size: 16px; font-weight: normal; letter-spacing: 1px;">%[3]s</span>
      </div>
      <div class="player-controls">
        <button class="ctrl-btn" id="%[4]s-status-btn">
          <span>⏳</span>
          <span>Loading</span>
        </button>
        <button class="ctrl-btn" id="%[4]s-fullscreen-btn">
          <span>⛶</span>
          <span>Full</span>
        </button>
        <a href="%[2]s" download class="ctrl-btn">
          <span>💿</span>
          <span>ISO</span>
        </a>
      </div>
    </div>
  </div>
  
  <div class="player-info">
    <div class="player-info-item">
      <span class="label">Format:</span>
      <span class="value">%[5]s</span>
    </div>
    <div class="player-info-item">
      <span class="label">System:</span>
      <span class="value">PSX</span>
    </div>
    <div class="player-info-item">
      <span class="label">Core:</span>
      <span class="value">PCSX-ReARMed</span>
    </div>
    <div class="player-info-item">
      <span class="label">Region:</span>
      <span class="value">NTSC/PAL</span>
    </div>
  </div>

  <div class="player-tip">
    <span class="tip-icon">💡</span>
    <span>Press Tab for Menu • Supports Memory Cards • DualShock Compatible</span>
  </div>
  
  <div class="error-message" id="%[4]s-error"></div>
  
  <!-- EmulatorJS Container -->
  <div class="emulator-outer">
    <div class="emulator-container" id="%[4]s-container">
      <div class="power-indicator">
        <div class="power-led"></div>
        <span class="power-text">Power</span>
      </div>
      <div class="loading-overlay" id="%[4]s-loading">
        <div class="ps-loading">
          <div class="ps-symbols">
            <div class="ps-symbol triangle">△</div>
            <div class="ps-symbol circle">○</div>
            <div class="ps-symbol x">✕</div>
            <div class="ps-symbol square">□</div>
          </div>
        </div>
        <div class="loading-text">PlayStation</div>
      </div>
      <div id="game"></div>
    </div>
  </div>

  <!-- Game Options -->
  <div class="game-options">
    <div class="option-group">
      <label>
        <input type="checkbox" id="%[4]s-turbo" />
        <span>⚡ Turbo</span>
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
        <input type="checkbox" id="%[4]s-vibration" />
        <span>📳 Vibration</span>
      </label>
    </div>
  </div>

  <!-- Controls Information -->
  <div class="controls-info">
    <h4>🎮 DualShock Controls</h4>
    <div class="controls-grid">
      <div class="control-item">
        <span class="action">D-Pad</span>
        <span class="key">WASD</span>
      </div>
      <div class="control-item">
        <span class="action ps-x">✕</span>
        <span class="key">K</span>
      </div>
      <div class="control-item">
        <span class="action ps-circle">○</span>
        <span class="key">L</span>
      </div>
      <div class="control-item">
        <span class="action ps-square">□</span>
        <span class="key">J</span>
      </div>
      <div class="control-item">
        <span class="action ps-triangle">△</span>
        <span class="key">I</span>
      </div>
      <div class="control-item">
        <span class="action">L1/R1</span>
        <span class="key">U / O</span>
      </div>
      <div class="control-item">
        <span class="action">L2/R2</span>
        <span class="key">Y / P</span>
      </div>
      <div class="control-item">
        <span class="action">Start</span>
        <span class="key">Enter</span>
      </div>
      <div class="control-item">
        <span class="action">Select</span>
        <span class="key">Shift</span>
      </div>
      <div class="control-item">
        <span class="action">L Analog</span>
        <span class="key">TFGH</span>
      </div>
      <div class="control-item">
        <span class="action">R Analog</span>
        <span class="key">8456</span>
      </div>
      <div class="control-item">
        <span class="action">Menu</span>
        <span class="key">Tab</span>
      </div>
    </div>
    <p style="margin-top: 12px; color: #666; font-size: 11px; text-transform: uppercase; letter-spacing: 1px;">
      💾 Supports Memory Card Saves • 🎮 USB/Bluetooth Controller Compatible
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
      let crtFilter = false;
      let vibrationEnabled = false;
      
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
          setTimeout(() => {
            loadingEl.classList.add('hidden');
          }, 500);
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
	
          // 设置键盘映射 - PlayStation 布局
          updateStatus('Configuring...', 'loading');
          window.EJS_defaultControls = {
            0: {
              // WASD 方向
              4: { 'value': 'w', 'value2': 'DPAD_UP' },
              5: { 'value': 's', 'value2': 'DPAD_DOWN' },
              6: { 'value': 'a', 'value2': 'DPAD_LEFT' },
              7: { 'value': 'd', 'value2': 'DPAD_RIGHT' },
              
              // PlayStation 按键布局
              0: { 'value': 'k', 'value2': 'BUTTON_2' },    // X (Cross)
              8: { 'value': 'l', 'value2': 'BUTTON_1' },    // Circle
              1: { 'value': 'j', 'value2': 'BUTTON_4' },    // Square
              9: { 'value': 'i', 'value2': 'BUTTON_3' },    // Triangle
              
              // 系统键
              2: { 'value': 'shift', 'value2': 'SELECT' },
              3: { 'value': 'enter', 'value2': 'START' },
              
              // 肩键
              10: { 'value': 'u', 'value2': 'LEFT_TOP_SHOULDER' },    // L1
              11: { 'value': 'o', 'value2': 'RIGHT_TOP_SHOULDER' },   // R1
              12: { 'value': 'y', 'value2': 'LEFT_BOTTOM_SHOULDER' }, // L2
              13: { 'value': 'p', 'value2': 'RIGHT_BOTTOM_SHOULDER' },// R2
              
              // 摇杆按下
              14: { 'value': 'v', 'value2': 'LEFT_STICK' },  // L3
              15: { 'value': 'n', 'value2': 'RIGHT_STICK' }, // R3
              
              // 左摇杆
              16: { 'value': 'h', 'value2': 'LEFT_STICK_X:+1' },
              17: { 'value': 'f', 'value2': 'LEFT_STICK_X:-1' },
              18: { 'value': 'g', 'value2': 'LEFT_STICK_Y:+1' },
              19: { 'value': 't', 'value2': 'LEFT_STICK_Y:-1' },
              
              // 右摇杆
              20: { 'value': '6', 'value2': 'RIGHT_STICK_X:+1' },
              21: { 'value': '4', 'value2': 'RIGHT_STICK_X:-1' },
              22: { 'value': '5', 'value2': 'RIGHT_STICK_Y:+1' },
              23: { 'value': '8', 'value2': 'RIGHT_STICK_Y:-1' },
              
              // 功能键
              24: { 'value': '1' },      // Quick Save
              25: { 'value': '3' },      // Quick Load  
              27: { 'value': 'space' }   // Fast Forward
            },
            1: {}, 2: {}, 3: {}
          };

          updateLoadingText('LOADING SYSTEM...');
          
          // EmulatorJS configuration
          window.EJS_player = '#game';
          window.EJS_gameUrl = ROM_URL;
          window.EJS_core = 'psx'; // PlayStation core
          window.EJS_gameName = ROM_NAME;
          
          // Optional configurations
          window.EJS_pathtodata = 'https://cdn.emulatorjs.org/stable/data/';
          window.EJS_startOnLoaded = true;
          window.EJS_biosUrl = '/bios/scph5501.bin'; // PSX BIOS - scph5501.bin == us language
          window.EJS_DEBUG_XX = false;
          
          // Color and theme
          window.EJS_color = '#0080FF';
          
          // Features
          window.EJS_cheats = true;
          window.EJS_saveStateName = ROM_NAME + '.state';
          
          // Language
          window.EJS_language = 'en-US';
          
          // Settings menu
          window.EJS_settingsMenu = true;
          
          // Screen settings (PS1 native resolution)
          window.EJS_aspectRatio = '4:3';
          
          // Virtual gamepad for mobile
          window.EJS_VirtualGamepadSettings = {
            scale: 1.5,
            opacity: 0.8,
            enabled: true
          };
          
          // Callbacks
          window.EJS_onGameStart = function() {
            console.log('PlayStation game started:', ROM_NAME);
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
            updateStatus('Loaded', 'success');
            setTimeout(() => updateStatus('Running', 'running'), 2000);
          };
          
          window.EJS_onSaveState = function() {
            console.log('State saved');
            updateStatus('Saved', 'success');
            setTimeout(() => updateStatus('Running', 'running'), 2000);
          };
          
          // Setup button handlers
          setupButtonHandlers();
          
          // Setup option handlers
          setupOptionHandlers();
          
          // Load EmulatorJS loader script
          updateLoadingText('BOOTING...');
          const script = document.createElement('script');
          script.src = '/emulatorjs/loader.js';
          script.onerror = function() {
            showError('Failed to load EmulatorJS. Please check installation.');
            hideLoading();
          };
          script.onload = function() {
            console.log('EmulatorJS loader loaded successfully');
            updateLoadingText('SONY COMPUTER ENTERTAINMENT');
            updateStatus('Starting...', 'loading');
          };
          document.body.appendChild(script);
          
        } catch (error) {
          console.error('Failed to initialize emulator:', error);
          showError('Failed to initialize: ' + error.message);
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
        
        
        
        // Vibration toggle
        const vibrationCheckbox = document.getElementById(UNIQUE_ID + '-vibration');
        if (vibrationCheckbox) {
          vibrationCheckbox.addEventListener('change', function() {
            vibrationEnabled = this.checked;
            // EmulatorJS will handle vibration if supported
          });
        }
      }
      
      // Check ROM format
      function checkRomFormat() {
        const validFormats = ['cue', 'bin', 'img', 'iso', 'chd', 'pbp', 'ccd'];
        if (!validFormats.includes(FILE_EXT)) {
          console.warn('Unusual file format for PSX:', FILE_EXT);
          return true; // Still try to load
        }
        return true;
      }
      
      // Initialize when DOM is ready
      if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
          if (!checkRomFormat()) {
            showError('Warning: Unusual format. Expected: .cue, .bin, .iso, .chd');
          }
          initializeEmulator();
        });
      } else {
        if (!checkRomFormat()) {
          showError('Warning: Unusual format. Expected: .cue, .bin, .iso, .chd');
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
