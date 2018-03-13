# triton-service-groups

The following environment variables are still required to demo.

```
NOMAD_URL
TRITON_ACCOUNT
TRITON_URL
TRITON_KEY_ID
TRITON_KEY_MATERIAL
```

## Run

```sh
$ make build
$ bin/triton-sg agent --log-level=DEBUG
```

## Environment

The following environment variables override any configuration file values.

```sh
TSG_LOG_LEVEL=DEBUG
TSG_POSTGRESQL_DATABASE=triton
TSG_POSTGRESQL_HOST=127.0.0.1
TSG_POSTGRESQL_PORT=26257
TSG_POSTGRESQL_USER=root
TSG_POSTGRESQL_PASSWORD=database123
TSG_HTTP_BIND=127.0.0.1
TSG_HTTP_PORT=3000
TSG_GOPS_ENABLE=true
TSG_GOPS_BIND=127.0.0.1
TSG_GOPS_PORT=9090
TSG_PPROF_ENABLE=true
TSG_PPROF_BIND=127.0.0.1
TSG_PPROF_PORT=9191
```

## Configuration

```toml
[log]
level = "INFO"

[postgresql]
database = "triton"
host = "127.0.0.1"
port = 26257
password = ""
user = "root"

[run]
log-format = "auto"

[http]
bind = "127.0.0.1"
port = 3000

[gops]
enable = true
bind = "127.0.0.1"
port = 9191

[pprof]
enable = true
bind = "127.0.0.1"
port = 9090
```
