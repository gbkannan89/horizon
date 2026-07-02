package engine

import (
	"context"
	"time"
)

type ReportGenerator struct {
	provider DataProvider
}

func NewReportGenerator(provider DataProvider) *ReportGenerator {
	return &ReportGenerator{provider: provider}
}

func (g *ReportGenerator) Generate(ctx context.Context, schedule ReportSchedule) (*GeneratedReport, error) {
	data := make(map[string]interface{})

	for _, section := range schedule.Sections {
		sectionData, err := g.provider.FetchSectionData(ctx, schedule.UserID, section.Type)
		if err != nil {
			sectionData = map[string]interface{}{"error": err.Error()}
		}
		data[section.Type] = sectionData
	}

	now := time.Now().UTC().Format(time.RFC3339)
	return &GeneratedReport{
		ScheduleID:  schedule.ID,
		UserID:      schedule.UserID,
		Name:        schedule.Name,
		Frequency:   schedule.Frequency,
		Sections:    schedule.Sections,
		Data:        data,
		GeneratedAt: now,
	}, nil
}

func (g *ReportGenerator) DetermineSectionTypes() []ReportSection {
	return []ReportSection{
		{Title: "Net Worth Summary", Type: "net_worth"},
		{Title: "Income vs Expenses", Type: "income_expense"},
		{Title: "Budget Performance", Type: "budget"},
		{Title: "Goal Progress", Type: "goals"},
	}
}

func (g *ReportGenerator) DefaultSections() []ReportSection {
	return g.DetermineSectionTypes()
}
