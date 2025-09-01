package constants

// 🚨 FIRE SALAMANDER - HTTP Constants
// Zero Hardcoding Policy - All HTTP-related constants

// ===== HTTP STATUS CODES =====

// 1xx Informational responses
const (
	HTTPStatusContinue           = 100
	HTTPStatusSwitchingProtocols = 101
	HTTPStatusProcessing         = 102
	HTTPStatusEarlyHints         = 103
)

// 2xx Success (HTTPStatusOK already in report_constants)
const (
	HTTPStatusCreated              = 201
	HTTPStatusAccepted             = 202
	HTTPStatusNonAuthoritativeInfo = 203
	HTTPStatusNoContent            = 204
	HTTPStatusResetContent         = 205
	HTTPStatusPartialContent       = 206
	HTTPStatusMultiStatus          = 207
	HTTPStatusAlreadyReported      = 208
	HTTPStatusIMUsed               = 226
)

// 3xx Redirection
const (
	HTTPStatusMultipleChoices   = 300
	HTTPStatusMovedPermanently  = 301
	HTTPStatusFound             = 302
	HTTPStatusSeeOther          = 303
	HTTPStatusNotModified       = 304
	HTTPStatusUseProxy          = 305
	HTTPStatusTemporaryRedirect = 307
	HTTPStatusPermanentRedirect = 308
)

// 4xx Client errors (HTTPStatusBadRequest, HTTPStatusNotFound already in report_constants)  
const (
	HTTPStatusUnauthorized                   = 401
	HTTPStatusPaymentRequired                = 402
	HTTPStatusForbidden                      = 403
	HTTPStatusMethodNotAllowed               = 405
	HTTPStatusNotAcceptable                  = 406
	HTTPStatusProxyAuthRequired              = 407
	HTTPStatusRequestTimeout                 = 408
	HTTPStatusConflict                       = 409
	HTTPStatusGone                          = 410
	HTTPStatusLengthRequired                = 411
	HTTPStatusPreconditionFailed            = 412
	HTTPStatusRequestEntityTooLarge         = 413
	HTTPStatusRequestURITooLong             = 414
	HTTPStatusUnsupportedMediaType          = 415
	HTTPStatusRequestedRangeNotSatisfiable  = 416
	HTTPStatusExpectationFailed             = 417
	HTTPStatusTeapot                        = 418
	HTTPStatusMisdirectedRequest            = 421
	HTTPStatusUnprocessableEntity           = 422
	HTTPStatusLocked                        = 423
	HTTPStatusFailedDependency              = 424
	HTTPStatusTooEarly                      = 425
	HTTPStatusUpgradeRequired               = 426
	HTTPStatusPreconditionRequired          = 428
	HTTPStatusTooManyRequests               = 429
	HTTPStatusRequestHeaderFieldsTooLarge   = 431
	HTTPStatusUnavailableForLegalReasons    = 451
)

// 5xx Server errors
const (
	HTTPStatusInternalServerError           = 500
	HTTPStatusNotImplemented                = 501
	HTTPStatusBadGateway                    = 502
	HTTPStatusServiceUnavailable            = 503
	HTTPStatusGatewayTimeout                = 504
	HTTPStatusHTTPVersionNotSupported       = 505
	HTTPStatusVariantAlsoNegotiates         = 506
	HTTPStatusInsufficientStorage           = 507
	HTTPStatusLoopDetected                  = 508
	HTTPStatusNotExtended                   = 510
	HTTPStatusNetworkAuthenticationRequired = 511
)

// ===== HTTP METHODS ===== (basic ones already in web_server_constants)
const (
	HTTPMethodPATCH   = "PATCH"
	HTTPMethodHEAD    = "HEAD"
	HTTPMethodTRACE   = "TRACE"
	HTTPMethodCONNECT = "CONNECT"
)

// ===== CONTENT TYPES ===== (avoiding duplicates with debug_constants and web_server_constants)
const (
	ContentTypeXML             = "application/xml"
	ContentTypePlain           = "text/plain"
	ContentTypeCSS             = "text/css"
	ContentTypeJavaScript      = "application/javascript"
	ContentTypeFormURLEncoded  = "application/x-www-form-urlencoded"
	ContentTypeMultipartForm   = "multipart/form-data"
	ContentTypePNG             = "image/png"
	ContentTypeJPEG            = "image/jpeg"
	ContentTypeGIF             = "image/gif"
	ContentTypeSVG             = "image/svg+xml"
	ContentTypeWebP            = "image/webp"
	ContentTypeMP4             = "video/mp4"
	ContentTypeZIP             = "application/zip"
)

