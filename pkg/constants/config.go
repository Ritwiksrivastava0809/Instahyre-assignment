package constants

// server related constants
const (
	ConstantDb    = "db"
	Origin        = "Origin"
	ContentLength = "Content-Length"
	ContentType   = "Content-Type"
	Authorization = "Authorization"
)

// db related constants
const (
	UserTable    = "users"
	ConactsTable = "contacts"
	SpamReport   = "spam_reports"
	GlobalSpam   = "global_spam"
)

// env related constants
const (
	DefaultConfigType = "yaml"
	DefaultConfigPath = "environment"
)

// error messages
const (
	ExipredToken              = "token has expired"
	InvalidToken              = "token is invalid"
	JWTValidationErrorExpired = 512
)

// handler related constants
const (
	MinSecretKeyLen         = 32
	TokenMaker              = "tokenMaker"
	ConstantPayload         = "payload"
	Bearer                  = "Bearer"
	UserID                  = "userID"
	AuthorizationPayloadKey = "authorization_payload"
	Pending                 = "pending"
	Approved                = "approved"
	Spam                    = "spam"
	Name                    = "name"
	PhoneNumber             = "phone_number"
	Email                   = "email"
	SpamLikelihood          = "spam_likelihood"
)
