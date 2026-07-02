package infrastructure

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/smart-alert/internal/engine"
)

type PGRuleLoader struct {
	pool *pgxpool.Pool
}

func NewPGRuleLoader(pool *pgxpool.Pool) *PGRuleLoader {
	return &PGRuleLoader{pool: pool}
}

func (l *PGRuleLoader) LoadAlertRules(ctx context.Context, userID string) ([]engine.Rule, error) {
	rows, err := l.pool.Query(ctx,
		`SELECT id, name, conditions::text, actions::text FROM rules WHERE user_id=$1 AND enabled=true ORDER BY priority`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []engine.Rule
	for rows.Next() {
		var id, name, condsJSON, actsJSON string
		if err := rows.Scan(&id, &name, &condsJSON, &actsJSON); err != nil {
			return nil, err
		}
		rule, err := mapRule(id, name, condsJSON, actsJSON)
		if err != nil {
			continue
		}
		rules = append(rules, *rule)
	}
	return rules, nil
}

func mapRule(id, name, condsJSON, actsJSON string) (*engine.Rule, error) {
	var conds []struct {
		Field    string `json:"field"`
		Operator string `json:"operator"`
		Value    string `json:"value"`
	}
	if err := json.Unmarshal([]byte(condsJSON), &conds); err != nil {
		return nil, err
	}
	var acts []struct {
		Type   string                 `json:"type"`
		Params map[string]interface{} `json:"params"`
	}
	if err := json.Unmarshal([]byte(actsJSON), &acts); err != nil {
		return nil, err
	}

	ec := make([]engine.Condition, len(conds))
	for i, c := range conds {
		ec[i] = engine.Condition{Field: c.Field, Operator: engine.Operator(c.Operator), Value: c.Value}
	}
	ea := make([]engine.Action, len(acts))
	for i, a := range acts {
		ea[i] = engine.Action{Type: engine.ActionType(a.Type), Params: a.Params}
	}
	return &engine.Rule{ID: id, Name: name, Enabled: true, Conditions: ec, Actions: ea}, nil
}
