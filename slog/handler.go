package hbslog

import (
	"context"
	"log/slog"

	"github.com/honeybadger-io/honeybadger-go"
)

type preformattedAttr struct {
	groups []string
	attr   slog.Attr
}

type Handler struct {
	c         *honeybadger.Client
	eventType string
	preformat []preformattedAttr
	groups    []string
	level     slog.Leveler
}

func New(c *honeybadger.Client) *Handler {
	return &Handler{
		c:         c,
		eventType: "log",
		level:     slog.LevelInfo,
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	data := map[string]any{
		"level":   r.Level.String(),
		"message": r.Message,
	}

	for _, pf := range h.preformat {
		h.appendAttr(data, pf.groups, pf.attr)
	}

	eventType := h.eventType
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "event_type" {
			if et, ok := a.Value.Any().(string); ok {
				eventType = et
				return true
			}
		}
		h.appendAttr(data, h.groups, a)
		return true
	})

	return h.c.Event(eventType, data)
}

func (h *Handler) appendAttr(data map[string]any, groups []string, attr slog.Attr) {
	value := attr.Value.Resolve()

	if value.Kind() == slog.KindGroup {
		groupAttrs := value.Group()
		if len(groupAttrs) == 0 {
			return
		}

		nestedGroups := append(groups, attr.Key)
		for _, groupAttr := range groupAttrs {
			h.appendAttr(data, nestedGroups, groupAttr)
		}
		return
	}

	if len(groups) == 0 {
		data[attr.Key] = value.Any()
		return
	}

	current := data
	for _, group := range groups {
		if _, ok := current[group]; !ok {
			current[group] = make(map[string]any)
		}
		current = current[group].(map[string]any)
	}
	current[attr.Key] = value.Any()
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newPreformat := make([]preformattedAttr, len(h.preformat), len(h.preformat)+len(attrs))
	copy(newPreformat, h.preformat)

	for _, attr := range attrs {
		groupsCopy := make([]string, len(h.groups))
		copy(groupsCopy, h.groups)
		newPreformat = append(newPreformat, preformattedAttr{
			groups: groupsCopy,
			attr:   attr,
		})
	}

	return &Handler{
		c:         h.c,
		eventType: h.eventType,
		preformat: newPreformat,
		groups:    h.groups,
		level:     h.level,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newGroups := make([]string, 0, len(h.groups)+1)
	newGroups = append(newGroups, h.groups...)
	newGroups = append(newGroups, name)
	return &Handler{
		c:         h.c,
		eventType: h.eventType,
		preformat: h.preformat,
		groups:    newGroups,
		level:     h.level,
	}
}

func (h *Handler) WithEventType(eventType string) *Handler {
	if eventType == "" {
		eventType = "log"
	}
	return &Handler{
		c:         h.c,
		eventType: eventType,
		preformat: h.preformat,
		groups:    h.groups,
		level:     h.level,
	}
}

func (h *Handler) WithLevel(level slog.Leveler) *Handler {
	return &Handler{
		c:         h.c,
		eventType: h.eventType,
		preformat: h.preformat,
		groups:    h.groups,
		level:     level,
	}
}
