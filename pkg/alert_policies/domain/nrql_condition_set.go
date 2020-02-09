package domain

type NrqlConditionSet struct {
	conditions map[string]NrqlConditionBody
}

func NewNrqlConditionSet(conditions NrqlConditionList) *NrqlConditionSet {
	set := &NrqlConditionSet{
		conditions: make(map[string]NrqlConditionBody),
	}
	for _, condition := range conditions.Condition {
		set.put(condition)
	}

	return set
}

func NewNrqlConditionSetFromSlice(conditions []*NrqlCondition) *NrqlConditionSet {
	set := &NrqlConditionSet{
		conditions: make(map[string]NrqlConditionBody),
	}
	for _, condition := range conditions {
		set.put(condition.Condition)
	}

	return set
}

func (set NrqlConditionSet) put(condition NrqlConditionBody) {
	set.conditions[condition.getHashKey()] = condition
}

func (set NrqlConditionSet) Contains(condition NrqlConditionBody) bool {
	key := condition.getHashKey()

	_, ok := set.conditions[key]
	return ok
}