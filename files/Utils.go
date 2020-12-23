package files

import (
	"bytes"
	"fmt"
	"github.com/saichler/syncit/model"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

var timestamp = time.Now().Unix()
var MIN_FILE_SIZE = int64(1024 * 1024 * 200)

func Scan(filename string) *model.File {
	s, e := os.Stat(filename)
	if e != nil {
		return nil
	}
	root := &model.File{}
	root.Name = filename
	timestamp = time.Now().Unix()
	if !s.IsDir() {
		root.Size = s.Size()
	} else {
		files, e := ioutil.ReadDir(filename)
		if e != nil {
			return root
		}
		for _, file := range files {
			seek(root, getFilename(filename, file.Name()))
		}
	}
	return root
}

func getFilename(path, name string) string {
	buff := &bytes.Buffer{}
	buff.WriteString(path)
	buff.WriteString("/")
	buff.WriteString(name)
	return buff.String()
}

func seek(parent *model.File, filename string) {
	if time.Now().Unix()-timestamp > 1 {
		fmt.Println("In dir:" + parent.Name)
		timestamp = time.Now().Unix()
	}
	s, e := os.Stat(filename)
	if e != nil {
		return
	}

	if s.Mode().IsDir() {
		ls, _ := os.Lstat(filename)
		if ls.Mode()&os.ModeSymlink != 0 {
			return
		}
	}

	fe := &model.File{}
	fe.Name = filename
	if parent.Files == nil {
		parent.Files = make([]*model.File, 0)
	}
	if !s.IsDir() {
		fe.Size = s.Size()
		fe.Date = s.ModTime().Unix()
		parent.Size += fe.Size
	} else {
		files, e := ioutil.ReadDir(filename)
		if e != nil {
			return
		}
		for _, file := range files {
			seek(fe, getFilename(filename, file.Name()))
		}
		sort.Slice(parent.Files, func(i, j int) bool {
			return parent.Files[i].Size > parent.Files[j].Size
		})
		parent.Size += fe.Size
	}
	parent.Files = append(parent.Files, fe)
}

func Print(fe *model.File, dept int, incFile, incLessThanBlock bool) {
	print(fe, 0, dept, incFile, incLessThanBlock)
}

func print(fe *model.File, lvl, dept int, incFiles, incLessThanBlock bool) {
	if lvl > dept {
		return
	}
	if fe.Files == nil && !incFiles {
		return
	}
	if fe.Size < MIN_FILE_SIZE && !incLessThanBlock {
		return
	}
	buff := bytes.Buffer{}
	buff.WriteString(fe.Name)
	sizeStr := sizeIt(fe.Size)
	buff.WriteString(" ")
	buff.WriteString(sizeStr)
	fmt.Println(buff.String())
	if fe.Files != nil {
		for _, child := range fe.Files {
			lvl++
			print(child, lvl, dept, incFiles, incLessThanBlock)
			lvl--
		}
	}
}

func sizeIt(size int64) string {
	buff := bytes.Buffer{}
	if size/1024 == 0 {
		buff.WriteString(strconv.Itoa(int(size)))
		buff.WriteString("b")
	} else if size/1024/1024 == 0 {
		buff.WriteString(strconv.Itoa(int(size / 1024)))
		buff.WriteString("k")
	} else if size/1024/1024/1024 == 0 {
		buff.WriteString(strconv.Itoa(int(size / 1024 / 1024)))
		buff.WriteString("m")
	} else {
		buff.WriteString(strconv.Itoa(int(size / 1024 / 1024 / 1024)))
		buff.WriteString("g")
	}
	return buff.String()
}
