// Package nats is a FrankenPHP extension exposing the NATS messaging
// system to PHP via the official Go client (github.com/nats-io/nats.go).
//
// All exported PHP symbols live under the "Abderrahim\Nats" namespace.
// Connections are stored in a process-wide registry keyed by name and
// persist across requests and worker reboots, mirroring the pattern
// established by frankenphp-etcd and frankenphp-grpc.
package nats

// #include <Zend/zend_types.h>
import "C"

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/dunglas/frankenphp"
	"github.com/google/uuid"
	natsgo "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

//export_php:namespace Abderrahim\Nats

// ----------------------------------------------------------------------------
// Time-unit constants (nanoseconds, matching nats.go's time.Duration values).
// ----------------------------------------------------------------------------

//export_php:const
const NANOSECOND = 1

//export_php:const
const MICROSECOND = 1000

//export_php:const
const MILLISECOND = 1000000

//export_php:const
const SECOND = 1000000000

//export_php:const
const MINUTE = 60000000000

// ----------------------------------------------------------------------------
// Process-wide registries.
// ----------------------------------------------------------------------------

var (
	connRegistry   = make(map[string]*natsgo.Conn)
	connRegistryMu sync.RWMutex

	// Serializes natsConnect to prevent two concurrent callers from both
	// dialing under the same name and silently leaking the loser.
	connectMu sync.Mutex

	subRegistry   = make(map[string]*subEntry)
	subRegistryMu sync.RWMutex
)

type subEntry struct {
	sub      *natsgo.Subscription
	connName string
}

var (
	errUnknownConnection   = errors.New("nats: unknown connection name")
	errUnknownSubscription = errors.New("nats: unknown subscription id")
)

// ----------------------------------------------------------------------------
// Connection lifecycle.
// ----------------------------------------------------------------------------

//export_php:function connect(string $name, array $servers, string $username = "", string $password = "", string $token = "", string $credsFile = "", string $nkeyFile = "", bool $tls = false, int $timeout = 2000000000, int $reconnectAttempts = 60, int $reconnectWait = 2000000000, int $pingInterval = 120000000000, int $maxPingsOut = 2): void
func natsConnect(name *C.zend_string, servers *C.zend_array, username *C.zend_string, password *C.zend_string, token *C.zend_string, credsFile *C.zend_string, nkeyFile *C.zend_string, useTLS bool, timeoutNs int64, reconnectAttempts int64, reconnectWaitNs int64, pingIntervalNs int64, maxPingsOut int64) {
	connName := frankenphp.GoString(unsafe.Pointer(name))

	// Serialize all connects to close the TOCTOU window between the
	// existence check and the registry insert. natsgo.Connect can take
	// seconds; do it inside the lock so concurrent callers see a single
	// writer per name. Connects are infrequent — this is not a hot path.
	connectMu.Lock()
	defer connectMu.Unlock()

	connRegistryMu.RLock()
	_, exists := connRegistry[connName]
	connRegistryMu.RUnlock()
	if exists {
		// Idempotent — calling connect() with an existing name is a no-op,
		// matching frankenphp-etcd's getOrCreate semantics.
		return
	}

	urls, err := serversFromPHPArray(servers)
	if err != nil {
		logError("connect", err.Error(), zap.String("name", connName))
		return
	}
	if len(urls) == 0 {
		logError("connect", "at least one server URL is required", zap.String("name", connName))
		return
	}

	user := frankenphp.GoString(unsafe.Pointer(username))
	pass := frankenphp.GoString(unsafe.Pointer(password))
	tok := frankenphp.GoString(unsafe.Pointer(token))
	creds := frankenphp.GoString(unsafe.Pointer(credsFile))
	nkey := frankenphp.GoString(unsafe.Pointer(nkeyFile))

	// Reject ambiguous auth: nats.go silently last-wins, which surprises
	// callers who pass two methods expecting validation. Allow only one of
	// {user/pass, token, credsFile, nkeyFile}. user without password (or
	// vice versa) is treated as no-auth and falls through.
	authMethods := 0
	if user != "" && pass != "" {
		authMethods++
	}
	if tok != "" {
		authMethods++
	}
	if creds != "" {
		authMethods++
	}
	if nkey != "" {
		authMethods++
	}
	if authMethods > 1 {
		logError("connect", "multiple auth methods provided; specify only one of (username+password, token, credsFile, nkeyFile)", zap.String("name", connName))
		return
	}

	if useTLS {
		urls = forceTLSScheme(urls)
	}

	opts := []natsgo.Option{
		natsgo.Name(connName),
		natsgo.Timeout(time.Duration(timeoutNs)),
		natsgo.MaxReconnects(int(reconnectAttempts)),
		natsgo.ReconnectWait(time.Duration(reconnectWaitNs)),
		natsgo.PingInterval(time.Duration(pingIntervalNs)),
		natsgo.MaxPingsOutstanding(int(maxPingsOut)),
	}

	if useTLS {
		opts = append(opts, natsgo.Secure(&tls.Config{MinVersion: tls.VersionTLS12}))
	}

	switch {
	case user != "" && pass != "":
		opts = append(opts, natsgo.UserInfo(user, pass))
	case tok != "":
		opts = append(opts, natsgo.Token(tok))
	case creds != "":
		opts = append(opts, natsgo.UserCredentials(creds))
	case nkey != "":
		nkOpt, err := natsgo.NkeyOptionFromSeed(nkey)
		if err != nil {
			logError("connect", "nkey: "+err.Error(), zap.String("name", connName))
			return
		}
		opts = append(opts, nkOpt)
	}

	urlList := joinURLs(urls)
	conn, err := natsgo.Connect(urlList, opts...)
	if err != nil {
		logError("connect", err.Error(), zap.String("name", connName))
		return
	}

	connRegistryMu.Lock()
	connRegistry[connName] = conn
	connRegistryMu.Unlock()
}

