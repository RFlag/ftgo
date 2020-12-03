package ftgo

import (
	"context"
	stdlog "log"
	"net/http"
	"os"
	"strings"

	"ftgo/ftconf"
	"ftgo/safeclose"

	"github.com/gin-gonic/gin"
)

func Run(addr string, router func(*gin.Engine)) {
	log := stdlog.New(os.Stdout, "[ftgo] ", stdlog.LstdFlags)

	gin.SetMode(gin.ReleaseMode)

	g := gin.New()

	g.Any("/health", ok)

	if ftconf.Mode == ftconf.Debug {
		g.Use(gin.Logger())
	}
	g.Use(gin.Recovery())
	g.Use(ErrorLog(stdlog.New(os.Stdout, "[log] ", stdlog.LstdFlags)))

	router(g)

	server := &http.Server{
		Addr:    addr,
		Handler: g,
	}

	safeclose.DoContext(func(ctx context.Context) {
		<-ctx.Done()
		server.Shutdown(context.Background())
	})

	log.Print("启动 " + server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		safeclose.Cancel()
		if err == http.ErrServerClosed {
			log.Print("http 服务器已关闭")
		} else {
			log.Print("服务器异常退出", err)
		}
	}

	safeclose.Wait()
}
func ErrorLog(log *stdlog.Logger) gin.HandlerFunc {
	if ftconf.Mode == ftconf.Debug {
		return func(c *gin.Context) {
			c.Next()
			errors := c.Errors.ByType(^gin.ErrorTypeBind).Errors()
			for _, v := range errors {
				log.Print(v)
			}
			bindError := c.Errors.ByType(gin.ErrorTypeBind).Last()
			if bindError != nil {
				c.JSON(400, struct {
					Code       int      `json:"code"`
					Error      string   `json:"error"`
					ParamError []string `json:"paramError"`
				}{
					Code:       ResultParamError.Code,
					Error:      ResultParamError.Error,
					ParamError: strings.Split(bindError.Error(), "\n"),
				})
				return
			}
			publicError := c.Errors.ByType(gin.ErrorTypePublic).Last()
			if publicError != nil {
				rce := publicError.Meta.(resultCodeError)
				c.JSON(-1, struct {
					Code         int    `json:"code"`
					Error        string `json:"error"`
					PrivateError string `json:"privateError"`
				}{
					Code:         rce.Code,
					Error:        rce.Error,
					PrivateError: publicError.Error(),
				})
				return
			}
		}
	} else {
		return func(c *gin.Context) {
			c.Next()
			errors := c.Errors.ByType(^gin.ErrorTypeBind).Errors()
			for _, v := range errors {
				log.Print(v)
			}
			bindErrors := c.Errors.ByType(gin.ErrorTypeBind)
			if len(bindErrors) > 0 {
				c.JSON(400, ResultParamError)
				return
			}
			publicError := c.Errors.ByType(gin.ErrorTypePublic).Last()
			if publicError != nil {
				c.JSON(-1, publicError.Meta.(resultCodeError))
				return
			}
		}
	}
}

func ok(c *gin.Context) {
	c.Status(200)
}
