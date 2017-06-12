# kick-kick-go

## Quick start on Docker

### Run with example config and templates on http://localhost:8000/.
```
docker run --rm -p 8000:8000 amane/kick-kick-go
```

### Run with your own config or/and templates.
Firstly, copy [kick-kick-go/example](https://github.com/amane-katagiri/kick-kick-go/tree/master/example) to `/path/to/your/app`.

Replace all of config, template, static files with yours.
```
docker run --rm -p 8000:8000 -v /path/to/your/app:/app:ro amane/kick-kick-go
```

Or mount `config`, `templates` or `static` directory individually.
```
docker run --rm -p 8000:8000 -v /path/to/your/app/config:/app/config:ro \
    -v /path/to/your/app/templates:/app/templates:ro \
    -v /path/to/your/app/static:/app/static:ro amane/kick-kick-go
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

If you want to use different websocket url, set `ws_url` on config.json or use `-wsurl`.

Example: You can use websocket url `ws://localhost.amane.moe:9999/dev/chair` as templates variable when using the following `config.json`. `Origin` header will be expected `http://localhost.amane.moe:9999` if `server.check_origin` is `true`.

```
{
    "server": {
        "key": "",
        "cert": "",
        "host": "localhost",
        "port": 8000,
        "check_origin": false,
        "ws_path": "ws://localhost.amane.moe:9999/dev/chair"
    },
    "static_dir": "static",
    "template_files": ["templates/index.tmpl"]
}
```

