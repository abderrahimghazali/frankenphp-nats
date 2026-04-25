package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	_ "github.com/abderrahimghazali/frankenphp-nats"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/dunglas/frankenphp/caddy"
)

func main() {
	caddycmd.Main()
}