//export_php:function close(string $name): void
func natsClose(name *C.zend_string) {
	connName := frankenphp.GoString(unsafe.Pointer(name))

	connRegistryMu.Lock()
	conn, ok := connRegistry[connName]
	delete(connRegistry, connName)
	connRegistryMu.Unlock()
	if !ok {
		return
	}

	subRegistryMu.Lock()
	for id, entry := range subRegistry {
		if entry.connName == connName {
			_ = entry.sub.Unsubscribe()
			delete(subRegistry, id)
		}
	}
	subRegistryMu.Unlock()

	conn.Close()
}

//export_php:function isConnected(string $name): bool
func natsIsConnected(name *C.zend_string) bool {
	conn, ok := lookupConn(name)
	if !ok {
		return false
	}
	return conn.IsConnected()
}

//export_php:function flush(string $name, int $timeout = 5000000000): void
func natsFlush(name *C.zend_string, timeoutNs int64) {
	conn, ok := lookupConn(name)
	if !ok {
		logError("flush", errUnknownConnection.Error(), zap.String("name", frankenphp.GoString(unsafe.Pointer(name))))
		return
	}
	if err := conn.FlushTimeout(time.Duration(timeoutNs)); err != nil {
		logError("flush", err.Error())
	}
}

//export_php:function stats(string $name): array
func natsStats(name *C.zend_string) unsafe.Pointer {
	conn, ok := lookupConn(name)
	if !ok {
		logError("stats", errUnknownConnection.Error(), zap.String("name", frankenphp.GoString(unsafe.Pointer(name))))
		return frankenphp.PHPMap(map[string]any{})
	}
	s := conn.Stats()
	return frankenphp.PHPMap(map[string]any{
		"in_msgs":    int64(s.InMsgs),
		"out_msgs":   int64(s.OutMsgs),
		"in_bytes":   int64(s.InBytes),
		"out_bytes":  int64(s.OutBytes),
		"reconnects": int64(s.Reconnects),
	})
}

// ----------------------------------------------------------------------------
// Pub/sub.
// ----------------------------------------------------------------------------

