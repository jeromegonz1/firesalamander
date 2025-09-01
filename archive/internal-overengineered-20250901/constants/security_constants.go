package constants

// ========================================
// DELTA-5 SECURITY CONSTANTS
// Security Analysis and Vulnerability Assessment Constants
// ========================================

// ========================================
// SECURITY VULNERABILITIES
// ========================================

// Common Web Vulnerabilities
const (
	SecurityVulnerabilityXSS                    = "xss"
	SecurityVulnerabilitySQLInjection           = "sql_injection"
	SecurityVulnerabilityCSRF                   = "csrf"
	SecurityVulnerabilityClickjacking           = "clickjacking"
	SecurityVulnerabilityDirectoryTraversal     = "directory_traversal"
	SecurityVulnerabilityCommandInjection       = "command_injection"
	SecurityVulnerabilityLDAPInjection          = "ldap_injection"
	SecurityVulnerabilityXPathInjection         = "xpath_injection"
	SecurityVulnerabilityXXE                    = "xxe"
	SecurityVulnerabilitySSRF                   = "ssrf"
	SecurityVulnerabilityRFI                    = "rfi"
	SecurityVulnerabilityLFI                    = "lfi"
)

// Application Security Vulnerabilities
const (
	SecurityVulnerabilityDeserialization           = "deserialization"
	SecurityVulnerabilityBufferOverflow            = "buffer_overflow"
	SecurityVulnerabilityRaceCondition             = "race_condition"
	SecurityVulnerabilityPrivilegeEscalation       = "privilege_escalation"
	SecurityVulnerabilityInformationDisclosure     = "information_disclosure"
	SecurityVulnerabilityDenialOfService           = "denial_of_service"
	SecurityVulnerabilityBrokenAuthentication      = "broken_authentication"
	SecurityVulnerabilitySessionFixation           = "session_fixation"
	SecurityVulnerabilitySensitiveDataExposure     = "sensitive_data_exposure"
	SecurityVulnerabilitySecurityMisconfiguration  = "security_misconfiguration"
	SecurityVulnerabilityInsecureDeserialization   = "insecure_deserialization"
)

// ========================================
// SECURITY HEADERS
// ========================================

// Content Security Headers
const (
	SecurityHeaderCSP                     = "Content-Security-Policy"
	SecurityHeaderXFrameOptions           = "X-Frame-Options"
	SecurityHeaderXXSSProtection          = "X-XSS-Protection"
	SecurityHeaderXContentTypeOptions     = "X-Content-Type-Options"
	SecurityHeaderReferrerPolicy          = "Referrer-Policy"
	SecurityHeaderFeaturePolicy          = "Feature-Policy"
	SecurityHeaderPermissionsPolicy       = "Permissions-Policy"
)

// Transport Security Headers
const (
	SecurityHeaderHSTS                    = "Strict-Transport-Security"
	SecurityHeaderExpectCT                = "Expect-CT"
	SecurityHeaderPKP                     = "Public-Key-Pins"
)

// Cross-Origin Headers
const (
	SecurityHeaderCOEP                               = "Cross-Origin-Embedder-Policy"
	SecurityHeaderCOOP                               = "Cross-Origin-Opener-Policy"
	SecurityHeaderCORP                               = "Cross-Origin-Resource-Policy"
	SecurityHeaderXPermittedCrossDomainPolicies      = "X-Permitted-Cross-Domain-Policies"
)

// ========================================
// SECURITY PROTOCOLS
// ========================================

// Network Security Protocols
const (
	SecurityProtocolHTTPS   = "https"
	SecurityProtocolTLS     = "tls"
	SecurityProtocolSSL     = "ssl"
	SecurityProtocolWSS     = "wss"
	SecurityProtocolSFTP    = "sftp"
	SecurityProtocolFTPS    = "ftps"
	SecurityProtocolSSH     = "ssh"
)

// Cryptographic Protocols
const (
	SecurityProtocolPGP     = "pgp"
	SecurityProtocolGPG     = "gpg"
	SecurityProtocolAES     = "aes"
	SecurityProtocolDES     = "des"
	SecurityProtocolRSA     = "rsa"
	SecurityProtocolECDSA   = "ecdsa"
	SecurityProtocolHMAC    = "hmac"
)

