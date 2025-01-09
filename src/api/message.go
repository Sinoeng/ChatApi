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

// new_message godoc
// @Summary send message to a server
// @Description send a message to server with serverid if you are a member
// @Schemes
// @Param Authorization header string true "jwt token"
// @Param message body string true "message"
// @Tags auth message
// @Accept json
// @Produce json
// @Success 200
// @Router /protected/message/:serverid [post]
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

// delete_message godoc
// @Summary delete messaage by id
// @Description delete message with id messageid if you are admin
// @Schemes
// @Param Authorization header string true "jwt token"
// @Tags auth message
// @Produce json
// @Success 200
// @Router /protected/message/:messageid [delete]
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

// get_message_by_server godoc
// @Summary get message by server
// @Description get message by serverid if you are sever member or admin
// @Schemes
// @Param Authorization header string true "jwt token"
// @Tags auth message
// @Produce json
// @Success 200
// @Router /protected/message/byserver/:serverid [get]
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

// get_message_by_user godoc
// @Summary get message by user
// @Description get message by userid if you are same user or admin
// @Schemes
// @Param Authorization header string true "jwt token"
// @Tags auth message
// @Produce json
// @Success 200
// @Router /protected/message/byuser/:userid [get]
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
		authorization.AuthorizeMiddleware(c, db, authorization.CheckServerMember)
	}, newMessageHandler)
	grp.DELETE("/:messageid", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin)
	}, deleteMessageHandler)
	grp.GET("/byserver/:serverid", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckServerMember)
	}, getByServerHandler)
	grp.GET("/byuser/:userid", func(c *gin.Context) {
		authorization.AuthorizeMiddleware(c, db, authorization.CheckGlobalAdmin, authorization.CheckSameUser)
	}, getByUserHandler)
}
