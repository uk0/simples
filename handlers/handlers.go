package handlers

import (
	_ "encoding/json"
	"log"
	"net/http"
	"retroblog/config"
	"retroblog/utils"
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
		time.Sleep(100 * time.Millisecond)
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
