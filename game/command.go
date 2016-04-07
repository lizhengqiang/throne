package game

import (
	"fmt"
	"throne/game/command"
)

type Command struct {
	Belong *Player
	Type   command.Type
	Plus   bool
}

func InitCommand(t command.Type, plus bool, player *Player) *Command {
	return &Command{
		Belong: player,
		Type:   t,
		Plus:   plus,
	}
}

func (c *Command) String() string {
	return fmt.Sprintf("(类型:%d,星:%v)", c.Type, c.Plus)
}
