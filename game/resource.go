package game

import "fmt"

type Resource struct {
	Type int64
}

func InitResource(typ int64) *Resource {
	return &Resource{
		Type: typ,
	}
}

func (r *Resource) String() string {
	return fmt.Sprintf("(类型:%d)", r.Type)
}
