package version

var (
	buildDate = "" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	goVersion = "" // go version go1.17 linux/amd64
	gitCommit = "" // short commit id
)
