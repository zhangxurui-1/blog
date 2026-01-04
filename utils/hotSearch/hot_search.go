package hotSearch

import "server/model/other"

type Source interface {
	GetHotSearchData(maxNum int) (hotSearchData other.HotSearchData, err error)
}

func NewSource(sourceStr string) Source {
	switch sourceStr {
	case "baidu":
		return &Baidu{}
	case "kuaishou":
		return &Kuaishou{}
	case "toutiao":
		return &Toutiao{}
	case "zhihu":
		return &Zhihu{}
	default:
		return nil
	}
}
