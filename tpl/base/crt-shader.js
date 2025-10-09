/**
 * Game CRT Shader Filter
 * 自动应用 CRT 效果到 .ejs_canvas 元素
 *
 * @author uk0
 * @version 1.0.0
 * @date 2025-01-10
 *
 * @example
 * // 基本使用
 * <script src="crt-shader.js"></script>
 * <script>
 *   const crt = new CRTShader();
 *   crt.enable();
 * </script>
 *
 * // 高级配置
 * const crt = new CRTShader({
 *   autoEnable: true,
 *   targetSelector: '.ejs_canvas',
 *   scanlineStrength: 0.8,
 *   curvature: 0.12
 * });
 */

class CRTShader {
    constructor(config = {}) {
        this.version = '1.0.0';
        this.enabled = false;
        this.canvasOverlays = new Map();
        this.observer = null;
        this.isProcessing = false;
        this.fpsCounter = 0;
        this.lastFpsTime = 0;

        // 默认配置
        this.config = {
            autoEnable: false,                    // 是否自动启用
            targetSelector: 'canvas.ejs_canvas',  // 目标 canvas 选择器
            logEnabled: true,                     // 是否启用日志
            ...config
        };

        // CRT 参数
        this.params = {
            intensity: 1.0,
            scanlineStrength: 0.8,
            scanlineCount: 1080,
            phosphorGlow: 0.3,
            rgbShift: 0.5,
            curvature: 0.12,
            vignette: 0.25,
            brightness: 1.1,
            contrast: 1.2,
            saturation: 1.1,
            gamma: 2.2,
            noiseAmount: 0.02,
            flickerAmount: 0.0,
            maskStrength: 0.4,
            // 从配置中覆盖参数
            ...this._extractParams(config)
        };

        // 预设配置
        this.presets = {
            arcade: {
                scanlineStrength: 1.0,
                phosphorGlow: 0.4,
                rgbShift: 0.3,
                contrast: 1.3,
                saturation: 1.2,
                maskStrength: 0.5,
                curvature: 0.08
            },
            console: {
                scanlineStrength: 0.7,
                phosphorGlow: 0.35,
                curvature: 0.15,
                vignette: 0.3,
                contrast: 1.15,
                maskStrength: 0.3
            },
            monitor: {
                scanlineStrength: 0.5,
                phosphorGlow: 0.2,
                curvature: 0.08,
                vignette: 0.15,
                maskStrength: 0.6
            },
            pvm: {
                scanlineStrength: 0.6,
                phosphorGlow: 0.25,
                curvature: 0.05,
                contrast: 1.25,
                saturation: 1.15,
                maskStrength: 0.7
            },
            gameboy: {
                scanlineStrength: 0.4,
                phosphorGlow: 0.15,
                curvature: 0.02,
                saturation: 0.3,
                gamma: 2.8,
                maskStrength: 0.2
            }
        };

        this._init();
    }

    // 提取 CRT 参数
    _extractParams(config) {
        const paramKeys = [
            'intensity', 'scanlineStrength', 'scanlineCount', 'phosphorGlow',
            'rgbShift', 'curvature', 'vignette', 'brightness', 'contrast',
            'saturation', 'gamma', 'noiseAmount', 'flickerAmount', 'maskStrength'
        ];

        const extracted = {};
        paramKeys.forEach(key => {
            if (config[key] !== undefined) {
                extracted[key] = config[key];
            }
        });

        return extracted;
    }

    _init() {
        this._log('init', `CRT Shader v${this.version} initialized`);

        if (this.config.autoEnable) {
            // 等待 DOM 加载完成后自动启用
            if (document.readyState === 'loading') {
                document.addEventListener('DOMContentLoaded', () => this.enable());
            } else {
                this.enable();
            }
        }

        // 页面卸载时清理
        window.addEventListener('beforeunload', () => this.destroy());
    }

