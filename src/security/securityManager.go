package security

import (
	"clinic/sqlc"
	"log"

	"clinic/src/entity"

	"github.com/gin-gonic/gin"
)

var RolePermissions = map[string]map[string]map[string]bool{
	"doctor": {
		"doctor": {
			"create":false,
			"read":  true,	
			"update":true,
			"delete":false,
		},
		"receptionist": {
			"create":false,
			"read":  false,
			"update":false,
			"delete":false,
		},
		"patient": {
			"create":false,
			"read":  false,
			"update":false,
			"delete":false,
		},
		"appointment": {
			"create":false,
			"read":  true,
			"update":true,
			"delete":true,
		},
		"application": {
			"create":false,
			"read":  false,
			"update":false,
			"delete":false,
		}

		// contniue and make for patient and receptionist
		//receptionist can read application and create appointment
		//patient can read and create and update application and read appointment
		//doctor can read and create and update appointment and cant do shit with application
	}
}

type SecurityManager struct {
	RolePermissions map[string]map[string]map[string]bool
}

func NewSecurityManager(db *sqlc.Queries) *SecurityManager {
	return &SecurityManager{
		RolePermissions: RolePermissions,
	}
}

func (sm *SecurityManager) CheckPermission(sve *entity.SessionValueEntity, target string, action string) bool {

	role := sve.Role

	// Get permissions for this target
	targetPermissions, exists := role[target]
	if !exists {
		log.Printf("Target %s not found for role: %s", target, role)
		return false
	}

	// Get permission for this action
	permission, exists := targetPermissions[action]
	if !exists {
		log.Printf("Action %s not found for target %s in role: %s", action, target, role)
		return false
	}

	return permission
}

// PermissionMiddleware creates a middleware for checking permissions in Gin HTTP handlers
func (sm *SecurityManager) GeneralPermissionMiddleware(target string, action string, relatedData interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session from context
		session, exists := c.Get("session")
		if !exists {
			log.Println("Session not found in context")
			c.JSON(403, gin.H{"error": "Session not found"})
			c.Abort()
			return
		}

		// Type assertion to get the session value entity
		sve, ok := session.(*entity.SessionValueEntity)
		if !ok {
			log.Println("Invalid session data in context")
			c.JSON(403, gin.H{"error": "Invalid session data"})
			c.Abort()
			return
		}

		// Check permission
		if !sm.CheckPermission(sve, target, action) {
			log.Printf("Permission denied for user with RoleID: %d on target: %s action: %s", sve.Role, target, action)
			c.JSON(403, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
