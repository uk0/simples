/**
 * Game CRT Shader Filter (Fixed Double Vision)
 * 自动应用 CRT 效果到 .ejs_canvas 元素
 *
 * @author uk0
 * @version 1.3.0
 * @date 2025-01-10
 *
 * @example
 * const crt = new CRTShader({ hideOriginal: true });
 * crt.enable();
 */

class CRTShader {
    constructor(config = {}) {
        this.version = '1.3.0';
        this.enabled = false;
        this.canvasOverlays = new Map();
        this.observer = null;
        this.isProcessing = false;
        this.fpsCounter = 0;
        this.lastFpsTime = 0;

        // 默认配置
        this.config = {
            autoEnable: false,
            targetSelector: 'canvas.ejs_canvas',
            logEnabled: true,
            hideOriginal: true,            // 🔧 隐藏原始canvas
            replaceMode: 'opacity',        // 🔧 替换模式: 'opacity', 'visibility', 'replace'
            ...config
        };

        // CRT 参数
        this.params = {
            intensity: 1.0,
            scanlineStrength: 0.6,
            scanlineCount: 800,
            phosphorGlow: 0.1,
            rgbShift: 0.15,
            curvature: 0.08,
            vignette: 0.15,
            brightness: 1.05,
            contrast: 1.15,
            saturation: 1.05,
            gamma: 2.2,
            noiseAmount: 0.0,
            flickerAmount: 0.0,
            maskStrength: 0.15,
            ...this._extractParams(config)
        };

        // 预设配置
        this.presets = {
            arcade: {
                scanlineStrength: 0.9,
                phosphorGlow: 0.25,
                rgbShift: 0.2,
                contrast: 1.25,
                saturation: 1.15,
                maskStrength: 0.35,
                curvature: 0.1
            },
            console: {
                scanlineStrength: 0.6,
                phosphorGlow: 0.2,
                curvature: 0.12,
                vignette: 0.2,
                contrast: 1.1,
                maskStrength: 0.2
            },
            monitor: {
                scanlineStrength: 0.4,
                phosphorGlow: 0.1,
                curvature: 0.05,
                vignette: 0.1,
                maskStrength: 0.4
            },
            sharp: {
                scanlineStrength: 0.5,
                phosphorGlow: 0.05,
                rgbShift: 0.1,
                curvature: 0.03,
                vignette: 0.05,
                contrast: 1.2,
                maskStrength: 0.1,
                noiseAmount: 0.0
            },
            pvm: {
                scanlineStrength: 0.5,
                phosphorGlow: 0.15,
                curvature: 0.04,
                contrast: 1.2,
                saturation: 1.1,
                maskStrength: 0.5
            }
        };

        this._init();
    }

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
        this._log('init', `CRT Shader v${this.version} (Replace Mode: ${this.config.replaceMode})`);

        if (this.config.autoEnable) {
            if (document.readyState === 'loading') {
                document.addEventListener('DOMContentLoaded', () => this.enable());
            } else {
                this.enable();
            }
        }

