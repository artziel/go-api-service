package rest

import "errors"

var ErrFormToStructPtrExpected = errors.New("expected error FormToStruct function")
var ErrGracefullShutdown = errors.New("service stopped gracefully")
