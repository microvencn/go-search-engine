package types

import (
	"github.com/gin-gonic/gin"
	"go-search-engine/src/controller"
	"go-search-engine/src/controller/search"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	//TEST
	g.GET("/ping", controller.Ping)

	// 成员管理
	g.POST("/member/create", controller.CreateMember)
	g.GET("/member", controller.GetMember)
	g.GET("/member/list", controller.GetMemberList)
	g.POST("/member/update", controller.UpdateMember)
	g.POST("/member/delete", controller.DeleteMember)

	// 登录
	g.POST("/auth/login", controller.Login)
	g.POST("/auth/logout", controller.Logout)
	g.GET("/auth/whoami", controller.WhoAmI)

	g.GET("/search", search.SimpleSearch)
	g.GET("/auto", search.AutoComplete)
	// 数据管理
	//g.GET("/data", controller.GetData)
	//g.GET("/data/list", controller.GetDataList)

}
