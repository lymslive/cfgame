package cmdline

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// 默认服务器地址
const (
	defHost string = "localhost"
	defPort int    = 8000
)

// 配置的数据结构定义
type ConfigT struct {
	WebDir  string
	CfgFile string
	Server  ServerCfg
}

type ServerCfg struct {
	Host   string
	Port   int
	Static []string
}

// 定一个全局的配置变量
var config *ConfigT = new(ConfigT)
var parsed = false

func GetConfig() *ConfigT {
	if !parsed {
		return nil
	}
	return config
}

// 读取 toml 配置文件
func loadConfig(file string) error {
	if _, err := toml.DecodeFile(file, config); err != nil {
		log.Fatalln("fail to load config file:", file, err)
		return err
	}
	return nil
}

// 解析命令行及配置文件
func ParseConfig() *ConfigT {
	// 解析命令行参数
	var host *string = flag.String("host", "", "server host, ip")
	var port *int = flag.Int("port", 0, "the default port")
	var dir *string = flag.String("dir", "", "web server base dir")
	var conf *string = flag.String("conf", "conf.toml", "config file")
	flag.Parse()

	fmt.Printf("Arg Parsed: host=%s, port=%d, dir=%s, conf=%s\n",
		*host, *port, *dir, *conf)

	if *dir == "" {
		*dir, _ = os.Getwd()
	}

	path := *dir + string(os.PathSeparator) + *conf
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalln("config file not exist: ", path, err)
	}

	// 解析配置文件
	loadConfig(path)

	// 解析 toml 时 unmarshal 会初始化目标数据 config
	// 故要在解析后再根据命令行参数覆盖某些值

	config.WebDir = *dir
	config.CfgFile = *conf
	// Lookup 并不能判断用户未是否输入了某选项
	if flag.Lookup("host") != nil && *host != "" {
		fmt.Println("host found, overrides toml", config.Server.Host)
		config.Server.Host = *host
	}
	if flag.Lookup("port") != nil && *port != 0 {
		fmt.Println("port found, overrides toml", config.Server.Port)
		config.Server.Port = *port
	}

	// 如果地址依然零值，使用默认地址
	if config.Server.Host == "" {
		config.Server.Host = defHost
	}
	if config.Server.Port == 0 {
		config.Server.Port = defPort
	}

	parsed = true
	return config
}
