package datastore

var (
	db *Database
)

// NewBaseRepo returns a base repository.
func NewBaseRepo(db *Database) *BaseRepo {
	return &BaseRepo{db: db}
}

// BaseRepo is a base repository.
type BaseRepo struct {
	db *Database
}

type Database struct{}

func GetDatabase() *Database {
	if db != nil {
		return db
	}

	return db
}
