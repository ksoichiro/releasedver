module c

go 1.14

require (
	github.com/gostaticanalysis/modfile v0.0.0-20200119053902-fbaf1a3ff141 // want `github.com/gostaticanalysis/modfile`
	golang.org/x/mod v0.3.0
	golang.org/x/tools v0.0.0-20200515010526-7d3b6ebf133d // want `golang.org/x/tools must use released version`
)