        window.addEventListener('beforeunload', () => this.destroy());
    }

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

    disable() {
        if (!this.enabled) {
            this._log('warn', 'CRT Shader already disabled');
            return this;
        }

        this.enabled = false;
        this._log('info', 'Disabling CRT Shader...');

        if (this.observer) {
            this.observer.disconnect();
            this.observer = null;
        }

        // 🔧 恢复原始canvas的显示
        this.canvasOverlays.forEach(processor => {
            if (processor.animationId) {
                cancelAnimationFrame(processor.animationId);
            }

            // 恢复原始canvas
            if (processor.sourceCanvas) {
                this._restoreOriginalCanvas(processor.sourceCanvas);
            }

            // 移除覆盖层
            if (processor.overlay && processor.overlay.parentElement) {
                processor.overlay.parentElement.removeChild(processor.overlay);
            }

            // 如果是替换模式，恢复原始canvas到DOM
            if (this.config.replaceMode === 'replace' && processor.originalParent && processor.sourceCanvas) {
                processor.originalParent.appendChild(processor.sourceCanvas);
            }
        });

        this.canvasOverlays.clear();

        this._log('success', 'CRT Shader disabled');

        return this;
    }

    toggle() {
        return this.enabled ? this.disable() : this.enable();
    }

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

    setParam(key, value) {
        if (this.params[key] === undefined) {
            this._log('error', `Unknown parameter: ${key}`);
            return this;
        }

        this.params[key] = value;
        return this;
    }

    setParams(params) {
        Object.keys(params).forEach(key => {
            if (this.params[key] !== undefined) {
                this.params[key] = params[key];
            }
        });

        return this;
    }

    getParam(key) {
        return this.params[key];
    }

    getParams() {
        return { ...this.params };
    }

    reset() {
        this.params = {
            intensity: 1.0,
            scanlineStrength: 0.6,
            scanlineCount: 800,
            phosphorGlow: 0.1,
            rgbShift: 0.15,
            curvature: 0.08,
            vignette: 0.15,
            brightness: 1.05,
            contrast: 1.15,
            saturation: 1.05,
            gamma: 2.2,
            noiseAmount: 0.0,
            flickerAmount: 0.0,
            maskStrength: 0.15
        };

        this._log('success', 'Parameters reset to defaults');
        return this;
    }

    getFPS() {
        return this.fpsCounter;
    }

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

    // 🔧 隐藏原始canvas
    _hideOriginalCanvas(canvas) {
        if (!this.config.hideOriginal) return;

        // 保存原始样式
        canvas.setAttribute('data-original-opacity', canvas.style.opacity || '');
        canvas.setAttribute('data-original-visibility', canvas.style.visibility || '');
        canvas.setAttribute('data-original-filter', canvas.style.filter || '');

        switch (this.config.replaceMode) {
            case 'opacity':
                // 🔧 使用opacity=0隐藏，保持布局
                canvas.style.opacity = '0';
                break;

            case 'visibility':
                // 🔧 使用visibility隐藏，保持布局
                canvas.style.visibility = 'hidden';
                break;

            case 'replace':
                // 🔧 完全移除原始canvas（最彻底）
                // 稍后在处理中执行
                break;

            default:
                canvas.style.opacity = '0';
        }
    }

    // 🔧 恢复原始canvas
    _restoreOriginalCanvas(canvas) {
        if (!canvas) return;

        const originalOpacity = canvas.getAttribute('data-original-opacity') || '';
        const originalVisibility = canvas.getAttribute('data-original-visibility') || '';
        const originalFilter = canvas.getAttribute('data-original-filter') || '';

        canvas.style.opacity = originalOpacity;
        canvas.style.visibility = originalVisibility;
        canvas.style.filter = originalFilter;

        canvas.removeAttribute('data-original-opacity');
        canvas.removeAttribute('data-original-visibility');
        canvas.removeAttribute('data-original-filter');
    }

    _processCanvas(sourceCanvas) {
        if (this.canvasOverlays.has(sourceCanvas)) return;

        this._log('process', `Processing canvas: ${sourceCanvas.id || 'unnamed'}`);

        try {
            const overlay = document.createElement('canvas');
            overlay.setAttribute('data-crt-overlay', 'true');
            overlay.width = sourceCanvas.width;
            overlay.height = sourceCanvas.height;

            // 🔧 复制原始canvas的样式
            const sourceStyle = window.getComputedStyle(sourceCanvas);
            overlay.style.cssText = `
                position: ${sourceStyle.position === 'static' ? 'relative' : sourceStyle.position};
                top: ${sourceStyle.top};
                left: ${sourceStyle.left};
                width: ${sourceStyle.width};
                height: ${sourceStyle.height};
                pointer-events: ${sourceStyle.pointerEvents};
                z-index: ${parseInt(sourceStyle.zIndex || '0') + 1};
                image-rendering: -webkit-optimize-contrast;
                image-rendering: crisp-edges;
                display: ${sourceStyle.display};
                margin: ${sourceStyle.margin};
                padding: ${sourceStyle.padding};
            `;

            const wrapper = sourceCanvas.parentElement;
            let originalParent = null;

            if (this.config.replaceMode === 'replace') {
                // 🔧 替换模式：完全替换原始canvas
                originalParent = wrapper;
                wrapper.insertBefore(overlay, sourceCanvas);
                wrapper.removeChild(sourceCanvas);

                // 将原始canvas保存在内存中，但不在DOM中
                overlay.sourceCanvas = sourceCanvas;
            } else {
                // 🔧 覆盖模式：在原始canvas上方添加覆盖层
                if (wrapper) {
                    const currentPosition = window.getComputedStyle(wrapper).position;
                    if (currentPosition === 'static') {
                        wrapper.style.position = 'relative';
                    }

                    // 确保overlay在sourceCanvas之后插入
                    if (sourceCanvas.nextSibling) {
                        wrapper.insertBefore(overlay, sourceCanvas.nextSibling);
                    } else {
                        wrapper.appendChild(overlay);
                    }

                    // 使overlay完全覆盖sourceCanvas
                    overlay.style.position = 'absolute';
                    overlay.style.top = '0';
                    overlay.style.left = '0';
                }

                // 🔧 隐藏原始canvas
                this._hideOriginalCanvas(sourceCanvas);
            }

            // 🔧 创建WebGL上下文（完全不透明）
            const gl = overlay.getContext('webgl', {
                alpha: false,                      // 🔧 关闭alpha通道
                antialias: false,
                premultipliedAlpha: false,
                preserveDrawingBuffer: false,
                powerPreference: 'high-performance',
                desynchronized: true,
                failIfMajorPerformanceCaveat: false
            }) || overlay.getContext('experimental-webgl', {
                alpha: false,
                antialias: false,
                premultipliedAlpha: false,
                preserveDrawingBuffer: false
            });

            if (!gl) {
                this._log('error', 'WebGL not supported');
                return;
            }

            // 🔧 设置WebGL混合模式为不透明
            gl.disable(gl.BLEND);
            gl.disable(gl.DEPTH_TEST);
            gl.disable(gl.STENCIL_TEST);

            const program = this._compileShaders(gl);
            if (!program) return;

            const processor = {
                sourceCanvas: this.config.replaceMode === 'replace' ? overlay.sourceCanvas : sourceCanvas,
                overlay,
                gl,
                program,
                texture: this._createTexture(gl),
                vertexBuffer: this._createVertexBuffer(gl, program),
                animationId: null,
                originalParent
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

        processor.animationId = requestAnimationFrame(render);
    }

    _renderFrame(processor) {
        const { gl, program, sourceCanvas, overlay, texture } = processor;

        // 🔧 设置视口
        gl.viewport(0, 0, overlay.width, overlay.height);

        // 🔧 使用不透明的黑色清除
        gl.clearColor(0, 0, 0, 1);
        gl.clear(gl.COLOR_BUFFER_BIT);

        gl.bindTexture(gl.TEXTURE_2D, texture);

        // 🔧 优化纹理上传
        gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, false);
        gl.pixelStorei(gl.UNPACK_PREMULTIPLY_ALPHA_WEBGL, false);
        gl.pixelStorei(gl.UNPACK_ALIGNMENT, 4);

        gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGB, gl.RGB, gl.UNSIGNED_BYTE, sourceCanvas);

        gl.useProgram(program);

        // 设置uniforms
        gl.uniform1i(gl.getUniformLocation(program, 'u_texture'), 0);
        gl.uniform2f(gl.getUniformLocation(program, 'u_resolution'), overlay.width, overlay.height);
        gl.uniform1f(gl.getUniformLocation(program, 'u_time'), performance.now());

        Object.keys(this.params).forEach(param => {
            gl.uniform1f(gl.getUniformLocation(program, `u_${param}`), this.params[param]);
        });

        gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4);

        // 🔧 强制刷新
        gl.flush();
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

        // 🔧 优化的片段着色器（完全不透明输出）
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
                highp float a = 12.9898;
                highp float b = 78.233;
                highp float c = 43758.5453;
                highp float dt = dot(co.xy, vec2(a,b));
                highp float sn = mod(dt, 3.14);
                return fract(sin(sn) * c);
            }

            vec2 curveScreen(vec2 uv) {
                if (u_curvature <= 0.001) return uv;
                vec2 offset = uv - 0.5;
                float r2 = offset.x * offset.x + offset.y * offset.y;
                offset *= 1.0 + u_curvature * r2;
                return offset + 0.5;
            }

            float scanlines(vec2 coord) {
                if (u_scanlineStrength <= 0.001) return 1.0;
                float line = sin(coord.y * PI * u_scanlineCount);
                return mix(1.0, line * line, u_scanlineStrength);
            }

            vec3 rgbShift(sampler2D tex, vec2 coord) {
                if (u_rgbShift <= 0.001) return texture2D(tex, coord).rgb;
                
                float shift = u_rgbShift * 0.0015;
                vec3 color;
                color.r = texture2D(tex, coord + vec2(shift, 0.0)).r;
                color.g = texture2D(tex, coord).g;
                color.b = texture2D(tex, coord - vec2(shift, 0.0)).b;
                return color;
            }

            vec3 phosphorGlow(vec3 color) {
                if (u_phosphorGlow <= 0.001) return color;
                return color * (1.0 + u_phosphorGlow * 0.3);
            }

            float vignette(vec2 coord) {
                if (u_vignette <= 0.001) return 1.0;
                float dist = length(coord - 0.5) * 1.414;
                return 1.0 - u_vignette * smoothstep(0.5, 1.0, dist);
            }

            vec3 rgbMask(vec2 coord, vec3 color) {
                if (u_maskStrength <= 0.001) return color;
                
                float x = coord.x * u_resolution.x;
                float mask = mod(x, 3.0);
                vec3 maskColor = vec3(1.0);
                
                if (mask < 1.0) {
                    maskColor.r = 1.0;
                    maskColor.g = 1.0 - u_maskStrength * 0.5;
                    maskColor.b = 1.0 - u_maskStrength * 0.5;
                } else if (mask < 2.0) {
                    maskColor.r = 1.0 - u_maskStrength * 0.5;
                    maskColor.g = 1.0;
                    maskColor.b = 1.0 - u_maskStrength * 0.5;
                } else {
                    maskColor.r = 1.0 - u_maskStrength * 0.5;
                    maskColor.g = 1.0 - u_maskStrength * 0.5;
                    maskColor.b = 1.0;
                }
                
                return color * maskColor;
            }

            vec3 colorGrade(vec3 color) {
                color = mix(vec3(0.5), color, u_contrast);
                color *= u_brightness;
                
                float gray = dot(color, vec3(0.299, 0.587, 0.114));
                color = mix(vec3(gray), color, u_saturation);
                
                return pow(max(color, vec3(0.0)), vec3(1.0 / u_gamma));
            }

            void main() {
                vec2 coord = curveScreen(v_texCoord);

                // 🔧 边界处理：使用黑色而不是透明
                if (coord.x < 0.0 || coord.x > 1.0 || coord.y < 0.0 || coord.y > 1.0) {
                    gl_FragColor = vec4(0.0, 0.0, 0.0, 1.0);
                    return;
                }

                vec3 color = rgbShift(u_texture, coord);
                color = phosphorGlow(color);
                color *= scanlines(coord);
                color *= vignette(coord);
                color = rgbMask(coord, color);
                color = colorGrade(color);
                
                if (u_noiseAmount > 0.001) {
                    float noise = (random(coord + u_time * 0.001) - 0.5) * u_noiseAmount;
                    color += vec3(noise);
                }
                
                color = clamp(color, 0.0, 1.0);
                
                // 🔧 完全不透明输出
                gl_FragColor = vec4(color, 1.0);
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
            this._log('error', `Shader compilation failed: ${gl.getShaderInfoLog(shader)}`);
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

// ==================== Auto Loading ====================

function crt_shader_init() {
    console.log('%c🔌 Initializing CRT Shader (Fixed Double Vision v1.3.0)......',
        'color: #00ff00; font-weight: bold; text-shadow: 0 0 5px #00ff00');

    var canvas_count = document.getElementsByClassName("ejs_canvas").length;
    console.log(`%c🔌 Game Canvas count: ${canvas_count}`,
        'color: #00ff00; font-weight: bold; text-shadow: 0 0 5px #00ff00');

    if (canvas_count >= 1) {
        game_crt();
    }

    function game_crt() {
        // 🔧 使用隐藏原始canvas模式
        const crt = new CRTShader({
            hideOriginal: true,
            replaceMode: 'opacity'  // 可选: 'opacity', 'visibility', 'replace'
        });

        crt.enable();

        // 优化的参数设置
        crt.setParams({
            intensity: 1.0,
            scanlineStrength: 0.6,
            scanlineCount: 600,
            phosphorGlow: 0.1,
            rgbShift: 0.15,
            curvature: 0.08,
            vignette: 0.15,
            brightness: 1.05,
            contrast: 1.15,
            saturation: 1.05,
            gamma: 2.2,
            noiseAmount: 0.0,
            flickerAmount: 0.0,
            maskStrength: 0.15
        });

        console.log('%c✅ CRT Shader initialized (Original Canvas Hidden)',
            'color: #00ff00; font-weight: bold');
    }
}

// 启动
setTimeout(() => {
    crt_shader_init();
}, 1000);