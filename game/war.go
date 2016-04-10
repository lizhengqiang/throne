package game

import (
	"encoding/json"
	"fmt"
	"qiniupkg.com/x/errors.v7"
	"qiniupkg.com/x/log.v7"
	"throne/utils"
)

type War struct {
	Game                 *Game `json:"-"`
	Area                 *Area
	Attacker             *Player
	Defender             *Player
	Winner               *Player
	AttackerSoldiers     []*Soldier
	DefenderSoldiers     []*Soldier // 失败方战败退出的士兵
	Helpers              map[*Area]interface{} `json:"-"`
	AttackerHelpers      map[*Area]interface{} `json:"-"`
	DefenderHelpers      map[*Area]interface{} `json:"-"`
	AttackerSrc          *Area
	DefenderDst          *Area      // 为空表示防守方胜利,否则为失败方将要退向的地区
	HelperWaitGroup      *utils.WaitGroup `json:"-"`
	DefenderDstWaitGroup *utils.WaitGroup `json:"-"`
}

func (w *War) String() string {
	return fmt.Sprintln(w.Area) + fmt.Sprintln(w.Helpers) + fmt.Sprintln(w.Attacker) + fmt.Sprintln(w.AttackerSoldiers) + fmt.Sprintln(w.AttackerHelpers)
}

func (w *War) SetDst(area *Area) {
	// 死亡
	if area == w.Area {
		w.DefenderDst = w.Defender.PreArea
		w.DefenderDstWaitGroup.Done()
		return
	}
	// 退到目的地
	w.DefenderDst = area
	w.DefenderDstWaitGroup.Done()
	return
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
	defer func() {
		w.Game.War = nil
	}()
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
func (w *War) Handle() {
	log.Println("有一个战争正在处理")
	w.HelperWaitGroup.Wait()
	w.Fight()
	if w.Winner == w.Attacker {
		log.Println("进攻方胜利")
		w.DefenderDstWaitGroup.Add(1)
		w.DefenderDstWaitGroup.Wait()
	}
	w.Finish()
}

func (w *War) Marshal() (m map[string]interface{}, err error) {
	m = map[string]interface{}{}
	bytes, err := json.Marshal(w)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return
	}

	m["Helpers"] = Areas(w.Helpers)
	m["AttackerHelpers"] = Areas(w.AttackerHelpers)
	m["DefenderHelpers"] = Areas(w.DefenderHelpers)
	return
}

func InitWar(player *Player, src *Area, dst *Area, soldiers []*Soldier) (w *War) {

	// 有一个战争要处理
	log.Println("有一个战争要处理")
	dst.MoveWaitGroup.Add(1)
	w = &War{
		Game:                 player.Game,
		Area:                 dst,
		Attacker:             player,
		Defender:             dst.Belong,
		AttackerSoldiers:     soldiers,
		DefenderSoldiers:     dst.Soldiers,
		AttackerSrc:          src,
		Helpers:              map[*Area]interface{}{},
		AttackerHelpers:      map[*Area]interface{}{},
		DefenderHelpers:      map[*Area]interface{}{},
		HelperWaitGroup:      utils.InitWaitGroup(),
		DefenderDstWaitGroup: utils.InitWaitGroup(),
	}
	// 找到周边区域
	w.Helpers = dst.Helpers()
	// 等待支援确认
	w.HelperWaitGroup.Add(len(w.Helpers))
	w.AttackerHelpers = map[*Area]interface{}{}
	w.DefenderHelpers = map[*Area]interface{}{}
	// 存入游戏
	w.Game.War = w
	return

}
