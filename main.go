package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/jfernand/todo/renderer"
	"github.com/jfernand/todo/todo"
	"os"
	"strings"
	"time"
)

var defStyle tcell.Style

type context = struct {
	escape      int
	addingNew   bool
	highlighted int
}

var ctx context

func addNewTodo(s tcell.Screen, newTodo string) {
	blue := tcell.StyleDefault.Foreground(tcell.ColorBlue)
	renderer.EmitStr(s, 0, 2, blue, "New todo: "+newTodo)
}

func tickTodos(x int, y int, todos todo.List) todo.List {
	sortedTodos := todos.SortTodo()
	goal := sortedTodos[0]
	important := sortedTodos[1]
	text := sortedTodos[2]
	shopping := sortedTodos[3]
	checked := sortedTodos[4]

	var allTodos todo.List

	index := 1
	for _, el := range goal {
		if index == y {
			el["done"] = true
		}
		index += 1
		allTodos = append(allTodos, el)
	}

	for _, el := range important {
		if index == y {
			el["done"] = true
		}
		index += 1
		allTodos = append(allTodos, el)
	}

	for _, el := range text {
		if index == y {
			el["done"] = true
		}
		index += 1
		allTodos = append(allTodos, el)
	}

	for _, el := range shopping {
		if index == y {
			el["done"] = true
		}
		index += 1
		allTodos = append(allTodos, el)
	}

	for _, el := range checked {
		if index == y {
			name := el["name"].(string)
			if (len(name) + 6) != x {
				el["done"] = false
				allTodos = append(allTodos, el)
				index += 1
			}
		} else {
			allTodos = append(allTodos, el)
			index += 1
		}
	}

	allTodos.SaveTodos()
	return allTodos
}

func main() {

	s, e := tcell.NewScreen()
	encoding.Register()

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	escapeKeypressCount := 0
	addNew := false
	newTodo := ""

	todos, err := todo.LoadTodos()
	if err != nil {
		s.Fini()
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(0)
	}

	renderer.RenderTodos(s, todos)
	s.Show()

	defer s.Fini()
	events := make(chan tcell.Event)
	go func() {
		for {
			ev := s.PollEvent()
			events <- ev
		}
	}()
	go func() {
		for {
			ev := <-events
			switch ev := ev.(type) {
			case *tcell.EventKey:
				addNew, newTodo = handleKey(ev, escapeKeypressCount, s, addNew, newTodo)
			case *tcell.EventMouse:
				x, y := ev.Position()
				switch ev.Buttons() {
				case tcell.Button1, tcell.Button2, tcell.Button3:
					allTodos, _ := todo.LoadTodos()
					l := len(allTodos)
					if y < (l + 2) {
						s.Clear()
						tickTodos(x, y, allTodos)
						todos, _ := todo.LoadTodos()
						renderer.RenderTodos(s, todos)
						s.Show()
					} else if y == l+2 {
						s.Clear()
						addNewTodo(s, newTodo)
						addNew = true
						s.Show()
					}
				}
			}
		}
	}()

	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			if !addNew {
				s.Clear()
				allTodos, _ := todo.LoadTodos()
				renderer.RenderTodos(s, allTodos)
				s.Sync()
				s.Show()
			}
		}
	}
}

func handleKey(ev *tcell.EventKey, escapeKeypressCount int, s tcell.Screen, addNew bool, newTodo string) (bool, string) {
	switch ev.Key() {
	case tcell.KeyEscape:
	case tcell.KeyHome:
		{
			escapeKeypressCount++
			if escapeKeypressCount > 1 {
				s.Fini()
				os.Exit(0)
			}
		}
	case tcell.KeyEnter:
		{
			if addNew == true {
				todos, _ := todo.LoadTodos()
				newValue := make(map[string]interface{})
				newValue["name"] = newTodo
				newValue["done"] = false
				addNew = false
				newTodo = ""
				todos = append(todos, newValue)
				todos.SaveTodos()
				allTodos, _ := todo.LoadTodos()
				s.Clear()
				renderer.RenderTodos(s, allTodos)
				s.Show()
			}
		}
	case tcell.KeyRune:
		{
			if addNew {
				keyValue := strings.Replace(strings.Replace(ev.Name(), "Rune[", "", 1), "]",
					"", 1)
				newTodo += keyValue
				addNewTodo(s, newTodo)
				addNew = true
				s.Show()
			}
		}
	case tcell.KeyBackspace2:
	case tcell.KeyBackspace:
		{
			if len(newTodo) > 0 {
				newTodo = strings.TrimSuffix(newTodo, newTodo[len(newTodo)-1:])
				s.Clear()
				addNewTodo(s, newTodo)
				s.Show()
			} else {
				addNew = false
				allTodos, _ := todo.LoadTodos()
				s.Clear()
				renderer.RenderTodos(s, allTodos)
				s.Show()
			}
		}
	}
	return addNew, newTodo
}
