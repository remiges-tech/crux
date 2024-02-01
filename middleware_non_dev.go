//go:build !dev

package main

import "github.com/gin-gonic/gin"

// WARNING!!
// This is dangerous. Do not modify it.
// Notice the build tag at the top of the file. It is: //go:build !dev
// Do not modify that as well.
// This middleware function is a no-op in production and it has to be a no-op.
// Incorrect CORS headers can lead to security vulnerabilities.
func corsMiddleware() gin.HandlerFunc {
	return nil
}
