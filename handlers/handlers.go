package handlers

import (
	_ "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"retroblog/config"
	"retroblog/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UnlockRequest represents password unlock request
type UnlockRequest struct {
	Path     string `json:"path"`
	Password string `json:"password"`
}

// HandleUnlock handles password verification and content unlock
func HandleUnlock(c *gin.Context) {
	var req UnlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req.Path = strings.Trim(req.Path, "/")

	// Wait for any ongoing rebuild
	config.RebuildingMu.Lock()
	isRebuilding := config.Rebuilding
	config.RebuildingMu.Unlock()

	if isRebuilding {
		time.Sleep(500 * time.Millisecond)
	}

	config.ProtectedMu.RLock()
	expectedHash, ok := config.ProtectedPosts[req.Path]
	content, hasContent := config.FullContent[req.Path]
	config.ProtectedMu.RUnlock()

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post not protected"})
		return
	}

	if utils.HashPassword(req.Password) != expectedHash {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	if !hasContent {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Content not found"})
		return
	}

	// Generate auth token and set cookie
	authToken := utils.GenerateAuthToken(req.Path, req.Password)
	c.SetCookie(
		"auth_"+utils.HashString(req.Path),
		authToken,
		86400, // 24 hours
		"/",
		"firsh.me",
		true, // Secure - set to true in production with HTTPS
		true, // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"content": content,
		"status":  "success",
	})
}

// SharedArrayBufferMiddleware 添加 VLC.js 所需的安全头
func SharedArrayBufferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 必需：启用 SharedArrayBuffer
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Embedder-Policy", "require-corp")

		// 处理 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// AutoVideoMiddleware 自动处理视频文件的万能中间件
// 特性：自动识别视频 + Range 支持 + CORS + 所有必需 Headers
// 使用：r.Use(middleware.AutoVideoMiddleware())

