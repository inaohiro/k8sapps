package token

import "k8soperation/token/internal/controller"

var (
	Routes               = controller.Controller()
	TokenParseMiddleware = middleware
	GetNamespace         = getNamespace
)
