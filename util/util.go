// Copy from peco/peco/internal/util/util.go

package util

type causer interface {
	Cause() error
}

type ignorable interface {
	Ignorable() bool
}

type collectResults interface {
	CollectResults() bool
}

func IsIgnorableError(err error) bool {
	for e := err; e != nil; {
		switch e.(type) {
		case ignorable:
			return e.(ignorable).Ignorable()
		case causer:
			e = e.(causer).Cause()
		default:
			return false
		}
	}
	return false
}

func IsCollectResultsError(err error) bool {
	for e := err; e != nil; {
		switch e.(type) {
		case collectResults:
			return e.(collectResults).CollectResults()
		case causer:
			e = e.(causer).Cause()
		default:
			return false
		}
	}
	return false
}
