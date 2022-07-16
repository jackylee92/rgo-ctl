package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
	fmt.Println("验证成功...")
	if err = c.clone(); err != nil {
		return err
	}
	fmt.Println("clone成功...")
	if err = c.replaceName(c.pwd); err != nil {
		log.Println("err", err)
		return err
	}
	fmt.Println("初始化成功...")
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
	fmt.Println("准备clone>", "git clone " + rgoTemplateUrl + " " + c.pwd)
	command.Dir = c.pwd
	fmt.Println("准备clone>>", "git clone " + rgoTemplateUrl + " " + c.pwd)
	err = command.Run()
	fmt.Println("准备clone>>>", "git clone " + rgoTemplateUrl + " " + c.pwd)
	if err == nil {
		fmt.Println("准备clone>>>>", "git clone " + rgoTemplateUrl + " " + c.pwd)
		_ = command.Wait()
	}
	fmt.Println("准备clone>>>>>", "git clone " + rgoTemplateUrl + " " + c.pwd)
	time.Sleep(20*time.Second)
	return errors.New("git clone失败：" + err.Error())
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