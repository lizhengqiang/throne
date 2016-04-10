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
	"encoding/json"
)

type Area struct {
	Game          *Game `json:"-"`
	Id            string
	Name          string
	Type          int64                            // 地面类型
	Around        map[*Area]interface{} `json:"-"` // 临近地区
	Resources     []*Resource                      // 资源
	Soldiers      []*Soldier                       // 士兵单位
	Command       *Command                         // 指令
	Belong        *Player                          // 指示物
	Home          *Player

	MoveWaitGroup *utils.WaitGroup
}

func AreasSlice(m []*Area) []interface{} {
	r := []interface{}{}

	for _, d := range m {
		dMap, err := d.Marshal()
		if err != nil {
			continue
		}
		r = append(r, dMap)
	}
	return r
}

func Areas(m map[*Area]interface{}) []interface{} {
	r := []interface{}{}

	for d, _ := range m {
		dMap, err := d.marshal()
		if err != nil {
			continue
		}
		r = append(r, dMap)
	}
	return r
}

func (a *Area) Marshal() (m map[string]interface{}, err error) {
	m, err = a.marshal()
	m["Around"] = Areas(a.Around)
	return m, err
}
func (a *Area)marshal() (m map[string]interface{}, err error) {
	m = map[string]interface{}{}
	bytes, err := json.Marshal(a)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return
	}

	return
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
var WrongCommand error = errors.New("指令不正确")
var NotTurn error = errors.New("还没轮到")
var NotStage error = errors.New("不是移动的阶段")
// 结算指令
func (a *Area) ConsumeCommand(typ command.Type) error {
	// 没有指令
	if a.Command == nil {
		return NoCommand
	}
	// 指令错误
	if a.Command.Type != typ {
		return WrongCommand
	}
	// 不是移动阶段
	if a.Game.Stage != stage.Move {
		return NotStage
	}
	// 还没有到执行顺序, 仅当进攻或者偷袭时检查
	if a.Belong.Orders[player.OrderA] != a.Game.Order && (a.Command.Type == command.Attack || a.Command.Type == command.Steal) {
		return NotTurn
	}
	return nil
}

func (a *Area) ConsumedCommand(typ command.Type) error {
	if typ == command.Steal {
		a.Game.WaitGroup.Done();
	}

	if typ == command.Attack {
		a.Game.WaitGroup.Done();
	}
	a.Command = nil
	return nil
}
func (a *Area) ConsumeSteal(areas ...*Area) error {
	// 检查命令是否正常
	if err := a.ConsumeCommand(command.Steal); err != nil {
		log.Println(err.Error())
		return err
	}
	// 取消选中的地区的指令
	for _, ar := range areas {
		if ar.Command.Type == command.Steal || ar.Command.Type == command.Help || ar.Command.Type == command.Money {
			ar.Command = nil
		}
	}
	// 完成这个指令
	return a.ConsumedCommand(command.Steal)
}

func (a *Area) ConsumeMoney() error {
	// 检查命令是否正常
	if err := a.ConsumeCommand(command.Money); err != nil {
		log.Println(err.Error())
		return err
	}
	// 增加金钱
	for _, res := range a.Resources {
		if res.Type == resource.Money {
			a.Belong.Money++
		}
	}
	// 最少增加一个
	a.Belong.Money++
	// 完成指令
	return a.ConsumedCommand(command.Money)
}

func (a *Area) ConsumeAttack(ar *Area, soldiers []*Soldier, stayControl bool) error {
	// 检查命令是否正常
	if err := a.ConsumeCommand(command.Attack); err != nil {
		log.Println(err.Error())
		return err
	}
	// 移动
	if err := a.Belong.Move(a, ar, soldiers, stayControl); err != nil {
		// 没法移动
		log.Println(err.Error())
		return err
	}
	// 处理战争
	log.Println("等待战争处理完毕")
	ar.MoveWaitGroup.Wait()
	log.Println("战争处理完毕")

	return a.ConsumedCommand(command.Attack)
}

func (a *Area) Calc() (r int64) {
	return CalcSoldiers(a.Soldiers)
}

func (a *Area)AroundCanMove() (r []*Area) {
	r = []*Area{}
	for area, _ := range a.AroundMove() {
		r = append(r, area)
	}
	return r
}
func (a *Area) LandHelpers(has map[*Area]interface{}) (r map[*Area]interface{}) {
	r = map[*Area]interface{}{}
	for ar, _ := range a.Around {
		// 已经有了
		if _, ok := has[ar]; ok {
			continue;
		}
		// 不是自己的
		if ar.Belong != a.Belong {
			continue
		}
		// 陆地支援
		if ar.Type == area.Land {
			r[ar] = nil
			has[ar] = nil
			continue
		}
		// 跨海支援
		has[ar] = nil
		for arr, _ := range ar.LandHelpers(has) {
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
		for are, _ := range ar.LandHelpers(map[*Area]interface{}{}) {
			r[are] = nil
		}
	}

	return

}

func (a *Area) AroundMoveLand(hasUsedAreas map[*Area]interface{}) (result map[*Area]interface{}) {
	result = map[*Area]interface{}{}
	for anAroundArea, _ := range a.Around {
		// 已经有了
		if _, ok := hasUsedAreas[anAroundArea]; ok {
			continue;
		}
		// 不是自己的
		if anAroundArea.Belong != a.Belong {
			continue
		}
		// 陆地支援
		if anAroundArea.Type == area.Land {
			result[anAroundArea] = nil
			hasUsedAreas[anAroundArea] = nil
			continue
		}
		// 如果是一片海洋
		// 跨海支援,这片海已经过去过了
		hasUsedAreas[anAroundArea] = nil
		// 在海上时候周围能移动过去的陆地
		for anAroundAreaInTheSea, _ := range anAroundArea.AroundMoveLand(hasUsedAreas) {
			result[anAroundAreaInTheSea] = nil
		}

	}
	return result
}
func (a *Area) AroundMove() (result map[*Area]interface{}) {
	result = map[*Area]interface{}{}
	// 自己属于自己周围可以移动的地方
	result[a] = nil
	// 海上
	if a.Type == area.Sea {
		for ar, _ := range a.Around {
			if ar.Type == area.Sea {
				result[ar] = nil
			}
		}
		return
	}
	// 陆地
	// 全部周围均可支援
	for anAroundArea, _ := range a.Around {
		// 陆地
		if anAroundArea.Type == area.Land {
			result[anAroundArea] = nil
			continue
		}

		// anAroundArea.Type == area.Sea

		// 不是自己的海洋
		if anAroundArea.Belong != a.Belong {
			continue
		}
		// 跨海支援
		for landAreasCanMove, _ := range anAroundArea.AroundMoveLand(map[*Area]interface{}{}) {
			result[landAreasCanMove] = nil
		}
	}

	return

}

var CannotMove error = errors.New("不能移动到这个位置")

func (a *Area) canMove(target *Area) bool {
	if _, ok := a.AroundMove()[target]; ok {
		return true
	}
	return false
}
func (a *Area) CanMove(target *Area) bool {
	return a.canMove(target)
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

func (a *Area)ReliveAll() error {
	for _, s := range a.Soldiers {
		s.Alive = true
	}
	return nil
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
