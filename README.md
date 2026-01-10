# Zig Playground

This is a rudimentary online compiler for the [Zig](https://ziglang.org) programming language. It
is inspired by the [Go](https://play.golang.org) playground.

It's currently served from this [page](https://zig-play.dev).

### Setup
The main server is a Go binary that serves up a single HTML page that allows you to enter your Zig
code and then run it.

Installation Requirements:
- [zvm](https://www.zvm.app/)
- [go](https://go.dev)

You'll also want a shell script called `zrun.sh` accessible from the home directory of the user you
run the playground as:

```bash
#!/bin/bash

set -eu

VERSION=$1
SOURCE=$2
TIMEOUT=$YOUR_TIMEOUT_IN_SECONDS

/usr/bin/timeout $TIMEOUT zvm run "$VERSION" run "$SOURCE"
```

### Hosting
I currently am using a VPS and have [Caddy](https://caddyserver.com) as a reverse proxy.

### FAQ
> What can this playground do?

It is currently set up to simply run and format a single Zig source file. (i.e. `zig run source.zig` & `zig fmt source.zig`)

> Can you share an example request?

```bash
# main.zig is the file containing the source you want to execute.
curl -v -H "X-Zig-Version: 0.15.2" https://zig-play.dev/server/run -d @main.zig
```

> Are there any timeouts?

If your code doesn't build within 30 seconds, the server will quit your request.

> Why am I getting rate-limited?

You're allowed five executions per minute which I think is fairly generous.

> Is it secure?

Go read the source. I do not collect logs of any kind and am not interested in your data. Unless it
is causing issues to the service.

> Will this always be available?

To the best of my ability, I will try and keep this online.

### Contact
Feel free to write to hello@zig-play.dev with any questions or comments.

### License
MIT
