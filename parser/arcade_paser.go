package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
)

// GenerateArcadePlayerHTML 生成用于在Web上运行arcade游戏的HTML播放器（使用EmulatorJS）
func GenerateArcadePlayerHTML(playerID, romPath, fileName, fileExt string) string {
	// 确保所有用户输入都被正确转义
	safePlayerID := template.HTMLEscapeString(playerID)
	safeRomPath := template.HTMLEscapeString(romPath)
	safeFileName := template.HTMLEscapeString(fileName)
	safeFileExt := template.HTMLEscapeString(fileExt)

	// 生成唯一ID
	h := md5.Sum([]byte(playerID))
	uniqueID := fmt.Sprintf("arc_%x", h)[:12]

	return fmt.Sprintf(`<div class="arcade-player-wrapper" id="%[1]s-wrapper">
  <style>
    @import url('https://fonts.googleapis.com/css2?family=Press+Start+2P&display=swap');
    
    #%[1]s-wrapper {
      margin: 20px auto;
      padding: 0;
      background: linear-gradient(180deg, #1a1a1a 0%%, #000000 100%%);
      border-radius: 12px;
      max-width: 850px;
      width: 100%%;
      font-family: 'Press Start 2P', monospace;
      box-shadow: 
        0 30px 60px rgba(0, 0, 0, 0.9),
        inset 0 2px 0 rgba(255, 255, 255, 0.1),
        inset 0 -2px 0 rgba(0, 0, 0, 0.8);
      overflow: hidden;
      border: 4px solid #FFD700;
      position: relative;
    }

    /* Arcade Cabinet Header */
    #%[1]s-wrapper .arcade-header {
      background: linear-gradient(135deg, #FFD700 0%%, #FFA500 50%%, #FF6347 100%%);
      padding: 15px 20px;
      border-bottom: 5px solid #000;
      position: relative;
      overflow: hidden;
      text-align: center;
    }

    /* Arcade Marquee Effect */
    #%[1]s-wrapper .arcade-header::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: repeating-linear-gradient(
        90deg,
        transparent,
        transparent 10px,
        rgba(255, 255, 255, 0.1) 10px,
        rgba(255, 255, 255, 0.1) 20px
      );
      animation: marqueeMove 5s linear infinite;
    }

    @keyframes marqueeMove {
      0%% { transform: translateX(0); }
      100%% { transform: translateX(20px); }
    }

    /* Neon Glow Effect */
    @keyframes neonGlow {
      0%%, 100%% {
        text-shadow: 
          0 0 10px #FF00FF,
          0 0 20px #FF00FF,
          0 0 30px #FF00FF,
          0 0 40px #FF00FF;
      }
      50%% {
        text-shadow: 
          0 0 20px #00FFFF,
          0 0 30px #00FFFF,
          0 0 40px #00FFFF,
          0 0 50px #00FFFF;
      }
    }

    #%[1]s-wrapper .player-header-top {
      position: relative;
      z-index: 1;
    }
    
    #%[1]s-wrapper .player-title {
      font-size: 20px;
      font-weight: normal;
      color: #000;
      text-transform: uppercase;
      letter-spacing: 3px;
      animation: neonGlow 3s ease-in-out infinite;
      margin: 10px 0;
    }

    /* Arcade Logo Style */
    #%[1]s-wrapper .arcade-logo {
      display: inline-block;
      background: linear-gradient(180deg, #FF00FF 0%%, #00FFFF 100%%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
      font-size: 28px;
      font-weight: bold;
      text-transform: uppercase;
      letter-spacing: 4px;
      margin-bottom: 10px;
    }

    #%[1]s-wrapper .game-name {
      color: #fff;
      font-size: 12px;
      background: #000;
      padding: 5px 10px;
      border-radius: 4px;
      display: inline-block;
      margin-top: 10px;
    }
    
    #%[1]s-wrapper .player-controls {
      position: absolute;
      top: 15px;
      right: 20px;
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }
    
    /* Arcade System Info */
    #%[1]s-wrapper .player-info {
      background: repeating-linear-gradient(
        180deg,
        #2a2a2a 0px,
        #2a2a2a 20px,
        #1a1a1a 20px,
        #1a1a1a 40px
      );
      padding: 12px 20px;
      display: flex;
      gap: 30px;
      font-size: 10px;
      color: #FFD700;
      border-bottom: 2px solid #FFD700;
      text-transform: uppercase;
      letter-spacing: 1px;
      justify-content: center;
    }

    #%[1]s-wrapper .player-info-item {
      display: flex;
      align-items: center;
      gap: 10px;
    }

    #%[1]s-wrapper .player-info-item .label {
      color: #999;
      font-size: 9px;
    }

    #%[1]s-wrapper .player-info-item .value {
      color: #FFD700;
      font-weight: normal;
    }
    
    /* Insert Coin Section */
    #%[1]s-wrapper .player-tip {
      font-size: 10px;
      color: #000;
      background: linear-gradient(90deg, #00FF00 0%%, #FFFF00 50%%, #00FF00 100%%);
      padding: 12px 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 15px;
      text-transform: uppercase;
      letter-spacing: 2px;
      animation: insertCoinBlink 1s ease-in-out infinite;
    }

    @keyframes insertCoinBlink {
      0%%, 100%% { opacity: 1; }
      50%% { opacity: 0.6; }
    }
    
    #%[1]s-wrapper .player-tip .tip-icon {
      font-size: 14px;
      animation: coinRotate 2s linear infinite;
    }

    @keyframes coinRotate {
      from { transform: rotateY(0deg); }
      to { transform: rotateY(360deg); }
    }
    
    /* Arcade Style Buttons */
    #%[1]s-wrapper .ctrl-btn {
      background: linear-gradient(135deg, #FF0000 0%%, #8B0000 100%%);
      color: #fff;
      border: 2px solid #FFD700;
      padding: 8px 14px;
      border-radius: 50px;
      cursor: pointer;
      font-size: 9px;
      transition: all 0.1s;
      text-decoration: none;
      display: inline-flex;
      align-items: center;
      gap: 6px;
      text-transform: uppercase;
      font-weight: normal;
      letter-spacing: 1px;
      box-shadow: 
        0 4px 0 #660000,
        0 6px 8px rgba(0,0,0,0.5);
      transform: translateY(-2px);
    }
    
    #%[1]s-wrapper .ctrl-btn:active {
      transform: translateY(2px);
      box-shadow: 
        0 1px 0 #660000,
        0 2px 4px rgba(0,0,0,0.5);
    }
    
    #%[1]s-wrapper .ctrl-btn:hover {
      background: linear-gradient(135deg, #FF6347 0%%, #FF0000 100%%);
      border-color: #FFFF00;
    }
    
    #%[1]s-wrapper .ctrl-btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    
    #%[1]s-wrapper .ctrl-btn.success {
      background: linear-gradient(135deg, #00FF00 0%%, #008800 100%%);
      box-shadow: 
        0 4px 0 #006600,
        0 6px 8px rgba(0,0,0,0.5);
    }

    #%[1]s-wrapper .ctrl-btn.primary {
      background: linear-gradient(135deg, #0080FF 0%%, #0050AA 100%%);
      box-shadow: 
        0 4px 0 #003388,
        0 6px 8px rgba(0,0,0,0.5);
    }
    
    /* Arcade Screen Container */
    #%[1]s-wrapper .emulator-outer {
      background: #000;
      padding: 30px;
      position: relative;
      background-image: 
        repeating-linear-gradient(
          0deg,
          rgba(255, 255, 255, 0.03) 0px,
          transparent 1px,
          transparent 2px,
          rgba(255, 255, 255, 0.03) 3px
        );
    }

    #%[1]s-wrapper .screen-bezel {
      background: linear-gradient(135deg, #333 0%%, #111 100%%);
      padding: 15px;
      border-radius: 20px;
      box-shadow: 
        inset 0 5px 15px rgba(0,0,0,0.8),
        0 0 30px rgba(255, 215, 0, 0.3);
    }

    #%[1]s-wrapper .emulator-container {
      background: #000;
      border-radius: 10px;
      padding: 0;
      display: block;
      position: relative;
      max-width: 720px;
      margin: 0 auto;
      aspect-ratio: 4/3;
      box-shadow: 
        inset 0 0 50px rgba(0, 255, 255, 0.2),
        inset 0 0 20px rgba(255, 0, 255, 0.1);
      border: 2px solid #222;
    }
    
    #%[1]s-wrapper #game {
      width: 100%%;
      height: 100%%;
      max-width: 100%%;
      border-radius: 8px;
      image-rendering: pixelated;
      image-rendering: -moz-crisp-edges;
      image-rendering: crisp-edges;
    }

    /* Credit Display */
    #%[1]s-wrapper .credit-display {
      position: absolute;
      top: 10px;
      right: 10px;
      background: #000;
      color: #00FF00;
      padding: 5px 10px;
      border-radius: 4px;
      font-size: 10px;
      border: 1px solid #00FF00;
      z-index: 10;
      text-transform: uppercase;
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

    /* Arcade Loading Animation */
    #%[1]s-wrapper .arcade-loading {
      text-align: center;
    }

    #%[1]s-wrapper .coin-slot {
      width: 60px;
      height: 60px;
      border: 3px solid #FFD700;
      border-radius: 50%%;
      margin: 0 auto 20px;
      position: relative;
      animation: pulse 1.5s ease-in-out infinite;
    }

    #%[1]s-wrapper .coin-slot::before {
      content: '¢';
      position: absolute;
      top: 50%%;
      left: 50%%;
      transform: translate(-50%%, -50%%);
      font-size: 30px;
      color: #FFD700;
      animation: coinRotate 2s linear infinite;
    }

    @keyframes pulse {
      0%%, 100%% { transform: scale(1); opacity: 1; }
      50%% { transform: scale(1.1); opacity: 0.8; }
    }
    
    #%[1]s-wrapper .loading-text {
      color: #FFD700;
      font-size: 12px;
      text-transform: uppercase;
      letter-spacing: 3px;
      margin-top: 20px;
      animation: textBlink 1s step-start infinite;
    }

    @keyframes textBlink {
      0%%, 50%% { opacity: 1; }
      51%%, 100%% { opacity: 0; }
    }

    /* Control Instructions */
    #%[1]s-wrapper .controls-info {
      margin: 20px;
      padding: 20px;
      background: linear-gradient(135deg, #1a1a1a 0%%, #000 100%%);
      border-radius: 8px;
      border: 2px solid #FFD700;
      color: #FFD700;
      font-size: 10px;
    }

    #%[1]s-wrapper .controls-info h4 {
      margin: 0 0 15px 0;
      color: #fff;
      font-size: 12px;
      text-transform: uppercase;
      letter-spacing: 3px;
      text-align: center;
    }

    #%[1]s-wrapper .controls-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
      gap: 10px;
      margin-top: 15px;
    }

    #%[1]s-wrapper .control-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      background: rgba(255, 215, 0, 0.1);
      padding: 8px 12px;
      border-radius: 4px;
      border: 1px solid rgba(255, 215, 0, 0.3);
    }

    #%[1]s-wrapper .control-item .action {
      color: #FFD700;
      font-size: 9px;
      text-transform: uppercase;
    }

    #%[1]s-wrapper .key {
      background: linear-gradient(135deg, #FFD700 0%%, #FFA500 100%%);
      color: #000;
      padding: 4px 8px;
      border-radius: 4px;
      font-size: 9px;
      font-weight: bold;
      text-transform: uppercase;
    }

    /* Arcade Button Colors */
    #%[1]s-wrapper .btn-a { background: #FF0000; color: #fff; }
    #%[1]s-wrapper .btn-b { background: #FFFF00; color: #000; }
    #%[1]s-wrapper .btn-c { background: #00FF00; color: #000; }
    #%[1]s-wrapper .btn-d { background: #0080FF; color: #fff; }
    
    /* Error message */
    #%[1]s-wrapper .error-message {
      background: #000;
      border: 2px solid #FF0000;
      color: #FF0000;
      padding: 10px 15px;
      border-radius: 6px;
      margin: 10px 20px;
      display: none;
      text-align: center;
      font-size: 10px;
      text-transform: uppercase;
    }
    
    #%[1]s-wrapper .error-message.show {
      display: block;
      animation: textBlink 0.5s step-start infinite;
    }

    /* Game options panel */
    #%[1]s-wrapper .game-options {
      display: flex;
      gap: 20px;
      margin: 20px;
      padding: 15px;
      background: linear-gradient(135deg, #2a2a2a 0%%, #1a1a1a 100%%);
      border-radius: 8px;
      flex-wrap: wrap;
      align-items: center;
      border: 1px solid #FFD700;
      justify-content: center;
    }

    #%[1]s-wrapper .option-group {
      display: flex;
      align-items: center;
      gap: 8px;
      color: #FFD700;
      font-size: 10px;
      text-transform: uppercase;
      letter-spacing: 1px;
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
      width: 80px;
      accent-color: #FFD700;
    }
    
    @media (max-width: 768px) {
      #%[1]s-wrapper {
        margin: 10px;
        border-radius: 8px;
      }
      
      #%[1]s-wrapper .emulator-container {
        max-width: calc(100%% - 40px);
      }
      
      #%[1]s-wrapper .player-title {
        font-size: 14px;
      }
      
      #%[1]s-wrapper .controls-grid {
        grid-template-columns: 1fr;
      }

      #%[1]s-wrapper .player-info {
        flex-wrap: wrap;
        gap: 10px;
        font-size: 9px;
      }
    }
  </style>
  
  <div class="arcade-header">
    <div class="player-header-top">
      <div class="arcade-logo">ARCADE</div>
      <div class="player-title">%[3]s</div>
      <div class="player-controls">
        <button class="ctrl-btn" id="%[4]s-status-btn">
          <span>🕹️</span>
          <span>Ready</span>
        </button>
        <button class="ctrl-btn" id="%[4]s-fullscreen-btn">
          <span>⛶</span>
        </button>
        <a href="%[2]s" download class="ctrl-btn">
          <span>💾</span>
        </a>
      </div>
    </div>
  </div>
  
  <div class="player-info">
    <div class="player-info-item">
      <span class="label">ROM:</span>
      <span class="value">%[5]s</span>
    </div>
    <div class="player-info-item">
      <span class="label">System:</span>
      <span class="value">MAME</span>
    </div>
    <div class="player-info-item">
      <span class="label">Core:</span>
      <span class="value">FBNeo</span>
    </div>
    <div class="player-info-item">
      <span class="label">Players:</span>
      <span class="value">1-2P</span>
    </div>
  </div>

  <div class="player-tip">
    <span class="tip-icon">🪙</span>
    <span>Insert Coin to Start • Press 5 for Credit</span>
  </div>
  
  <div class="error-message" id="%[4]s-error">Game Over</div>
  
  <!-- EmulatorJS Container -->
  <div class="emulator-outer">
    <div class="screen-bezel">
      <div class="emulator-container" id="%[4]s-container">
        <div class="credit-display">Credit: 0</div>
        <div class="loading-overlay" id="%[4]s-loading">
          <div class="arcade-loading">
            <div class="coin-slot"></div>
            <div class="loading-text">Insert Coin</div>
          </div>
        </div>
        <div id="game"></div>
      </div>
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
        <span>🔊</span>
        <input type="range" id="%[4]s-volume" min="0" max="100" value="70" />
      </label>
    </div>
    <div class="option-group">
      <label>
        <input type="checkbox" id="%[4]s-scanlines" />
        <span>📺 Scanlines</span>
      </label>
    </div>
    <div class="option-group">
      <label>
        <input type="checkbox" id="%[4]s-smooth" />
        <span>🎨 Smooth</span>
      </label>
    </div>
  </div>

  <!-- Controls Information -->
  <div class="controls-info">
    <h4>🕹️ Arcade Controls</h4>
    <div class="controls-grid">
      <div class="control-item">
        <span class="action">Move</span>
        <span class="key">WASD</span>
      </div>
      <div class="control-item">
        <span class="action">Button 1</span>
        <span class="key btn-a">J</span>
      </div>
      <div class="control-item">
        <span class="action">Button 2</span>
        <span class="key btn-b">K</span>
      </div>
      <div class="control-item">
        <span class="action">Button 3</span>
        <span class="key btn-c">L</span>
      </div>
      <div class="control-item">
        <span class="action">Button 4</span>
        <span class="key btn-d">I</span>
      </div>
      <div class="control-item">
        <span class="action">Button 5</span>
        <span class="key">U</span>
      </div>
      <div class="control-item">
        <span class="action">Button 6</span>
        <span class="key">O</span>
      </div>
      <div class="control-item">
        <span class="action">Coin</span>
        <span class="key">5</span>
      </div>
      <div class="control-item">
        <span class="action">Start</span>
        <span class="key">1</span>
      </div>
      <div class="control-item">
        <span class="action">2P Start</span>
        <span class="key">2</span>
      </div>
      <div class="control-item">
        <span class="action">Menu</span>
        <span class="key">Tab</span>
      </div>
      <div class="control-item">
        <span class="action">Pause</span>
        <span class="key">P</span>
      </div>
    </div>
    <p style="margin-top: 15px; color: #999; font-size: 9px; text-transform: uppercase; text-align: center;">
      🎮 USB Controller Support • 💾 Save States Available
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
      let scanlinesEnabled = false;
      let smoothingEnabled = false;
      let creditCount = 0;
      
      // Status update function
      function updateStatus(message, type = 'loading') {
        const statusBtn = document.getElementById(UNIQUE_ID + '-status-btn');
        if (statusBtn) {
          const icons = {
            loading: '⏳',
            success: '✅',
            error: '❌',
            running: '🕹️',
            paused: '⏸️'
          };
          statusBtn.innerHTML = '<span>' + (icons[type] || '🕹️') + '</span><span>' + message + '</span>';
          statusBtn.className = 'ctrl-btn';
          if (type === 'success' || type === 'running') {
            statusBtn.classList.add('success');
          } else if (type === 'loading') {
            statusBtn.classList.add('primary');
          }
        }
      }
      
      // Update credit display
      function updateCredit() {
        const creditDisplay = document.querySelector('#' + UNIQUE_ID + '-container .credit-display');
        if (creditDisplay) {
          creditDisplay.textContent = 'Credit: ' + creditCount;
        }
      }
      
      // Show error message
      function showError(message) {
        const errorEl = document.getElementById(UNIQUE_ID + '-error');
        if (errorEl) {
          errorEl.textContent = message;
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
          updateStatus('Loading', 'loading');
	
          // 设置键盘映射 - Arcade 6按钮布局
          updateStatus('Setup', 'loading');
		  window.EJS_defaultControls = {
			  0: {
				// WASD 方向
				4: { 'value': 'w', 'value2': 'DPAD_UP' },
				5: { 'value': 's', 'value2': 'DPAD_DOWN' },
				6: { 'value': 'a', 'value2': 'DPAD_LEFT' },
				7: { 'value': 'd', 'value2': 'DPAD_RIGHT' },
				
				// 6按钮格斗游戏布局
				0: { 'value': 'j', 'value2': 'BUTTON_2' },    // Button 1
				8: { 'value': 'k', 'value2': 'BUTTON_1' },    // Button 2
				1: { 'value': 'l', 'value2': 'BUTTON_4' },    // Button 3
				9: { 'value': 'i', 'value2': 'BUTTON_3' },    // Button 4
				10: { 'value': 'u', 'value2': 'LEFT_TOP_SHOULDER' },   // Button 5
				11: { 'value': 'o', 'value2': 'RIGHT_TOP_SHOULDER' },  // Button 6
				
				// 系统键
				2: { 'value': '5', 'value2': 'SELECT' },      // Coin
				3: { 'value': '1', 'value2': 'START' },       // Start
				
				// 2P控制（可选）
				12: { 'value': '6', 'value2': 'LEFT_BOTTOM_SHOULDER' },
				13: { 'value': '2', 'value2': 'RIGHT_BOTTOM_SHOULDER' },
				
				// 功能键
				24: { 'value': 'f1' },
				25: { 'value': 'f4' },
				27: { 'value': 'space' }
			  },
			  1: {}, 2: {}, 3: {}
			};

          updateLoadingText('BOOTING ROM...');
          
          // EmulatorJS configuration
          window.EJS_player = '#game';
          window.EJS_gameUrl = ROM_URL;
          window.EJS_gameParentUrl = ROM_URL;
          window.EJS_core = 'fbneo'; // mame2003，fbneo
          window.EJS_gameName = ROM_NAME;
 		 
		  window.EJS_biosUrl = '/bios/arcade.7z'; // 如果需要bios请放在此路径
          window.EJS_lightgun = false; // 除非是光枪游戏

          // Optional configurations
          window.EJS_pathtodata = 'https://cdn.emulatorjs.org/stable/data/';
          window.EJS_startOnLoaded = true;
          window.EJS_DEBUG_XX = false;
          
          // Color and theme
          window.EJS_color = '#FFD700';
          
          // Features
          window.EJS_cheats = false;
          window.EJS_saveStateName = ROM_NAME + '.state';
          
          // Language
          window.EJS_language = 'en-US';
          
          // Settings menu
          window.EJS_settingsMenu = true;
          
          // Screen settings (Arcade standard)
          window.EJS_aspectRatio = '4:3';
          
          // Virtual gamepad for mobile
          window.EJS_VirtualGamepadSettings = {
            scale: 1.5,
            opacity: 0.8,
            enabled: true
          };
          
          // Callbacks
          window.EJS_onGameStart = function() {
            console.log('Arcade game started:', ROM_NAME);
            updateStatus('Playing', 'running');
            hideLoading();
            emulatorReady = true;
            creditCount = 99; // Free play mode
            updateCredit();
            
            // Initialize custom controls
            initializeCustomControls();
          };
          
          window.EJS_onLoadState = function() {
            console.log('State loaded');
            updateStatus('Loaded', 'success');
            setTimeout(() => updateStatus('Playing', 'running'), 2000);
          };
          
          window.EJS_onSaveState = function() {
            console.log('State saved');
            updateStatus('Saved', 'success');
            setTimeout(() => updateStatus('Playing', 'running'), 2000);
          };
          
          // Setup button handlers
          setupButtonHandlers();
          
          // Setup option handlers
          setupOptionHandlers();
          
          // Load EmulatorJS loader script
          updateLoadingText('LOADING...');
          const script = document.createElement('script');
          script.src = '/emulatorjs/loader.js';
          script.onerror = function() {
            showError('SYSTEM ERROR - CHECK ROM');
            hideLoading();
          };
          script.onload = function() {
            console.log('EmulatorJS loader loaded successfully');
            updateLoadingText('READY PLAYER ONE');
            updateStatus('Starting', 'loading');
          };
          document.body.appendChild(script);
          
        } catch (error) {
          console.error('Failed to initialize emulator:', error);
          showError('BOOT ERROR: ' + error.message);
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
        
        // Scanlines filter toggle
        const scanlinesCheckbox = document.getElementById(UNIQUE_ID + '-scanlines');
        if (scanlinesCheckbox) {
          scanlinesCheckbox.addEventListener('change', function() {
            scanlinesEnabled = this.checked;
            if (window.EJS_emulator && window.EJS_emulator.setShader) {
              window.EJS_emulator.setShader(scanlinesEnabled ? 'crt' : 'none');
            }
          });
        }
        
        // Smoothing toggle
        const smoothCheckbox = document.getElementById(UNIQUE_ID + '-smooth');
        if (smoothCheckbox) {
          smoothCheckbox.addEventListener('change', function() {
            smoothingEnabled = this.checked;
            const gameCanvas = document.querySelector('#game canvas');
            if (gameCanvas) {
              gameCanvas.style.imageRendering = smoothingEnabled ? 'auto' : 'pixelated';
            }
          });
        }
      }
      
      // Check ROM format
      function checkRomFormat() {
        const validFormats = ['zip', '7z', 'chd', 'bin', 'iso'];
        if (!validFormats.includes(FILE_EXT)) {
          console.warn('Unusual file format for Arcade:', FILE_EXT);
          return true; // Still try to load
        }
        return true;
      }
      
      // Initialize when DOM is ready
      if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
          if (!checkRomFormat()) {
            showError('INVALID ROM FORMAT');
          }
          initializeEmulator();
        });
      } else {
        if (!checkRomFormat()) {
          showError('INVALID ROM FORMAT');
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
