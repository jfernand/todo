package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var defStyle tcell.Style

func todosPath() string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	return filepath.Join(dir, "/.todos.yaml")
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
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
func sortList(todos []map[string]interface{}) {
	sort.Slice(todos, func(i, j int) bool {
		return todos[i]["name"].(string) < todos[j]["name"].(string)
	})
}

func sortTodo(todos []map[string]interface{}) [][]map[string]interface{} {
	var goal []map[string]interface{}
	var important []map[string]interface{}
	var todo []map[string]interface{}
	var shopping []map[string]interface{}
	var done []map[string]interface{}

	for _, el := range todos {
		name := el["name"].(string)
		checked := el["done"].(bool)
		if checked == true {
			done = append(done, el)
		} else if strings.Contains(name, "Goal:") {
			goal = append(goal, el)
		} else if strings.Contains(name, "*") {
			important = append(important, el)
		} else if strings.Contains(name, "Shopping:") {
			shopping = append(shopping, el)
		} else {
			todo = append(todo, el)
		}
	}

	sortList(goal)
	sortList(important)
	sortList(shopping)
	sortList(todo)
	sortList(done)

	var allTodos [][]map[string]interface{}

	allTodos = append(allTodos, goal)
	allTodos = append(allTodos, important)
	allTodos = append(allTodos, todo)
	allTodos = append(allTodos, shopping)
	allTodos = append(allTodos, done)

	return allTodos
}

func renderTodos(s tcell.Screen, todos []map[string]interface{}) {
	green := tcell.StyleDefault.Foreground(tcell.ColorLawnGreen)
	yellow := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	seaGreen := tcell.StyleDefault.Foreground(tcell.ColorDarkSeaGreen)
	purple := tcell.StyleDefault.Foreground(tcell.ColorPurple)
	grey := tcell.StyleDefault.Foreground(tcell.ColorGrey)
	orange := tcell.StyleDefault.Foreground(tcell.ColorOrange)
	blue := tcell.StyleDefault.Foreground(tcell.ColorBlue)

	emitStr(s, 0, 0, green, "TODO")

	index := 1

	allTodos := sortTodo(todos)
	goal := allTodos[0]
	important := allTodos[1]
	todo := allTodos[2]
	shopping := allTodos[3]
	done := allTodos[4]

	for _, el := range goal {
		name := el["name"].(string)
		emitStr(s, 0, index, seaGreen, "-- "+name+" --")
		index += 1
	}

	for _, el := range important {
		name := el["name"].(string)
		emitStr(s, 0, index, yellow, "[ ] "+name)
		index += 1
	}

	for _, el := range todo {
		name := el["name"].(string)
		emitStr(s, 0, index, blue, "[ ] "+name)
		index += 1
	}

	for _, el := range shopping {
		name := el["name"].(string)
		emitStr(s, 0, index, purple, "[ ] "+name)
		index += 1
	}

	for _, el := range done {
		name := el["name"].(string)
		emitStr(s, 0, index, grey, "[x] "+name+" (-)")
		index += 1
	}

	emitStr(s, 0, index, blue, "")
	index += 1
	emitStr(s, 0, index, orange, "Add +")
}

func addNewTodo(s tcell.Screen, newTodo string) {
	blue := tcell.StyleDefault.Foreground(tcell.ColorBlue)
	emitStr(s, 0, 2, blue, "New Todo: "+newTodo)
}

func tickTodos(x int, y int, todos []map[string]interface{}) []map[string]interface{} {

	sortedTodos := sortTodo(todos)
	goal := sortedTodos[0]
	important := sortedTodos[1]
	todo := sortedTodos[2]
	shopping := sortedTodos[3]
	checked := sortedTodos[4]

	var allTodos []map[string]interface{}

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

	for _, el := range todo {
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

	saveTodos(allTodos)
	return allTodos
}

func saveTodos(todos []map[string]interface{}) {
	b := make(map[string]interface{})
	b["todos"] = todos
	d, err := yaml.Marshal(b)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	f, err := os.Create(todosPath())
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	f.Write(d)
}

func getTodos() ([]map[string]interface{}, error) {
	yamlFile, err := ioutil.ReadFile(todosPath())

	type AllToDos struct {
		ToDos []map[string]interface{}
	}

	todos := AllToDos{}
	if err == nil {
		err = yaml.Unmarshal(yamlFile, &todos)
	}

	return todos.ToDos, err
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

	todos, err := getTodos()
	if err != nil {
		s.Fini()
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(0)
	}

	renderTodos(s, todos)
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
				if ev.Key() == tcell.KeyEscape {
					escapeKeypressCount++
					if escapeKeypressCount > 1 {
						s.Fini()
						os.Exit(0)
					}
				} else if ev.Key() == tcell.KeyEnter {
					if addNew == true {
						todos, _ := getTodos()
						newValue := make(map[string]interface{})
						newValue["name"] = newTodo
						newValue["done"] = false
						addNew = false
						newTodo = ""
						todos = append(todos, newValue)
						saveTodos(todos)
						allTodos, _ := getTodos()
						s.Clear()
						renderTodos(s, allTodos)
						s.Show()
					}
				} else if ev.Key() == tcell.KeyRune {
					if addNew {
						keyValue := strings.Replace(strings.Replace(ev.Name(), "Rune[", "", 1), "]",
							"", 1)
						newTodo += keyValue
						addNewTodo(s, newTodo)
						addNew = true
						s.Show()
					}
				} else if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBackspace {
					if len(newTodo) > 0 {
						newTodo = strings.TrimSuffix(newTodo, newTodo[len(newTodo)-1:])
						s.Clear()
						addNewTodo(s, newTodo)
						s.Show()
					} else {
						addNew = false
						allTodos, _ := getTodos()
						s.Clear()
						renderTodos(s, allTodos)
						s.Show()
					}
				}
			case *tcell.EventMouse:
				x, y := ev.Position()
				switch ev.Buttons() {
				case tcell.Button1, tcell.Button2, tcell.Button3:
					allTodos, _ := getTodos()
					l := len(allTodos)
					if y < (l + 2) {
						s.Clear()
						tickTodos(x, y, allTodos)
						todos, _ := getTodos()
						renderTodos(s, todos)
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
				allTodos, _ := getTodos()
				renderTodos(s, allTodos)
				s.Sync()
				s.Show()
			}
		}

	}
}
