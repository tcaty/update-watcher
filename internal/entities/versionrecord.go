package entities

// VersionRecord is the main entity in this project
type VersionRecord struct {
	// Target means some target to watch updates for
	// It could be grafana dashboard id, docker image name etc.
	Target string
	// Version means the version of the current target
	// It could be revision in case of grafana dashboard or
	// a tag in docker image case
	Version string
}
