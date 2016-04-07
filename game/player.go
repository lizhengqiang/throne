package game

import (
	"fmt"
	"qiniupkg.com/x/errors.v7"
	"throne/game/player"
	"throne/game/stage"
)

type Player struct {
	Game   *Game
	Name   string
	Money  int64
	Orders map[player.OrderType]int64
	Wars   map[*War]interface{}
}

func InitPlayer(game *Game, name string) *Player {
	return &Player{
		Game:   game,
		Name:   name,
		Money:  5,
		Wars:   map[*War]interface{}{},
		Orders: map[player.OrderType]int64{},
	}
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

func (p *Player) Move(src *Area, dst *Area, soldiers []*Soldier, stayControl bool) error {
	if src.Belong != p {
		return InvalidSrcArea
	}

	if e := src.CanMove(dst); e != nil {
		return e
	}

	if src.Type != dst.Type {
		return InvalidDstArea
	}
	// 都是自己的土地,或者进入空地
	if src.Belong == p && (dst.Belong == p || dst.Belong == nil) {
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
