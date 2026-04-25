<?php

use function Abderrahim\Nats\connect;
use function Abderrahim\Nats\publish;
use function Abderrahim\Nats\subscribe;
use function Abderrahim\Nats\nextMessage;
use function Abderrahim\Nats\unsubscribe;
use function Abderrahim\Nats\flush as nats_flush;
use function Abderrahim\Nats\isConnected;
use function Abderrahim\Nats\stats;
use function Abderrahim\Nats\close as nats_close;
use const Abderrahim\Nats\SECOND;
use const Abderrahim\Nats\MILLISECOND;

header('Content-Type: text/plain; charset=utf-8');

$conn = 'frankenphp-test';

connect($conn, ['nats://127.0.0.1:4222']);

echo 'connected: ' . (isConnected($conn) ? 'yes' : 'no') . "\n";

$sub = subscribe($conn, 'frankenphp.test');

publish($conn, 'frankenphp.test', 'hello from frankenphp', ['x-trace' => 'abc123']);

nats_flush($conn, 2 * SECOND);

$msg = nextMessage($sub, 2 * SECOND);
if ($msg === null) {
    echo "no message received\n";
} else {
    echo "subject: {$msg['subject']}\n";
    echo "data: {$msg['data']}\n";
    echo 'header x-trace: ' . ($msg['headers']['X-Trace'][0] ?? '(none)') . "\n";
}

unsubscribe($sub);

$s = stats($conn);
echo "in_msgs >= 1: " . ($s['in_msgs'] >= 1 ? 'yes' : 'no') . "\n";
echo "out_msgs >= 1: " . ($s['out_msgs'] >= 1 ? 'yes' : 'no') . "\n";

nats_close($conn);
echo "closed\n";
