package main

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"throne/game"
	"throne/v1"
	"throne/game/command"
)

var Game *game.Game;

func PlayerHandler(ctx *macaron.Context, sess session.Store) {
	if p := Game.FindPlayer(ctx.Params("player")); p != nil {
		ctx.Map(p)
		return
	}
	player, ok := sess.Get("Player").(*game.Player)

	if !ok {
		ctx.Map((*game.Player)(nil))
		return
	}

	ctx.Map(player)
	return
}

func RegisterCommands(m *macaron.Macaron) {
	g := Game
	m.Get("/player", func(ctx *macaron.Context, sess session.Store) {
		player, ok := sess.Get("Player").(*game.Player)
		if !ok || player == nil {
			player = g.GetPlayer()
			sess.Set("Player", player)
		}
		if player != nil {
			ctx.JSON(200, player)
			return
		}
		ctx.JSON(400, nil)
		return

	})

	m.Get("/areas", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) {
		ctx.JSON(200, game.AreasSlice(player.Game.Map.Areas))
	})

	m.Get("/war", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) {
		if g.War == nil {
			ctx.JSON(200, nil)
			return
		}
		wm, _ := g.War.Marshal()
		ctx.JSON(200, wm)
	})

	m.Get("/war/help/attacker", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {
		war := g.War
		war.HelpAttacker(player)
		return "ok"
	})

	m.Get("/war/help//defender", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {
		war := g.War
		war.HelpDefender(player)
		return "ok"
	})

	m.Get("/war/back/:dst", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {
		war := g.War
		dst := g.FindArea(ctx.Params("dst"))
		war.SetDst(dst)
		return "ok"
	})
	m.Get("/:src/attack/:dst", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {
		//soldiers := ctx.Query("soldiers")
		stayControl := (ctx.Query("stayControl") == "true")
		src := g.FindArea(ctx.Params("src"))
		dst := g.FindArea(ctx.Params("dst"))
		go src.ConsumeAttack(dst, src.SelectAllSoldiers(), stayControl)
		return "ok"
	})
	m.Get("/:src/select", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) {
		src := g.FindArea(ctx.Params("src"))
		ctx.JSON(200, src.AroundCanMove())
	})
	m.Get("/set/:src/:cmd", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {
		src := g.FindArea(ctx.Params("src"))
		err := player.Set(src, game.InitCommand(command.Type(ctx.ParamsInt64("cmd")), false, player))
		if err != nil {
			return err.Error()
		}
		return "ok"
	})

	m.Get("/set/finish", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {
		err := player.FinishSet()
		if err != nil {
			return err.Error()
		}

		return "ok"
	})
}

func main() {

	g := v1.NewGame()
	Game = g

	m := macaron.Classic()
	m.Use(session.Sessioner())
	m.Use(macaron.Static("./templates"))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Delims:macaron.Delims{"<<", ">>"},
	}))
	m.Get("/game", PlayerHandler, func(ctx *macaron.Context) {
		ctx.JSON(200, g)
	})
	m.Get("/players", func(ctx *macaron.Context) {
		ctx.JSON(200, game.Players(g.Players))
	})
	m.Group("/:player", func() {
		RegisterCommands(m)
	})

	m.Run()

	//go func() {
	//	time.Sleep(5 * time.Second)
	//	sLand.ConsumeAttack(tLand, []*game.Soldier{}, true)
	//	time.Sleep(5 * time.Second)
	//	mainLand.ConsumeAttack(sLand, mainLand.SelectAllSoldiers(), true)
	//	time.Sleep(5 * time.Second)
	//	// 结算战争
	//	for w, _ := range Game.Wars {
	//		w.HelpAttacker(mainPlayer)
	//		w.HelpDefender(secondPlayer)
	//		w.DefenderDst = tLand
	//		fmt.Println(w.Fight(), w.Finish())
	//	}
	//	time.Sleep(5 * time.Second)
	//
	//}()
	//
	//go func() {
	//	for {
	//		time.Sleep(5 * time.Second)
	//	}
	//}()

	//Game.Next()
	//// 结算指令
	//
	//go func() {
	//	for {
	//		time.Sleep(5 * time.Second)
	//		// 结算战争
	//		for w, _ := range Game.Wars {
	//			w.HelpAttacker(mainPlayer)
	//			w.HelpDefender(secondPlayer)
	//			w.DefenderDst = tLand
	//			fmt.Println(w.Fight(), w.Finish())
	//		}
	//		time.Sleep(5 * time.Second)
	//	}
	//}()
	//
	//sLand.ConsumeAttack(tLand, []*game.Soldier{}, true)
	//mainLand.ConsumeAttack(sLand, mainLand.SelectAllSoldiers(), true)
	//
	//fmt.Println(Game)

}
