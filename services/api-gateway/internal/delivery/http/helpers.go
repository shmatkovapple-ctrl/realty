package http

import (
"net/http"

"google.golang.org/grpc/codes"
)

func grpcStatusToHTTP(code codes.Code) int {
switch code {
case codes.OK:
return http.StatusOK
case codes.InvalidArgument:
return http.StatusBadRequest
case codes.NotFound:
return http.StatusNotFound
case codes.AlreadyExists:
return http.StatusConflict
case codes.PermissionDenied:
return http.StatusForbidden
case codes.Unauthenticated:
return http.StatusUnauthorized
case codes.ResourceExhausted:
return http.StatusTooManyRequests
case codes.Unimplemented:
return http.StatusNotImplemented
default:
return http.StatusInternalServerError
}
}
