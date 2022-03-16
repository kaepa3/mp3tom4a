package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Print(err)
		return
	}
	done := make(chan struct{})
	fileChan := make(chan string)
	go dirwalk(done, fileChan, dir)
	wg := &sync.WaitGroup{}
Loop:
	for {
		select {
		case fileName := <-fileChan:
			use := filter(fileName)
			if use {
				log.Print("convert!!! " + fileName)
				wg.Add(1)
				convert(fileName)
				wg.Done()
			}
		case <-done:
			break Loop
		}
	}
	wg.Wait()
	log.Println("end convert")
}

func convert(src string) {
	e := getFileNameWithoutExt(src)
	dst := e + ".m4a"
	err := exec.Command(
		"ffmpeg", "-i", src,
		"-acodec", "aac",
		"-ab", "64k",
		dst,
	).Run()
	if err != nil {
		log.Print("err: " + err.Error())
	}
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func filter(file string) bool {
	e := filepath.Ext(file)
	if e == ".mp3" {
		return true
	}
	return false
}

func dirwalk(done chan<- struct{}, fileChan chan<- string, dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileChan <- file.Name()
	}
	close(done)
}
