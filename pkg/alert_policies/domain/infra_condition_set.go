package domain

type InfraConditionSet struct {
	conditions map[string]InfraConditionBody
}

func NewInfraConditionSet(conditions InfraConditionList) *InfraConditionSet {
	set := &InfraConditionSet{
		conditions: make(map[string]InfraConditionBody),
	}
	for _, condition := range conditions.Condition {
		set.put(condition)
	}

	return set
}

func NewInfraConditionSetFromSlice(conditions []*InfraCondition) *InfraConditionSet {
	set := &InfraConditionSet{
		conditions: make(map[string]InfraConditionBody),
	}
	for _, condition := range conditions {
		set.put(condition.Condition)
	}

	return set
}

func (set InfraConditionSet) put(condition InfraConditionBody) {
	set.conditions[condition.getHashKey()] = condition
}

func (set InfraConditionSet) Contains(condition InfraConditionBody) bool {
	key := condition.getHashKey()

	_, ok := set.conditions[key]
	return ok
}
