package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

type config struct {
	action      string
	name        string
	projectName string
	pwd         string
	sysType     string
	version     string
	module      string
}

var allowFiles = []string{
	"readme.md",
	"readme",
}

func main() {
	if err := start(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func start() (err error) {
	cfg, err := parseConfig(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}
	if err := cfg.checkEnv(); err != nil {
		return err
	}
	if err := cfg.switchDo(); err != nil {
		return err
	}
	cfg.outMessage()
	return err
}

func parseConfig(args []string) (*config, error) {
	var (
		flagAction = flag.String("tool", "", "执行动作")
		flagName   = flag.String("name", "", "动作参数")
	)
	if err := flag.CommandLine.Parse(args); err != nil {
		return nil, err
	}

	if flag.NFlag() == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return nil, flag.ErrHelp
	}
	sysType := runtime.GOOS
	pwd, err := getPwd(sysType)
	if err != nil {
		return nil, err
	}
	fileInfo, err := os.Stat(pwd)
	if err != nil {
		return nil, err
	}
	dirName := fileInfo.Name()
	cfg := &config{
		action:      *flagAction,
		name:        *flagName,
		projectName: dirName,
		sysType:     runtime.GOOS,
		pwd:         pwd,
	}
	return cfg, nil
}

func (cfg *config) switchDo() (err error) {
	switch cfg.action {
	case "init": // 初始化项目
		client := projectTool{*cfg}
		if err = client.create(); err != nil {
			return err
		}
	case "check":
	case "createMysql":
	}
	return err
}

func (cfg *config) outMessage() {
	fmt.Println("项目创建中: ")
	fmt.Println("    系统:" + cfg.sysType)
	fmt.Println("    go 版本:" + cfg.version)
	fmt.Println("    go Mod:" + cfg.module)
	fmt.Println("    项目名:" + cfg.name)
}
