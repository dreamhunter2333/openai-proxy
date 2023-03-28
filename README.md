# openai proxy

此项目由 chatgpt 生成

- [x] 新建 api_key 并配置 token 限制
- [x] 代理 openai
- [ ] check 请求是否超 token

## docker 部署

### 仅代理

```yaml
version: '3'

services:
  proxy:
    build:
      context: ..
      dockerfile: docker/Dockerfile
    container_name: proxy
    ports:
      - "8080:8080"
```

### 代理加计费

修改 `docker/docker-compose.yaml` 文件

```bash
docker compose -f docker/docker-compose.yaml up
```