// Hash Protocols
const (
	SecurityProtocolSHA1    = "sha1"
	SecurityProtocolSHA256  = "sha256"
	SecurityProtocolSHA512  = "sha512"
	SecurityProtocolMD5     = "md5"
	SecurityProtocolBCrypt  = "bcrypt"
	SecurityProtocolSCrypt  = "scrypt"
	SecurityProtocolPBKDF2  = "pbkdf2"
)

// Authentication Protocols
const (
	SecurityProtocolOAuth   = "oauth"
	SecurityProtocolSAML    = "saml"
	SecurityProtocolJWT     = "jwt"
	SecurityProtocolOpenID  = "openid"
)

// ========================================
// SECURITY LEVELS
// ========================================

// Risk Levels
const (
	SecurityLevelLow        = "low"
	SecurityLevelMedium     = "medium"
	SecurityLevelHigh       = "high"
	SecurityLevelCritical   = "critical"
)

// Log Levels
const (
	SecurityLevelInfo       = "info"
	SecurityLevelWarning    = "warning"
	SecurityLevelError      = "error"
	SecurityLevelFatal      = "fatal"
	SecurityLevelTrace      = "trace"
	SecurityLevelDebug      = "debug"
)

// Security Maturity Levels
const (
	SecurityLevelNone           = "none"
	SecurityLevelBasic          = "basic"
	SecurityLevelIntermediate   = "intermediate"
	SecurityLevelAdvanced       = "advanced"
	SecurityLevelExpert         = "expert"
)

// ========================================
// SECURITY ACTIONS
// ========================================

// Access Control Actions
const (
	SecurityActionAllow     = "allow"
	SecurityActionDeny      = "deny"
	SecurityActionBlock     = "block"
	SecurityActionPermit    = "permit"
	SecurityActionReject    = "reject"
	SecurityActionIgnore    = "ignore"
)

// Monitoring Actions
const (
	SecurityActionLog       = "log"
	SecurityActionAlert     = "alert"
	SecurityActionWarn      = "warn"
	SecurityActionMonitor   = "monitor"
	SecurityActionAudit     = "audit"
)

// Testing Actions
const (
	SecurityActionScan      = "scan"
	SecurityActionTest      = "test"
	SecurityActionValidate  = "validate"
	SecurityActionVerify    = "verify"
)

// Authentication Actions
const (
	SecurityActionAuthenticate = "authenticate"
	SecurityActionAuthorize    = "authorize"
)

// Cryptographic Actions
const (
	SecurityActionEncrypt         = "encrypt"
	SecurityActionDecrypt         = "decrypt"
	SecurityActionHash            = "hash"
	SecurityActionSign            = "sign"
	SecurityActionVerifySignature = "verify_signature"
)

// Input Processing Actions
const (
	SecurityActionSanitize    = "sanitize"
	SecurityActionEscape      = "escape"
	SecurityActionFilter      = "filter"
	SecurityActionWhitelist   = "whitelist"
	SecurityActionBlacklist   = "blacklist"
	SecurityActionQuarantine  = "quarantine"
)

// ========================================
// SECURITY CATEGORIES
// ========================================

// Authentication & Authorization
const (
	SecurityCategoryAuthentication    = "authentication"
	SecurityCategoryAuthorization     = "authorization"
	SecurityCategorySessionManagement = "session_management"
)

// Input & Output Security
const (
	SecurityCategoryInputValidation  = "input_validation"
	SecurityCategoryOutputEncoding   = "output_encoding"
	SecurityCategoryErrorHandling    = "error_handling"
)

// Technical Security Categories
const (
	SecurityCategoryLogging       = "logging"
	SecurityCategoryCryptography  = "cryptography"
	SecurityCategoryCommunication = "communication"
	SecurityCategoryConfiguration = "configuration"
)

// Application Security Categories
const (
	SecurityCategoryBusinessLogic  = "business_logic"
	SecurityCategoryFileUpload     = "file_upload"
	SecurityCategoryXML            = "xml"
	SecurityCategoryJSON           = "json"
	SecurityCategoryAPI            = "api"
	SecurityCategoryWebService     = "web_service"
)

// Infrastructure Security Categories
const (
	SecurityCategoryDatabase       = "database"
	SecurityCategoryInfrastructure = "infrastructure"
	SecurityCategoryNetwork        = "network"
	SecurityCategoryApplication    = "application"
	SecurityCategorySystem         = "system"
)

// ========================================
// OWASP TOP 10 CATEGORIES
// ========================================

