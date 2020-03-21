package widget

import "sort"

type WidgetList []Widget

func (list WidgetList) Equals (other WidgetList) bool {
	if len(list) != len(other) {
		return false
	}

	sort.Sort(list)
	sort.Sort(other)

	for i, _ := range list {
		if !list[i].Equals(other[i]) {
			return false
		}
	}

	return true
}

func (list WidgetList) Len() int {
	return len(list)
}

func (list WidgetList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list WidgetList) Less(i int, j int) bool {
	return list[i].getHashKey() < list[j].getHashKey()
}