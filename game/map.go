package game

import "fmt"

type Map struct {
	Game      *Game
	Areas     []*Area
	currentId byte
}

func InitMap(game *Game) *Map {
	return &Map{
		Game:  game,
		Areas: []*Area{},
		currentId: 'A',
	}
}

func (m *Map) String() string {
	result := ""
	for _, area := range m.Areas {
		result = result + fmt.Sprintln(area)
	}
	return result
}

func (m *Map) AddArea(area *Area) *Map {
	m.Areas = append(m.Areas, area)
	// 添加ID
	area.Id = string(append([]byte{}, m.currentId))
	m.currentId ++
	return m
}
