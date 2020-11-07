package widget

import (
	"sort"
)

type WidgetList []Widget

func (list WidgetList) Equals(other WidgetList) bool {
	if len(list) != len(other) {
		return false
	}

	listCopy := list.copy()
	otherCopy := other.copy()

	sort.Slice(listCopy, listCopy.comparer)
	sort.Slice(otherCopy, otherCopy.comparer)

	for i, _ := range list {
		if !listCopy[i].Equals(otherCopy[i]) {
			return false
		}
	}

	return true
}

func (list WidgetList) comparer(i int, j int) bool {
	return list[i].getComparisonKey() < list[j].getComparisonKey()
}

func (list WidgetList) copy() WidgetList {
	listCopy := make([]Widget, len(list))
	copy(listCopy, list)

	return listCopy
}
