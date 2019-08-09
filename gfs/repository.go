package gfs

const (
	// NCEPRepoType get files from NCEP
	NCEPRepoType RepositoryType = "NCEP"
	// NCDCRepoType get files from NCDC
	NCDCRepoType RepositoryType = "NCDC"
)

// RepositoryType the type of repository being accessed
type RepositoryType string

// Repository interface for different NOMADS file servers
type Repository interface {
	GetBaseURL() (string, error)
	GetURIs() ([]string, error)
	GetURIsForDate(date string) ([]string, error)
	GetURIsForDateAndTime(date string, timeFrame TimeFrame) ([]string, error)
	LoadParams(*Params) error
}

// NewRepository NewRepository
func NewRepository(rt RepositoryType) Repository {
	if rt == NCEPRepoType {
		return new(NCEPRepository)
	} else if rt == NCDCRepoType {
		return nil
	}
	return nil
}
