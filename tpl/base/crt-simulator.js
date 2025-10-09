/**
 * CRT WebGL FW900 Complete Edition - 完整Sony FW900 CRT模拟
 * 基于Sony FW900专业CRT显示器的精确WebGL实现
 * @author uk0
 * @version 20.0.0
 * @date 2025-01-10
 */

(function(window) {
    'use strict';

    class CRTWebGLFW900 {
        constructor() {
            this.version = '20.0.0';
            this.enabled = false;
            this.gl = null;
            this.canvas = null;
            this.program = null;
            this.frameTexture = null;
            this.startTime = Date.now();
            this.frameCount = 0;

            // Sony FW900 技术规格
            this.fw900Specs = {
                model: 'Sony GDM-FW900',
                type: 'Trinitron FD',
                size: '24 inch',
                maxResolution: '2304x1440',
                refreshRate: 85,  // 默认85Hz
                dotPitch: 0.24,    // mm
                phosphorType: 'P22',
                tubeType: 'Flat Trinitron'
            };

            // 可调参数
            this.config = {
                // 基础显示参数
                refreshRate: 85,
                resolution: 1.0,

                // CRT效果强度控制
                crtIntensity: 1.0,  // 总体CRT效果强度 (0-2, 1为标准)

                // 扫描线
                scanlineIntensity: 0.15,
                scanlineCount: 1080,
                scanlineThickness: 0.45,
                interlaced: false,

                // Trinitron特性
                trinitronWires: true,
                wirePositions: [0.333, 0.667],
                wireIntensity: 0.08,

                // 磷光特性
                phosphorPersistence: 0.15,
                phosphorBloom: 0.25,
                subpixelLayout: 'rgb', // rgb, bgr, grille

                // 电子束
                beamIntensity: 0.03,
                beamSpeed: 85,
                beamWidth: 0.02,

                // 画质调整
                brightness: 1.02,
                contrast: 1.18,
                saturation: 0.98,
                sharpness: 0.92,
                gamma: 2.2,

                // 几何畸变
                curvature: 0.02,      // FW900几乎是平的
                cornerRadius: 0.03,
                vignette: 0.25,

                // 色彩
                colorTemperature: 6500,
                tint: 0,
                convergence: 0.001,

                // 信号特性
                noiseAmount: 0.01,
                jitter: 0.001,
                ghosting: 0.02,
                halation: 0.05
            };

            // 增强的顶点着色器
            this.vertexShader = `
                attribute vec2 a_position;
                attribute vec2 a_texCoord;
                varying vec2 v_texCoord;
                varying vec2 v_screenCoord;
                
                void main() {
                    gl_Position = vec4(a_position, 0.0, 1.0);
                    v_texCoord = a_texCoord;
                    v_screenCoord = a_position * 0.5 + 0.5;
                }
            `;

            // 完整的FW900片段着色器
            this.fragmentShader = `
                precision highp float;
                
                // Uniforms
                uniform float u_time;
                uniform vec2 u_resolution;
                uniform sampler2D u_texture;
                uniform sampler2D u_frameBuffer;
                
                // 配置参数
                uniform float u_scanlineIntensity;
                uniform float u_scanlineCount;
                uniform float u_phosphorPersistence;
                uniform float u_brightness;
                uniform float u_contrast;
                uniform float u_saturation;
                uniform float u_gamma;
                uniform float u_curvature;
                uniform float u_vignette;
                uniform float u_noiseAmount;
                uniform float u_convergence;
                uniform float u_halation;
                uniform float u_crtIntensity;
                
                // Varyings
                varying vec2 v_texCoord;
                varying vec2 v_screenCoord;
                
                // 常量
                const float PI = 3.14159265359;
                const vec3 kColorTemperature6500K = vec3(1.0, 0.9549, 0.9171);
                
                // CRT屏幕曲率变换（FW900是近乎平坦的）
                vec2 curveRemapUV(vec2 uv) {
                    vec2 offset = uv - 0.5;
                    float curveFactor = u_curvature;
                    
                    // FW900的轻微曲率
                    offset *= vec2(
                        1.0 + curveFactor * offset.y * offset.y,
                        1.0 + curveFactor * offset.x * offset.x
                    );
                    
                    // 角落圆角
                    float cornerRadius = 0.03;
                    offset *= 1.0 - cornerRadius * dot(offset, offset);
                    
                    return offset + 0.5;
                }
                
                // 噪声生成
                float random(vec2 co) {
                    return fract(sin(dot(co.xy + u_time * 0.001, vec2(12.9898, 78.233))) * 43758.5453);
                }
                
                // 2D噪声
                float noise(vec2 p) {
                    vec2 i = floor(p);
                    vec2 f = fract(p);
                    f = f * f * (3.0 - 2.0 * f);
                    
                    float a = random(i);
                    float b = random(i + vec2(1.0, 0.0));
                    float c = random(i + vec2(0.0, 1.0));
                    float d = random(i + vec2(1.0, 1.0));
                    
                    return mix(mix(a, b, f.x), mix(c, d, f.x), f.y);
                }
                
                // Trinitron RGB荧光粉排列
                vec3 trinitronPhosphor(vec2 uv, vec3 color) {
                    float pixelX = uv.x * u_resolution.x;
                    float subpixel = mod(pixelX, 3.0);
                    
                    vec3 mask = vec3(0.0);
                    
                    // RGB垂直条纹（Trinitron特征）
                    if (subpixel < 1.0) {
                        mask.r = 1.0;
                    } else if (subpixel < 2.0) {
                        mask.g = 1.0;
                    } else {
                        mask.b = 1.0;
                    }
                    
                    // 添加磷光粉间隙
                    mask *= 0.9 + 0.1 * sin(pixelX * PI * 2.0);
                    
                    return color * mask;
                }
                
                // 扫描线生成（支持隔行扫描）
                float scanline(vec2 uv) {
                    float y = uv.y * u_scanlineCount;
                    float scanline = sin(y * PI * 2.0);
                    scanline = (scanline + 1.0) * 0.5;
                    
                    // 模拟电子束宽度变化
                    float beamWidth = 0.7 + 0.3 * sin(u_time * 5.0 + y);
                    scanline = pow(scanline, beamWidth);
                    
                    return mix(1.0 - u_scanlineIntensity, 1.0, scanline);
                }
                
                // Trinitron支撑线
                float trinitronWires(vec2 uv) {
                    float wire = 1.0;
                    
                    // FW900有两条支撑线，在1/3和2/3位置
                    if (abs(uv.y - 0.333) < 0.001 || abs(uv.y - 0.667) < 0.001) {
                        wire = 1.0 - 0.08; // 支撑线强度
                    }
                    
                    return wire;
                }
                
                // 色彩会聚误差（轻微的RGB偏移）
                vec3 convergenceError(vec2 uv, sampler2D tex) {
                    vec3 color;
                    float offset = u_convergence;
                    
                    color.r = texture2D(tex, uv + vec2(offset, 0.0)).r;
                    color.g = texture2D(tex, uv).g;
                    color.b = texture2D(tex, uv - vec2(offset, 0.0)).b;
                    
                    return color;
                }
                
                // 光晕效果（亮区溢出）- 修复循环问题
                vec3 halation(vec2 uv, vec3 color) {
                    vec3 halo = vec3(0.0);
                    const int SAMPLES = 5;  // 使用常量定义循环次数
                    
                    for (int i = 0; i < SAMPLES; i++) {
                        float angle = (float(i) / float(SAMPLES)) * PI * 2.0;
                        vec2 offset = vec2(cos(angle), sin(angle)) * u_halation * 0.01;
                        vec3 sample = texture2D(u_texture, uv + offset).rgb;
                        
                        // 只对亮区产生光晕
                        float brightness = dot(sample, vec3(0.299, 0.587, 0.114));
                        if (brightness > 0.8) {
                            halo += sample * (brightness - 0.8);
                        }
                    }
                    
                    return color + (halo / float(SAMPLES)) * 0.3;
                }
                
                // 暗角效果
                float vignette(vec2 uv) {
                    vec2 center = uv - 0.5;
                    float dist = length(center);
                    float vignette = smoothstep(0.8, 0.4, dist);
                    return mix(1.0 - u_vignette, 1.0, vignette);
                }
                
                // 电子束扫描效果
                float electronBeam(vec2 uv) {
                    float beamY = mod(u_time * 0.1, 1.0);
                    float beam = 1.0 - abs(uv.y - beamY) * 50.0;
                    beam = clamp(beam, 0.0, 1.0);
                    return 1.0 + beam * 0.02;
                }
                
                // 色彩调整
                vec3 colorGrading(vec3 color) {
                    // 对比度
                    color = (color - 0.5) * u_contrast + 0.5;
                    
                    // 亮度
                    color *= u_brightness;
                    
                    // 饱和度
                    vec3 gray = vec3(dot(color, vec3(0.299, 0.587, 0.114)));
                    color = mix(gray, color, u_saturation);
                    
                    // Gamma校正
                    color = pow(color, vec3(1.0 / u_gamma));
                    
                    // 色温调整（FW900标准6500K）
                    color *= kColorTemperature6500K;
                    
                    return color;
                }
                
                // 信号噪声和干扰
                vec3 signalNoise(vec2 uv, vec3 color) {
                    // 随机噪点
                    float staticNoise = random(uv) * u_noiseAmount;
                    
                    // 水平干扰线
                    float lineNoise = 0.0;
                    if (random(vec2(0.0, floor(uv.y * 100.0))) > 0.99) {
                        lineNoise = random(uv) * 0.03;
                    }
                    
                    // 信号抖动
                    vec2 jitter = vec2(
                        sin(u_time * 50.0 + uv.y * 30.0) * 0.0005,
                        0.0
                    );
                    
                    color += staticNoise + lineNoise;
                    
                    return color;
                }
                
                // 磷光余辉（模拟P22磷光粉特性）
                vec3 phosphorPersistence(vec2 uv, vec3 currentColor) {
                    // 获取上一帧
                    vec3 previousFrame = texture2D(u_frameBuffer, uv).rgb;
                    
                    // P22磷光粉的衰减特性
                    vec3 decay = vec3(0.95, 0.93, 0.90); // RGB不同的衰减率
                    previousFrame *= decay * u_phosphorPersistence;
                    
                    // 混合当前帧和上一帧
                    return max(currentColor, previousFrame);
                }
                
                void main() {
                    // 应用CRT曲率
                    vec2 uv = curveRemapUV(v_texCoord);
                    
                    // 边界检查
                    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0) {
                        gl_FragColor = vec4(0.0, 0.0, 0.0, 1.0);
                        return;
                    }
                    
                    // 获取基础颜色（带会聚误差）
                    vec3 color = convergenceError(uv, u_texture);
                    
                    // 如果没有输入纹理，生成测试图案
                    if (color == vec3(0.0)) {
                        // 生成测试网格
                        float grid = step(0.99, sin(uv.x * 50.0) * sin(uv.y * 50.0));
                        color = vec3(grid) * vec3(0.2, 0.8, 0.2);
                        
                        // 添加颜色条
                        if (uv.y > 0.7) {
                            float bar = floor(uv.x * 8.0) / 8.0;
                            color = vec3(
                                step(0.125, bar) * step(bar, 0.375),
                                step(0.375, bar) * step(bar, 0.625),
                                step(0.625, bar) * step(bar, 0.875)
                            );
                        }
                    }
                    
                    // 保存原始颜色
                    vec3 originalColor = color;
                    
                    // 应用Trinitron磷光栅格（带强度控制）
                    color = mix(originalColor, trinitronPhosphor(uv, color), u_crtIntensity);
                    
                    // 应用扫描线（带强度控制）
                    float scanlineEffect = scanline(uv);
                    color *= mix(1.0, scanlineEffect, u_crtIntensity);
                    
                    // Trinitron支撑线（带强度控制）
                    float wireEffect = trinitronWires(uv);
                    color *= mix(1.0, wireEffect, u_crtIntensity * 0.5);
                    
                    // 电子束扫描（带强度控制）
                    float beamEffect = electronBeam(uv);
                    color *= mix(1.0, beamEffect, u_crtIntensity * 0.3);
                    
                    // 光晕效果（带强度控制）
                    vec3 haloColor = halation(uv, color);
                    color = mix(color, haloColor, u_crtIntensity);
                    
                    // 暗角（带强度控制）
                    float vignetteEffect = vignette(uv);
                    color *= mix(1.0, vignetteEffect, u_crtIntensity);
                    
                    // 信号噪声（带强度控制）
                    vec3 noisyColor = signalNoise(uv, color);
                    color = mix(color, noisyColor, u_crtIntensity * 0.5);
                    
                    // 色彩调整
                    color = colorGrading(color);
                    
                    // 磷光余辉（如果有帧缓冲）
                    // color = phosphorPersistence(uv, color);
                    
                    // 最终裁剪
                    color = clamp(color, 0.0, 1.0);
                    
                    gl_FragColor = vec4(color, 1.0);
                }
            `;

            this.init();
        }

        init() {
            this.printQuickStart();
        }

        padRight(str, length) {
            return str + ' '.repeat(Math.max(0, length - str.length));
        }


        printQuickStart() {
            console.log(`%c
⚡ 快速开始:
━━━━━━━━━━━━━━━━━━━━━━━━━━━
CRT.on()          - 开启CRT效果
CRT.off()         - 关闭CRT效果
CRT.preset()      - 查看预设
CRT.quickSet()    - 快速设置强度
CRT.adjust()      - 调整参数
CRT.status()      - 查看状态
━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
                'color: #ffff00; font-family: monospace');
        }

        createCanvas() {
            // 移除旧canvas
            if (this.canvas) {
                this.canvas.remove();
            }

            this.canvas = document.createElement('canvas');
            this.canvas.id = 'crt-webgl-fw900';
            this.canvas.style.cssText = `
                position: fixed !important;
                top: 0 !important;
                left: 0 !important;
                width: 100vw !important;
                height: 100vh !important;
                pointer-events: none !important;
                z-index: 2147483647 !important;
                image-rendering: high-quality;
                image-rendering: -webkit-optimize-contrast;
            `;

            const scale = this.config.resolution;
            this.canvas.width = window.innerWidth * scale;
            this.canvas.height = window.innerHeight * scale;

            // 获取WebGL上下文
            this.gl = this.canvas.getContext('webgl', {
                alpha: true,
                depth: false,
                stencil: false,
                antialias: false,
                premultipliedAlpha: false,
                preserveDrawingBuffer: false,
                powerPreference: 'high-performance',
                failIfMajorPerformanceCaveat: false
            });

            if (!this.gl) {
                console.error('❌ WebGL不可用，无法启动CRT效果');
                return false;
            }

            document.body.appendChild(this.canvas);

            // 监听窗口大小变化
            this.resizeHandler = () => {
                const scale = this.config.resolution;
                this.canvas.width = window.innerWidth * scale;
                this.canvas.height = window.innerHeight * scale;
            };
            window.addEventListener('resize', this.resizeHandler);

            this.setupWebGL();
            return true;
        }

        setupWebGL() {
            const gl = this.gl;

            // 编译着色器
            const vs = this.compileShader(gl.VERTEX_SHADER, this.vertexShader);
            const fs = this.compileShader(gl.FRAGMENT_SHADER, this.fragmentShader);

            if (!vs || !fs) {
                console.error('❌ 着色器编译失败');
                return;
            }

            // 创建程序
            this.program = gl.createProgram();
            gl.attachShader(this.program, vs);
            gl.attachShader(this.program, fs);
            gl.linkProgram(this.program);

            if (!gl.getProgramParameter(this.program, gl.LINK_STATUS)) {
                console.error('❌ 着色器程序链接失败:', gl.getProgramInfoLog(this.program));
                return;
            }

            // 设置顶点数据（全屏四边形）
            const positions = new Float32Array([
                -1, -1,  // 左下
                1, -1,  // 右下
                -1,  1,  // 左上
                1,  1,  // 右上
            ]);

            const texCoords = new Float32Array([
                0, 1,    // 左下
                1, 1,    // 右下
                0, 0,    // 左上
                1, 0,    // 右上
            ]);

            // 创建并绑定缓冲区
            const positionBuffer = gl.createBuffer();
            gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer);
            gl.bufferData(gl.ARRAY_BUFFER, positions, gl.STATIC_DRAW);

            const texCoordBuffer = gl.createBuffer();
            gl.bindBuffer(gl.ARRAY_BUFFER, texCoordBuffer);
            gl.bufferData(gl.ARRAY_BUFFER, texCoords, gl.STATIC_DRAW);

            // 使用着色器程序
            gl.useProgram(this.program);

            // 设置顶点属性
            const positionLocation = gl.getAttribLocation(this.program, 'a_position');
            gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer);
            gl.enableVertexAttribArray(positionLocation);
            gl.vertexAttribPointer(positionLocation, 2, gl.FLOAT, false, 0, 0);

            const texCoordLocation = gl.getAttribLocation(this.program, 'a_texCoord');
            gl.bindBuffer(gl.ARRAY_BUFFER, texCoordBuffer);
            gl.enableVertexAttribArray(texCoordLocation);
            gl.vertexAttribPointer(texCoordLocation, 2, gl.FLOAT, false, 0, 0);

            // 创建空纹理（用于测试）
            this.createTestTexture();

            // 开始渲染循环
            this.render();
        }

        createTestTexture() {
            const gl = this.gl;
            const texture = gl.createTexture();
            gl.bindTexture(gl.TEXTURE_2D, texture);

            // 创建一个1x1的透明纹理
            gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE,
                new Uint8Array([0, 0, 0, 0]));

            gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
            gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
            gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
            gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
        }

        compileShader(type, source) {
            const gl = this.gl;
            const shader = gl.createShader(type);
            gl.shaderSource(shader, source);
            gl.compileShader(shader);

            if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
                const info = gl.getShaderInfoLog(shader);
                console.error('❌ 着色器编译错误:', info);
                gl.deleteShader(shader);
                return null;
            }

            return shader;
        }

        render() {
            if (!this.enabled || !this.gl || !this.program) return;

            const gl = this.gl;
            const time = (Date.now() - this.startTime) * 0.001;

            // 设置视口
            gl.viewport(0, 0, this.canvas.width, this.canvas.height);

            // 清空画布
            gl.clearColor(0, 0, 0, 0);
            gl.clear(gl.COLOR_BUFFER_BIT);

            // 更新uniform变量
            this.updateUniforms(time);

            // 绘制全屏四边形
            gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4);

            // 更新帧计数
            this.frameCount++;

            // 继续渲染循环
            requestAnimationFrame(() => this.render());
        }

        updateUniforms(time) {
            const gl = this.gl;
            const program = this.program;

            // 时间和分辨率
            gl.uniform1f(gl.getUniformLocation(program, 'u_time'), time);
            gl.uniform2f(gl.getUniformLocation(program, 'u_resolution'),
                this.canvas.width, this.canvas.height);

            // CRT参数
            gl.uniform1f(gl.getUniformLocation(program, 'u_crtIntensity'),
                this.config.crtIntensity);
            gl.uniform1f(gl.getUniformLocation(program, 'u_scanlineIntensity'),
                this.config.scanlineIntensity);
            gl.uniform1f(gl.getUniformLocation(program, 'u_scanlineCount'),
                this.config.scanlineCount);
            gl.uniform1f(gl.getUniformLocation(program, 'u_phosphorPersistence'),
                this.config.phosphorPersistence);
            gl.uniform1f(gl.getUniformLocation(program, 'u_brightness'),
                this.config.brightness);
            gl.uniform1f(gl.getUniformLocation(program, 'u_contrast'),
                this.config.contrast);
            gl.uniform1f(gl.getUniformLocation(program, 'u_saturation'),
                this.config.saturation);
            gl.uniform1f(gl.getUniformLocation(program, 'u_gamma'),
                this.config.gamma);
            gl.uniform1f(gl.getUniformLocation(program, 'u_curvature'),
                this.config.curvature);
            gl.uniform1f(gl.getUniformLocation(program, 'u_vignette'),
                this.config.vignette);
            gl.uniform1f(gl.getUniformLocation(program, 'u_noiseAmount'),
                this.config.noiseAmount);
            gl.uniform1f(gl.getUniformLocation(program, 'u_convergence'),
                this.config.convergence);
            gl.uniform1f(gl.getUniformLocation(program, 'u_halation'),
                this.config.halation);
        }

        on() {
            if (this.enabled) {
                console.warn('⚠️ CRT效果已经开启');
                return this;
            }

            if (this.createCanvas()) {
                this.enabled = true;
                this.startTime = Date.now();
                // 添加页面滤镜
                this.applyPageFilters();
                // 开机动画
                this.powerOnAnimation();
            }

            return this;
        }

        off() {
            if (!this.enabled) {
                console.warn('⚠️ CRT效果已经关闭');
                return this;
            }

            this.enabled = false;

            // 关机动画
            this.powerOffAnimation(() => {
                // 移除canvas
                if (this.canvas) {
                    this.canvas.remove();
                    this.canvas = null;
                }

                // 清理WebGL
                if (this.gl) {
                    const loseContext = this.gl.getExtension('WEBGL_lose_context');
                    if (loseContext) {
                        loseContext.loseContext();
                    }
                    this.gl = null;
                }

                // 移除事件监听
                if (this.resizeHandler) {
                    window.removeEventListener('resize', this.resizeHandler);
                    this.resizeHandler = null;
                }

                // 移除页面滤镜
                this.removePageFilters();
            });

            return this;
        }

        powerOnAnimation() {
            if (!this.canvas) return;

            let progress = 0;
            const animate = () => {
                progress += 0.05;
                if (progress >= 1) {
                    this.canvas.style.transform = 'scaleY(1) scaleX(1)';
                    this.canvas.style.filter = 'brightness(1) blur(0px)';
                    return;
                }

                // 模拟CRT开机效果
                const scaleY = 0.001 + progress * 0.999;
                const scaleX = 0.5 + progress * 0.5;
                const brightness = 1 + (1 - progress) * 10;
                const blur = (1 - progress) * 5;

                this.canvas.style.transform = `scaleY(${scaleY}) scaleX(${scaleX})`;
                this.canvas.style.filter = `brightness(${brightness}) blur(${blur}px)`;
                this.canvas.style.transition = 'none';

                requestAnimationFrame(animate);
            };

            this.canvas.style.transformOrigin = 'center center';
            requestAnimationFrame(animate);
        }

        powerOffAnimation(callback) {
            if (!this.canvas) {
                callback();
                return;
            }

            let progress = 1;
            const animate = () => {
                progress -= 0.08;
                if (progress <= 0) {
                    callback();
                    return;
                }

                // 模拟CRT关机效果
                const scaleY = progress * progress * progress;
                const scaleX = 1 + (1 - progress) * 0.2;
                const brightness = 1 + (1 - progress) * 20;

                this.canvas.style.transform = `scaleY(${scaleY}) scaleX(${scaleX})`;
                this.canvas.style.filter = `brightness(${brightness})`;

                requestAnimationFrame(animate);
            };

            requestAnimationFrame(animate);
        }

        applyPageFilters() {
            const style = document.createElement('style');
            style.id = 'crt-fw900-page-filters';
            style.textContent = `
                /* FW900 页面滤镜效果 */
                body.crt-fw900-active {
                    /* 轻微的色彩调整 */
                    filter: 
                        contrast(${this.config.contrast * 0.95})
                        brightness(${this.config.brightness * 0.98})
                        saturate(${this.config.saturation});
                    
                    /* 防止滚动条被覆盖 */
                    overflow: auto !important;
                    
                    /* 启用 macOS 风格的滚动条 */
                    scrollbar-width: thin;
                    scrollbar-color: rgba(128, 128, 128, 0.3) transparent;
                }
                
                /* 文字渲染优化 */
                body.crt-fw900-active * {
                    text-rendering: geometricPrecision;
                    -webkit-font-smoothing: none;
                    -moz-osx-font-smoothing: grayscale;
                }
                
                /* macOS 风格滚动条 - WebKit (Chrome, Safari, Edge) */
                body.crt-fw900-active::-webkit-scrollbar {
                    width: 8px;  /* 更细的宽度 */
                    height: 8px; /* 水平滚动条高度 */
                }
                
                body.crt-fw900-active::-webkit-scrollbar-track {
                    background: transparent; /* 透明轨道 */
                    border-radius: 10px;
                }
                
                body.crt-fw900-active::-webkit-scrollbar-thumb {
                    background: rgba(128, 128, 128, 0.3); /* 半透明灰色 */
                    border-radius: 10px;
                    border: 2px solid transparent; /* 边距 */
                    background-clip: padding-box;
                    transition: background-color 0.2s ease;
                }
                
                body.crt-fw900-active::-webkit-scrollbar-thumb:hover {
                    background: rgba(128, 128, 128, 0.5); /* 悬停时稍微深一点 */
                    background-clip: padding-box;
                }
                
                body.crt-fw900-active::-webkit-scrollbar-thumb:active {
                    background: rgba(128, 128, 128, 0.6); /* 点击时更深 */
                    background-clip: padding-box;
                }
                
                /* 角落（当两个滚动条都显示时） */
                body.crt-fw900-active::-webkit-scrollbar-corner {
                    background: transparent;
                }
                
                /* 隐藏滚动条按钮 (macOS 风格不显示按钮) */
                body.crt-fw900-active::-webkit-scrollbar-button {
                    display: none;
                }
                
                /* Firefox 滚动条样式 */
                @supports (scrollbar-width: thin) {
                    body.crt-fw900-active {
                        scrollbar-width: thin;
                        scrollbar-color: rgba(128, 128, 128, 0.3) transparent;
                    }
                }
                
                /* 自动隐藏效果（可选）- 仅在滚动时显示 */
                body.crt-fw900-active::-webkit-scrollbar-thumb {
                    opacity: 0;
                    transition: opacity 0.3s ease;
                }
                
                body.crt-fw900-active:hover::-webkit-scrollbar-thumb {
                    opacity: 1;
                }
            `;

            document.head.appendChild(style);
            document.body.classList.add('crt-fw900-active');
        }

        removePageFilters() {
            const style = document.getElementById('crt-fw900-page-filters');
            if (style) style.remove();
            document.body.classList.remove('crt-fw900-active');
        }

        toggle() {
            return this.enabled ? this.off() : this.on();
        }

        // 设置CRT效果强度
        intensity(value) {
            if (value === undefined) {
                console.log(`%c当前CRT效果强度: ${this.config.crtIntensity}`, 'color: #00ff00');
                console.log('%c调整范围: 0-2 (0=无效果, 1=标准, 2=强烈)', 'color: #ffff00');
                return this;
            }

            this.config.crtIntensity = Math.max(0, Math.min(2, value));
            console.log(`%c✅ CRT效果强度已设置为: ${this.config.crtIntensity}`, 'color: #00ff00; font-weight: bold');

            if (this.config.crtIntensity === 0) {
                console.log('%c💡 提示: 强度为0时CRT效果将完全关闭', 'color: #ffff00');
            } else if (this.config.crtIntensity > 1.5) {
                console.log('%c⚠️ 警告: 强度过高可能影响视觉体验', 'color: #ff8800');
            }

            return this;
        }

        // 快捷预设方法
        quickSet(mode) {
            const modes = {
                'subtle': {
                    name: '轻微效果',
                    crtIntensity: 0.3,
                    scanlineIntensity: 0.05,
                    vignette: 0.1,
                    noiseAmount: 0.005
                },
                'moderate': {
                    name: '中等效果',
                    crtIntensity: 0.7,
                    scanlineIntensity: 0.1,
                    vignette: 0.2,
                    noiseAmount: 0.008
                },
                'standard': {
                    name: '标准效果',
                    crtIntensity: 1.0,
                    scanlineIntensity: 0.15,
                    vignette: 0.25,
                    noiseAmount: 0.01
                },
                'intense': {
                    name: '强烈效果',
                    crtIntensity: 1.5,
                    scanlineIntensity: 0.25,
                    vignette: 0.35,
                    noiseAmount: 0.02
                },
                'extreme': {
                    name: '极限效果',
                    crtIntensity: 2.0,
                    scanlineIntensity: 0.35,
                    vignette: 0.45,
                    noiseAmount: 0.03
                }
            };

            if (!mode) {
                console.log('%c🎚️ CRT效果强度快速设置:', 'color: #00ff00; font-weight: bold');
                for (const [key, value] of Object.entries(modes)) {
                    console.log(`  %c${key}%c - ${value.name} (强度: ${value.crtIntensity})`,
                        'color: #ffff00; font-weight: bold',
                        'color: #ffffff');
                }
                console.log('\n使用: CRT.quickSet("moderate")');
                return this;
            }

            const preset = modes[mode];
            if (!preset) {
                console.error(`❌ 未找到模式: ${mode}`);
                this.quickSet();
                return this;
            }

            // 应用设置
            Object.assign(this.config, preset);
            console.log(`✅ 已应用CRT效果: %c${preset.name}`, 'color: #00ff00; font-weight: bold');
            console.log(`   强度: ${preset.crtIntensity}, 扫描线: ${preset.scanlineIntensity}`);

            return this;
        }

        // 预设配置
        preset(name) {
            const presets = {
                'fw900-85hz': {
                    name: 'FW900 标准 (85Hz)',
                    refreshRate: 85,
                    scanlineIntensity: 0.15,
                    scanlineCount: 1080,
                    brightness: 1.02,
                    contrast: 1.18,
                    sharpness: 0.92
                },
                'fw900-96hz': {
                    name: 'FW900 高刷 (96Hz)',
                    refreshRate: 96,
                    scanlineIntensity: 0.12,
                    scanlineCount: 1200,
                    brightness: 1.01,
                    contrast: 1.20,
                    sharpness: 0.94
                },
                'fw900-gaming': {
                    name: 'FW900 游戏模式',
                    refreshRate: 120,
                    scanlineIntensity: 0.10,
                    phosphorPersistence: 0.08,
                    brightness: 1.05,
                    contrast: 1.25,
                    saturation: 1.1
                },
                'fw900-movie': {
                    name: 'FW900 影院模式',
                    refreshRate: 72,
                    scanlineIntensity: 0.18,
                    vignette: 0.35,
                    brightness: 1.0,
                    contrast: 1.22,
                    saturation: 0.92,
                    gamma: 2.4
                },
                'fw900-text': {
                    name: 'FW900 文本模式',
                    refreshRate: 85,
                    scanlineIntensity: 0.08,
                    scanlineCount: 1440,
                    brightness: 1.0,
                    contrast: 1.15,
                    sharpness: 0.98,
                    noiseAmount: 0.005
                }
            };

            if (!name) {
                console.log('%c📺 可用的FW900预设:', 'color: #00ff00; font-weight: bold');
                for (const [key, value] of Object.entries(presets)) {
                    console.log(`  %c${key}%c - ${value.name}`,
                        'color: #ffff00; font-weight: bold',
                        'color: #ffffff');
                }
                console.log('\n使用: CRT.preset("fw900-85hz")');
                return this;
            }

            const preset = presets[name];
            if (!preset) {
                console.error(`❌ 未找到预设: ${name}`);
                this.preset();
                return this;
            }

            // 应用预设
            Object.assign(this.config, preset);

            console.log(`✅ 已应用预设: %c${preset.name}`, 'color: #00ff00; font-weight: bold');
            return this;
        }

        // 参数调整
        adjust(param, value) {
            if (param === undefined) {
                console.log('%c可调整的参数:', 'color: #00ff00; font-weight: bold');
                console.table({
                    crtIntensity: { 当前值: this.config.crtIntensity, 范围: '0-2', 说明: 'CRT总体强度' },
                    scanlineIntensity: { 当前值: this.config.scanlineIntensity, 范围: '0-1' },
                    scanlineCount: { 当前值: this.config.scanlineCount, 范围: '240-2160' },
                    brightness: { 当前值: this.config.brightness, 范围: '0.5-2' },
                    contrast: { 当前值: this.config.contrast, 范围: '0.5-2' },
                    saturation: { 当前值: this.config.saturation, 范围: '0-2' },
                    sharpness: { 当前值: this.config.sharpness, 范围: '0-1' },
                    curvature: { 当前值: this.config.curvature, 范围: '0-0.1' },
                    vignette: { 当前值: this.config.vignette, 范围: '0-1' },
                    noiseAmount: { 当前值: this.config.noiseAmount, 范围: '0-0.1' }
                });
                return this;
            }

            if (this.config.hasOwnProperty(param)) {
                this.config[param] = value;
                console.log(`✅ ${param} = ${value}`);
            } else {
                console.error(`❌ 未知参数: ${param}`);
            }

            return this;
        }

        // 状态查看
        status() {
            console.log(`%c
╔════════════════════════════════════════════════════╗
║            Sony FW900 CRT 状态                     ║
╠════════════════════════════════════════════════════╣
║  状态: ${this.enabled ? '✅ 运行中' : '❌ 已停止'}                                 ║
║  版本: v${this.version}                                    ║
║  分辨率: ${this.canvas ? this.canvas.width + 'x' + this.canvas.height : 'N/A'}                         ║
║  刷新率: ${this.config.refreshRate}Hz                              ║
║  扫描线: ${this.config.scanlineCount} lines                       ║
║  CRT强度: ${this.config.crtIntensity}                                ║
║  帧数: ${this.frameCount}                                  ║
╚════════════════════════════════════════════════════╝`,
                'color: #00ff00; font-family: monospace');
            return this;
        }

        help() {
            console.log(`%c
╔════════════════════════════════════════════════════╗
║      Sony FW900 CRT WebGL - 帮助                   ║
╠════════════════════════════════════════════════════╣
║                                                    ║
║  基础命令:                                          ║
║    CRT.on()              - 开启CRT效果             ║
║    CRT.off()             - 关闭CRT效果             ║
║    CRT.toggle()          - 切换开关                ║
║                                                    ║
║  效果强度:                                          ║
║    CRT.intensity()       - 查看当前强度            ║
║    CRT.intensity(1.5)    - 设置强度(0-2)           ║
║    CRT.quickSet()        - 查看快速预设            ║
║    CRT.quickSet('mode')  - 应用快速预设            ║
║                                                    ║
║  预设配置:                                          ║
║    CRT.preset()          - 查看所有预设            ║
║    CRT.preset('name')    - 应用预设                ║
║                                                    ║
║  参数调整:                                          ║
║    CRT.adjust()          - 查看可调参数            ║
║    CRT.adjust(p, v)      - 调整参数                ║
║                                                    ║
║  信息查看:                                          ║
║    CRT.status()          - 查看运行状态            ║
║    CRT.help()            - 显示此帮助              ║
║                                                    ║
╚════════════════════════════════════════════════════╝`,
                'color: #00ff00; font-family: monospace; line-height: 1.2');
            return this;
        }
    }

    // 创建全局实例
    window.CRT = new CRTWebGLFW900();
    window.CRTWebGL = window.CRT; // 兼容别名

    // 启动提示
    setTimeout(() => {
        console.log(
            '%c💡 输入 %cCRT.on()%c 开启效果, %cCRT.quickSet("moderate")%c 调整强度',
            'color: #888888',
            'color: #00ff00; background: #000; padding: 2px 4px; border-radius: 3px; font-weight: bold',
            'color: #888888',
            'color: #00ff00; background: #000; padding: 2px 4px; border-radius: 3px; font-weight: bold',
            'color: #888888'
        );
    }, 500);

})(window);


setTimeout(function (){
    // CRT.on()
    // CRT.off()
    // 调整参数
    // CRT.adjust('scanlineIntensity', 0.05)
    // CRT.adjust('scanlineCount', 1080)
    // CRT.adjust('crtIntensity', 0.2)
    // CRT.adjust('brightness', 0.9)

},600)