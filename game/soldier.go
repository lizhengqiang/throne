package game

import (
	"fmt"
	"throne/game/soldier"
)

type Soldier struct {
	Type  int64
	Alive bool
}

func InitSoldier(typ int64) *Soldier {
	return &Soldier{
		Type:  typ,
		Alive: true,
	}
}

func CalcSoldiers(soldiers []*Soldier) (r int64) {
	for _, s := range soldiers {
		switch s.Type {
		case soldier.Ship:
			r = r + 1
		case soldier.Foot:
			r = r + 1
		case soldier.Cavalry:
			r = r + 2
		}
	}
	return
}

func RemoveSoldiers(slice []*Soldier, elems ...*Soldier) (r []*Soldier) {
	r = make([]*Soldier, len(slice))
	copy(r, slice)
	fmt.Println(slice, elems, r)
	for _, e := range elems {
		for i, s := range r {

			if e == s {

				// 找到了这个元素
				r = append(r[:i], r[i+1:]...)
				break
			}
		}
	}
	return
}

func (s *Soldier) String() string {
	alive := ""
	if s.Alive {
		alive = "正常"
	} else {
		alive = "战败"
	}
	return fmt.Sprintf("(类型%d,%s)", s.Type, alive)
}