    // 启用 CRT 滤镜
    enable() {
        if (this.enabled) {
            this._log('warn', 'CRT Shader already enabled');
            return this;
        }

        this._log('info', 'Enabling CRT Shader...');
        this.enabled = true;

        this._startAutoDetection();

        const canvases = document.querySelectorAll(this.config.targetSelector);
        canvases.forEach(canvas => this._processCanvas(canvas));

        this._log('success', `CRT Shader enabled (${canvases.length} canvas found)`);

        return this;
    }

    // 禁用 CRT 滤镜
    disable() {
        if (!this.enabled) {
            this._log('warn', 'CRT Shader already disabled');
            return this;
        }

        this.enabled = false;
        this._log('info', 'Disabling CRT Shader...');

        // 断开观察器
        if (this.observer) {
            this.observer.disconnect();
            this.observer = null;
        }

        // 清理所有渲染循环和 DOM
        this.canvasOverlays.forEach(processor => {
            if (processor.animationId) {
                cancelAnimationFrame(processor.animationId);
            }
            if (processor.overlay && processor.overlay.parentElement) {
                processor.overlay.parentElement.removeChild(processor.overlay);
            }
        });

        this.canvasOverlays.clear();

        this._log('success', 'CRT Shader disabled');

        return this;
    }

    // 切换启用/禁用
    toggle() {
        return this.enabled ? this.disable() : this.enable();
    }

    // 应用预设
    applyPreset(name) {
        const preset = this.presets[name];
        if (!preset) {
            this._log('error', `Unknown preset: ${name}`);
            return this;
        }

        this._log('info', `Applying preset: ${name}`);

        Object.keys(preset).forEach(key => {
            if (this.params[key] !== undefined) {
                this.params[key] = preset[key];
            }
        });

        return this;
    }

    // 设置参数
    setParam(key, value) {
        if (this.params[key] === undefined) {
            this._log('error', `Unknown parameter: ${key}`);
            return this;
        }

        this.params[key] = value;
        this._log('info', `Set ${key} = ${value}`);

        return this;
    }

    // 批量设置参数
    setParams(params) {
        Object.keys(params).forEach(key => {
            if (this.params[key] !== undefined) {
                this.params[key] = params[key];
            }
        });

        return this;
    }

    // 获取参数
    getParam(key) {
        return this.params[key];
    }

    // 获取所有参数
    getParams() {
        return { ...this.params };
    }

    // 重置为默认值
    reset() {
        this.params = {
            intensity: 1.0,
            scanlineStrength: 0.8,
            scanlineCount: 1080,
            phosphorGlow: 0.3,
            rgbShift: 0.5,
            curvature: 0.12,
            vignette: 0.25,
            brightness: 1.1,
            contrast: 1.2,
            saturation: 1.1,
            gamma: 2.2,
            noiseAmount: 0.02,
            flickerAmount: 0.0,
            maskStrength: 0.4
        };

        this._log('success', 'Parameters reset to defaults');

        return this;
    }

    // 获取 FPS
    getFPS() {
        return this.fpsCounter;
    }

    // 销毁实例
    destroy() {
        this._log('info', 'Destroying CRT Shader...');
        this.disable();
        this._log('success', 'CRT Shader destroyed');
    }

    // ==================== 内部方法 ====================

    _startAutoDetection() {
        this._detectCanvases();

        this.observer = new MutationObserver((mutations) => {
            if (this.isProcessing) return;

            let hasNewCanvas = false;
            for (const mutation of mutations) {
                for (const node of mutation.addedNodes) {
                    if (node.nodeType === 1) {
                        if (node.matches && node.matches(this.config.targetSelector)) {
                            hasNewCanvas = true;
                            break;
                        } else if (node.querySelectorAll) {
                            const canvases = node.querySelectorAll(this.config.targetSelector);
                            if (canvases.length > 0) {
                                hasNewCanvas = true;
                                break;
                            }
                        }
                    }
                }
                if (hasNewCanvas) break;
            }

            if (hasNewCanvas) {
                setTimeout(() => this._detectCanvases(), 50);
            }
        });

        this.observer.observe(document.body, {
            childList: true,
            subtree: true
        });
    }

