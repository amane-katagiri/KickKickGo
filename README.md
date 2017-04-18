# kick-kick-go

## Quick start on Docker

* Run with example config and templates on http://localhost:8000/.
```
docker run --rm -p 8000:8000 amane/kick-kick-go
```

* Run with your own config and templates.
```
docker run --rm -p 8000:8000 -v /path/to/your/app:/app:ro amane/kick-kick-go
```

## Save count permanently

Use Redis binding to save count.

```
docker run --name kick-kick-go-redis -d redis
docker run -p 8000:8000 --link kick-kick-go-redis:redis -d amane/kick-kick-go -redis.address redis:6379
```

## Options

Use `-h` to see all options.

```
docker run --rm amane/kick-kick-go -h
```

### Run with proxy

If you want to use different scheme(use tls or not), host, port or path, set `ws_url` on config.json or use `-wsurl.*`. explicit `ws_url` parameters overwrites Origin.

Example: You can use websocket url `ws://localhost.amane.moe:9999/dev/chair` as templates variable when using the following `config.json`.

* The Origin will be `http://localhost.amane.moe:9999`.

* You must pass `/dev/chair` to `/ws` in proxy.

```
{
    "server": {
        "key": "",
        "cert": "",
        "host": "localhost",
        "port": 8000,
        "ws_path": "/ws"
    },
    "static_dir": "static",
    "template_files": ["templates/index.tmpl"],
    "ws_url": {
        "ssl": false,
        "host": "localhost.amane.moe",
        "port": 9999,
        "path": "/dev/chair"
    }
}
```
