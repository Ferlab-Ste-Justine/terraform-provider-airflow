package airflow

import (
	"context"
)

type Logger interface {
	Debug(context.Context, string)
	Info(context.Context, string)
	Warn(context.Context, string)
	Error(context.Context, string)
}