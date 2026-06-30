package engine

import "sort"

// Engine composes the notification center from domain events.
type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

// BuildCenter composes the notification center view.
func (e *Engine) BuildCenter(inputs Inputs) *NotificationCenter {
	prefs := e.buildPreferencesMap(inputs.Preferences)

	var items []NotificationCard
	p1, p2, p3, p4 := 0, 0, 0, 0

	for _, n := range inputs.Notifications {
		pref, ok := prefs[Category(n.Category)]
		if ok && !pref.Enabled {
			continue
		}
		if ok && e.priorityWeight(Priority(n.Priority)) < e.priorityWeight(Priority(pref.MinPriority)) {
			continue
		}

		card := NotificationCard{
			NotifID:     n.NotifID,
			Category:    Category(n.Category),
			Priority:    Priority(n.Priority),
			Title:       n.Title,
			Summary:     n.Summary,
			ActionLabel: n.ActionLabel,
			ActionRoute: n.ActionRoute,
			State:       NotificationState(n.State),
			GeneratedAt: n.GeneratedAt,
			ExpiresAt:   n.ExpiresAt,
		}
		items = append(items, card)

		switch n.Priority {
		case "P1": p1++ ; case "P2": p2++ ; case "P3": p3++ ; default: p4++
		}
	}

	// Sort by priority (P1 first), then by generated time descending
	sort.Slice(items, func(i, j int) bool {
		wi := e.priorityWeight(items[i].Priority)
		wj := e.priorityWeight(items[j].Priority)
		if wi != wj { return wi < wj }
		return items[i].GeneratedAt > items[j].GeneratedAt
	})

	return &NotificationCenter{
		UnreadCount: len(items),
		P1Count:     p1, P2Count: p2, P3Count: p3, P4Count: p4,
		Items: items,
	}
}

func (e *Engine) buildPreferencesMap(inputs []PreferenceInput) map[Category]Preference {
	m := make(map[Category]Preference)
	for _, p := range inputs {
		m[Category(p.Category)] = Preference{
			Category: Category(p.Category), InApp: p.InApp, Push: p.Push,
			Email: p.Email, Digest: p.Digest, MinPriority: p.MinPriority, Enabled: p.Enabled,
		}
	}
	return m
}

func (e *Engine) priorityWeight(p Priority) int {
	switch p { case P1Critical: return 1; case P2Urgent: return 2; case P3Important: return 3; default: return 4 }
}
