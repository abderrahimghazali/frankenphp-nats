# FrankenPHP Extension for NATS

A high-performance [NATS](https://nats.io) client for [PHP](https://php.net), designed to work with [FrankenPHP](https://frankenphp.dev).
It leverages [the official NATS Go client](https://github.com/nats-io/nats.go).

This extension enables the creation of global, shared NATS connections that persist across requests and worker script instances,
giving PHP applications first-class access to NATS messaging without the per-request reconnect cost of pure-PHP clients.

> [!NOTE]
>
> `frankenphp-nats` follows the pattern established by [`dunglas/frankenphp-etcd`](https://github.com/dunglas/frankenphp-etcd)
> and is built using FrankenPHP's [Extension Generator](https://github.com/php/frankenphp/blob/main/docs/extensions.md).
> The v0.1.0 scope is **core pub/sub** (publish, subscribe, request/reply, headers, basic auth, TLS).
>
> JetStream, Key-Value, Object Store, push subscriptions and the NATS Micro framework are planned for follow-up releases —
> see [Roadmap](#roadmap).

## Status

| Feature                                        | Version  |
| ---------------------------------------------- | -------- |
| Core pub/sub, request/reply, headers           | v0.1.0   |
| JetStream (streams, consumers, JS publish)     | v0.2.0   |
| Key-Value & Object Store                       | v0.3.0   |
| Push subscriptions, Caddyfile-declared conns   | v0.4.0   |
| NATS Micro / Services framework                | v0.5.0   |

## Installation

First, if not already done, follow [the instructions to install a ZTS version of libphp and `xcaddy`](https://frankenphp.dev/docs/compile/#install-php).
Then, use [`xcaddy`](https://github.com/caddyserver/xcaddy) to build FrankenPHP with the `frankenphp-nats` module:

```console
CGO_ENABLED=1 \
CGO_CFLAGS=$(php-config --includes) \
CGO_LDFLAGS="$(php-config --ldflags) $(php-config --libs)" \
xcaddy build \
    --output frankenphp \
    --with github.com/abderrahimghazali/frankenphp-nats/build=./build \
    --with github.com/dunglas/frankenphp/caddy
    # Add extra Caddy modules and FrankenPHP extensions here
```

That's all — your custom FrankenPHP build now exposes the `Abderrahim\Nats` namespace to PHP.

## Usage

```php
<?php

use function Abderrahim\Nats\connect;
use function Abderrahim\Nats\publish;
use function Abderrahim\Nats\subscribe;
use function Abderrahim\Nats\nextMessage;
use function Abderrahim\Nats\unsubscribe;
use function Abderrahim\Nats\request;
use const Abderrahim\Nats\SECOND;

// Open (or reuse) a globally registered NATS connection. Persists across
// requests and worker reboots.
connect('default', ['nats://127.0.0.1:4222']);

// Fire-and-forget publish, with optional headers.
publish('default', 'orders.created', json_encode(['id' => 42]), [
    'X-Trace-Id' => 'abc123',
]);

// Subscribe and pull one message synchronously.
$sub = subscribe('default', 'orders.>');
$msg = nextMessage($sub, 5 * SECOND);
if ($msg !== null) {
    echo "got {$msg['subject']}: {$msg['data']}\n";
}
unsubscribe($sub);

// Synchronous request/reply.
$reply = request('default', 'svc.echo', 'ping', 1 * SECOND);
echo $reply['data'] ?? 'timeout';
```

### API surface (v0.1.0)

All symbols live under `Abderrahim\Nats`:

| Function | Purpose |
| --- | --- |
| `connect(name, servers, …auth…)` | Create or reuse a globally registered connection. |
| `publish(name, subject, data, ?headers)` | Fire-and-forget publish. |
| `request(name, subject, data, timeout)` | Synchronous request/reply, returns `array\|null`. |
| `subscribe(name, subject, ?queue)` | Returns a subscription ID (string). |
| `nextMessage(subId, timeout)` | Pulls one message synchronously. |
| `unsubscribe(subId)` | Removes a subscription. |
| `flush(name, timeout)` | Block until pending publishes are flushed. |
| `isConnected(name)` | True if currently connected. |
| `stats(name)` | Returns counters: `in_msgs`, `out_msgs`, `in_bytes`, `out_bytes`, `reconnects`. |
| `close(name)` | Close and remove from the global registry. |

Time-unit constants are exposed as `NANOSECOND`, `MICROSECOND`, `MILLISECOND`, `SECOND`, `MINUTE`.

Message arrays returned by `request()` and `nextMessage()` have the shape:

```php
[
    'subject' => string,
    'data'    => string,
    'reply'   => ?string,
    'headers' => array<string, string[]>,
]
```

### Error handling

Because the [Extension Generator](https://github.com/php/frankenphp/blob/main/docs/extensions.md)
does not yet expose a Go-callable API for raising PHP exceptions, failures (connection errors,
publish errors, unknown connection names) are logged via the FrankenPHP error log and surface
to PHP as zero values:

| Function | On failure |
| --- | --- |
| `connect()`, `publish()`, `flush()`, `unsubscribe()`, `close()` | Logs and returns. Subsequent calls (`isConnected()`, `stats()`) reveal state. |
| `request()`, `nextMessage()` | Returns `null` (also returned on legitimate timeout). |
| `subscribe()` | Returns an empty string. |
| `stats()` | Returns an empty array. |
| `isConnected()` | Returns `false`. |

A future release will introduce explicit exception types once the upstream generator gains
exception-throwing primitives.

### Authentication

`connect()` accepts the full set of NATS auth options:

| Argument | Description |
| --- | --- |
| `username` + `password` | Basic auth. |
| `token` | Token auth. |
| `credsFile` | Path to a NATS `.creds` file (NGS / decentralised auth). |
| `nkeyFile` | Path to an NKey seed file. |
| `tls` | Enable TLS with sane defaults (TLS 1.2+). |

## Development

This extension is built using FrankenPHP's
[Extension Generator](https://github.com/php/frankenphp/blob/main/docs/extensions.md).
After editing `nats.go`, regenerate the C/PHP boilerplate:

```console
GEN_STUB_SCRIPT=path/to/php-src/build/gen_stub.php \
  frankenphp extension-init nats.go
```

This refreshes `build/nats.c`, `build/nats.h`, `build/nats_arginfo.h`, `build/nats.stub.php`, and `build/nats_generated.go`.
The generated `build/` directory is committed so consumers of this extension don't need a PHP build environment.

Run the Go tests against a real `nats-server`:

```console
docker run -d --name nats-test -p 4222:4222 nats:latest
go test -race -tags nobadger,nomysql,nopgx,nowatcher,nobrotli -v ./...
docker rm -f nats-test
```

## Credits

Created by [Abderrahim Ghazali](https://github.com/abderrahimghazali), inspired by
[`dunglas/frankenphp-etcd`](https://github.com/dunglas/frankenphp-etcd) and
[`dunglas/frankenphp-grpc`](https://github.com/dunglas/frankenphp-grpc).
