# openai proxy

此项目由 chatgpt 生成

- [x] 新建 api_key 并配置 token 限制
- [x] 代理 openai
- [ ] check 请求是否超 token


## docker 部署

修改 `docker/docker-compose.yaml` 文件

```bash
docker compose -f docker/docker-compose.yaml up
```

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
    volumes:
      - ./config.yaml:/config.yaml
    environment:
      - API_KEY=xxx
      - CONF_PATH=/
      # - HTTP_PROXY=xxx
      - REDIS_HOST=localhost:6379
      # - REDIS_PASS=xxx

  redis:
    image: redis
    container_name: redis
    network_mode: service:proxy
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

volumes:
  redis_data:
```
