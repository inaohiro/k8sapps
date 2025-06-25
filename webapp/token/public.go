package token

import "k8soperation/token/internal/controller"

var (
	Routes               = controller.Controller()
	TokenParseMiddleware = controller.Middleware
	GetNamespace         = controller.GetNamespace
)
