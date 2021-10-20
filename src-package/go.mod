module github.com/panlq/src-pkg

go 1.16

require internal/foo v1.0.0
replace internal/foo => ./internal/foo
require internal/bar v1.0.0
replace internal/bar => ./internal/bar
