package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey    = "user_id"
	UserEmailKey = "user_email"
	UserRoleKey  = "user_role"
)

// AuthRequired reads X-User-Id, X-User-Email and X-User-Role injected by the
// upstream API Gateway + Lambda Authorizer. Returns 401 if any header is absent.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-Id")
		userEmail := c.GetHeader("X-User-Email")
		userRole := c.GetHeader("X-User-Role")

		if userID == "" || userEmail == "" || userRole == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing required user headers"})
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Set(UserEmailKey, userEmail)
		c.Set(UserRoleKey, userRole)

		c.Next()
	}
}

// RoleRequired returns a middleware that grants access if the request carries
// one of the required roles. Returns 401 if claims are missing, 403 if role is insufficient.
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(UserRoleKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing identity claims"})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing identity claims"})
			c.Abort()
			return
		}

		for _, r := range roles {
			if roleStr == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient role"})
		c.Abort()
	}
}

// GetUserID reads the user ID from the Gin context (set by AuthRequired).
func GetUserID(c *gin.Context) string {
	id, exists := c.Get(UserIDKey)
	if !exists {
		return ""
	}
	idStr, _ := id.(string)
	return idStr
}

// GetUserEmail reads the user email from the Gin context.
func GetUserEmail(c *gin.Context) string {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return ""
	}
	emailStr, _ := email.(string)
	return emailStr
}

// GetUserRole reads the user role from the Gin context.
func GetUserRole(c *gin.Context) string {
	role, exists := c.Get(UserRoleKey)
	if !exists {
		return ""
	}
	roleStr, _ := role.(string)
	return roleStr
}
