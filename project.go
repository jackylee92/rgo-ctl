package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
 * @Content : main
 * @Author  : LiJunDong
 * @Time    : 2022-07-14$
 */

type projectTool struct {
	config
}

const rgoTemplateUrl = "http://git.ruigushop.com/golang/rgo-template.git"

func (c *projectTool)create() (err error) {
	if err = c.check(); err != nil {
		return err
	}
	if err = c.clone(); err != nil {
		return err
	}
	if err = c.replaceName(c.pwd); err != nil {
		log.Println("err", err)
		return err
	}
	return err
}

func (c *projectTool)check()(err error){
	if err := c.checkEnv(); err != nil {
		return err
	}
	if err := c.checkEmpty(); err != nil {
		return err
	}
	return err
}

func (c *projectTool)clone() (err error){
	command := exec.Command("git", "clone", rgoTemplateUrl, c.pwd)
	command.Dir = c.pwd
	err = command.Run()
	return err
}

func (c *projectTool)replaceName(path string) (err error) {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	if f.IsDir() {
		if list, err := ioutil.ReadDir(path); err == nil {
			for _, item := range list {
				name := item.Name()
				if strings.Index(name, ".") == 0 {
					continue
				}
				if err = c.replaceName(filepath.Join(path, name)); err != nil {
					return err
				}
			}
		}
	} else {
		err = c.replaceOne(path)
		if err != nil {
			return err
		}
	}
	return err
}

func (c *projectTool)replaceOne(fileName string) (err error) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	newBody := strings.Replace(string(body), "rgo-template", c.name, -1)
	_, err = io.WriteString(f, newBody)
	if err != nil {
		return err
	}
	return err
}