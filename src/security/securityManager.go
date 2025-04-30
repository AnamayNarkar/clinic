package security

import (
	"log"

	"clinic/src/entity"

	"github.com/gin-gonic/gin"
)

var RolePermissions = map[string]map[string]map[string]bool{

	"admin": {
		"doctor": {
			"read":   true,
			"create": true,
			"update": true,
			"delete": true,
		},
		"patient": {
			"read":   true,
			"create": true,
			"update": true,
			"delete": true,
		},
		"receptionist": {
			"read":   true,
			"create": true,
			"update": true,
			"delete": true,
		},
		"application": {
			"read":   true,
			"create": true,
			"update": true,
			"delete": true,
		},
		"appointment": {
			"read":   true,
			"create": true,
			"update": true,
			"delete": true,
		},
	},

	"doctor": {
		"doctor": {
			"read":   true,
			"create": false,
			"update": true,
			"delete": false,
		},
		"patient": {
			"read":   false,
			"create": false,
			"update": false,
			"delete": false,
		},
		"receptionist": {
			"read":   false,
			"create": false,
			"update": false,
			"delete": false,
		},
		"application": {
			"read":   false,
			"create": false,
			"update": false,
			"delete": false,
		},
		"appointment": {
			"read":   true,
			"create": false,
			"update": true,
			"delete": true,
		},
	},
	"receptionist": {
		"doctor": {
			"read":   false,
			"create": false,
			"update": false,
			"delete": false,
		},
		"patient": {
			"read":   true,
			"create": true,
			"update": true,
			"delete": false,
		},
		"receptionist": {
			"read":   true,
			"create": false,
			"update": true,
			"delete": false,
		},
		"application": {
			"read":   true,
			"create": false,
			"update": true,
			"delete": false,
		},
		"appointment": {
			"read":   true,
			"create": true,
			"update": false,
			"delete": true,
		},
	},
	"patient": {
		"doctor": {
			"read":   false,
			"create": false,
			"update": false,
			"delete": false,
		},
		"patient": {
			"read":   true,
			"create": false,
			"update": true,
			"delete": false,
		},
		"receptionist": {
			"read":   false,
			"create": false,
			"update": false,
			"delete": false,
		},
		"application": {
			"read":   true,
			"create": true,
			"update": true,
			"delete": false,
		},
		"appointment": {
			"read":   true,
			"create": false,
			"update": false,
			"delete": false,
		},
	},
}

type SecurityManager struct {
	RolePermissions map[string]map[string]map[string]bool
}

func NewSecurityManager() *SecurityManager {
	return &SecurityManager{
		RolePermissions: RolePermissions,
	}
}

func (sm *SecurityManager) CheckPermission(sve *entity.SessionValueEntity, target string, action string) bool {

	role := sve.Role

	targetPermissions, exists := sm.RolePermissions[sve.Role][target]
	if !exists {
		log.Printf("Target %s not found for role: %s", target, role)
		return false
	}

	permission, exists := targetPermissions[action]
	if !exists {
		log.Printf("Action %s not found for target %s in role: %s", action, target, role)
		return false
	}

	return permission
}

func (sm *SecurityManager) GeneralPermissionMiddleware(target string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, exists := c.Get("session")
		if !exists {
			log.Println("Session not found in context")
			c.JSON(403, gin.H{"error": "Session not found"})
			c.Abort()
			return
		}

		sve, ok := session.(*entity.SessionValueEntity)
		if !ok {
			log.Println("Invalid session data in context")
			c.JSON(403, gin.H{"error": "Invalid session data"})
			c.Abort()
			return
		}

		if !sm.CheckPermission(sve, target, action) {
			log.Printf("Permission denied for user with Role: %s on target: %s action: %s", sve.Role, target, action)
			c.JSON(403, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
