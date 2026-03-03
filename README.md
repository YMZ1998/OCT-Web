# OCT-Web

## GitHub Actions 自动部署

仓库已配置 `.github/workflows/deploy.yml`，在 `main` 分支 push 后会自动执行：

1. 后端 `go test ./...`
2. 前端 `npm install && npm run build`
3. 通过 SSH 登录部署机并执行 `docker compose up -d --build`

### 需要配置的仓库 Secrets

- `DEPLOY_HOST`：部署服务器地址
- `DEPLOY_USER`：SSH 用户名
- `DEPLOY_SSH_KEY`：私钥（建议专用 deploy key）
- `DEPLOY_PORT`：SSH 端口（可选，默认 22）
- `DEPLOY_PATH`：服务器上的项目目录（目录内需已初始化 git 且可执行 `docker compose`）
