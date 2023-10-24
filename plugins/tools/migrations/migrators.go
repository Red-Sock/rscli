package migrations

type MigrationTool interface {
	Install() error
	Version() (version string, err error)
	GetLatestVersion() (version string, err error)
}
