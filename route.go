package main

import "ProjectName/framework"

//注册路由规则
func registerRouter(core *framework.Core) {
	//需求1+2：HTTP方法+静态路由匹配

	core.Get("foo", FooControllerHandler)
}
