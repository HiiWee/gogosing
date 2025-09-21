package util

import (
	"log/slog"
	"os"
)

func MustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		slog.Warn("missing env %s", k)
	}
	return v
}
