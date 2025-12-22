package slogconsole

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"
)

type Handler struct {
	out       io.Writer
	level     slog.Level
	addSource bool
}

type Option func(*Handler)

func New(out io.Writer, opts ...Option) *Handler {
	h := &Handler{
		out:   out,
		level: slog.LevelInfo,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func WithLevel(level slog.Level) Option {
	return func(h *Handler) {
		h.level = level
	}
}

func WithSource(flag bool) Option {
	return func(h *Handler) {
		h.addSource = flag
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder

	// time
	b.WriteString(r.Time.Format(time.RFC3339))
	b.WriteByte(' ')

	// level
	b.WriteString(colorizeLevel(r.Level))
	b.WriteByte(' ')

	// message
	b.WriteString(r.Message)

	// attrs
	r.Attrs(func(a slog.Attr) bool {
		b.WriteByte(' ')
		b.WriteString(a.Key)
		b.WriteByte('=')
		b.WriteString(fmt.Sprint(a.Value.Any()))
		return true
	})

	// source
	if h.addSource && r.Source() != nil {
		b.WriteString(" (")
		b.WriteString(r.Source().File)
		b.WriteByte(':')
		b.WriteString(fmt.Sprint(r.Source().Line))
		b.WriteByte(')')
	}

	b.WriteByte('\n')

	_, err := h.out.Write([]byte(b.String()))
	return err
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h
}

func colorizeLevel(level slog.Level) string {
	var color string

	switch {
	case level == slog.LevelDebug:
		color = colorGray
	case level == slog.LevelInfo:
		color = colorBlue
	case level == slog.LevelWarn:
		color = colorYellow
	case level >= slog.LevelError:
		color = colorRed
	default:
		color = colorReset
	}

	return fmt.Sprintf("%s%-5s%s", color, level.String(), colorReset)
}
