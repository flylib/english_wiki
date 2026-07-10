package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"

	"go-admin/app/jobs"
)

// StartEmbedded 以非阻塞方式启动 API 服务,供桌面端(Wails)在进程内嵌入调用。
// 复用 cobra `server` 命令的 setup()/initRouter()/AppRouters,但不做信号阻塞。
// 调用方需保证工作目录下存在 configPath 指向的配置(及其引用的 sqlite/日志路径)。
func StartEmbedded(configPath string) {
	configYml = configPath
	setup()

	if config.ApplicationConfig.Mode == pkg.ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}
	initRouter()
	for _, f := range AppRouters {
		f()
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.ApplicationConfig.Host, config.ApplicationConfig.Port),
		Handler:      sdk.Runtime.GetEngine().(*gin.Engine),
		ReadTimeout:  time.Duration(config.ApplicationConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.ApplicationConfig.WriterTimeout) * time.Second,
	}

	go func() {
		jobs.InitJob()
		jobs.Setup(sdk.Runtime.GetDb())
	}()
	go func() {
		log.Infof("embedded api server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("embedded api server error: %s", err.Error())
		}
	}()
}
