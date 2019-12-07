package boot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/xinjiayu/SimServerUnicom/app/api/analyse"
	"github.com/xinjiayu/SimServerUnicom/app/api/collect"
	"github.com/xinjiayu/SimServerUnicom/app/api/notice"
	"github.com/xinjiayu/SimServerUnicom/app/api/operate"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
// 统一路由注册
func routerInit() {
	s := g.Server()
	s.SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
	s.BindMiddlewareDefault(MiddlewareCORS)

	s.BindObject("/unicom", new(collect.Controller))
	s.BindObject("/unicom/op",new(operate.Controller))
	s.BindObject("/unicom/analyse",new(analyse.Controller))
	s.BindObject("/unicom/notice",new(notice.Controller))

}