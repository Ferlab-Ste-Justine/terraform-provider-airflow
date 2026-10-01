package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Logger struct {
	ctx context.Context
}

func NewLogger(ctx context.Context) *Logger {
	return &Logger{ctx: ctx}
}

func (log *Logger) Debug(ctx context.Context, msg string) {
	tflog.Debug(ctx, msg)
}

func (log *Logger) Info(ctx context.Context, msg string) {
	tflog.Info(ctx, msg)
}

func (l *Logger) Warn(ctx context.Context, msg string) {
	tflog.Warn(ctx, msg)
}

func (l *Logger) Error(ctx context.Context, msg string) {
	tflog.Error(ctx, msg)
}
