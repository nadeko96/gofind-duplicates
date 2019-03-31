package main

import (
	"bufio"
	"crypto/sha512"
	"flag"
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
var dir string
var skipRename bool

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
		if !skipRename {
			checkBadNames(path, hash)
		}

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
	reader := bufio.NewReader(os.Stdin)

	var input string
	fmt.Printf("%q is a shitty filename, open and change it? \n", file+ext)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")
	if strings.ToUpper("YES") == input {
		openImage(path)
		fmt.Printf("What is a good name for it fam? Press \"s\" to skip")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		newPath := path[0:len(path)-len(file+ext)] + input + ext

		fmt.Println(newPath)

		if len(input) > 2 {
			os.Rename(path, newPath)
		}
	}

	fmt.Println()
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

	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	flag.StringVar(&dir, "p", "current directory", "path to search")
	flag.BoolVar(&skipRename, "s", false, "Skips renaming files")
	flag.Parse()
	if dir == "current directory" {
		cwd, err := os.Getwd()
		handle(err)
		dir = cwd
	}

	err := filepath.Walk(dir, checkDepulicates)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
