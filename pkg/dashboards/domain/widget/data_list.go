package widget

type DataList [1]Data

func (list DataList) getSortKey() string {
	// The list contains 1 item so no need to come up with a sort key
	return ""
}

func (list DataList) Equals(other DataList) bool {
	return list[0].Equals(other[0])
}
