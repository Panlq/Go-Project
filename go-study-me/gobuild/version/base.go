package version

var (
	buildDate  = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	goVersion  = ""
	gitCommit  = ""
	gitVersion = ""
	platform   = ""
)