// OWASP Top 10 2021 (A01-A10)
const (
	SecurityOWASPA01 = "A01"
	SecurityOWASPA02 = "A02"
	SecurityOWASPA03 = "A03"
	SecurityOWASPA04 = "A04"
	SecurityOWASPA05 = "A05"
	SecurityOWASPA06 = "A06"
	SecurityOWASPA07 = "A07"
	SecurityOWASPA08 = "A08"
	SecurityOWASPA09 = "A09"
	SecurityOWASPA10 = "A10"
)

// OWASP Vulnerability Categories
const (
	SecurityOWASPInjectionType                  = "injection"
	SecurityOWASPBrokenAuthentication           = "broken_authentication"
	SecurityOWASPSensitiveData                  = "sensitive_data"
	SecurityOWASPXMLExternalEntities            = "xml_external_entities"
	SecurityOWASPBrokenAccessControl            = "broken_access_control"
	SecurityOWASPSecurityMisconfiguration       = "security_misconfiguration"
	SecurityOWASPCrossSiteScripting             = "cross_site_scripting"
	SecurityOWASPInsecureDeserialization        = "insecure_deserialization"
	SecurityOWASPComponentsVulnerabilities      = "components_vulnerabilities"
	SecurityOWASPLoggingMonitoring              = "logging_monitoring"
)

// ========================================
// SECURITY TESTS
// ========================================

// Security Testing Types
const (
	SecurityTestPenetration     = "penetration"
	SecurityTestVulnerability   = "vulnerability"
	SecurityTestSecurity        = "security"
	SecurityTestCompliance      = "compliance"
	SecurityTestAudit           = "audit"
	SecurityTestAssessment      = "assessment"
	SecurityTestScan            = "scan"
	SecurityTestAnalysis        = "analysis"
	SecurityTestReview          = "review"
	SecurityTestInspection      = "inspection"
	SecurityTestValidation      = "validation"
	SecurityTestVerification    = "verification"
)

// ========================================
// SECURITY TOOLS
// ========================================

// Network Security Tools
const (
	SecurityToolNmap        = "nmap"
	SecurityToolNessus      = "nessus"
	SecurityToolBurp        = "burp"
	SecurityToolOWASPZAP    = "owasp_zap"
	SecurityToolMetasploit  = "metasploit"
	SecurityToolSQLMap      = "sqlmap"
	SecurityToolNikto       = "nikto"
	SecurityToolDirb        = "dirb"
	SecurityToolGobuster    = "gobuster"
)

// Password Security Tools
const (
	SecurityToolHydra       = "hydra"
	SecurityToolJohn        = "john"
	SecurityToolHashcat     = "hashcat"
)

// Network Analysis Tools
const (
	SecurityToolWireshark   = "wireshark"
	SecurityToolTcpdump     = "tcpdump"
	SecurityToolNcat        = "ncat"
	SecurityToolSocat       = "socat"
)

// Cryptographic Tools
const (
	SecurityToolOpenSSL     = "openssl"
	SecurityToolGPG         = "gpg"
	SecurityToolSSHKeygen   = "ssh_keygen"
)

// ========================================
// ENCRYPTION ALGORITHMS
// ========================================

// AES Variants
const (
	SecurityEncryptionAES128 = "AES-128"
	SecurityEncryptionAES192 = "AES-192"
	SecurityEncryptionAES256 = "AES-256"
)

// Legacy Encryption
const (
	SecurityEncryptionDES    = "DES"
	SecurityEncryption3DES   = "3DES"
)

// RSA Key Sizes
const (
	SecurityEncryptionRSA1024 = "RSA-1024"
	SecurityEncryptionRSA2048 = "RSA-2048"
	SecurityEncryptionRSA4096 = "RSA-4096"
)

// ECDSA Curves
const (
	SecurityEncryptionECDSA256 = "ECDSA-256"
	SecurityEncryptionECDSA384 = "ECDSA-384"
	SecurityEncryptionECDSA521 = "ECDSA-521"
)

// Hash Algorithms
const (
	SecurityHashSHA1    = "SHA-1"
	SecurityHashSHA224  = "SHA-224"
	SecurityHashSHA256  = "SHA-256"
	SecurityHashSHA384  = "SHA-384"
	SecurityHashSHA512  = "SHA-512"
	SecurityHashMD4     = "MD4"
	SecurityHashMD5     = "MD5"
)

// HMAC Variants
const (
	SecurityHashHMACSHA1   = "HMAC-SHA1"
	SecurityHashHMACSHA256 = "HMAC-SHA256"
	SecurityHashHMACSHA512 = "HMAC-SHA512"
)

