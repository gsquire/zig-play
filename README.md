# Zig Playground

This is a rudimentary online compiler for the [Zig](https://ziglang.org) programming language. It
is inspired by the [Go](https://play.golang.org) playground.

### Setup
The main server is a Go binary that serves up a single HTML page that allows you to enter your Zig
code and then run it. To run it yourself, you will need a Go tool chain which can be installed via
`brew` on a Mac. If you wish to run it locally, you must compile it for your `GOOS` and `GOARCH`
but I have included a small shell script to make a Linux binary that Docker can use as well. You
should also have Zig installed and accessible from within your `$PATH` on the host.

### Hosting
In theory this could be run anywhere that a Docker container can execute. Google's Cloud Run may be
a cheap option considering their generous free tier.

### FAQ
> What can this playground do?

It is currently set up to simply run a single Zig source file. (i.e. `zig run source.zig`)

> Why does the page look so bad?

I'm clearly more of a server-side programmer. Any styling changes are very much appreciated!

> Why am I getting rate-limited?

You're allowed five compilations per minute which I think is fairly generous.

> Is it secure?

Go read the source. I do not collect logs of any kind and am not interested in your data. Unless it
is causing issues to the service.

> Will this always be available?

To the best of my ability, I will try and keep this online. I may even buy a domain for it.

### Contact
Feel free to email me at my listed address with any other questions or comments.

### License
MIT
