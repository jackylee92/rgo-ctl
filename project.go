package main

import (
	"bufio"
	"errors"
	"fmt"
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
	command1 := exec.Command("rm", "-rf", c.pwd + ".*")
	command1.Dir = c.pwd
	err = command1.Run()
	if err != nil {
		return err
	}
	command1.Wait()

	command2 := exec.Command("git", "clone", rgoTemplateUrl, c.pwd)
	command2.Dir = c.pwd
	err = command2.Run()
	if err != nil {
		return errors.New("git clone失败：" + err.Error())
	}
	command2.Wait()
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
		err = c.replaceFile(path)
		if err != nil {
			return err
		}
	}
	return err
}


func (c *projectTool)replaceFile(fileName string)(err error){
	content, err := c.getNewFile(fileName)
	if err != nil {
		return err
	}
	err = c.writeToFile(fileName, content)
	return err
}
func (c *projectTool)getNewFile(fileName string) (content []byte, err error) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return content, err
	}
	defer f.Close()
	content, err = ioutil.ReadAll(f)
	if err != nil {
		return content, err
	}
	return content, err
}

func (c *projectTool)writeToFile(filePath string, outPut []byte) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	_, err = writer.Write(outPut)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}