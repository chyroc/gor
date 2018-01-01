package middlerware

import (
	"fmt"

	"github.com/Chyroc/gor"
)

// Recover default gor Recover middleware
var Recover = func(req *gor.Req, res *gor.Res, next gor.Next) {
	defer func() {
		if rev := recover(); rev != nil {
			res.Error(fmt.Sprintf("panic: %s\n", rev))
		}
	}()
	next()
}
