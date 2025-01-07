package api

import (
	"net/http"
	"primary/api/middleware/authorization"
	"primary/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Message string `form:"message" json:"message" xml:"message" binding:"required"`
}

func newMessageHandler(c *gin.Context) {
	claims, err := utils.GetClaims(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "jwt error"})
		return
	}
	userID := claims.Userid

	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "serverid is not a number"})
		return
	}

	var msg Message
	if err := c.Bind(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get message payload"})
		return
	}

	msgID, err := db.NewMessage(msg.Message, userID, uint(serverID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "message sent", "data": msgID})
}

func deleteMessageHandler(c *gin.Context) {
	messageID, err := strconv.ParseUint(c.Param("messageid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "messageid is not a number"})
		return
	}

	if err := db.DeleteMessage(uint(messageID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "message deleted"})
}

func getByIdHandler(c *gin.Context) {
	messageID, err := strconv.ParseUint(c.Param("messageid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "messageid is not a number"})
		return
	}

	message, err := db.GetMessageByID(uint(messageID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "message not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "message returned", "data": gin.H{"messageid": message}})
}

func getByServerHandler(c *gin.Context) {
	serverID, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "serverid is not a number"})
		return
	}

	messages, err := db.GetMessagesByServerID(uint(serverID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "message not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "messages returned", "data": messages})
}

func getByUserHandler(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "userid is not a number"})
		return
	}

	messages, err := db.GetMessagesByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "message not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "messages returned", "data": messages})
}

func AddMessageRoutes(grp *gin.RouterGroup) { //TODO: add authorization
	grp.POST("/:serverid", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckSameUser)
	}, newMessageHandler)
	grp.DELETE("/:messageid", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckServerAdmin, authorization.CheckSameUser)
	}, deleteMessageHandler)
	grp.GET("/byserver/:serverid", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckServerMember)
	}, getByServerHandler)
	grp.GET("/byuser/:userid", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckSameUser)
	}, getByUserHandler)
}
