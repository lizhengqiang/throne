package v1

import (
	"throne/game"
	"throne/game/player"
	"throne/game/soldier"
	"fmt"
	"throne/game/command"
	"throne/game/area"
	"throne/game/resource"
)

func NewGame() *game.Game {
	g := game.InitGame()
	mainPlayer := game.InitPlayer(g, "正强")
	secondPlayer := game.InitPlayer(g, "电脑")
	// 初始化地图
	mainLand := game.InitArea(g, "主城", area.Land, mainPlayer, []*game.Resource{game.InitResource(resource.City)})
	sLand := game.InitArea(g, "副城", area.Land, secondPlayer, []*game.Resource{game.InitResource(resource.City)})
	tLand := game.InitArea(g, "小城", area.Land, secondPlayer, []*game.Resource{game.InitResource(resource.City)})
	mainSea := game.InitArea(g, "主海", area.Sea, nil, []*game.Resource{})
	secondSea := game.InitArea(g, "副海", area.Sea, nil, []*game.Resource{})
	// 关联
	mainLand.AddAround(mainSea)
	sLand.AddAround(secondSea)
	sLand.AddAround(tLand)
	mainSea.AddAround(secondSea)
	// 添加
	g.Map.AddArea(mainLand)
	g.Map.AddArea(sLand)
	g.Map.AddArea(tLand)
	g.Map.AddArea(mainSea)
	g.Map.AddArea(secondSea)
	// 初始化角色
	mainPlayer.SetOrder(player.OrderA, 2)
	secondPlayer.SetOrder(player.OrderA, 1)
	g.AddPlayer(mainPlayer)
	g.AddPlayer(secondPlayer)
	// 初始化兵力
	mainLand.Enter(mainPlayer, []*game.Soldier{game.InitSoldier(soldier.Cavalry), game.InitSoldier(soldier.Foot)})
	mainSea.Enter(mainPlayer, []*game.Soldier{game.InitSoldier(soldier.Ship), game.InitSoldier(soldier.Ship)})
	secondSea.Enter(mainPlayer, []*game.Soldier{game.InitSoldier(soldier.Ship), game.InitSoldier(soldier.Ship)})
	sLand.Enter(secondPlayer, []*game.Soldier{game.InitSoldier(soldier.Cavalry), game.InitSoldier(soldier.Cavalry)})

	fmt.Println(g)

	// 放置指令
	mainLand.PutCommand(game.InitCommand(command.Attack, false, mainPlayer))
	mainSea.PutCommand(game.InitCommand(command.Help, false, mainPlayer))
	secondSea.PutCommand(game.InitCommand(command.Help, false, mainPlayer))
	sLand.PutCommand(game.InitCommand(command.Attack, false, secondPlayer))
	tLand.PutCommand(game.InitCommand(command.Help, false, secondPlayer))

	fmt.Println(g)

	go g.Begin()

	return g

}
