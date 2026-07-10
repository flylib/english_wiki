# english_wiki — 启动 / 构建
# 用法:make(看帮助) / make dev(前后端一起) / make backend / make frontend
SHELL := /bin/bash

SERVER_DIR := server
WEB_DIR    := web
CONFIG     := config/settings.sqlite.yml
GO_TAGS    := sqlite3
BIN        := english-wiki
BACKEND_PORT  := 739
FRONTEND_PORT := 1798

.DEFAULT_GOAL := help
.PHONY: help dev backend run-server frontend build-server deps stop

help: ## 显示可用命令
	@echo "english_wiki 命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	  | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

# ---------- 一起启动 ----------
dev: ## 同时启动后端(:739)+ 前端(:1798),Ctrl+C 一起退出
	@echo "▶ 启动 后端:$(BACKEND_PORT) + 前端:$(FRONTEND_PORT)  (Ctrl+C 退出)"
	@trap 'kill 0' EXIT INT TERM; \
	  $(MAKE) --no-print-directory backend & \
	  $(MAKE) --no-print-directory frontend & \
	  wait

# ---------- 分别启动 ----------
backend: build-server ## 仅启动后端(先构建再运行,:739)
	@cd $(SERVER_DIR) && mkdir -p temp/logs && ./$(BIN) server -c $(CONFIG)

run-server: ## 仅运行后端(不重新构建)
	@cd $(SERVER_DIR) && mkdir -p temp/logs && ./$(BIN) server -c $(CONFIG)

frontend: ## 仅启动前端(:1798,代理 /api/v1 -> :739)
	@cd $(WEB_DIR) && corepack pnpm dev

# ---------- 构建 / 其他 ----------
build-server: ## 构建后端二进制(必带 sqlite3 tag + CGO)
	@cd $(SERVER_DIR) && CGO_ENABLED=1 go build -tags $(GO_TAGS) -o $(BIN) .

deps: ## 安装依赖(go mod download + pnpm install)
	@cd $(SERVER_DIR) && go mod download
	@cd $(WEB_DIR) && corepack pnpm install

stop: ## 停掉占用 739/1798 端口的前后端进程
	@for p in $(BACKEND_PORT) $(FRONTEND_PORT); do \
	  pid=$$(lsof -ti tcp:$$p -sTCP:LISTEN 2>/dev/null); \
	  if [ -n "$$pid" ]; then echo "✖ kill :$$p (pid $$pid)"; kill $$pid; else echo "· :$$p 没有进程"; fi; \
	done
