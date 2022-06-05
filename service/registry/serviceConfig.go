package registry

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Config struct {
	*RedisConf `yaml:"Redis"`
}

type RedisConf struct {
	Addr     string `yaml:"Addr"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

var (
	config *Config
	rds    *redis.Client
	ctx    context.Context
	once   sync.Once
)

func InitRegistry(yamlPath string) {
	once.Do(func() {
		initConfig(yamlPath)
		initRedis()
	})
}

func initConfig(yamlPath string) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func initRedis() {
	rds = redis.NewClient(&redis.Options{
		Addr:     config.RedisConf.Addr,
		Password: config.RedisConf.Password, // no password set
		DB:       config.RedisConf.DB,       // use default DB
	})
	ctx = context.Background()
	val, err := rds.Get(ctx, "key").Result()
	if err == redis.Nil {
		fmt.Println("服务注册Redis测试： key does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("服务注册Redis测试： key", val)
	}
}
