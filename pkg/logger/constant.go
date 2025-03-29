package logger

type ctxKeyLocalization int

const (
	EnvironmentProduction                    = "production"
	EnvironmentDevelop                       = "develop"
	CtxKeyLocalization    ctxKeyLocalization = 0
)
