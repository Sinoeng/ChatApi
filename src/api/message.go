package api

import (
	"net/http"
	"primary/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Message string `form:"message" json:"message" xml:"message" binding:"required"`
}

func NewMessageHandler(c *gin.Context) {
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

func DeleteMessageHandler(c *gin.Context) {
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

func GetByIdHandler(c *gin.Context) {
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

func GetByServerHandler(c *gin.Context) {
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

func GetByUserHandler(c *gin.Context) {
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
	grp.POST("/:serverid", NewMessageHandler)
	grp.DELETE("/:messageid", DeleteMessageHandler)
	grp.GET("/byid/:messageid", GetByIdHandler)
	grp.GET("/byserver/:serverid", GetByServerHandler)
	grp.GET("/byuser/:userid", GetByUserHandler)
}
