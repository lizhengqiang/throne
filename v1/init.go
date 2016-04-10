package v1

import (
	"throne/game"
	"throne/game/player"
	"throne/game/soldier"
	"throne/game/area"
	"throne/game/resource"
)

func NewGame() *game.Game {
	g := game.InitGame()
	GreyjoyPlayer := game.InitPlayer(g, "Greyjoy")
	TyrellPlayer := game.InitPlayer(g, "Tyrell")
	// 初始化地图
	PaikeIsland := game.InitArea(g, "派克岛", area.Land, GreyjoyPlayer, []*game.Resource{game.InitResource(resource.City), game.InitResource(resource.Supply), game.InitResource(resource.Money)})
	Gaoting := game.InitArea(g, "高庭", area.Land, TyrellPlayer, []*game.Resource{game.InitResource(resource.City), game.InitResource(resource.Supply), game.InitResource(resource.Supply)})
	HewanLand := game.InitArea(g, "河湾地", area.Land, nil, []*game.Resource{game.InitResource(resource.City)})
	TieminSea := game.InitArea(g, "铁民湾", area.Sea, nil, []*game.Resource{game.InitResource(resource.Town)})
	XixiaSea := game.InitArea(g, "西夏海", area.Sea, nil, []*game.Resource{})
	RiluoSea := game.InitArea(g, "日落之海", area.Sea, nil, []*game.Resource{})
	// 关联
	RiluoSea.AddAround(TieminSea).AddAround(XixiaSea)
	PaikeIsland.AddAround(TieminSea)
	Gaoting.AddAround(XixiaSea).AddAround(HewanLand)
	// 添加
	g.Map.AddArea(PaikeIsland).AddArea(Gaoting).AddArea(TieminSea).AddArea(XixiaSea).AddArea(HewanLand).AddArea(RiluoSea)
	// 初始化角色
	GreyjoyPlayer.SetOrder(player.OrderA, 1)
	TyrellPlayer.SetOrder(player.OrderA, 2)
	g.AddPlayer(GreyjoyPlayer)
	g.AddPlayer(TyrellPlayer)
	// 初始化兵力
	PaikeIsland.Enter(GreyjoyPlayer, []*game.Soldier{game.InitSoldier(soldier.Cavalry), game.InitSoldier(soldier.Foot)})
	TieminSea.Enter(GreyjoyPlayer, []*game.Soldier{game.InitSoldier(soldier.Ship), game.InitSoldier(soldier.Ship)})
	Gaoting.Enter(TyrellPlayer, []*game.Soldier{game.InitSoldier(soldier.Cavalry), game.InitSoldier(soldier.Foot)})
	XixiaSea.Enter(TyrellPlayer, []*game.Soldier{game.InitSoldier(soldier.Ship)})
	go g.Begin()

	return g

}
