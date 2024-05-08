package wfinstance

import (
	"github.com/remiges-tech/crux/db/sqlc-gen"
)

// To Get ParentList from WFInstanceList data respone
func getParentList(data []sqlc.Wfinstance) []int32 {
	var parentsMap = make(map[int32]struct{})
	var parents []int32

	for _, val := range data {
		if val.Parent.Valid {
			parentValue := val.Parent.Int32
			_, exists := parentsMap[parentValue]
			if !exists {
				parents = append(parents, parentValue)
				parentsMap[parentValue] = struct{}{}
			}
		}
	}
	return parents

}
