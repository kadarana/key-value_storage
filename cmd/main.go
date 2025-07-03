package main

import (
	"fmt"
	"myproj/internal/pkg/storage"
)

func main() {
	s, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	s.Set("key1", "value1")
	s.Set("key2", "1337")

	fmt.Println(*s.Get("key1"))
	fmt.Println(*s.Get("key2"))
	fmt.Println(s.Get("key3"))

	fmt.Println(s.GetType("key1"))
	fmt.Println(s.GetType("key2"))
	fmt.Println(s.GetType("key3"))

}
