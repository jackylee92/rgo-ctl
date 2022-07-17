package main

import (
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
	if err = c.mv(); err != nil {
		log.Println("err", err)
		return err
	}
	c.cleanTemplate()
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
	command2 := exec.Command("git", "clone", rgoTemplateUrl)
	command2.Dir = c.pwd
	err = command2.Run()
	if err != nil {
		return errors.New("git clone失败：" + err.Error())
	}
	command2.Wait()
	return err
}

func (c *projectTool)cleanTemplate() (err error){
	command2 := exec.Command("rm", "-rf", c.pwd + "rgo-template")
	command2.Dir = c.pwd
	err = command2.Run()
	if err != nil {
		return errors.New("git clone失败：" + err.Error())
	}
	command2.Wait()
	return err
}

// mv
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-11
func (cfg *config) mv() (err error) {
	from := cfg.pwd + "rgo-template"
	to := cfg.pwd
	if cfg.sysType == "windows" {
		from = "\\util\\rgtemplate\\code"
	}
	err = copy(cfg.projectName, from, to)
	if err != nil {
		return errors.New("移动模版文件失败，" + err.Error())
	}
	return err
}

// <LiJunDong : 2022-06-12 00:18:32> --- 移动的时候需要将文件后缀tmp去掉 将文件内容中项目名替换
func copy(project, from, to string) error {
	fmt.Println("from:", from, strings.Index(from, ".git/"))
	if strings.Index(from, ".git/") != -1 {
		return nil
	}
	f, e := os.Stat(from)
	if e != nil {
		return e
	}
	if f.IsDir() {
		//from是文件夹，那么定义to也是文件夹
		if list, e := ioutil.ReadDir(from); e == nil {
			for _, item := range list {
				if e = copy(project, filepath.Join(from, item.Name()), filepath.Join(to, item.Name())); e != nil {
					return e
				}
			}
		}
	} else {
		//from是文件，那么创建to的文件夹
		p := filepath.Dir(to)
		if _, e = os.Stat(p); e != nil {
			if e = os.MkdirAll(p, 0777); e != nil {
				return e
			}
		}
		//读取源文件
		file, e := os.Open(from)
		if e != nil {
			return e
		}
		defer file.Close()
		fileContent, e := ioutil.ReadAll(file)
		if e != nil {
			return e
		}
		newFileContent := strings.Replace(string(fileContent), "rgo-template", project, -1)
		out, e := os.Create(to)
		if e != nil {
			return e
		}
		defer out.Close()
		_, e = out.WriteString(newFileContent)
	}
	return e
}