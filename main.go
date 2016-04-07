package main

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"qiniupkg.com/x/log.v7"
	"throne/game"
	"throne/v1"
	"fmt"
)

func PlayerHandler(ctx *macaron.Context, sess session.Store) {
	player, ok := sess.Get("Player").(*game.Player)

	if !ok {
		ctx.Map((*game.Player)(nil))
		return
	}

	ctx.Map(player)
	return
}

func main() {

	g := v1.NewGame()

	m := macaron.Classic()
	m.Use(session.Sessioner())

	m.Get("/player", func(ctx *macaron.Context, sess session.Store) string {

		player, ok := sess.Get("Player").(*game.Player)

		if !ok || player == nil {
			player = g.GetPlayer()
			sess.Set("Player", player)
		}

		log.Println(player)
		if player != nil {

			return "角色:" + player.Name
		}

		return "没有角色"

	})

	m.Get("/areas", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {

		return fmt.Sprintf("%+v", player.Game.Map.Areas)
	})

	m.Get("/wars", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {
		return fmt.Sprintf("%+v", player.Wars)
	})

	m.Get("/help/:dst/attacker", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {
		dst := g.FindArea(ctx.Params("dst"))
		player.Wars
	})

	m.Get("/help/:dst/defender", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {
		return fmt.Sprintf("%+v", player.Wars)
	})
	m.Get("/:src/attack/:dst", PlayerHandler, func(ctx *macaron.Context, sess session.Store, player *game.Player) string {
		//soldiers := ctx.Query("soldiers")
		stayControl := (ctx.Query("stayControl") == "true")
		src := g.FindArea(ctx.Params("src"))
		dst := g.FindArea(ctx.Params("dst"))
		go src.ConsumeAttack(dst, src.SelectAllSoldiers(), stayControl)
		return "OK"
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
