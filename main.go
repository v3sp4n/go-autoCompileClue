
package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "regexp"
    "os/exec"
    "time"
    "strings"
)

var PATH string

func main() {
	PATH = `D:\gta_sa_sborki\main\moonloader\clue`
	clues := map[string]map[string]string{}

	fmt.Println("start")

	clock := time.Now().UnixMilli()
	for {
		if time.Now().UnixMilli()-clock >= 100 {
			clock = time.Now().UnixMilli()

			dir, _ := ioutil.ReadDir(PATH)
			for _,f := range dir {
				if f.IsDir() {
					dir, _ := ioutil.ReadDir(PATH+"\\"+f.Name()+"\\")
					nameDir := f.Name()
					if _, ok := clues[nameDir]; !ok {
						fmt.Println("new folder",nameDir)
						clues[nameDir] = map[string]string{}
					}
					for _, f := range dir {
						lua, _ := ioutil.ReadFile(PATH+"\\"+nameDir+"\\"+f.Name())
						if isClue, _ := regexp.Match(".clue",[]byte(f.Name())); isClue {
							if _, ok := clues[nameDir][f.Name()]; !ok {
								clues[nameDir][f.Name()] = string(lua)
								fmt.Println(">>new file",nameDir,f.Name())
							} else if clues[nameDir][f.Name()] != string(lua) {
								fmt.Println(f.Name())
								clues[nameDir][f.Name()] = string(lua)
								compile(nameDir)
							}
						}
					}
				} else {
					lua, _ := ioutil.ReadFile(PATH+"\\"+f.Name())
					if _, ok := clues["."]; !ok {
						clues["."] = map[string]string{}//root
					}
					if isClue, _ := regexp.Match(".clue",[]byte(f.Name())); isClue {
						if _, ok := clues["."][f.Name()]; !ok {
							clues["."][f.Name()] = string(lua)
						} else if clues["."][f.Name()] != string(lua) {
							clues["."][f.Name()] = string(lua)
							compile(f.Name())
						}
					}
				}
			}
		}
	}

}
func compile(who string) {
	// fmt.Println("[compile]",who)
	fmt.Println("[compile]",who)
	os.Remove(PATH+who+".lua")

	cmd := exec.Command("clue","-t=luajit",who,who+".lua")
	cmd.Dir = PATH
	cmd.Start()

	for {
		if _,err := os.Stat(PATH+who+".lua"); err != nil {
			break
		}
	}

	newLua,_ := ioutil.ReadFile(PATH+"\\"+who+".lua")
	ioutil.WriteFile(strings.TrimRight(PATH,"clue") + who+".lua", newLua, 0644)

}