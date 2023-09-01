package pkg

import (
	"encoding/json"

	dtrack "github.com/DependencyTrack/client-go"
)

type PolicyCondition struct {
	Operator dtrack.PolicyConditionOperator `json:"operator"`
	Subject  dtrack.PolicyConditionSubject  `json:"subject"`
	Value    string                         `json:"value"`
}

type PolicyConditions []PolicyCondition

func NewPolicyConditions(b []byte) (pp PolicyConditions, err error) {
	if err := json.Unmarshal(b, &pp); err != nil {
		return nil, err
	}
	return pp, nil
}
