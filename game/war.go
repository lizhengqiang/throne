package game

import (
	"fmt"
	"qiniupkg.com/x/errors.v7"
	"qiniupkg.com/x/log.v7"
	"throne/utils"
)

type War struct {
	Game                 *Game
	Area                 *Area
	Attacker             *Player
	Defender             *Player
	Winner               *Player
	AttackerSoldiers     []*Soldier
	DefenderSoldiers     []*Soldier // 失败方战败退出的士兵
	Helpers              map[*Area]interface{}
	AttackerHelpers      map[*Area]interface{}
	DefenderHelpers      map[*Area]interface{}
	AttackerSrc          *Area
	DefenderDst          *Area      // 为空表示防守方胜利,否则为失败方将要退向的地区
	HelperWaitGroup      *utils.WaitGroup
	DefenderDstWaitGroup *utils.WaitGroup
}

func (w *War) String() string {
	return fmt.Sprintln(w.Area) + fmt.Sprintln(w.Helpers) + fmt.Sprintln(w.Attacker) + fmt.Sprintln(w.AttackerSoldiers) + fmt.Sprintln(w.AttackerHelpers)
}

func (w *War) CalcAttacker() (r int64) {
	r += CalcSoldiers(w.AttackerSoldiers)
	for helper, _ := range w.AttackerHelpers {
		r += helper.Calc()
	}
	return
}

func (w *War) CalcDefender() (r int64) {
	r += w.Area.Calc()
	for helper, _ := range w.DefenderHelpers {
		r += helper.Calc()
	}
	return
}

func (w *War) HelpAttacker(player *Player) {
	for a, _ := range w.Helpers {
		if a.Belong == player {
			w.HelperWaitGroup.Done()
			delete(w.Helpers, a)
			w.AttackerHelpers[a] = a
		}
	}
}

func (w *War) HelpDefender(player *Player) {
	for a, _ := range w.Helpers {
		if a.Belong == player {
			w.HelperWaitGroup.Done()
			delete(w.Helpers, a)
			w.DefenderHelpers[a] = nil
		}
	}
}

var HelpersNotReady error = errors.New("相邻位置尚未确认支援")

func (w *War) Fight() error {
	// 支援尚未就绪
	if len(w.Helpers) > 0 {
		return HelpersNotReady
	}
	// 战败的队伍为进攻方
	if w.CalcAttacker() < w.CalcDefender() {
		w.Winner = w.Defender
		for _, s := range w.AttackerSoldiers {
			s.Alive = false
		}
		return nil
	}
	// 战败的队伍为防守方
	w.Winner = w.Attacker
	w.DefenderSoldiers = append(w.Area.Soldiers)
	for _, s := range w.DefenderSoldiers {
		s.Alive = false
	}
	return nil
}

var DefenderDstNotSet error = errors.New("防守方尚未设置撤退目的地")

func (w *War) Finish() error {
	// 处理了一个战争

	defer w.Area.MoveWaitGroup.Done()
	defer log.Println("有一个战争处理完了")

	if w.Winner == w.Attacker && w.DefenderDst == nil {
		return DefenderDstNotSet
	}

	// 清空指令
	w.AttackerSrc.Command = nil
	w.Area.Command = nil

	// 防守方胜利
	if w.Winner == w.Defender {
		w.AttackerSrc.Enter(w.Attacker, w.AttackerSoldiers) // 进攻方撤退到后方
		return nil
	}

	// 进攻方胜利
	w.Area.Leave(w.Defender, w.DefenderSoldiers, false) // 离开战场
	w.DefenderDst.Enter(w.Defender, w.DefenderSoldiers) // 撤退到后方
	w.Area.Enter(w.Attacker, w.AttackerSoldiers)        // 进攻方占领战场
	return nil

}
func (w *War)Handle() {
	log.Println("有一个战争正在处理")
	w.HelperWaitGroup.Wait()
	w.Fight()
	if w.Winner == w.Attacker {
		w.DefenderDstWaitGroup.Wait()
	}
	w.Finish()
}

func InitWar(player *Player, src *Area, dst *Area, soldiers []*Soldier) (w *War) {



	// 有一个战争要处理
	log.Println("有一个战争要处理")
	dst.MoveWaitGroup.Add(1)

	w = &War{
		Game:             player.Game,
		Area:             dst,
		Attacker:         player,
		Defender:         dst.Belong,
		AttackerSoldiers: soldiers,
		DefenderSoldiers: dst.Soldiers,
		AttackerSrc:      src,
	}
	// 找到周边区域
	w.Helpers = dst.Helpers()
	// 等待支援确认
	w.HelperWaitGroup.Add(len(w.Helpers))
	w.AttackerHelpers = map[*Area]interface{}{}
	w.DefenderHelpers = map[*Area]interface{}{}
	// 战争通知到进攻方
	player.Wars[w] = w
	// 战争通知到防守方
	w.Area.Belong.Wars[w] = w
	// 通知支援方
	for a, _ := range w.Helpers {
		a.Belong.Wars[w] = w
	}
	// 存入游戏
	w.Game.Wars[w] = nil
	return

}
