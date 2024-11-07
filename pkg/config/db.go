package config

type DB struct {
	URL      string
	DBtype   string
	Database string
	Username string
	Password string
}

type DBtype string

const (
	RDS   DBtype = "rds"
	LOCAL DBtype = "local"
)

func (d DBtype) String() string {
	switch d {
	case RDS:
		return "rds"

	case LOCAL:
		return "local"
	}
	return "invalid"
}
