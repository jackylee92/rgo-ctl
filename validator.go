package main

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

/*
 * @Content : main
 * @Author  : LiJunDong
 * @Time    : 2022-07-14$
 */

func validate(cfg *config) (err error) {
	return err
}


// checkEnv
// @Param   :
// @Return  : err error
// @Author  : LiJunDong
// @Time    : 2022-06-11
func (cfg *config) checkEnv() (err error) {
	//  <LiJunDong : 2022-06-11 15:57:09> --- go version
	v := runtime.Version()
	cfg.version = v
	vNum := v[2:]
	vArr := strings.Split(vNum, ".")
	for index, item := range vArr {
		if index == 0 && item < "1" {
			return errors.New("go版本过低，最低要求1.16")
		}
		if index == 1 && item < "16" {
			return errors.New("go版本过低，最低要求1.16")
		}
	}
	//  <LiJunDong : 2022-06-11 15:57:19> --- go module
	gomod := os.Getenv("GO111MODULE")
	if cfg.sysType == "windows" {
		gomod, err = getWinMod()
		if err != nil {
			return err
		}
	}
	cfg.module = gomod
	gomod = strings.ToLower(gomod)
	if gomod != "auto" && gomod != "on" {
		return errors.New("go环境未开启go mod，需要auto/on")
	}
	return err
}

// checkContent
// @Param   :
// @Return  : err error
// @Author  : LiJunDong
// @Time    : 2022-06-11
func (cfg *config) checkEmpty() (err error) {
	f, err := os.Open(cfg.pwd)
	if err != nil {
		return errors.New("打开" + cfg.pwd + "失败，" + err.Error())
	}

	dirs, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return err
	}
	files := make([]string, 0)
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	for _, item := range dirs {
		files = append(files, item.Name())
	}
OuterLoop:
	for key, item := range files {
		if key == 0 {
			continue
		}
		if item[0:1] == "." {
			continue
		}
		fileName := strings.ToLower(item)
		for _, allowFileName := range allowFiles {
			if fileName == allowFileName {
				continue OuterLoop
			}
		}
		return errors.New(cfg.pwd + "不为空:" + item)
	}
	return err
}
func getWinMod() (mod string, err error){
	command := exec.Command("go", "env")
	outBt, err := command.Output()
	if err != nil {
		return mod, errors.New("获取Windows go env 数据失败，" + err.Error())
	}
	err = command.Run()
	goEnvStr := string(outBt)
	goEnvArr := strings.Split(goEnvStr, "\n")
	goEnvMap := make(map[string]string,0)
	// 可以获取所有的go配置，方便以后备用
	for _, item := range goEnvArr {
		if item == "" {
			continue
		}
		cutEndIndex := strings.Index(item, "=")
		if cutEndIndex == 0 {
			continue
		}
		envField := item[4:cutEndIndex]
		envValue := item[cutEndIndex+1:]
		goEnvMap[envField] = envValue
	}
	mod, ok := goEnvMap["GO111MODULE"]
	if !ok {
		return mod, errors.New("go env中未查询到GO111MODULE配置")
	}
	return mod, nil
}

// getPwd 获取项目地址
// @Param   :
// @Return  : path string
// @Author  : LiJunDong
// @Time    : 2022-07-01
func getPwd(sysType string) (path string, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return path, errors.New("获取当前路径失败，" + err.Error())
	}
	if sysType == "windows" {
		path = pwd + "\\"
	} else {
		path = pwd + "/"
	}
	return path, err
}
