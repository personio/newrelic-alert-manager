package widget

import (
	"fmt"
	"sort"
	"strings"
)

type ApmMetric struct {
	Duration  int        `json:"duration"`
	EntityIds []int      `json:"entity_ids"`
	Metrics   MetricList `json:"metrics"`
	Facet     string     `json:"facet,omitempty"`
	OrderBy   string     `json:"order_by,omitempty"`
}

func (m ApmMetric) Equals(other *ApmMetric) bool {
	if other == nil {
		return false
	}

	return m.equalEntityIds(other) &&
		m.Metrics.Equals(other.Metrics) &&
		m.Duration == other.Duration &&
		m.Facet == other.Facet &&
		m.OrderBy == other.OrderBy
}

func (m ApmMetric) equalEntityIds(other *ApmMetric) bool {
	if len(m.EntityIds) != len(other.EntityIds) {
		return false
	}

	sort.Ints(m.EntityIds)
	sort.Ints(other.EntityIds)
	for i, _ := range m.EntityIds {
		if m.EntityIds[i] != other.EntityIds[i] {
			return false
		}
	}

	return true
}

type MetricList []Metric

func (list MetricList) Equals(other MetricList) bool {
	if len(list) != len(other) {
		return false
	}

	sort.Slice(list, list.comparer)
	sort.Slice(other, other.comparer)

	for i, _ := range list {
		if !list[i].Equals(other[i]) {
			return false
		}
	}

	return true
}

func (list MetricList) comparer(i int, j int) bool {
	return list[i].getSortKey() < list[j].getSortKey()
}

type Metric struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func (m Metric) Equals(other Metric) bool {
	if len(m.Values) != len(other.Values) {
		return false
	}

	sort.Strings(m.Values)
	sort.Strings(other.Values)
	for i, _ := range m.Values {
		if m.Values[i] != other.Values[i] {
			return false
		}
	}

	return m.Name == other.Name
}

func (m Metric) getSortKey() string {
	return fmt.Sprintf("%s-%s", m.Name, strings.Join(m.Values, ";"))
}
