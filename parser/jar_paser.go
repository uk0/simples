package parser

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"strings"
)

// GenerateJARPlayerHTML 生成用于在Web上运行JAR文件的HTML播放器
func GenerateJARPlayerHTML(playerID, jarPath, fileName, fileExt string) string {
	// 确保所有用户输入都被正确转义
	safePlayerID := template.HTMLEscapeString(playerID)
	safeJarPath := template.HTMLEscapeString(jarPath)
	safeFileName := template.HTMLEscapeString(fileName)

	// 生成唯一ID
	h := md5.Sum([]byte(playerID))
	uniqueID := fmt.Sprintf("jar_%x", h)[:12]

	// 确保文件名包含扩展名
	if !strings.HasSuffix(safeFileName, ".jar") {
		safeFileName = safeFileName + ".jar"
	}

	return fmt.Sprintf(`<div class="jar-player-wrapper" id="%[1]s-wrapper">
  <style>
    #%[1]s-wrapper {
      margin: 20px auto;
      padding: 20px;
      background: linear-gradient(135deg, #1a1a2e 0%%, #16213e 100%%);
      border-radius: 12px;
      max-width: fit-content;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
    }
    
    #%[1]s-wrapper .player-header {
      color: #fff;
      margin-bottom: 20px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 20px;
    }
    
    #%[1]s-wrapper .player-title {
      font-size: 18px;
      font-weight: 500;
      opacity: 0.9;
    }
    
    #%[1]s-wrapper .player-controls {
      display: flex;
      gap: 8px;
    }
    
    #%[1]s-wrapper .ctrl-btn {
      background: rgba(255, 255, 255, 0.1);
      color: #fff;
      border: 1px solid rgba(255, 255, 255, 0.2);
      padding: 6px 12px;
      border-radius: 6px;
      cursor: pointer;
      font-size: 13px;
      transition: all 0.2s;
      text-decoration: none;
    }
    
    #%[1]s-wrapper .ctrl-btn:hover {
      background: rgba(255, 255, 255, 0.2);
    }
    
    #%[1]s-wrapper .ctrl-btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    
    #%[1]s-wrapper .ctrl-btn.success {
      background: rgba(76, 175, 80, 0.3);
      border-color: rgba(76, 175, 80, 0.5);
    }
    
    /* Main display container */
    #%[1]s-wrapper #display-container {
      background: #000;
      border-radius: 8px;
      padding: 10px;
      display: inline-block;
      position: relative;
    }
    
    #%[1]s-wrapper #display {
      position: relative;
      background: #000;
    }
    
    #%[1]s-wrapper #canvas {
      display: block;
      image-rendering: pixelated;
      image-rendering: -moz-crisp-edges;
      image-rendering: crisp-edges;
      background: #000;
    }
    
    /* J2ME.js screens */
    #%[1]s-wrapper .screen {
      display: none;
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0,0,0,0.95);
      color: #fff;
      text-align: center;
      padding: 20px;
      z-index: 10;
    }
    
    #%[1]s-wrapper .screen.active {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
    }
    
    #%[1]s-wrapper .screen .title {
      font-size: 16px;
      margin-bottom: 15px;
    }
    
    #%[1]s-wrapper .screen progress {
      width: 200px;
      height: 4px;
    }
    
    /* Virtual keypad */
    #%[1]s-wrapper #keypad {
      margin-top: 15px;
      display: flex;
      justify-content: center;
      gap: 30px;
      user-select: none;
    }
    
    #%[1]s-wrapper #joypad,
    #%[1]s-wrapper #numpad {
      display: inline-block;
    }
    
    #%[1]s-wrapper .key {
      background: #2a2a2a;
      border: 1px solid #444;
      border-radius: 6px;
      color: #fff;
      width: 40px;
      height: 40px;
      margin: 2px;
      cursor: pointer;
      font-size: 14px;
      font-weight: 500;
      transition: all 0.1s;
      vertical-align: top;
    }
    
    #%[1]s-wrapper .key:active {
      background: #444;
      transform: scale(0.95);
    }
    
    #%[1]s-wrapper #choice { background: #4a5568; }
    #%[1]s-wrapper #back { background: #e53e3e; }
    #%[1]s-wrapper #ok { background: #38a169; }
    
    #%[1]s-wrapper .numpad2 { display: none; }
    #%[1]s-wrapper.gamepad-2 .numpad { display: none; }
    #%[1]s-wrapper.gamepad-2 .numpad2 { display: inline-block; }
    
    #%[1]s-wrapper .lcdui-alert {
      position: fixed;
      top: 50%%;
      left: 50%%;
      transform: translate(-50%%, -50%%);
      background: #fff;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 0 10px 30px rgba(0,0,0,0.5);
      z-index: 1000;
      min-width: 280px;
      color: #333;
    }
    
    #%[1]s-wrapper #header {
      display: none !important;
    }
    
    #%[1]s-wrapper .loading-message {
      position: absolute;
      top: 50%%;
      left: 50%%;
      transform: translate(-50%%, -50%%);
      color: #fff;
      font-size: 14px;
      z-index: 100;
      background: rgba(0,0,0,0.8);
      padding: 10px 20px;
      border-radius: 8px;
    }
    
    @media (max-width: 640px) {
      #%[1]s-wrapper {
        margin: 10px;
        padding: 15px;
      }
      
      #%[1]s-wrapper #canvas {
        max-width: calc(100vw - 60px);
        height: auto;
      }
      
      #%[1]s-wrapper .key {
        width: 36px;
        height: 36px;
        font-size: 12px;
      }
    }
  </style>
  
  <div class="player-header">
    <div class="player-title">📱 %[3]s</div>
    <div class="player-controls">
      <button class="ctrl-btn" id="%[4]s-status-btn">⏳ Loading...</button>
      <button class="ctrl-btn" onclick="document.getElementById('display-container').requestFullscreen?.()">⛶</button>
      <button class="ctrl-btn" onclick="location.reload()">↻</button>
      <a href="%[2]s" download class="ctrl-btn">⬇</a>
    </div>
  </div>
  
  <!-- Hidden file input required by J2ME.js -->
  <input id="File" type="file" style="visibility:hidden;position:absolute;width:0;height:0">
  
  <!-- Main display matching J2ME.js structure -->
  <div id="display-container">
    <div id="display">
      <div id="header"></div>
      <canvas id="canvas" width="240" height="320"></canvas>
      <div id="download-screen" class="screen">
        <h1 class="title">Downloading midlet…</h1>
        <progress id="download-bar" value="0" max="100"></progress>
      </div>
      <div id="splash-screen" class="screen">
        <h1 class="title">Starting midlet…</h1>
        <progress></progress>
      </div>
      <div id="exit-screen" class="screen"></div>
      <div id="background-screen"></div>
      <div class="loading-message" id="%[4]s-loading-msg">Initializing...</div>
    </div>
  </div>
  <!-- LCDUI Alert Dialog -->
  <form role="dialog" data-type="confirm" class="lcdui-alert" id="lcdui-alert" style="display:none">
    <section>
      <h1 class="title"></h1>
      <p class="text"></p>
      <p><progress style="display:none"></progress></p>
    </section>
    <menu>
      <button class="button0" style="display:none"></button>
      <button class="button1" style="display:none"></button>
    </menu>
  </form>
  <!-- Virtual Keypad -->
  <div id="keypad" class="gamepad-3 row-reverse">
    <div id="joypad">
      <button id="choice" class="key">&nbsp;</button>
      <button id="up" class="key">↑</button>
      <button id="back" class="key">&nbsp;</button><br>
      <button id="left" class="key">←</button>
      <button id="ok" class="key">&nbsp;</button>
      <button id="right" class="key">→</button><br>
      <button id="down" class="key">↓</button>
    </div>
    <div id="numpad">
      <button id="num1" class="key">1</button>
      <button id="num2" class="key">2</button>
      <button id="num3" class="key">3</button>
      <button id="star" class="key numpad2">*</button><br>
      <button id="num4" class="key">4</button>
      <button id="num5" class="key">5</button>
      <button id="num6" class="key">6</button>
      <button id="num0" class="key numpad2">0</button><br>
      <button id="num7" class="key">7</button>
      <button id="num8" class="key">8</button>
      <button id="num9" class="key">9</button>
      <button id="pound" class="key numpad2">#</button><br>
      <button id="star" class="key numpad">*</button>
      <button id="num0" class="key numpad">0</button>
      <button id="pound" class="key numpad">#</button>
    </div>
  </div>
  <!-- Load J2ME.js dependencies first -->
  <script src="/jar_em/libs/encoding.js"></script>
  <script src="/jar_em/libs/webaudio-tinysynth.min.js"></script>
  <script src="/jar_em/libs/BenzAMRRecorder.min.js"></script>
  <!-- J2ME.js core files -->
  <script src="/jar_em/bld/native.js" defer></script>
  <script src="/jar_em/config/default.js" defer></script>
  <script src="/jar_em/config/build.js" defer></script>
  <script src="/jar_em/config/urlparams.js" defer></script>
  <script src="/jar_em/bld/j2me.js" defer></script>
  <script src="/jar_em/bld/main-all.js" defer></script>
  <script src="/jar_em/keymap.js" defer></script>
  <script>
  // Initialize J2ME.js with proper configuration
  (function() {
    const JAR_URL = '%[2]s';
    const JAR_NAME = '%[3]s';
    const UNIQUE_ID = '%[4]s';
    // Status update function
    function updateStatus(message, isSuccess) {
      const statusBtn = document.getElementById(UNIQUE_ID + '-status-btn');
      if (statusBtn) {
        statusBtn.textContent = message;
        if (isSuccess) {
          statusBtn.classList.add('success');
        }
      }
    }
    // Pre-configure global config object before J2ME.js loads
    config = {};
    config.localjar = JAR_NAME;
    config.enginemode = 'enginemode2-classes2.jar';
    config.gamepadSize = 'gamepad-3';
    config.gamepad = 1;
    config.canvasSize = 'size-240x320';
    config.gameresize = 'resize-1';
    config.platform = 'nokia/5800';
    config.autosize = false;
    config.rotate = false;
    config.scanvas = 0;
    config.jars = [];
    config.jad = null;
    config.midletClassName = '';
    // Override urlParams to use our config
	window.config = config; // Make sure config is global
    window.urlParams = {
      parse: function(url) {
        console.log('urlParams.parse called, using pre-configured settings');
        return config;
      }
    };
    // Function to download and install JAR using JARStore
    async function downloadAndInstallJAR() {
      const loadingMsg = document.getElementById(UNIQUE_ID + '-loading-msg');
      try {
        updateStatus('⏳ Downloading...', false);
        if (loadingMsg) loadingMsg.textContent = 'Downloading JAR file...';
        console.log('Downloading JAR from:', JAR_URL);
        const response = await fetch(JAR_URL);
        if (!response.ok) {
          throw new Error('Failed to download JAR: ' + response.status);
        }
        const arrayBuffer = await response.arrayBuffer();
        console.log('JAR downloaded, size:', arrayBuffer.byteLength);
        updateStatus('⏳ Installing...', false);
        if (loadingMsg) loadingMsg.textContent = 'Installing JAR file...';
        // Wait for JARStore to be available and install
        return new Promise((resolve, reject) => {
          let attempts = 0;
          const maxAttempts = 100; // 10 seconds
          const checkAndInstall = setInterval(() => {
            attempts++;
            if (typeof JARStore !== 'undefined' && JARStore.installJAR) {
              clearInterval(checkAndInstall);
              // Install JAR using JARStore
              JARStore.installJAR(JAR_NAME, arrayBuffer).then(
                function() {
                  console.log('JAR installed successfully:', JAR_NAME);
                  updateStatus('✅ Installed', true);
                  if (loadingMsg) loadingMsg.textContent = 'JAR installed, starting...';
                  resolve();
                },
                function(error) {
                  console.error('Failed to install JAR:', error);
                  reject(new Error('JAR installation failed: ' + error));
                }
              );
            } else if (attempts >= maxAttempts) {
              clearInterval(checkAndInstall);
              reject(new Error('JARStore not available (timeout)'));
            }
          }, 100);
        });
      } catch (error) {
        console.error('Failed to download/install JAR:', error);
        updateStatus('❌ Error', false);
        if (loadingMsg) {
          loadingMsg.textContent = 'Error: ' + error.message;
          loadingMsg.style.color = '#ff4444';
        }
        throw error;
      }
    }
    // Auto start the game after JAR is ready
    async function autoStart() {
      const loadingMsg = document.getElementById(UNIQUE_ID + '-loading-msg');
      try {
        // Wait a bit for main function to be ready
        await new Promise(resolve => {
          let checkCount = 0;
          const checkMain = setInterval(() => {
            checkCount++;
            if (typeof start === 'function') {
              clearInterval(checkMain);
              resolve();
            } else if (checkCount > 50) { // 5 seconds timeout
              clearInterval(checkMain);
              resolve(); // Continue anyway
            }
          }, 100);
        });
        // Call main() to start the game
        if (typeof start === 'function') {
          console.log('Calling main() function with config:', config);
          updateStatus('🎮 Running', true);
          if (loadingMsg) loadingMsg.style.display = 'none';
          // Call main after a small delay to ensure everything is ready
          setTimeout(() => {
            start();
          }, 100);
        } else {
          throw new Error('main() function not found');
        }
      } catch (error) {
        console.error('Failed to start game:', error);
        updateStatus('❌ Start failed', false);
        if (loadingMsg) {
          loadingMsg.textContent = 'Failed to start: ' + error.message;
          loadingMsg.style.color = '#ff4444';
        }
      }
    }
    // Main initialization function
    async function initialize() {
      try {
        console.log('Starting J2ME initialization...');
        // Step 1: Download and install JAR
        await downloadAndInstallJAR();
        // Step 2: Auto start the game
        await autoStart();
      } catch (error) {
        console.error('Initialization failed:', error);
        const loadingMsg = document.getElementById(UNIQUE_ID + '-loading-msg');
        if (loadingMsg) {
          loadingMsg.innerHTML = 'Failed to load game<br><button onclick="location.reload()" style="margin-top:10px;padding:5px 10px;">Retry</button>';
          loadingMsg.style.color = '#ff4444';
        }
      }
    }
    // Start when page is fully loaded
    window.addEventListener('load', function() {
      console.log('Page loaded, starting J2ME initialization...');
      // Delay to ensure all scripts are loaded
      setTimeout(initialize, 1000);
    });
  })();
  </script>
</div>`,
		safePlayerID, // %[1]s
		safeJarPath,  // %[2]s
		safeFileName, // %[3]s - 包含 .jar 扩展名
		uniqueID,     // %[4]s
	)
}
