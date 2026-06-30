package engine

import (
	"sort"
	"strings"
)

// Engine composes timeline view models from event data.
type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

// BuildFeed composes the paginated, filtered, grouped timeline feed.
func (e *Engine) BuildFeed(inputs Inputs) *TimelineOutput {
	// Apply filters
	filtered := e.applyFilters(inputs.Events, inputs.Filters)

	// Sort by date descending (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].EventDate > filtered[j].EventDate
	})

	// Paginate
	limit := inputs.Limit
	if limit <= 0 { limit = 50 }
	hasMore := len(filtered) > limit
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}

	// Build timeline events
	events := e.toTimelineEvents(filtered)

	// Compute daily summaries
	summaries := e.computeDailySummaries(filtered)

	return &TimelineOutput{
		View:          inputs.View,
		Feed:          TimelineFeed{Events: events, Cursor: inputs.Cursor, HasMore: hasMore, Total: len(filtered)},
		DailySummaries: summaries,
		Filters:       inputs.Filters,
		TotalEvents:   len(filtered),
	}
}

func (e *Engine) applyFilters(events []EventInput, filter FilterState) []EventInput {
	var result []EventInput
	for _, ev := range events {
		// Category filter
		if len(filter.Categories) > 0 {
			found := false
			for _, c := range filter.Categories {
				if ev.SourceType == c { found = true; break }
			}
			if !found { continue }
		}

		// Entity filter
		if filter.EntityID != "" && ev.SourceEntityID != filter.EntityID {
			continue
		}

		// Search text filter
		if filter.SearchText != "" {
			search := strings.ToLower(filter.SearchText)
			if !strings.Contains(strings.ToLower(ev.Title), search) &&
				!strings.Contains(strings.ToLower(ev.Description), search) &&
				!strings.Contains(strings.ToLower(ev.EventType), search) {
				continue
			}
		}

		result = append(result, ev)
	}
	return result
}

func (e *Engine) toTimelineEvents(inputs []EventInput) []TimelineEvent {
	events := make([]TimelineEvent, len(inputs))
	for i, ev := range inputs {
		events[i] = TimelineEvent{
			EventID:         ev.EventID,
			SourceType:      ev.SourceType,
			EventType:       ev.EventType,
			Title:           ev.Title,
			Description:     ev.Description,
			Amount:          ev.Amount,
			EventDate:       ev.EventDate,
			IsCelebration:   ev.IsCelebration,
			IsWarning:       ev.IsWarning,
			IsHistoricalOnly: ev.IsHistorical,
			SourceEntityID:  ev.SourceEntityID,
			SourceRoute:     ev.SourceRoute,
		}
	}
	return events
}

func (e *Engine) computeDailySummaries(events []EventInput) []DailySummary {
	dayMap := map[string]*DailySummary{}

	for _, ev := range events {
		date := ev.EventDate[:10] // YYYY-MM-DD
		if _, ok := dayMap[date]; !ok {
			dayMap[date] = &DailySummary{Date: date}
		}
		dayMap[date].EventCount++
		dayMap[date].NetFinancialImpact += ev.Amount
		if ev.Title != "" {
			dayMap[date].MostSignificant = ev.Title
		}
	}

	var summaries []DailySummary
	for _, s := range dayMap {
		summaries = append(summaries, *s)
	}
	sort.Slice(summaries, func(i, j int) bool { return summaries[i].Date > summaries[j].Date })
	if len(summaries) > 7 {
		summaries = summaries[:7]
	}
	return summaries
}
