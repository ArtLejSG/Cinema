package middleware

//func RoleCheck(requiredRole string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		role, exists := c.Get("role")
//		if !exists || role != requiredRole {
//			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
//			c.Abort()
//			return
//		}
//		c.Next()
//	}
//}