// Password Hashing
const (
	SecurityHashPBKDF2 = "PBKDF2"
	SecurityHashBCrypt = "BCrypt"
	SecurityHashSCrypt = "SCrypt"
)

// ========================================
// SECURITY STATUS CODES
// ========================================

// Success Codes
const (
	SecurityStatusOK      = "200"
	SecurityStatusCreated = "201"
)

// Redirection Codes
const (
	SecurityStatusMovedPermanently = "301"
	SecurityStatusFound           = "302"
)

// Client Error Codes
const (
	SecurityStatusBadRequest          = "400"
	SecurityStatusUnauthorized        = "401"
	SecurityStatusForbidden           = "403"
	SecurityStatusNotFound            = "404"
	SecurityStatusMethodNotAllowed    = "405"
	SecurityStatusNotAcceptable       = "406"
	SecurityStatusConflict            = "409"
	SecurityStatusGone                = "410"
	SecurityStatusPayloadTooLarge     = "413"
	SecurityStatusURITooLong          = "414"
	SecurityStatusUnsupportedMediaType = "415"
	SecurityStatusUnprocessableEntity = "422"
	SecurityStatusTooManyRequests     = "429"
)

// Server Error Codes
const (
	SecurityStatusInternalServerError = "500"
	SecurityStatusNotImplemented      = "501"
	SecurityStatusBadGateway          = "502"
	SecurityStatusServiceUnavailable  = "503"
	SecurityStatusGatewayTimeout      = "504"
)

// ========================================
// SECURITY ENDPOINTS
// ========================================

// Authentication Endpoints
const (
	SecurityEndpointAuth     = "/auth"
	SecurityEndpointLogin    = "/login"
	SecurityEndpointLogout   = "/logout"
	SecurityEndpointRegister = "/register"
	SecurityEndpointReset    = "/reset"
	SecurityEndpointVerify   = "/verify"
	SecurityEndpoint2FA      = "/2fa"
	SecurityEndpointMFA      = "/mfa"
)

// OAuth Endpoints
const (
	SecurityEndpointOAuth = "/oauth"
	SecurityEndpointSAML  = "/saml"
	SecurityEndpointJWT   = "/jwt"
)

// API Security Endpoints
const (
	SecurityEndpointAPISecurity = "/api/security"
	SecurityEndpointAPIAuth     = "/api/auth"
	SecurityEndpointAPIValidate = "/api/validate"
	SecurityEndpointAPIScan     = "/api/scan"
	SecurityEndpointAPIAudit    = "/api/audit"
)

// ========================================
// SECURITY CONFIG KEYS
// ========================================

// General Security Configuration
const (
	SecurityConfigSecurity    = "security"
	SecurityConfigAuth        = "auth"
	SecurityConfigEncryption  = "encryption"
)

// SSL/TLS Configuration
const (
	SecurityConfigSSL         = "ssl"
	SecurityConfigTLS         = "tls"
	SecurityConfigCertificate = "certificate"
	SecurityConfigKey         = "key"
)

// Secrets Configuration
const (
	SecurityConfigSecret   = "secret"
	SecurityConfigToken    = "token"
	SecurityConfigPassword = "password"
	SecurityConfigHash     = "hash"
	SecurityConfigSalt     = "salt"
	SecurityConfigPepper   = "pepper"
	SecurityConfigNonce    = "nonce"
	SecurityConfigIV       = "iv"
)

// Cryptographic Configuration
const (
	SecurityConfigCipher    = "cipher"
	SecurityConfigAlgorithm = "algorithm"
	SecurityConfigProtocol  = "protocol"
)

// Session Configuration
const (
	SecurityConfigTimeout     = "timeout"
	SecurityConfigMaxAttempts = "max_attempts"
	SecurityConfigLockout     = "lockout"
	SecurityConfigSession     = "session"
	SecurityConfigCookie      = "cookie"
)

// Security Headers Configuration
const (
	SecurityConfigCORS = "cors"
	SecurityConfigCSP  = "csp"
	SecurityConfigHSTS = "hsts"
)

// ========================================
// COMPLIANCE STANDARDS
// ========================================

// Financial Compliance
const (
	SecurityCompliancePCIDSS = "PCI-DSS"
	SecurityComplianceSOX    = "SOX"
)

// Healthcare Compliance
const (
	SecurityComplianceHIPAA = "HIPAA"
)

// Privacy Compliance
const (
	SecurityComplianceGDPR = "GDPR"
)

