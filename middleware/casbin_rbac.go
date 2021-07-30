package middleware

import (
	"casbin-gorm/util"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的URL
		obj := c.Request.URL.RequestURI()
		// 获取请求的方法
		act := c.Request.Method
		// 获取用户的角色
		sub := "1"

		adapter, err := gormAdapter.NewAdapter("mysql", "root:0000@tcp(127.0.0.1:3306)/gorm1", true)
		if err != nil {
			fmt.Println("0", err)
		}

		enforcer, err := casbin.NewEnforcer("./rbac.conf", adapter)
		if err != nil {
			fmt.Println("1", err)
		}

		enforcer.AddFunction("ParamsMatch", util.ParamsMatchFunc)
		enforcer.LoadPolicy()

		ok, err := enforcer.Enforce(sub, obj, act)
		fmt.Println("ok", ok)
		if err != nil {
			fmt.Println("2", err)
		}
		if ok {
			c.Next()
		} else  {
			c.JSON(200, gin.H{
				"err_code": 1,
				"err_message": "没有权限",
			})
			c.Abort()
			return
		}
	}
}
