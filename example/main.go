package main

import (
	"fmt"
	"github.com/nocd5/pecoru"
)

type Item struct {
	Name  string
	和名    string
	Price int
}

func main() {
	fruits := []Item{
		{"Apple", "りんご", 90},
		{"Banana", "バナナ", 100},
		{"Grape", "ぶどう", 120},
		{"Mango", "太陽のタマゴ", 16000},
		{"Orange", "オレンジ", 80},
	}
	var list []string
	for _, v := range fruits {
		list = append(list, v.Name)
	}

	for content := range pecoru.Select(list) {
		if content.Error == nil {
			fmt.Printf("%s => %sは%d円です。\n", content.Label, fruits[content.Index].和名, fruits[content.Index].Price)
		}
	}
}
