package game

import (
	"errors"
	"fmt"
	"qiniupkg.com/x/log.v7"
	"throne/game/area"
	"throne/game/command"
	"throne/game/player"
	"throne/game/resource"
	"throne/game/stage"
	"throne/utils"
)

type Area struct {
	Game          *Game
	Id            string
	Name          string
	Type          int64                 // 地面类型
	Around        map[*Area]interface{} // 临近地区
	Resources     []*Resource           // 资源
	Soldiers      []*Soldier            // 士兵单位
	Command       *Command              // 指令
	Belong        *Player               // 指示物
	Home          *Player

	MoveWaitGroup *utils.WaitGroup
}

func (a *Area) String() string {
	return fmt.Sprintf("%s,%s,所属%v,指令%v,", a.Id, a.Name, a.Belong, a.Command) + fmt.Sprintf("资源%v,", a.Resources) + fmt.Sprintf("士兵%v;", a.Soldiers)
}

var CannotPutCommand error = errors.New("不能放置指令")

func (a *Area) PutCommand(command *Command) error {

	if a.Game.Stage != stage.Set {
		return CannotPutCommand
	}

	if a.Belong != command.Belong {
		return CannotPutCommand
	}

	if len(a.Soldiers) == 0 {
		return CannotPutCommand
	}

	a.Command = command
	return nil

}

var NoCommand error = errors.New("没有指令")
var NotTurn error = errors.New("还没轮到")

// 结算指令
func (a *Area) ConsumeCommand(typ command.Type) error {

	if a.Command == nil {
		return NoCommand
	}
	// 还没有到执行顺序
	if a.Belong.Orders[player.OrderA] != a.Game.Order || a.Game.Stage != stage.Move {
		return NotTurn
	}
	return nil
}

func (a *Area) ConsumedCommand(typ command.Type) error {
	defer a.Game.WaitGroup.Done()
	a.Command = nil
	return nil
}
func (a *Area) ConsumeSteal(areas ...*Area) error {
	if err := a.ConsumeCommand(command.Steal); err != nil {
		return err
	}
	for _, ar := range areas {
		if ar.Command.Type == command.Steal || ar.Command.Type == command.Help || ar.Command.Type == command.Money {
			ar.Command = nil
		}
	}
	return a.ConsumedCommand(command.Steal)
}

func (a *Area) ConsumeMoney() error {
	if err := a.ConsumeCommand(command.Money); err != nil {
		return err
	}
	for _, res := range a.Resources {
		if res.Type == resource.Money {
			a.Belong.Money++
		}
	}
	a.Belong.Money++
	return a.ConsumedCommand(command.Money)
}

func (a *Area) ConsumeAttack(ar *Area, soldiers []*Soldier, stayControl bool) error {
	// 等待战争处理完毕
	defer a.ConsumedCommand(command.Attack)
	defer log.Println("战争处理完毕")
	defer ar.MoveWaitGroup.Wait()
	defer log.Println("等待战争处理完毕")
	if err := a.ConsumeCommand(command.Attack); err != nil {
		return err
	}
	if err := a.Belong.Move(a, ar, soldiers, stayControl); err != nil {
		return err
	}

	return nil
}

func (a *Area) Calc() (r int64) {
	return CalcSoldiers(a.Soldiers)
}
func (a *Area) LandHelpers() (r map[*Area]interface{}) {
	r = map[*Area]interface{}{}
	for ar, _ := range a.Around {
		// 不是自己的
		if ar.Belong != a.Belong {
			continue
		}
		// 陆地支援
		if ar.Type == area.Land {
			r[ar] = nil
			continue
		}
		// 跨海支援
		for arr, _ := range ar.LandHelpers() {
			r[arr] = nil
		}

	}
	return r
}
func (a *Area) Helpers() (r map[*Area]interface{}) {
	r = map[*Area]interface{}{}
	// 海上
	if a.Type == area.Sea {
		for ar, _ := range a.Around {
			if ar.Belong == nil {
				continue
			}
			if ar.Type == area.Sea {
				r[ar] = nil
			}
		}
		return
	}
	// 陆地
	// 全部周围均可支援
	for ar, _ := range a.Around {
		if ar.Belong == nil {
			continue
		}
		r[ar] = nil
		// 陆地
		if ar.Type == area.Land {
			continue
		}
		// 不是自己的海洋
		if ar.Belong != a.Belong {
			continue
		}
		// 跨海支援
		for are, _ := range ar.LandHelpers() {
			r[are] = nil
		}
	}

	return

}

var CannotMove error = errors.New("不能移动到这个位置")

func (a *Area) canMove(target *Area, path map[*Area]interface{}) error {
	for s, _ := range a.Around {

		// 直接相邻
		if s == target {
			return nil
		}
		// 不是自己的地
		if s.Belong != a.Belong {
			continue
		}
		// 走过这个地儿
		if _, has := path[s]; has {
			continue
		}
		// 添加到路径
		path[s] = s
		// 递归
		if e := s.canMove(target, path); e == nil {
			return nil
		}
	}
	return CannotMove
}
func (a *Area) CanMove(target *Area) error {
	return a.canMove(target, map[*Area]interface{}{a: a})
}
func (a *Area) AppendAround(area *Area) {
	a.Around[area] = area
}

func (a *Area) AddAround(area *Area) *Area {
	area.AppendAround(a)
	a.AppendAround(area)
	return a
}

func (a *Area) SelectAllSoldiers() []*Soldier {
	return append(a.Soldiers)
}

func (a *Area) Leave(player *Player, soldiers []*Soldier, stayControl bool) *Area {
	a.Soldiers = RemoveSoldiers(a.Soldiers, soldiers...)
	// 主城
	if a.Home == player {
		a.Belong = player
		return a
	}
	// 有剩余兵力
	if len(a.Soldiers) > 0 {
		a.Belong = player
		return a
	}
	// 没有多余的权利标记
	if player.Money == 0 {
		a.Belong = nil
		return a
	}
	// 留下一个钱币
	if stayControl {
		player.Money--
		a.Belong = player
		return a
	}
	// 直接走了
	a.Belong = nil
	return a

}

func (a *Area) Enter(player *Player, soldiers []*Soldier) *Area {
	a.Soldiers = append(a.Soldiers, soldiers...)
	// 进入自己的主城
	if a.Home == player {
		a.Belong = player
		return a
	}
	// 进入有权利标记的地方
	if a.Belong == player {
		player.Money++
		return a
	}
	// 进入其他地方
	a.Belong = player
	return a
}

func InitArea(game *Game, name string, typ int64, home *Player, resources []*Resource) *Area {
	return &Area{
		Game:          game,
		Name:          name,
		Type:          typ,
		Around:        map[*Area]interface{}{},
		Soldiers:      []*Soldier{},
		Resources:     resources,
		Home:          home,
		Belong:        home,
		MoveWaitGroup: utils.InitWaitGroup(),
	}
}
