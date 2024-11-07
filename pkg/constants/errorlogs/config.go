package errorlogs

const (
	ParsingError              = "got error %v on parsong configuration file"
	InvalidKeySize            = "invalid key size : must be atleast %d characters long"
	BindingJsonError          = "got error while binding the request. error :: %s"
	InvalidEmailFormatError   = "got error while parsing the email. error :: %s"
	CheckUserExistenceError   = "got error while checking the user existence. error :: %s"
	GetDBError                = "got error while getting the db connection. error :: %s"
	DataBaseGetError          = "got error while getting the data from db. error :: %s"
	InsertTxnError            = "got error while inserting into table %s. error :: %s"
	GetUserByNumberError      = "got error while getting the user with phone number %s. error :: %s"
	BeginSQLTransactionError  = "got error while starting a transaction. error :: %s"
	HashedPasswordError       = "got error while hashing the password. error :: %v"
	InsertUserError           = "got error while inserting the user. error :: %s"
	CommitSQLTransactionError = "got error while committing the transaction for table %s. error :: %s"
	GetUserByPhoneNumberError = "got error while getting the user with phone number %s. error :: %s"
	TokenError                = "got error while generating the token. error :: %s"
)
