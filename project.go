package main

import (
	"io"
	"io/ioutil"
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
	command := exec.Command("git clone " + rgoTemplateUrl + " " + c.pwd)
	command.Dir = c.pwd
	err = command.Run()
	return err
}

func (c *projectTool)replaceName(path string) (err error) {
	f, e := os.Stat(path)
	if e != nil {
		return e
	}
	if f.IsDir() {
		if list, err := ioutil.ReadDir(path); e == nil {
			for _, item := range list {
				if err = c.replaceName(filepath.Join(path, item.Name())); err != nil {
					return e
				}
			}
		}
	} else {
		err = c.replaceOne(path)
		if err != nil {
			return err
		}
	}
	return e
}

func (c *projectTool)replaceOne(fileName string) (err error) {
	f, err := os.OpenFile(fileName, os.O_APPEND, 0666)
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