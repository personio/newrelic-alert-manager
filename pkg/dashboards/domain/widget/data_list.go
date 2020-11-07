package widget

type DataList [1]Data

func (list DataList) Equals(other DataList) bool {
	return list[0].Equals(other[0])
}
