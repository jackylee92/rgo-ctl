package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type config struct {
	name        string
	git         string
	projectPath string
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
	err = validate(cfg)
	if err != nil {
		return err
	}
	fmt.Println("项目创建中: ")
	fmt.Println("    系统:" + cfg.sysType)
	fmt.Println("    go 版本:" + cfg.version)
	fmt.Println("    go Mod:" + cfg.module)
	fmt.Println("    项目名:" + cfg.name)
	fmt.Println("    Git:" + cfg.git)
	err = cfg.mv()
	if err == nil {
		fmt.Println("    项目创建完成。。。")
	}
	return err
}

func parseConfig(args []string) (*config, error) {
	var (
		flagName = flag.String("name", "", "指定项目名称，确保app下有该项目")
		//flagGit = flag.String("git", "", "指定项目git地址")
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
	projectPath, err := os.Getwd()
	if err != nil {
		return nil, errors.New("获取当前路径失败，" + err.Error())
	}
	if sysType == "windows" {
		projectPath = projectPath + "\\app\\" + *flagName
	} else {
		projectPath = projectPath + "/app/" + *flagName
	}
	cfg := &config{
		name:        *flagName,
		sysType:     runtime.GOOS,
		projectPath: projectPath,
		//git: *flagGit
	}
	return cfg, nil
}

func validate(cfg *config) (err error) {
	// <LiJunDong : 2022-06-11 15:21:39> --- go环境检测 go module go版本
	if err := cfg.checkEnv(); err != nil {
		return err
	}
	// <LiJunDong : 2022-06-11 15:21:23> --- 必须在rgo.com目录下面
	//if err := cfg.checkDir(); err != nil {
	//	return err
	//}
	// <LiJunDong : 2022-06-11 15:22:46> --- 项目文件夹是否存在，是否已存在内容
	if err := cfg.checkContent(); err != nil {
		return err
	}

	// <LiJunDong : 2022-06-11 15:23:32> --- git地址是否正确 git submodule 会将子项目仓库自动克隆下来，所以不需要地址
	return err
}

func (cfg *config) create() (err error) {
	return err
}

// clone
// @Param   :
// @Return  : err error
// @Author  : LiJunDong
// @Time    : 2022-06-11
func (cfg *config) clone() (err error) {
	return err
}

// mv
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-11
func (cfg *config) mv() (err error) {
	from := "./util/rgtemplate/code"
	if cfg.sysType == "windows" {
		from = "\\util\\rgtemplate\\code"
	}
	err = copy(cfg.name, from, cfg.projectPath)
	if err != nil {
		return errors.New("移动模版文件失败，" + err.Error())
	}
	return err
}

// replace
// @Param   :
// @Return  :
// @Author  : LiJunDong
// @Time    : 2022-06-11
func (cfg *config) replace() (err error) {
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
	cfg.module = gomod
	gomod = strings.ToLower(gomod)
	if gomod != "auto" && gomod != "on" {
		return errors.New("go环境未开启go mod，需要auto/on")
	}
	return err
}

// checkDir
// @Param   :
// @Return  : err error
// @Author  : LiJunDong
// @Time    : 2022-06-11
func (cfg *config) checkDir() (err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return errors.New("检查路径失败，" + err.Error())
	}
	fileInfo, err := os.Stat(pwd)
	dirName := fileInfo.Name()
	if dirName != "rgo.com" {
		return errors.New("执行位置错误，仅可在rgo.com目录下执行")
	}
	return err
}

// checkContent
// @Param   :
// @Return  : err error
// @Author  : LiJunDong
// @Time    : 2022-06-11
func (cfg *config) checkContent() (err error) {
	dir, err := os.Stat(cfg.projectPath)
	if os.IsNotExist(err) {
		return errors.New(cfg.projectPath + "不存在，请联系项目管理员创建项目")
	}
	if err != nil {
		return errors.New("读取" + cfg.projectPath + "项目路径失败")
	}
	if !dir.IsDir() {
		return errors.New(cfg.projectPath + "不存在，请联系项目管理员创建项目")
	}
	gitModulesPath := "./.gitmodules"
	if cfg.sysType == "windows" {
		gitModulesPath = ".\\.gitmodules"
	}
	contentBt, err := ioutil.ReadFile(gitModulesPath)
	if err != nil {
		return errors.New("读取.gitmodules文件失败" + err.Error())
	}
	contentArr := strings.Split(string(contentBt), "\n")
	if len(contentArr) == 0 {
		return errors.New("读取.gitmodules文件失败" + err.Error())
	}
	startIndex := 0
	for key, item := range contentArr {
		if item == "[submodule \"app/"+cfg.name+"\"]" {
			startIndex = key
		}
	}
	if startIndex+2 >= len(contentArr) {
		return errors.New("gitmodules中未配置" + cfg.name + "项目")
	}
	gitPathTmp := contentArr[startIndex+2]
	gitPathTmp = strings.ReplaceAll(gitPathTmp, " ", "")
	gitPathTmp = strings.ReplaceAll(gitPathTmp, "	", "")
	if len(gitPathTmp) < 4 {
		return errors.New("gitmodules中读取" + cfg.name + " git失败")
	}
	cfg.git = gitPathTmp[4:]
	files := make([]string, 0)



	f, err := os.Open(cfg.projectPath)
	if err != nil {
		return err
	}
	dirs, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return err
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	for _, item := range dirs {
		files = append(files, item.Name())
	}
	hasMod := false
	if len(files) == 1 {
		command := exec.Command("git clone " + cfg.git + " ./")
		command.Dir = cfg.projectPath
		command.Run()
		command = exec.Command("go mod init")
		command.Dir = cfg.projectPath
		command.Run()
	} else {
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
			if item == "go.mod" {
				hasMod = true
			}
			return errors.New(cfg.projectPath + "不为空" + item)
		}
	}
	if !hasMod {
		command := exec.Command("go", "mod", "init")
		command.Dir = cfg.projectPath
		err = command.Run()
		if err != nil {
			return err
		}
	}
	return err
}

// <LiJunDong : 2022-06-12 00:18:32> --- 移动的时候需要将文件后缀tmp去掉 将文件内容中项目名替换
func copy(project, from, to string) error {
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
		newFileContent := strings.Replace(string(fileContent), "rgtemplate", project, -1)
		// 创建一个文件用于保存
		if strings.ToLower(to[len(to)-4:]) == ".tmp" {
			to = to[:len(to)-4]
		}
		out, e := os.Create(to)
		if e != nil {
			return e
		}
		defer out.Close()
		_, e = out.WriteString(newFileContent)
	}
	return e
}
