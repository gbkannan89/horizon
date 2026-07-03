package engine

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// evaluateCondition checks if a single condition holds for the payload
func evaluateCondition(cond Condition, payload map[string]interface{}) (bool, error) {
	val, ok := payload[cond.Field]
	if !ok {
		// If field is missing, condition fails
		return false, nil
	}

	// Cast payload value to string for easier comparison
	strVal := fmt.Sprintf("%v", val)

	switch cond.Operator {
	case OpEquals:
		return strings.EqualFold(strVal, cond.Value), nil
	case OpContains:
		return strings.Contains(strings.ToLower(strVal), strings.ToLower(cond.Value)), nil
	case OpGreaterThan:
		payloadFloat, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return false, nil
		}
		condFloat, err := strconv.ParseFloat(cond.Value, 64)
		if err != nil {
			return false, nil
		}
		return payloadFloat > condFloat, nil
	case OpLessThan:
		payloadFloat, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return false, nil
		}
		condFloat, err := strconv.ParseFloat(cond.Value, 64)
		if err != nil {
			return false, nil
		}
		return payloadFloat < condFloat, nil
	case OpMatchesRegex:
		matched, err := regexp.MatchString(cond.Value, strVal)
		if err != nil {
			return false, err
		}
		return matched, nil
	case OpBetween:
		parts := strings.Split(cond.Value, ",")
		if len(parts) != 2 {
			return false, nil
		}
		min, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		max, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		payloadFloat, err3 := strconv.ParseFloat(strVal, 64)
		if err1 != nil || err2 != nil || err3 != nil {
			return false, nil
		}
		return payloadFloat >= min && payloadFloat <= max, nil
	}

	return false, fmt.Errorf("unknown operator: %s", cond.Operator)
}
