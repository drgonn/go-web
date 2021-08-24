package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-web/global"
	"github.com/go-web/internal/service"
	"github.com/go-web/pkg/app"
	"github.com/go-web/pkg/errcode"
)

type Tag struct{}

func NewTag() Tag {
	return Tag{}
}

func (t Tag) Get(c *gin.Context) {}
func (t Tag) List(c *gin.Context) {
	// param := struct {
	// 	Name string `form:"name" binding:"max=100"`
	// 	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
	// }{}
	// param 是验证结构体的实例
	param := service.TagListRequest{}
	response := app.NewResponse(c)
	// 验证入参，并绑定参数到param当中
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	// 获取service，包含context和gorm
	svc := service.New(c.Request.Context())
	// 获取页数和页面条数放到Pager中
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
	// 数据库查到数据条数
	totalRows, err := svc.CountTag(&service.CountTagRequest{Name: param.Name, State: param.State})

	response.ToResponse(gin.H{})
	return
}

func (t Tag) Create(c *gin.Context) {}
func (t Tag) Update(c *gin.Context) {}
func (t Tag) Delete(c *gin.Context) {}
