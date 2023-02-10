package datastore

var (
	db *database
)

// NewBaseRepo returns a base repository.
func NewBaseRepo(db *database) *BaseRepo {
	return &BaseRepo{db: db}
}

// BaseRepo is a base repository.
type BaseRepo struct {
	db *database
}

type database struct{}

func GetDatabase() *database {
	if db != nil {
		return db
	}

	return db
}
