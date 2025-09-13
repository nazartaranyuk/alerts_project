#To build this project create config.yaml file in configs/ directory with following name config.(dev, stage or prod).yaml with following syntax
```yaml
env: "dev" // could be dev, stage, prod

server:
  host: "name of your host"
  port: // your port (int type)

client:
  api_base_url: "YOUR_API_BASE_URL"
  api_key: "YOUR_API_KEY"
  timeout_sec: // timeout (int type)
```
