package config

const (
	KeyDefaultProfile = "default"

	KeyRedisBase = "redis"
	KeyRedisHost = "host"
	KeyRedisPort = "port"
	KeyRedisPass = "pass"
	KeyRedisDb   = "db"

	KeyMongoBase        = "mongo"
	KeyMongoAddrs       = "addrs"
	KeyMongoPerformAuth = "with_auth"

	KeyMongoAuthBase      = "auth"
	KeyMongoAuthDatabase  = "db"
	KeyMongoAuthUser      = "user"
	KeyMongoAuthPass      = "pass"
	KeyMongoAuthMechanism = "mechanism"

	KeyMongoTimeoutSec = "timeout_sec"
	KeyMongoDatabase   = "db"

	KeyMailBase = "mail"

	KeyMailSMTPBase                     = "smtp"
	KeyMailSMTPHost                     = "host"
	KeyMailSMTPPort                     = "port"
	KeyMailSMTPSender                   = "sender"
	KeyMailSMTPUser                     = "user"
	KeyMailSMTPPass                     = "pass"
	KeyMailSMTPOption                   = "option"
	KeyMailSMTPOptionInsecureSkipVerify = "insecure_skip_verify"
)
