package game

import (
	"fmt"
	"qiniupkg.com/x/log.v7"
	"throne/game/command"
	"throne/game/stage"
	"throne/utils"
)

type Game struct {
	Round     int64
	Stage     int64
	Order     int64
	Map       *Map
	Players   map[*Player]interface{}
	Wars      map[*War]interface{}
	WaitGroup *utils.WaitGroup
}

func (g *Game) Notify() {
	for player, _ := range g.Players {
		player.Notify()
	}
}

func (g *Game) nextStage() {
	if g.Stage == stage.Move {
		g.Round++
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
	// 等待一个移动
	log.Println("需要消费一个指令")
	g.WaitGroup.Add(1)
	log.Println("等待消费一个指令")
	g.WaitGroup.Wait()
	log.Println("已经消费一个指令")
	g.Order++

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
	if g.Stage == stage.Move {
		// 移动一次
		g.move()
		// 如果还有其他指令
		if !g.isClear() {
			return
		}
	}
	g.nextStage()
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

func (g *Game)GetPlayer() *Player {
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
		Round:     0,
		Stage:     stage.Set,
		Players:   map[*Player]interface{}{},
		Wars:      map[*War]interface{}{},
		Order:     1,
		WaitGroup: utils.InitWaitGroup(),
	}
	g.Map = InitMap(g)
	return
}

func (g *Game) String() string {
	return fmt.Sprintln("----Game Begin----") + fmt.Sprintln(g.Map) + fmt.Sprintln(g.Players) + fmt.Sprintln("----Game End----")
}
