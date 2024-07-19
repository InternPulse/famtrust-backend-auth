package interfaces

import "github.com/gin-gonic/gin"

type Handlers interface {
	Users() UserHandlers
}

type UserHandlers interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Validate(c *gin.Context)
}