    _detectCanvases() {
        if (this.isProcessing) return;

        this.isProcessing = true;

        try {
            const canvases = document.querySelectorAll(this.config.targetSelector);

            canvases.forEach(canvas => {
                if (!this.canvasOverlays.has(canvas) &&
                    !canvas.hasAttribute('data-crt-overlay') &&
                    this.enabled) {
                    this._processCanvas(canvas);
                }
            });
        } finally {
            this.isProcessing = false;
        }
    }

    _processCanvas(sourceCanvas) {
        if (this.canvasOverlays.has(sourceCanvas)) return;

        this._log('process', `Processing canvas: ${sourceCanvas.id || 'unnamed'}`);

        try {
            const overlay = document.createElement('canvas');
            overlay.setAttribute('data-crt-overlay', 'true');
            overlay.width = sourceCanvas.width;
            overlay.height = sourceCanvas.height;
            overlay.style.cssText = `
                position: absolute;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                pointer-events: none;
                z-index: 10;
            `;

            const wrapper = sourceCanvas.parentElement;
            if (wrapper) {
                const currentPosition = window.getComputedStyle(wrapper).position;
                if (currentPosition === 'static') {
                    wrapper.style.position = 'relative';
                }
                wrapper.appendChild(overlay);
            }

            const gl = overlay.getContext('webgl', {
                alpha: true,
                antialias: false,
                preserveDrawingBuffer: false,
                powerPreference: 'high-performance'
            });

            if (!gl) {
                this._log('error', 'WebGL not supported');
                return;
            }

            const program = this._compileShaders(gl);
            if (!program) return;

            const processor = {
                sourceCanvas,
                overlay,
                gl,
                program,
                texture: this._createTexture(gl),
                vertexBuffer: this._createVertexBuffer(gl, program),
                animationId: null
            };

            this._startRenderLoop(processor);
            this.canvasOverlays.set(sourceCanvas, processor);

            this._log('success', 'Canvas processing started');

        } catch (error) {
            this._log('error', `Failed to process canvas: ${error.message}`);
        }
    }

    _startRenderLoop(processor) {
        const render = () => {
            if (!this.enabled) {
                processor.animationId = null;
                return;
            }

            try {
                this._renderFrame(processor);
                this.fpsCounter++;
            } catch (error) {
                this._log('error', `Render error: ${error.message}`);
                return;
            }

            processor.animationId = requestAnimationFrame(render);
        };

        render();
    }

    _renderFrame(processor) {
        const { gl, program, sourceCanvas, overlay, texture } = processor;

        gl.bindTexture(gl.TEXTURE_2D, texture);
        gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, sourceCanvas);

        gl.viewport(0, 0, overlay.width, overlay.height);
        gl.clearColor(0, 0, 0, 0);
        gl.clear(gl.COLOR_BUFFER_BIT);

        gl.useProgram(program);

        gl.uniform1i(gl.getUniformLocation(program, 'u_texture'), 0);
        gl.uniform2f(gl.getUniformLocation(program, 'u_resolution'), overlay.width, overlay.height);
        gl.uniform1f(gl.getUniformLocation(program, 'u_time'), performance.now());

        Object.keys(this.params).forEach(param => {
            gl.uniform1f(gl.getUniformLocation(program, `u_${param}`), this.params[param]);
        });

        gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4);
    }

    _compileShaders(gl) {
        const vertexShaderSource = `
            attribute vec2 a_position;
            attribute vec2 a_texCoord;
            varying vec2 v_texCoord;

            void main() {
                gl_Position = vec4(a_position, 0.0, 1.0);
                v_texCoord = a_texCoord;
            }
        `;

        const fragmentShaderSource = `
            precision highp float;

            uniform sampler2D u_texture;
            uniform vec2 u_resolution;
            uniform float u_time;
            uniform float u_intensity;
            uniform float u_scanlineStrength;
            uniform float u_scanlineCount;
            uniform float u_phosphorGlow;
            uniform float u_rgbShift;
            uniform float u_curvature;
            uniform float u_vignette;
            uniform float u_brightness;
            uniform float u_contrast;
            uniform float u_saturation;
            uniform float u_gamma;
            uniform float u_noiseAmount;
            uniform float u_flickerAmount;
            uniform float u_maskStrength;

            varying vec2 v_texCoord;
            const float PI = 3.14159265359;

            float random(vec2 co) {
                return fract(sin(dot(co.xy + u_time * 0.0001, vec2(12.9898, 78.233))) * 43758.5453);
            }

            vec2 curveScreen(vec2 uv) {
                vec2 offset = uv - 0.5;
                offset *= vec2(1.0 + u_curvature * offset.y * offset.y,
                              1.0 + u_curvature * offset.x * offset.x);
                return offset + 0.5;
            }

            float scanlines(vec2 coord) {
                float line = sin(coord.y * PI * u_scanlineCount);
                return 1.0 - u_scanlineStrength * (1.0 - line * line);
            }

            vec3 rgbShift(sampler2D tex, vec2 coord) {
                float shift = u_rgbShift * 0.003;
                return vec3(
                    texture2D(tex, coord + vec2(shift, 0.0)).r,
                    texture2D(tex, coord).g,
                    texture2D(tex, coord - vec2(shift, 0.0)).b
                );
            }

            vec3 phosphorGlow(vec3 color, vec2 coord) {
                float dist = length(coord - 0.5);
                float glow = (1.0 - smoothstep(0.0, 0.7, dist)) * u_phosphorGlow;
                return color + color * glow;
            }

            float vignette(vec2 coord) {
                float dist = length(coord - 0.5);
                return 1.0 - u_vignette * smoothstep(0.3, 0.9, dist);
            }

            float noise(vec2 coord) {
                return (random(coord * 100.0) - 0.5) * u_noiseAmount;
            }

            vec3 rgbMask(vec2 coord, vec3 color) {
                if (u_maskStrength <= 0.0) return color;

                float x = coord.x * u_resolution.x;
                float mask = mod(x, 3.0);
                vec3 maskColor = vec3(1.0);

                if (mask < 1.0) {
                    maskColor = vec3(1.0, 1.0 - u_maskStrength, 1.0 - u_maskStrength);
                } else if (mask < 2.0) {
                    maskColor = vec3(1.0 - u_maskStrength, 1.0, 1.0 - u_maskStrength);
                } else {
                    maskColor = vec3(1.0 - u_maskStrength, 1.0 - u_maskStrength, 1.0);
                }

                return color * mix(vec3(1.0), maskColor, u_intensity * 0.6);
            }

            vec3 colorGrade(vec3 color) {
                color = (color - 0.5) * u_contrast + 0.5;
                color *= u_brightness;
                vec3 gray = vec3(dot(color, vec3(0.299, 0.587, 0.114)));
                color = mix(gray, color, u_saturation);
                return pow(max(color, vec3(0.0)), vec3(1.0 / u_gamma));
            }

            void main() {
                vec2 coord = curveScreen(v_texCoord);

                if (coord.x < 0.0 || coord.x > 1.0 || coord.y < 0.0 || coord.y > 1.0) {
                    gl_FragColor = vec4(0.0, 0.0, 0.0, 0.0);
                    return;
                }

                vec3 color = rgbShift(u_texture, coord);
                color = phosphorGlow(color, coord);
                color *= scanlines(coord);
                color *= vignette(coord);
                color = rgbMask(coord, color);
                color = colorGrade(color);
                color += vec3(noise(coord));

                if (u_flickerAmount > 0.0) {
                    color *= 1.0 + sin(u_time * 0.06) * u_flickerAmount;
                }

                gl_FragColor = vec4(clamp(color, 0.0, 1.0), 1.0);
            }
        `;

        const vertexShader = this._createShader(gl, gl.VERTEX_SHADER, vertexShaderSource);
        const fragmentShader = this._createShader(gl, gl.FRAGMENT_SHADER, fragmentShaderSource);

        if (!vertexShader || !fragmentShader) return null;

        const program = gl.createProgram();
        gl.attachShader(program, vertexShader);
        gl.attachShader(program, fragmentShader);
        gl.linkProgram(program);

        if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
            this._log('error', 'Shader linking failed');
            return null;
        }

        return program;
    }

    _createShader(gl, type, source) {
        const shader = gl.createShader(type);
        gl.shaderSource(shader, source);
        gl.compileShader(shader);

        if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
            this._log('error', 'Shader compilation failed');
            gl.deleteShader(shader);
            return null;
        }

        return shader;
    }

    _createVertexBuffer(gl, program) {
        const vertices = new Float32Array([
            -1, -1,  0, 1,
            1, -1,  1, 1,
            -1,  1,  0, 0,
            1,  1,  1, 0
        ]);

        const buffer = gl.createBuffer();
        gl.bindBuffer(gl.ARRAY_BUFFER, buffer);
        gl.bufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW);

        const posLoc = gl.getAttribLocation(program, 'a_position');
        const texLoc = gl.getAttribLocation(program, 'a_texCoord');

        gl.enableVertexAttribArray(posLoc);
        gl.vertexAttribPointer(posLoc, 2, gl.FLOAT, false, 16, 0);

        gl.enableVertexAttribArray(texLoc);
        gl.vertexAttribPointer(texLoc, 2, gl.FLOAT, false, 16, 8);

        return buffer;
    }

    _createTexture(gl) {
        const texture = gl.createTexture();
        gl.bindTexture(gl.TEXTURE_2D, texture);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
        return texture;
    }

    _log(type, message) {
        if (!this.config.logEnabled) return;

        const prefix = {
            'init': '[CRT]',
            'info': '[CRT]',
            'warn': '[CRT WARN]',
            'error': '[CRT ERROR]',
            'success': '[CRT]',
            'process': '[CRT]'
        }[type] || '[CRT]';

        const method = type === 'error' ? 'error' : type === 'warn' ? 'warn' : 'log';
        console[method](`${prefix} ${message}`);
    }
}

