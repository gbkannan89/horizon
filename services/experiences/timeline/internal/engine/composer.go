package engine

import (
	"sort"
	"strings"
)

// Composer builds timeline view models from raw event data.
type Composer struct{}

func NewComposer() *Composer { return &Composer{} }

// BuildFeed composes the paginated, filtered timeline feed.
func (c *Composer) BuildFeed(inputs Inputs) *TimelineOutput {
	filtered := c.applyFilters(inputs.Items, inputs.Filters)

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp > filtered[j].Timestamp
	})

	limit := inputs.Limit
	if limit <= 0 { limit = 50 }
	hasMore := len(filtered) > limit
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}

	items := c.toItems(filtered)

	nextCursor := ""
	if hasMore && len(items) > 0 {
		nextCursor = items[len(items)-1].TimelineID
	}

	return &TimelineOutput{
		View: inputs.View,
		Feed: TimelineFeed{Items: items, Cursor: nextCursor, HasMore: hasMore, Total: len(filtered)},
		Filters: inputs.Filters,
		Total:  len(filtered),
	}
}

// GetByID returns a single timeline item by its ID.
func (c *Composer) GetByID(all []RawItem, id string) *TimelineItem {
	for _, r := range all {
		if r.TimelineID == id {
			item := c.toItem(r)
			return &item
		}
	}
	return nil
}

// GetRecent returns the N most recent items.
func (c *Composer) GetRecent(all []RawItem, n int) []TimelineItem {
	if n <= 0 { n = 10 }
	sorted := make([]RawItem, len(all))
	copy(sorted, all)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp > sorted[j].Timestamp
	})
	if len(sorted) > n { sorted = sorted[:n] }
	return c.toItems(sorted)
}

func (c *Composer) Search(all []RawItem, q string) []TimelineItem {
	q = strings.ToLower(q)
	var matched []RawItem
	for _, r := range all {
		if strings.Contains(strings.ToLower(r.Title), q) ||
			strings.Contains(strings.ToLower(r.Summary), q) ||
			strings.Contains(strings.ToLower(r.Description), q) ||
			strings.Contains(strings.ToLower(r.EventType), q) ||
			strings.Contains(r.TimelineID, q) {
			matched = append(matched, r)
		}
	}
	return c.toItems(matched)
}

func (c *Composer) applyFilters(items []RawItem, filter FilterState) []RawItem {
	var result []RawItem
	for _, it := range items {
		if len(filter.Categories) > 0 {
			found := false
			for _, cat := range filter.Categories {
				if it.Category == cat { found = true; break }
			}
			if !found { continue }
		}

		if filter.StartDate != "" && it.Timestamp < filter.StartDate { continue }
		if filter.EndDate != "" {
			end := filter.EndDate
			if len(end) == 10 { end += "T23:59:59Z" }
			if it.Timestamp > end { continue }
		}

		if filter.Severity != "" && it.Severity != filter.Severity { continue }
		if filter.Goal != "" && it.RelatedGoal != filter.Goal { continue }
		if filter.Account != "" && it.RelatedAcct != filter.Account { continue }
		if filter.Asset != "" && it.RelatedAsset != filter.Asset { continue }
		if filter.Liability != "" && it.RelatedEntity != filter.Liability { continue }
		if filter.Portfolio != "" && it.RelatedAgg != filter.Portfolio { continue }
		if filter.EntityID != "" && it.RelatedEntity != filter.EntityID {
			if it.RelatedAgg != filter.EntityID { continue }
		}

		if filter.SearchText != "" {
			q := strings.ToLower(filter.SearchText)
			if !strings.Contains(strings.ToLower(it.Title), q) &&
				!strings.Contains(strings.ToLower(it.Summary), q) &&
				!strings.Contains(strings.ToLower(it.Description), q) &&
				!strings.Contains(strings.ToLower(it.EventType), q) {
				continue
			}
		}

		result = append(result, it)
	}
	return result
}

func (c *Composer) toItems(raw []RawItem) []TimelineItem {
	items := make([]TimelineItem, len(raw))
	for i, r := range raw {
		items[i] = c.toItem(r)
	}
	return items
}

func (c *Composer) toItem(r RawItem) TimelineItem {
	icon, color := c.resolveIcon(r.Category, r.EventType, r.Severity)
	return TimelineItem{
		TimelineID:    r.TimelineID,
		Timestamp:     r.Timestamp,
		EventType:     r.EventType,
		Category:      r.Category,
		Title:         r.Title,
		Summary:       r.Summary,
		Description:   r.Description,
		Severity:      r.Severity,
		RelatedEntity: r.RelatedEntity,
		RelatedAgg:    r.RelatedAgg,
		RelatedGoal:   r.RelatedGoal,
		RelatedAcct:   r.RelatedAcct,
		RelatedAsset:  r.RelatedAsset,
		Amount:        r.Amount,
		Metadata:      r.Metadata,
		Icon:          icon,
		Color:         color,
	}
}

func (c *Composer) resolveIcon(cat EventCategory, evtType string, sev Severity) (string, string) {
	switch cat {
	case ECFinancial:
		if sev == SevWarning || sev == SevCritical { return "warning", "amber" }
		return "currency_rupee", "green"
	case ECGoal:
		if sev == SevSuccess { return "celebration", "green" }
		if sev == SevWarning { return "warning", "amber" }
		return "track_changes", "blue"
	case ECAccount:
		return "account_balance", "indigo"
	case ECAsset:
		return "trending_up", "teal"
	case ECLiability:
		if sev == SevCritical { return "error", "red" }
		return "credit_card", "purple"
	case ECPortfolio:
		return "pie_chart", "blue"
	case ECHealth:
		if sev == SevWarning { return "heart_broken", "red" }
		return "favorite", "pink"
	case ECRisk:
		if sev == SevCritical { return "gpp_bad", "red" }
		return "shield", "amber"
	case ECRecommend:
		return "lightbulb", "amber"
	case ECSimulation:
		return "science", "cyan"
	case ECOptimization:
		return "auto_graph", "cyan"
	case ECAchievement:
		return "emoji_events", "gold"
	case ECMilestone:
		return "flag", "green"
	case ECUser:
		return "person", "grey"
	default:
		return "circle", "grey"
	}
}
