package domain

type ApmConditionSet struct {
	conditions map[string]ApmConditionBody
}

func NewApmConditionSet(conditions ApmConditionList) *ApmConditionSet {
	set := &ApmConditionSet{
		conditions: make(map[string]ApmConditionBody),
	}
	for _, condition := range conditions.Condition {
		set.put(condition)
	}

	return set
}

func NewApmConditionSetFromSlice(conditions []*ApmCondition) *ApmConditionSet {
	set := &ApmConditionSet{
		conditions: make(map[string]ApmConditionBody),
	}
	for _, condition := range conditions {
		set.put(condition.Condition)
	}

	return set
}

func (set ApmConditionSet) put(condition ApmConditionBody) {
	set.conditions[condition.getHashKey()] = condition
}

func (set ApmConditionSet) Contains(condition ApmConditionBody) bool {
	key := condition.getHashKey()

	_, ok := set.conditions[key]
	return ok
}
