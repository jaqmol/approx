package run

// DataProvider ...
type DataProvider interface {
	Connection(hash uint32) *ConnItem
	FormationBasePath() string
}