//export_php:function publish(string $name, string $subject, string $data, array $headers = []): void
func natsPublish(name *C.zend_string, subject *C.zend_string, data *C.zend_string, headers *C.zend_array) {
	conn, ok := lookupConn(name)
	if !ok {
		logError("publish", errUnknownConnection.Error(), zap.String("name", frankenphp.GoString(unsafe.Pointer(name))))
		return
	}

	subj := frankenphp.GoString(unsafe.Pointer(subject))
	body := []byte(frankenphp.GoString(unsafe.Pointer(data)))

	hdr := headersFromPHPArray(headers)
	if len(hdr) == 0 {
		if err := conn.Publish(subj, body); err != nil {
			logError("publish", err.Error(), zap.String("subject", subj))
		}
		return
	}

	msg := &natsgo.Msg{
		Subject: subj,
		Data:    body,
		Header:  hdr,
	}
	if err := conn.PublishMsg(msg); err != nil {
		logError("publish", err.Error(), zap.String("subject", subj))
	}
}

//export_php:function request(string $name, string $subject, string $data, int $timeout = 1000000000): ?array
func natsRequest(name *C.zend_string, subject *C.zend_string, data *C.zend_string, timeoutNs int64) unsafe.Pointer {
	conn, ok := lookupConn(name)
	if !ok {
		logError("request", errUnknownConnection.Error(), zap.String("name", frankenphp.GoString(unsafe.Pointer(name))))
		return nil
	}

	ctx, cancel := context.WithTimeout(getContext(), time.Duration(timeoutNs))
	defer cancel()

	subj := frankenphp.GoString(unsafe.Pointer(subject))
	resp, err := conn.RequestWithContext(
		ctx,
		subj,
		[]byte(frankenphp.GoString(unsafe.Pointer(data))),
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, natsgo.ErrTimeout) {
			return nil
		}
		logError("request", err.Error(), zap.String("subject", subj))
		return nil
	}

	return frankenphp.PHPMap(messageToMap(resp))
}

//export_php:function subscribe(string $name, string $subject, string $queue = ""): string
func natsSubscribe(name *C.zend_string, subject *C.zend_string, queue *C.zend_string) unsafe.Pointer {
	connName := frankenphp.GoString(unsafe.Pointer(name))
	conn, ok := lookupConn(name)
	if !ok {
		logError("subscribe", errUnknownConnection.Error(), zap.String("name", connName))
		return frankenphp.PHPString("", false)
	}

	subj := frankenphp.GoString(unsafe.Pointer(subject))

	var (
		sub *natsgo.Subscription
		err error
	)
	if q := frankenphp.GoString(unsafe.Pointer(queue)); q != "" {
		sub, err = conn.QueueSubscribeSync(subj, q)
	} else {
		sub, err = conn.SubscribeSync(subj)
	}
	if err != nil {
		logError("subscribe", err.Error(), zap.String("subject", subj))
		return frankenphp.PHPString("", false)
	}

	id := uuid.NewString()
	subRegistryMu.Lock()
	subRegistry[id] = &subEntry{
		sub:      sub,
		connName: connName,
	}
	subRegistryMu.Unlock()

	return frankenphp.PHPString(id, false)
}

//export_php:function unsubscribe(string $subId): void
func natsUnsubscribe(subID *C.zend_string) {
	id := frankenphp.GoString(unsafe.Pointer(subID))

	subRegistryMu.Lock()
	entry, ok := subRegistry[id]
	delete(subRegistry, id)
	subRegistryMu.Unlock()
	if !ok {
		return
	}
	_ = entry.sub.Unsubscribe()
}

//export_php:function nextMessage(string $subId, int $timeout = 5000000000): ?array
func natsNextMessage(subID *C.zend_string, timeoutNs int64) unsafe.Pointer {
	id := frankenphp.GoString(unsafe.Pointer(subID))

	subRegistryMu.RLock()
	entry, ok := subRegistry[id]
	subRegistryMu.RUnlock()
	if !ok {
		logError("nextMessage", errUnknownSubscription.Error(), zap.String("subId", id))
		return nil
	}

	msg, err := entry.sub.NextMsg(time.Duration(timeoutNs))
	if err != nil {
		if errors.Is(err, natsgo.ErrTimeout) {
			return nil
		}
		// Subscription/connection became invalid (e.g. close() fired or the
		// server dropped us). Drop the entry so subscriptionValid() returns
		// false and a polling loop can detect the terminal state via
		// subscriptionValid($subId) instead of confusing it with a timeout.
		if errors.Is(err, natsgo.ErrBadSubscription) || errors.Is(err, natsgo.ErrConnectionClosed) {
			subRegistryMu.Lock()
			delete(subRegistry, id)
			subRegistryMu.Unlock()
		}
		logError("nextMessage", err.Error(), zap.String("subId", id))
		return nil
	}

	return frankenphp.PHPMap(messageToMap(msg))
}

