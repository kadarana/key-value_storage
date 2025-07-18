package main

import (
	"encoding/json"
	"fmt"
	"log"
	"myproj/livecode/pt3"
	"os"
	"path/filepath"
)

// var rootDir string

// такая штука помогает атомарно писать файлы в линуксе
func writeAtomic(path string, b []byte) error {
	dir := filepath.Dir(path)       // /mnt/c/ProjectsGOLANG/projgo/key-value_project/livecode/pt3/file_test
	filename := filepath.Base(path) // название файла показывает
	tmpPathName := filepath.Join(dir, filename+".tmp")
	err := os.WriteFile(tmpPathName, b, 755) // 755 - разрешение
	if err != nil {
		return err
	}

	defer func() {
		os.Remove(tmpPathName)
	}()

	return os.Rename(tmpPathName, path)
}

func main() {

	tasks := []pt3.Task{
		{
			ID:     2,
			Status: "started",
			Cfg: pt3.Config{
				StartAt:  "now",
				Timeline: "before 28.10",
			},
		},
		{
			ID:     3,
			Status: "not started",
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cwd)

	b, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if err = writeAtomic("tasks.json", b); err != nil {
		log.Fatal(err)
	}

	// filepath.Join("/home/student/", "tasks.json")
	// записываем в файл
	// err = os.WriteFile("tasks.json", b, 755) // 755 - разрешение
	// if err != nil {
	// 	log.Fatal(err)
	// } // прочитали ошибку

	// fromFile, err := os.ReadFile("tasks.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var decoded []pt3.Task
	// json.Unmarshal(fromFile, &decoded)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(decoded)

}