// 支持 CommonJS/AMD/UMD
if (typeof module !== 'undefined' && typeof module.exports !== 'undefined') {
    module.exports = CRTShader;
} else if (typeof define === 'function' && define.amd) {
    define([], function() {
        return CRTShader;
    });
} else {
    window.CRTShader = CRTShader;
}



// auto loading


//  Crt Shader Initialization
function crt_shader_init() {
    console.log('%c🔌 Initializing CRT Shader......',
        'color: #00ff00; font-weight: bold; text-shadow: 0 0 5px #00ff00');
    var canvas_count = document.getElementsByClassName("ejs_canvas").length
    console.log(`%c🔌 Game Canvas count: ${canvas_count}`,
        'color: #00ff00; font-weight: bold; text-shadow: 0 0 5px #00ff00');
    if (canvas_count === 1) {
        game_crt();
    }

    function game_crt() {
        // 创建实例并启用
        const crt = new CRTShader();
        crt.enable();
        // 批量设置
        crt.setParams({
            intensity: 1.0,
            scanlineStrength: 0.8,
            scanlineCount: 750,// 根据分辨率调整600线一般够用了
            phosphorGlow: 0.3,
            rgbShift: 0.5,
            curvature: 0.12,
            vignette: 0.23,
            brightness: 1.1,
            contrast: 1.2,
            saturation: 1.1,
            gamma: 2.2,
            noiseAmount: 0.02,
            flickerAmount: 0.0,
            maskStrength: 0.4
        });
    }
}


// 启动提示
setTimeout(() => {
    crt_shader_init()
}, 1000);