//export_php:function subscriptionValid(string $subId): bool
func natsSubscriptionValid(subID *C.zend_string) bool {
	id := frankenphp.GoString(unsafe.Pointer(subID))

	subRegistryMu.RLock()
	entry, ok := subRegistry[id]
	subRegistryMu.RUnlock()
	if !ok {
		return false
	}
	return entry.sub.IsValid()
}

// ----------------------------------------------------------------------------
// Internal helpers.
// ----------------------------------------------------------------------------

func lookupConn(name *C.zend_string) (*natsgo.Conn, bool) {
	connName := frankenphp.GoString(unsafe.Pointer(name))
	connRegistryMu.RLock()
	defer connRegistryMu.RUnlock()
	conn, ok := connRegistry[connName]
	return conn, ok
}

func serversFromPHPArray(arr *C.zend_array) ([]string, error) {
	if arr == nil {
		return nil, errors.New("servers must be a non-empty array of strings")
	}
	values, err := frankenphp.GoPackedArray[any](unsafe.Pointer(arr))
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0, len(values))
	for _, v := range values {
		s, ok := v.(string)
		if !ok || s == "" {
			return nil, errors.New("each server entry must be a non-empty string")
		}
		urls = append(urls, s)
	}
	return urls, nil
}

func joinURLs(urls []string) string {
	if len(urls) == 1 {
		return urls[0]
	}
	out := urls[0]
	for _, u := range urls[1:] {
		out += "," + u
	}
	return out
}

// forceTLSScheme rewrites server URLs so the nats.go client picks the TLS
// dialer. natsgo.Secure() alone is not enough: a `nats://` URL still triggers
// a plain TCP dial. Bare `host:port` strings get a `tls://` prefix; explicit
// schemes other than nats:// (e.g. ws://) are left alone for the caller to
// handle.
func forceTLSScheme(urls []string) []string {
	out := make([]string, len(urls))
	for i, u := range urls {
		switch {
		case strings.HasPrefix(u, "tls://"):
			out[i] = u
		case strings.HasPrefix(u, "nats://"):
			out[i] = "tls://" + strings.TrimPrefix(u, "nats://")
		case strings.Contains(u, "://"):
			out[i] = u
		default:
			out[i] = "tls://" + u
		}
	}
	return out
}

// headersFromPHPArray accepts both array<string,string> and array<string,string[]>.
// A scalar string value becomes a single-valued NATS header.
func headersFromPHPArray(arr *C.zend_array) natsgo.Header {
	out := natsgo.Header{}
	if arr == nil {
		return out
	}
	m, err := frankenphp.GoMap[any](unsafe.Pointer(arr))
	if err != nil {
		return out
	}
	for k, v := range m {
		switch typed := v.(type) {
		case string:
			out.Add(k, typed)
		case []any:
			for _, item := range typed {
				if s, ok := item.(string); ok {
					out.Add(k, s)
				}
			}
		}
	}
	return out
}

func messageToMap(msg *natsgo.Msg) map[string]any {
	headers := make(map[string]any, len(msg.Header))
	for k, vs := range msg.Header {
		out := make([]any, 0, len(vs))
		for _, v := range vs {
			out = append(out, v)
		}
		headers[k] = out
	}
	var reply any
	if msg.Reply != "" {
		reply = msg.Reply
	}
	return map[string]any{
		"subject": msg.Subject,
		"data":    string(msg.Data),
		"reply":   reply,
		"headers": headers,
	}
}

// logError records an error via the Caddy logger. The generator-based
// extension cannot raise PHP exceptions directly from Go, so failures
// surface to PHP as zero values (null, false, empty) accompanied by an
// entry in the FrankenPHP error log.
func logError(op, msg string, fields ...zap.Field) {
	logger := getLogger()
	if logger == nil {
		return
	}
	logger.Error("nats: "+op+": "+msg, fields...)
}
