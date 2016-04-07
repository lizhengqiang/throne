package command

type Type int64

const (
	Attack Type = iota
	Defend
	Help
	Money
	Steal
)
