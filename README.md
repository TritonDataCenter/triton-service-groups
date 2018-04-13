# triton-service-groups

Triton Service Groups (TSGs) provide a simple way to manage any number of instances running the same image with the same configuration. Simply define the instance template, then create the service group and set the number of instances that are needed. As needs change, you can easily scale the number of instances.

Triton Service Groups help you maximize application efficiency so that you only pay for resources that you actually need.

At this time, there is no integration with Triton CLI.

## API Usage

The API has 2 main endpoints:

* [groups](docs/groups/index.md)
* [templates](docs/templates/index.md)

All API calls to the API require an Authorization header. An example Authorization header may look as follows:

```
Authorization: Signature keyId=/demo/keys/foo,algorithm="rsa-sha256" ${Base64(sign($Date))}
```

The default value to sign for API requests is simply the value of the HTTP Date header. For more information on the Date header value, see [RFC 2616](http://tools.ietf.org/html/rfc2616#section-14.18). All requests to the API using the Signature authentication scheme must send a Date header.

### Using CURL with Triton Service Groups

```bash
function tsg() {
    local now=$(date -u '+%a, %d %h %Y %H:%M:%S GMT')
    local signature=$(echo -n "$now" | openssl dgst -sha256 -sign ~/.ssh/id_rsa | openssl enc -e -a | tr -d '\n')
    local url="$TSG_URL$1"
    shift

    curl -s -k -i \
        -H 'Accept: application/json' \
        -H "accept-version: ~8" \
        -H "Date: $now" \
        -H "Authorization: Signature keyId=\"/$TRITON_ACCOUNT/keys/id_rsa\",algorithm=\"rsa-sha256\" $signature" \
        "$@" "$url"
    echo
}
```

You may need to alter the path to your SSH key in the above function. With this function, you could just do:

```bash
TSG_URL=https://tsg.us-sw-1.svc.joyent.zone tsg /v1/tsg/templates
```

## Development

### Running the API

```sh
$ make build
$ bin/triton-sg agent --log-level=DEBUG
```

While developing anything except scaling actions, you can rely on `TSG_DEV_MODE=1` to skip authentication.

```sh
$ TSG_DEV_MODE=1 bin/triton-sg agent --log-level=DEBUG
```

When dev mode is enabled any request sent to the TSG API (regardless of headers) will be linked to the seed data we've provided within `./dev/setup_db.sh`. This data is only provided as a stub and will not work against CloudAPI.

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
TSG_NOMAD_URL=127.0.0.1
TSG_NOMAD_PORT=4646
TSG_TRITON_DC=us-east-1
TSG_TRITON_URL=https://us-east-1.api.joyent.com
```

## Configuration

```toml
[log]
level = "INFO"

[crdb]
database = "triton"
host = "127.0.0.1"
port = 26257
password = ""
user = "root"

[agent]
log-format = "auto"

[http]
bind = "127.0.0.1"
port = 3000
dc = "us-east-1"

[gops]
enable = true
bind = "127.0.0.1"
port = 9191

[pprof]
enable = true
bind = "127.0.0.1"
port = 9090

[nomad]
url = "127.0.0.1"
port = 4646

[triton]
dc = "us-east-1"
url = "https://us-east-1.api.joyent.com"
```
