/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - PaaS 平台 (BlueKing - PaaS System) available.
 * Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 *     http://opensource.org/licenses/MIT
 *
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package webserver

import (
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"

	dc "github.com/TencentBlueking/bkpaas/cnb-builder-shim/internal/devcontainer"
	"github.com/TencentBlueking/bkpaas/cnb-builder-shim/internal/devcontainer/webserver/service"
)

type envConfig struct {
	DevServerAddr string `env:"DEV_SERVER_ADDR" envDefault:":8000"`
	Token         string `env:"TOKEN" envDefault:"jwram1lpbnuugmcv"`
	UploadDir     string `env:"UPLOAD_DIR" envDefault:"/cnb/devcontainer/src"`
}

// WebServer 实现了 DevWatchServer 接口, 提供 Web 服务
type WebServer struct {
	server *gin.Engine
	lg     *logr.Logger
	ch     chan dc.AppReloadEvent
	env    envConfig
}

// New creates a new WebServer instance.
//
// It takes a pointer to a Logger object as a parameter and returns a pointer to a WebServer object.
func New(lg *logr.Logger) (*WebServer, error) {
	cfg := envConfig{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	r := gin.Default()
	r.Use(tokenAuthMiddleware(cfg.Token))

	s := &WebServer{
		server: r,
		lg:     lg,
		// 1 chan buffer, flow control
		ch:  make(chan dc.AppReloadEvent, 1),
		env: cfg,
	}

	mgr := service.NewDeployManager()
	r.POST("/deploys", DeployHandler(s, mgr))
	r.GET("/deploys/:deployID/results", ResultHandler(mgr))

	return s, nil
}

// Start starts the WebServer.
//
// It returns an error if the server fails to run.
func (s *WebServer) Start() error {
	return s.server.Run(s.env.DevServerAddr)
}

// AppReloadEvents 返回 reload 事件通道
func (s *WebServer) AppReloadEvents() <-chan dc.AppReloadEvent {
	return s.ch
}

func (s *WebServer) Clean() {
	close(s.ch)
}

func tokenAuthMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		reqToken := strings.TrimPrefix(authorizationHeader, "Bearer ")
		if reqToken == "" {
			c.String(http.StatusUnauthorized, "missing authorization token")
			c.Abort()
			return
		}
		if reqToken != token {
			c.String(http.StatusUnauthorized, "invalid authorization token")
			c.Abort()
			return
		}

		c.Next()
	}
}

// DeployHandler handles the deployment of a file to the web server.
func DeployHandler(s *WebServer, svc service.DeployServiceHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}

		fileName := filepath.Base(file.Filename)
		srcFilePath := path.Join(s.env.UploadDir, fileName)
		if err = c.SaveUploadedFile(file, srcFilePath); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}

		status, err := svc.Deploy(srcFilePath)
		if err != nil {
			c.String(http.StatusInternalServerError, "deploy error: %s", err.Error())
			return
		}

		if len(s.ch) == cap(s.ch) {
			c.String(http.StatusTooManyRequests, "app is deploying, please wait")
		} else {
			s.ch <- dc.AppReloadEvent{ID: status.DeployID, Rebuild: status.StepOpts.Rebuild, Relaunch: status.StepOpts.Relaunch}
			c.JSON(http.StatusOK, map[string]interface{}{"deployID": status.DeployID})
		}
	}
}

// ResultHandler is a function that get the result of a deployment.
func ResultHandler(svc service.DeployServiceHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		deployID := c.Param("deployID")
		withLog, _ := strconv.ParseBool(c.Query("log"))

		result, err := svc.Result(deployID, withLog)
		if err != nil {
			c.String(http.StatusBadRequest, "get result error: %s", err.Error())
			return
		}

		if withLog {
			c.JSON(http.StatusOK, map[string]interface{}{"status": result.Status, "log": result.Log})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{"status": result.Status})
		}
	}
}

var _ dc.DevWatchServer = (*WebServer)(nil)