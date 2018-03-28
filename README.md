# triton-service-groups

## Run

```sh
$ make build
$ bin/triton-sg agent --log-level=DEBUG
```

While developing eveything besides scaling actions you can rely on `TSG_DEV_MODE=1` to skip authentication.

```sh
$ TSG_DEV_MODE=1 bin/triton-sg agent --log-level=DEBUG
```

When dev mode is enabled any request sent to the TSG API (regardless of headerS) will be linked to the seed data we've provided within `./dev/setup_db.sh`. This data is only provided as a stub and will not work against Triton's CloudAPI.

### Whitelist

Authentication provides a whitelisting feature which only allows incoming requests to be authenticated if the account has been entered into the TSG database. If whitelisting is not enabled than all Triton accounts that can be authenticated with CloudAPI will generate a new account and key within the TSG API.

The following SQL is a snippet for adding your Triton account to TSG.

```sql
INSERT INTO tsg_accounts (account_name, triton_uuid, created_at, updated_at) VALUES ('demouser', 'd82a1f04-b9f6-4075-998f-af20e3d49de6', NOW(), NOW());
```

To turn off whitelisting requires changing a boolean within `server/handlers/auth/consts.go`. Set `isWhitelistOnly` to either `true` or `false`. `true` is the default. A new build of TSG must include this change as it is not yet configurable within the config file or env var.

## Environment

The following environment variables override any configuration file values.

```sh
TSG_LOG_LEVEL=DEBUG
TSG_CRDB_DATABASE=triton
TSG_CRDB_HOST=127.0.0.1
TSG_CRDB_PORT=26257
TSG_CRDB_USER=root
TSG_CRDB_PASSWORD=database123
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
