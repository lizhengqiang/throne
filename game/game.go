package game

import (
	"fmt"
	"qiniupkg.com/x/log.v7"
	"throne/game/command"
	"throne/game/stage"
	"throne/utils"
	"throne/game/player"
)

type Game struct {
	Round     int64
	Stage     int64
	Order     int64
	Map       *Map                    `json:"-"`
	Players   map[*Player]interface{} `json:"-"`
	War       *War                    `json:"-"`
	WaitGroup *utils.WaitGroup        `json:"-"`
}

func (g *Game) FindPlayer(name string) *Player {
	for player, _ := range g.Players {
		if player.Name == name {
			return player
		}
	}
	return nil
}
func (g *Game) Notify() {
	for player, _ := range g.Players {
		player.Notify()
	}
}

func (g *Game) nextStage() {
	if g.Stage == stage.Move {
		g.Round++
		g.Order = 1
		g.Stage = stage.Event
	}

	if g.Stage == stage.Set {
		g.Stage = stage.Move
	}

	if g.Stage == stage.Event {
		g.Stage = stage.Set
	}
}

func (g *Game) move() {

	if !g.isPlayerClear(g.GetOrderPlayer()) {
		// 等待一个移动
		log.Println("需要消费一个指令")
		g.WaitGroup.Add(1)
	}

	log.Println("等待消费一个指令")
	g.WaitGroup.Wait()

	log.Println("已经消费一个指令")

	// 循环1-len(g.Players)
	if g.Order == int64(len(g.Players)) {
		g.Order = 1
	}else {
		g.Order++
	}

}

func (g *Game)GetOrderPlayer() *Player {
	for p, _ := range g.Players {
		if p.Orders[player.OrderA] == g.Order {
			return p
		}
	}
	return nil
}

func (g *Game)isPlayerClear(player *Player) bool {
	for _, a := range g.Map.Areas {
		if a.Belong == player {
			if a.Command == nil {
				continue
			}
			if a.Command.Type == command.Steal {
				return false
			}

			if a.Command.Type == command.Attack {
				return false
			}
		}

	}
	return true
}

func (g *Game) isClear() bool {
	for _, a := range g.Map.Areas {
		if a.Command == nil {
			continue
		}
		if a.Command.Type == command.Steal {
			return false
		}

		if a.Command.Type == command.Attack {
			return false
		}
	}
	return true
}
func (g *Game) Next() {
	// 结束前通知动作
	g.Notify()
	if g.Stage == stage.Move && !g.isClear() {
		// 如果还有其他指令
		// 结算偷袭,进攻
		g.move()
		return
	}

	if g.Stage == stage.Move && g.isClear() {
		// 没有偷袭或者进攻命令了
		// 结算巩固权利
		for _, a := range g.Map.Areas {
			a.ConsumeMoney()
		}
		// 进入事件阶段
		g.nextStage()
		return
	}

	if g.Stage == stage.Set {
		// 放置指令
		for p, _ := range g.Players {
			p.CanSet = true;
			g.WaitGroup.Add(1)
		}
		// 复生所有的士兵
		for _, a := range g.Map.Areas {
			a.ReliveAll()
		}
		// 等待放置完成
		log.Println("等待放置指令")
		g.WaitGroup.Wait()
		log.Println("放置指令完成")
		// 进入移动阶段
		g.nextStage()
		return
	}

	if g.Stage == stage.Event {
		// 进入放置阶段
		g.nextStage()
		return
	}

	return
}
func (g *Game) AddPlayer(player *Player) {
	g.Players[player] = nil
}

func (g *Game) Begin() {
	for g.Round < 10 {
		g.Next()
	}
}

func (g *Game) FindArea(id string) *Area {
	for _, area := range g.Map.Areas {
		if area.Id == id {
			return area
		}
	}
	return nil
}

func (g *Game) GetPlayer() *Player {
	for player, status := range g.Players {
		// 已经领取
		if status != nil {
			continue
		}
		// 领取
		g.Players[player] = true
		return player
	}
	return nil

}
func InitGame() (g *Game) {
	g = &Game{
		Round:     1,
		Stage:     stage.Set,
		Players:   map[*Player]interface{}{},
		Order:     1,
		WaitGroup: utils.InitWaitGroup(),
	}
	g.Map = InitMap(g)
	return
}

func (g *Game) String() string {
	return fmt.Sprintln("----Game Begin----") + fmt.Sprintln(g.Map) + fmt.Sprintln(g.Players) + fmt.Sprintln("----Game End----")
}
