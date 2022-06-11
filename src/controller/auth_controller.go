package controller

import (
	global "go-search-engine/src/global"
	"go-search-engine/src/model"
	"go-search-engine/src/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var cookiesName string = "camp-session"

func Login(c *gin.Context) {
	loginRequest := global.LoginRequest{}
	if err := c.ShouldBind(&loginRequest); err != nil {
		c.JSON(http.StatusOK, global.ResponseMeta{Code: global.UnknownError})
		return
	}
	log.Println(loginRequest)

	//user, err := model.GetMemberByUsernameAndPassword(loginRequest.Username, loginRequest.Password)
	user, err := model.GetMemberByUsername(loginRequest.Username)
	//用户不存在或者密码错误
	if err != nil || user.Password != utils.Md5Encrypt(loginRequest.Password) {

		c.JSON(http.StatusOK, global.ResponseMeta{Code: global.WrongPassword})
		return
	}
	//用户已删除
	if user.IsDeleted {
		c.JSON(http.StatusOK, global.ResponseMeta{Code: global.UserHasDeleted})
		return
	}

	session := sessions.Default(c)
	var sessionId = getSessionId()

	log.Println(sessionId, user)
	v := global.TMember{
		UserID:   strconv.Itoa(user.UserID),
		Nickname: user.Nickname,
		Username: user.Username,
		UserType: user.UserType,
	}
	session.Set(sessionId, v)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, global.ResponseMeta{Code: global.UnknownError})
		return
	}

	c.SetCookie(cookiesName, sessionId, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, global.LoginResponse{
		Code: global.OK,
		Data: struct {
			UserID string
		}{strconv.Itoa(user.UserID)},
	})
}

func Logout(c *gin.Context) {
	sessionId, err := c.Cookie(cookiesName)
	if err != nil {
		c.JSON(http.StatusOK, global.ResponseMeta{Code: global.LoginRequired})
		return
	}

	session := sessions.Default(c)

	if session.Get(sessionId) == nil {
		c.JSON(http.StatusOK, global.ResponseMeta{Code: global.LoginRequired})
		return
	}

	session.Delete(sessionId)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, global.ResponseMeta{Code: global.UnknownError})
		return
	}

	c.SetCookie(cookiesName, sessionId, -1, "/", "", false, true)
	c.JSON(http.StatusOK, global.LogoutResponse{Code: global.OK})
}

func WhoAmI(c *gin.Context) {
	sessionId, err := c.Cookie(cookiesName)
	if err != nil {
		c.JSON(http.StatusOK, global.ResponseMeta{Code: global.LoginRequired})
		return
	}

	session := sessions.Default(c)
	v := session.Get(sessionId)
	if v == nil {
		c.JSON(http.StatusOK, global.ResponseMeta{Code: global.LoginRequired})
		return
	}
	log.Println(v)
	user := v.(global.TMember)
	c.JSON(http.StatusOK, global.WhoAmIResponse{Code: global.OK, Data: user})
}

func getSessionId() string {
	//b := make([]byte, 32)
	//if _, err := io.ReadFull(rand.Reader, b); err != nil {
	//	return ""
	//}
	//return base64.StdEncoding.EncodeToString(b)
	return uuid.NewV4().String()
}