// ===== HTTP HEADERS =====
const (
	HeaderAccept                          = "Accept"
	HeaderAcceptCharset                   = "Accept-Charset"
	HeaderAcceptEncoding                  = "Accept-Encoding"
	HeaderAcceptLanguage                  = "Accept-Language"
	HeaderAcceptRanges                    = "Accept-Ranges"
	HeaderAccessControlAllowCredentials   = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowHeaders       = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowMethods       = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowOrigin        = "Access-Control-Allow-Origin"
	HeaderAccessControlExposeHeaders      = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge             = "Access-Control-Max-Age"
	HeaderAccessControlRequestHeaders     = "Access-Control-Request-Headers"
	HeaderAccessControlRequestMethod      = "Access-Control-Request-Method"
	HeaderAge                             = "Age"
	HeaderAllow                           = "Allow"
	HeaderAuthorization                   = "Authorization"
	HeaderConnection                      = "Connection"
	HeaderContentEncoding                 = "Content-Encoding"
	HeaderContentLanguage                 = "Content-Language"
	HeaderContentLocation                 = "Content-Location"
	HeaderContentRange                    = "Content-Range"
	HeaderCookie                          = "Cookie"
	HeaderDate                            = "Date"
	HeaderETag                            = "ETag"
	HeaderExpect                          = "Expect"
	HeaderExpires                         = "Expires"
	HeaderFrom                            = "From"
	HeaderHost                            = "Host"
	HeaderIfMatch                         = "If-Match"
	HeaderIfModifiedSince                 = "If-Modified-Since"
	HeaderIfNoneMatch                     = "If-None-Match"
	HeaderIfRange                         = "If-Range"
	HeaderIfUnmodifiedSince               = "If-Unmodified-Since"
	HeaderLastModified                    = "Last-Modified"
	HeaderLink                            = "Link"
	HeaderLocation                        = "Location"
	HeaderMaxForwards                     = "Max-Forwards"
	HeaderOrigin                          = "Origin"
	HeaderPragma                          = "Pragma"
	HeaderProxyAuthenticate               = "Proxy-Authenticate"
	HeaderProxyAuthorization              = "Proxy-Authorization"
	HeaderPublicKeyPins                   = "Public-Key-Pins"
	HeaderRange                           = "Range"
	HeaderReferer                         = "Referer"
	HeaderRetryAfter                      = "Retry-After"
	HeaderSecWebSocketAccept              = "Sec-WebSocket-Accept"
	HeaderSecWebSocketExtensions          = "Sec-WebSocket-Extensions"
	HeaderSecWebSocketKey                 = "Sec-WebSocket-Key"
	HeaderSecWebSocketProtocol            = "Sec-WebSocket-Protocol"
	HeaderSecWebSocketVersion             = "Sec-WebSocket-Version"
	HeaderServer                          = "Server"
	HeaderSetCookie                       = "Set-Cookie"
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderTE                              = "TE"
	HeaderTrailer                         = "Trailer"
	HeaderTransferEncoding                = "Transfer-Encoding"
	HeaderUpgrade                         = "Upgrade"
	HeaderUpgradeInsecureRequests         = "Upgrade-Insecure-Requests"
	HeaderUserAgent                       = "User-Agent"
	HeaderVary                            = "Vary"
	HeaderVia                             = "Via"
	HeaderWarning                         = "Warning"
	HeaderWWWAuthenticate                 = "WWW-Authenticate"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderXForwardedFor                   = "X-Forwarded-For"
	HeaderXForwardedHost                  = "X-Forwarded-Host"
	HeaderXForwardedProto                 = "X-Forwarded-Proto"
	HeaderXRealIP                         = "X-Real-IP"
	HeaderXRequestedWith                  = "X-Requested-With"
	HeaderXXSSProtection                  = "X-XSS-Protection"
)

// ===== USER AGENT STRINGS =====
const (
	UserAgentFireSalamander = "FireSalamander/1.0 (SEO Crawler; +https://github.com/firesalamander)"
	UserAgentChrome         = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	UserAgentFirefox        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0"
	UserAgentSafari         = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15"
	UserAgentBot            = "Mozilla/5.0 (compatible; bot/1.0; +http://example.com/bot)"
)

// ===== URL SCHEMES =====
const (
	SchemeHTTP    = "http"
	SchemeHTTPS   = "https"
	SchemeFTP     = "ftp"
	SchemeFTPS    = "ftps"
	SchemeFile    = "file"
	SchemeMailto  = "mailto"
	SchemeTel     = "tel"
	SchemeSMS     = "sms"
	SchemeData    = "data"
)

// ===== DEFAULT PORTS =====
const (
	DefaultPortHTTP     = 80
	DefaultPortHTTPS    = 443
	DefaultPortFTP      = 21
	DefaultPortSSH      = 22
	DefaultPortTelnet   = 23
	DefaultPortSMTP     = 25
	DefaultPortDNS      = 53
	DefaultPortPOP3     = 110
	DefaultPortIMAP     = 143
	DefaultPortSNMP     = 161
	DefaultPortLDAP     = 389
	DefaultPortSMTPS    = 465
	DefaultPortIMAPS    = 993
	DefaultPortPOP3S    = 995
	DefaultPortMySQL    = 3306
	DefaultPortPostgres = 5432
	DefaultPortRedis    = 6379
	DefaultPortMongoDB  = 27017
	DefaultPortElastic  = 9200
)

// ===== TIMEOUT VALUES =====
const (
	DefaultRequestTimeoutSec     = 30
	DefaultConnectTimeoutSec     = 10
	DefaultReadTimeoutSec        = 30
	DefaultWriteTimeoutSec       = 30
	DefaultIdleTimeoutSec        = 120
	DefaultKeepAliveTimeoutSec   = 60
	DefaultHandshakeTimeoutSec   = 10
	DefaultDialTimeoutSec        = 30
)

// ===== RATE LIMITING =====
const (
	DefaultRateLimitPerSecond = 10
	DefaultRateLimitPerMinute = 600
	DefaultRateLimitPerHour   = 36000
	DefaultBurstSize          = 20
)

// ===== CACHE CONTROL VALUES =====
const (
	CacheControlNoCache       = "no-cache"
	CacheControlNoStore       = "no-store"
	CacheControlMustRevalidate = "must-revalidate"
	CacheControlPublic        = "public"
	CacheControlPrivate       = "private"
	CacheControlMaxAge        = "max-age"
)