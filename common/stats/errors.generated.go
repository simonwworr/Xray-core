package stats

import "github.com/xtls/xray-core/common/errors"

type errPathObjHolder struct{}

// newError creates a new error with the stats package path context.
func newError(values ...interface{}) *errors.Error {
	return errors.New(values...).WithPathObj(errPathObjHolder{})
}