// Security Standards
const (
	SecurityComplianceISO27001 = "ISO27001"
	SecurityComplianceNIST     = "NIST"
	SecurityComplianceOWASP    = "OWASP"
	SecurityComplianceCIS      = "CIS"
	SecurityComplianceSANS     = "SANS"
)

// IT Governance
const (
	SecurityComplianceCOBIT = "COBIT"
	SecurityComplianceITIL  = "ITIL"
)

// Cloud Compliance
const (
	SecurityComplianceSOC2    = "SOC2"
	SecurityComplianceFedRAMP = "FedRAMP"
	SecurityComplianceFISMA   = "FISMA"
)

// ========================================
// SECURITY FILE EXTENSIONS
// ========================================

// Certificate Files
const (
	SecurityFileExtensionKey        = ".key"
	SecurityFileExtensionCRT        = ".crt"
	SecurityFileExtensionPEM        = ".pem"
	SecurityFileExtensionP12        = ".p12"
	SecurityFileExtensionPFX        = ".pfx"
	SecurityFileExtensionJKS        = ".jks"
	SecurityFileExtensionKeystore   = ".keystore"
	SecurityFileExtensionTruststore = ".truststore"
	SecurityFileExtensionCER        = ".cer"
	SecurityFileExtensionDER        = ".der"
	SecurityFileExtensionCSR        = ".csr"
	SecurityFileExtensionP7B        = ".p7b"
	SecurityFileExtensionP7C        = ".p7c"
)

// ========================================
// SECURITY ERROR MESSAGES
// ========================================

// Authentication Errors
const (
	SecurityErrorUnauthorized        = "Accès non autorisé"
	SecurityErrorForbidden           = "Accès interdit"
	SecurityErrorInvalidCredentials  = "Identifiants invalides"
	SecurityErrorAuthenticationFailed = "Échec de l'authentification"
	SecurityErrorAuthorizationFailed  = "Échec de l'autorisation"
)

// Security Violation Errors
const (
	SecurityErrorSecurityViolation  = "Violation de sécurité détectée"
	SecurityErrorSuspiciousActivity = "Activité suspecte détectée"
	SecurityErrorAttackDetected     = "Attaque détectée"
	SecurityErrorVulnerabilityFound = "Vulnérabilité trouvée"
	SecurityErrorComplianceViolation = "Violation de conformité"
	SecurityErrorAuditFailed        = "Échec de l'audit de sécurité"
)

// ========================================
// SECURITY LOG MESSAGES
// ========================================

// Security Event Logs
const (
	SecurityLogSecurityScan          = "Scan de sécurité: %s"
	SecurityLogVulnerabilityAssessment = "Évaluation de vulnérabilité: %s"
	SecurityLogPenetrationTest       = "Test de pénétration: %s"
	SecurityLogSecurityAudit         = "Audit de sécurité: %s"
	SecurityLogComplianceCheck       = "Vérification de conformité: %s"
)

// Authentication Logs
const (
	SecurityLogAuthenticationAttempt = "Tentative d'authentification: %s"
	SecurityLogAuthorizationRequest  = "Demande d'autorisation: %s"
	SecurityLogSuspiciousRequest     = "Requête suspecte: %s"
	SecurityLogAttackAttempt         = "Tentative d'attaque: %s"
	SecurityLogSecurityEvent         = "Événement de sécurité: %s"
	SecurityLogSecurityAlert         = "Alerte de sécurité: %s"
)

// ========================================
// DEFAULT SECURITY VALUES
// ========================================

// Default Timeouts (in seconds)
const (
	SecurityDefaultAuthTimeout     = 300  // 5 minutes
	SecurityDefaultSessionTimeout  = 1800 // 30 minutes
	SecurityDefaultTokenTimeout    = 3600 // 1 hour
	SecurityDefaultScanTimeout     = 600  // 10 minutes
)

// Default Security Limits
const (
	SecurityDefaultMaxLoginAttempts = 5
	SecurityDefaultLockoutDuration  = 900  // 15 minutes
	SecurityDefaultPasswordLength   = 12
	SecurityDefaultKeySize          = 256
)

// Default Security Levels
const (
	SecurityDefaultRiskLevel      = SecurityLevelMedium
	SecurityDefaultLogLevel       = SecurityLevelInfo
	SecurityDefaultEncryption     = SecurityEncryptionAES256
	SecurityDefaultHashAlgorithm  = SecurityHashSHA256
)