func AutoVideoMiddleware(basePath string) gin.HandlerFunc {
	// 视频 MIME 类型映射
	videoExts := map[string]string{
		".mp4":  "video/mp4",
		".webm": "video/webm",
		".mkv":  "video/x-matroska",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".flv":  "video/x-flv",
		".wmv":  "video/x-ms-wmv",
		".m4v":  "video/x-m4v",
		".mpg":  "video/mpeg",
		".mpeg": "video/mpeg",
		".ogv":  "video/ogg",
		".3gp":  "video/3gpp",
		".ts":   "video/mp2t",
		".m3u8": "application/vnd.apple.mpegurl",
	}

	return func(c *gin.Context) {
		// 1. 检查是否为视频文件
		ext := strings.ToLower(filepath.Ext(c.Request.URL.Path))
		contentType, isVideo := videoExts[ext]
		if !isVideo {
			c.Next()
			return
		}

		// 2. 构建文件路径并检查是否存在
		filePath := filepath.Join(basePath, c.Request.URL.Path)
		log.Printf("[UK0] Request: %s %s → %s", c.Request.Method, c.Request.URL.Path, filePath)

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Printf("[UK0] File not found: %v", err)
			c.Next() // 交给后续处理器（如 404）
			return
		}

		fileSize := fileInfo.Size()
		log.Printf("[UK0] File: %s (%.2f MB)", fileInfo.Name(), float64(fileSize)/1024/1024)

		// 3. 设置通用响应头
		c.Header("Content-Type", contentType)
		c.Header("Accept-Ranges", "bytes")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Range, Content-Type, Cache-Control")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges")
		c.Header("Cache-Control", "public, max-age=86400")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Accel-Buffering", "no") // 禁用 Nginx 缓冲（如果使用反向代理）
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(filePath)))

		// 4. 处理 OPTIONS 预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 5. 处理 HEAD 请求（仅返回头部）
		if c.Request.Method == http.MethodHead {
			c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
			c.Status(http.StatusOK)
			c.Abort()
			return
		}

		// 6. 打开文件
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("[UK0] Failed to open file: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// 7. 检查 Range 请求头
		rangeHeader := c.GetHeader("Range")

		// 8. 无 Range：返回完整文件 (200 OK)
		if rangeHeader == "" {
			c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
			c.Status(http.StatusOK)

			log.Printf("[UK0] Serving full file (%d bytes)", fileSize)
			written, err := io.Copy(c.Writer, file)
			if err != nil {
				log.Printf("[UK0] ❌ Transfer failed: %v (sent %d bytes)", err, written)
			} else {
				log.Printf("[UK0] ✅ Sent %d bytes", written)
			}

			c.Abort()
			return
		}

		// 9. 解析 Range 请求头
		start, end, err := parseRange(rangeHeader, fileSize)
		if err != nil {
			log.Printf("[UK0] Invalid Range: %s (%v)", rangeHeader, err)
			c.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
			c.AbortWithStatus(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		// 10. 验证范围
		if start < 0 || start >= fileSize || end >= fileSize || start > end {
			c.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
			c.AbortWithStatus(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		contentLength := end - start + 1

		// 11. 设置 206 Partial Content 响应
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
		c.Status(http.StatusPartialContent)

		// 12. 定位文件并传输
		if _, err := file.Seek(start, io.SeekStart); err != nil {
			log.Printf("[UK0] ❌ Seek failed: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		log.Printf("[UK0] Range: bytes=%d-%d/%d (%d bytes)", start, end, fileSize, contentLength)
		written, err := io.CopyN(c.Writer, file, contentLength)
		if err != nil {
			log.Printf("[UK0] ⚠️ Range transfer interrupted: %v (sent %d/%d bytes)", err, written, contentLength)
		} else {
			log.Printf("[UK0] ✅ Range sent: %d bytes", written)
		}

		c.Abort()
	}
}

// parseRange 解析 Range 头，返回 (start, end, error)
// 支持：
// - bytes=100-200  → start=100, end=200
// - bytes=100-     → start=100, end=fileSize-1
// - bytes=-500     → start=fileSize-500, end=fileSize-1
func parseRange(rangeHeader string, fileSize int64) (start, end int64, err error) {
	// 移除 "bytes=" 前缀
	rangeStr := strings.TrimPrefix(rangeHeader, "bytes=")

	// 暂不支持多范围（如 "0-100,200-300"）
	if strings.Contains(rangeStr, ",") {
		return 0, 0, fmt.Errorf("multi-range not supported")
	}

	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format")
	}

	// Case 1: bytes=-500 (最后 500 字节)
	if parts[0] == "" && parts[1] != "" {
		suffix, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil || suffix <= 0 {
			return 0, 0, fmt.Errorf("invalid suffix length")
		}
		if suffix > fileSize {
			suffix = fileSize
		}
		return fileSize - suffix, fileSize - 1, nil
	}

	// Case 2: bytes=100- (从 100 到结尾)
	if parts[0] != "" && parts[1] == "" {
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil || start < 0 {
			return 0, 0, fmt.Errorf("invalid start")
		}
		return start, fileSize - 1, nil
	}

	// Case 3: bytes=100-200 (完整范围)
	start, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil || start < 0 {
		return 0, 0, fmt.Errorf("invalid start")
	}

	end, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil || end < 0 {
		return 0, 0, fmt.Errorf("invalid end")
	}

	return start, end, nil
}

// AttachmentAuthMiddleware checks if attachments require authentication
func AttachmentAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		log.Printf("[Middleware] Protected path: %s ", path)

		// Check if this is an attachment request
		if strings.Contains(path, "/attachments/") {
			fullPath := strings.TrimPrefix(path, "/")

			config.ProtectedMu.RLock()
			isProtected := config.ProtectedAttachments[fullPath]
			config.ProtectedMu.RUnlock()
			log.Printf("[Middleware] Protected path: %s fullPath: %s isProtected: %v", path, fullPath, isProtected)
			if isProtected {
				// Extract post path from URL
				parts := strings.Split(fullPath, "/")
				if len(parts) >= 2 {
					postPath := parts[0] + "/" + parts[1]

					cookie, err := c.Cookie("auth_" + utils.HashString(postPath))
					if err != nil || !utils.VerifyAuthCookie(postPath, cookie) {
						c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
							"error": "Unauthorized - Please unlock the post first",
						})
						return
					}
				}
			}
		}

		c.Next()
	}
}
