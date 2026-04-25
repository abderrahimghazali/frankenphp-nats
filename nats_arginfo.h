/* This is a generated file, edit the .stub.php file instead.
 * Stub hash: 5c4fdd2890a101b00196fb1f44732445192a3ae2 */

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_connect, 0, 2, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, name, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, servers, IS_ARRAY, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, username, IS_STRING, 0, "\"\"")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, password, IS_STRING, 0, "\"\"")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, token, IS_STRING, 0, "\"\"")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, credsFile, IS_STRING, 0, "\"\"")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, nkeyFile, IS_STRING, 0, "\"\"")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, tls, _IS_BOOL, 0, "false")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_LONG, 0, "2000000000")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, reconnectAttempts, IS_LONG, 0, "60")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, reconnectWait, IS_LONG, 0, "2000000000")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, pingInterval, IS_LONG, 0, "120000000000")
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, maxPingsOut, IS_LONG, 0, "2")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_close, 0, 1, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, name, IS_STRING, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_isConnected, 0, 1, _IS_BOOL, 0)
	ZEND_ARG_TYPE_INFO(0, name, IS_STRING, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_flush, 0, 1, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, name, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_LONG, 0, "5000000000")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_stats, 0, 1, IS_ARRAY, 0)
	ZEND_ARG_TYPE_INFO(0, name, IS_STRING, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_publish, 0, 3, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, name, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, subject, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, data, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, headers, IS_ARRAY, 0, "[]")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_request, 0, 3, IS_ARRAY, 1)
	ZEND_ARG_TYPE_INFO(0, name, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, subject, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, data, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_LONG, 0, "1000000000")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_subscribe, 0, 2, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, name, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO(0, subject, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, queue, IS_STRING, 0, "\"\"")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_unsubscribe, 0, 1, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, subId, IS_STRING, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_nextMessage, 0, 1, IS_ARRAY, 1)
	ZEND_ARG_TYPE_INFO(0, subId, IS_STRING, 0)
	ZEND_ARG_TYPE_INFO_WITH_DEFAULT_VALUE(0, timeout, IS_LONG, 0, "5000000000")
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_Abderrahim_Nats_subscriptionValid, 0, 1, _IS_BOOL, 0)
	ZEND_ARG_TYPE_INFO(0, subId, IS_STRING, 0)
ZEND_END_ARG_INFO()

ZEND_FUNCTION(Abderrahim_Nats_connect);
ZEND_FUNCTION(Abderrahim_Nats_close);
ZEND_FUNCTION(Abderrahim_Nats_isConnected);
ZEND_FUNCTION(Abderrahim_Nats_flush);
ZEND_FUNCTION(Abderrahim_Nats_stats);
ZEND_FUNCTION(Abderrahim_Nats_publish);
ZEND_FUNCTION(Abderrahim_Nats_request);
ZEND_FUNCTION(Abderrahim_Nats_subscribe);
ZEND_FUNCTION(Abderrahim_Nats_unsubscribe);
ZEND_FUNCTION(Abderrahim_Nats_nextMessage);
ZEND_FUNCTION(Abderrahim_Nats_subscriptionValid);

static const zend_function_entry ext_functions[] = {
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "connect"), zif_Abderrahim_Nats_connect, arginfo_Abderrahim_Nats_connect, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "close"), zif_Abderrahim_Nats_close, arginfo_Abderrahim_Nats_close, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "isConnected"), zif_Abderrahim_Nats_isConnected, arginfo_Abderrahim_Nats_isConnected, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "flush"), zif_Abderrahim_Nats_flush, arginfo_Abderrahim_Nats_flush, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "stats"), zif_Abderrahim_Nats_stats, arginfo_Abderrahim_Nats_stats, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "publish"), zif_Abderrahim_Nats_publish, arginfo_Abderrahim_Nats_publish, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "request"), zif_Abderrahim_Nats_request, arginfo_Abderrahim_Nats_request, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "subscribe"), zif_Abderrahim_Nats_subscribe, arginfo_Abderrahim_Nats_subscribe, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "unsubscribe"), zif_Abderrahim_Nats_unsubscribe, arginfo_Abderrahim_Nats_unsubscribe, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "nextMessage"), zif_Abderrahim_Nats_nextMessage, arginfo_Abderrahim_Nats_nextMessage, 0, NULL, NULL)
	ZEND_RAW_FENTRY(ZEND_NS_NAME("Abderrahim\\Nats", "subscriptionValid"), zif_Abderrahim_Nats_subscriptionValid, arginfo_Abderrahim_Nats_subscriptionValid, 0, NULL, NULL)
	ZEND_FE_END
};

static void register_nats_symbols(int module_number)
{
	REGISTER_LONG_CONSTANT("Abderrahim\\Nats\\NANOSECOND", 1, CONST_PERSISTENT);
	REGISTER_LONG_CONSTANT("Abderrahim\\Nats\\MICROSECOND", 1000, CONST_PERSISTENT);
	REGISTER_LONG_CONSTANT("Abderrahim\\Nats\\MILLISECOND", 1000000, CONST_PERSISTENT);
	REGISTER_LONG_CONSTANT("Abderrahim\\Nats\\SECOND", 1000000000, CONST_PERSISTENT);
	REGISTER_LONG_CONSTANT("Abderrahim\\Nats\\MINUTE", 60000000000, CONST_PERSISTENT);
}
