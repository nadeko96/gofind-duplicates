package main

import (
	"crypto/sha512"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var files = make(map[[sha512.Size]byte]string)

func checkDepulicates(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
	}

	if info.IsDir() {
		return nil
	}

	ext := filepath.Ext(path)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return nil
	}

	data, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	hash := sha512.Sum512(data)
	// fmt.Printf("%x\n", hash)

	if f, ok := files[hash]; ok {
		fmt.Printf("%q is a duplicate of %q\n", path, f)
		promptDelete(path)
	} else {
		checkBadNames(path, hash)
		files[hash] = path
	}

	return nil
}

func promptDelete(path string) {
	var input string
	fmt.Printf("Do you want to delete %q (type yes or no): ", path)
	fmt.Scanf("%s", &input)

	if strings.ToUpper(input) == "YES" {
		err := os.Remove(path)
		handle(err)
	}
}

func checkBadNames(path string, hash [64]uint8) {

	_, file := filepath.Split(path)
	ext := filepath.Ext(path)
	file = file[0 : len(file)-len(ext)]
	_, err := strconv.ParseInt(file, 10, 64)
	if file == string(hash[:]) || err == nil {
		promptNameChange(path, file, ext)
	}

}

func promptNameChange(path string, file string, ext string) {
	var input string
	fmt.Printf("%q is a shitty filename, open and change it? \n", file+"."+ext)
	fmt.Scanf("%s", input)
	openImage(path)
	fmt.Printf("What is a good name for it fam? Press \"s\" to skip")
	fmt.Scanf("%s", &input)

	newPath := path[0:len(path)-len(file+ext)] + input + ext

	if len(input) > 2 {
		os.Rename(path, newPath)
	}
}

func openImage(p string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", p).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", p).Start()
	case "darwin":
		err = exec.Command("open", p).Start()
	default:
		err = fmt.Errorf("Can't open platform")
	}

	if err != nil {
		log.Fatal(err)
	}
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var dir string
	if len(os.Args) != 2 {
		cwd, err := os.Getwd()
		handle(err)
		dir = cwd
	} else {
		dir = os.Args[1] // get the target directory
	}

	err := filepath.Walk(dir, checkDepulicates)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
