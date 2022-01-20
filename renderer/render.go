package renderer

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jfernand/todo/todo"
	"github.com/mattn/go-runewidth"
)

func EmitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func RenderTodos(s tcell.Screen, todos todo.List) {
	green := tcell.StyleDefault.Foreground(tcell.ColorLawnGreen)
	yellow := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	seaGreen := tcell.StyleDefault.Foreground(tcell.ColorDarkSeaGreen)
	purple := tcell.StyleDefault.Foreground(tcell.ColorPurple)
	grey := tcell.StyleDefault.Foreground(tcell.ColorGrey)
	orange := tcell.StyleDefault.Foreground(tcell.ColorOrange)
	blue := tcell.StyleDefault.Foreground(tcell.ColorBlue)

	EmitStr(s, 0, 0, green, "TODO")

	index := 1

	allTodos := todos.SortTodo()
	goal := allTodos[0]
	important := allTodos[1]
	todo := allTodos[2]
	shopping := allTodos[3]
	done := allTodos[4]

	for _, el := range goal {
		name := el["name"].(string)
		EmitStr(s, 0, index, seaGreen, "-- "+name+" --")
		index += 1
	}

	for _, el := range important {
		name := el["name"].(string)
		EmitStr(s, 0, index, yellow, "[ ] "+name)
		index += 1
	}

	for _, el := range todo {
		name := el["name"].(string)
		EmitStr(s, 0, index, blue, "[ ] "+name)
		index += 1
	}

	for _, el := range shopping {
		name := el["name"].(string)
		EmitStr(s, 0, index, purple, "[ ] "+name)
		index += 1
	}

	for _, el := range done {
		name := el["name"].(string)
		EmitStr(s, 0, index, grey, "[x] "+name+" (-)")
		index += 1
	}

	EmitStr(s, 0, index, blue, "")
	index += 1
	EmitStr(s, 0, index, orange, "Add +")
}
