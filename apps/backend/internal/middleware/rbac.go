package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gms/backend/internal/models"
)

// RequireRoles is the baseline RBAC gate for module endpoints.
// This can be replaced later by a policy/permission engine.
func RequireRoles(roles ...models.Role) gin.HandlerFunc {
	allowed := make(map[models.Role]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		user, ok := CurrentUser(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
			return
		}

		if _, exists := allowed[user.Role]; !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
			return
		}

		c.Next()
	}
}
