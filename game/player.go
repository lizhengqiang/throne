package game

import (
	"encoding/json"
	"fmt"
	"qiniupkg.com/x/errors.v7"
	"throne/game/player"
	"throne/game/stage"
	"throne/game/area"
)

type Player struct {
	Game    *Game `json:"-"`
	CanSet  bool
	Name    string
	Money   int64
	Orders  map[player.OrderType]int64 `json:"-"`
	PreArea *Area`json:"-"`
}

func Players(m map[*Player]interface{}) []interface{} {
	r := []interface{}{}
	for d, _ := range m {
		d1, e := d.Marshal()
		if e != nil {
			continue
		}
		r = append(r, d1)
	}
	return r
}

func (a *Player) Marshal() (m map[string]interface{}, err error) {
	m = map[string]interface{}{}
	bytes, err := json.Marshal(a)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return
	}

	for t, o := range a.Orders {
		m[player.OrderKey(t)] = o
	}

	return
}

func InitPlayer(game *Game, name string) *Player {
	p := &Player{
		Game:   game,
		Name:   name,
		CanSet: false,
		Money:  5,
		Orders: map[player.OrderType]int64{},

	}
	p.PreArea = InitArea(game, "预备", area.Land, p, []*Resource{})
	return p
}

func (p *Player) String() string {
	return fmt.Sprintf("[%s,金钱:%d]", p.Name, p.Money)
}

func (p *Player) SetOrder(t player.OrderType, order int64) {
	p.Orders[t] = order
}

func (p *Player) Notify() {
	if p.Game.Order == p.Orders[player.OrderA] && p.Game.Stage == stage.Move {
		fmt.Printf("%s:可以执行一个移动指令!\n", p.Name)
	}
}

var InvalidSrcArea error = errors.New("无效起始地区")
var InvalidDstArea error = errors.New("无效目标地区")

func (p *Player) Set(area *Area, cmd *Command) error {
	if p.CanSet == false {
		return CannotSet
	}

	area.PutCommand(cmd)
	return nil
}

var CannotSet error = errors.New("还不能放置")

func (p *Player) FinishSet() error {
	if p.CanSet == false {
		return CannotSet
	}

	// 完成一个角色
	p.CanSet = false
	p.Game.WaitGroup.Done()
	return nil
}

func (p *Player) Move(src *Area, dst *Area, soldiers []*Soldier, stayControl bool) error {
	if src.Belong != p {
		return InvalidSrcArea
	}

	if !src.CanMove(dst) {
		return CannotMove
	}

	if src.Type != dst.Type {
		return InvalidDstArea
	}
	// 都是自己的土地,或者进入空地
	if dst.Belong == p || dst.Belong == nil {
		src.Leave(p, soldiers, stayControl)
		dst.Enter(p, soldiers)
		return nil
	}

	if len(dst.Soldiers) == 0 {
		src.Leave(p, soldiers, stayControl)
		dst.Enter(p, soldiers)
		return nil
	}
	// 是别人的地,发起战争
	src.Leave(p, soldiers, stayControl)
	w := InitWar(p, src, dst, soldiers)
	go w.Handle()

	return nil
}
