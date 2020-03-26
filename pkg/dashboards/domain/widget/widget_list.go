package widget

import "sort"

type WidgetList []Widget

func (list WidgetList) Equals (other WidgetList) bool {
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

func (list WidgetList) comparer (i int, j int) bool {
	return list[i].getSortKey() < list[j].getSortKey()
}
