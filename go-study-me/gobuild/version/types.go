package version

type Info struct {
	GitCommit string `json:"gitCommit"` // sha1 from git, output of $(git rev-parse HEAD)
	BuildDate string `json:"buildDate"` // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	GoVersion string `json:"goVersion"` // output of go version
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return info.GitCommit
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in ./base.go
	return Info{
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: goVersion,
	}
}
