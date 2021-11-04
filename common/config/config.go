package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var Config config

type config struct {
	Redis struct {
		DBAddress     string `yaml:"dbAddress"`
		DBMaxIdle     int    `yaml:"dbMaxIdle"`
		DBMaxActive   int    `yaml:"dbMaxActive"`
		DBIdleTimeout int    `yaml:"dbIdleTimeout"`
		DBPassWord    string `yaml:"dbPassWord"`
	}

	Mysql struct {
		DBAddress      string `yaml:"dbAddress"`
		DBUserName     string `yaml:"dbUserName"`
		DBPassword     string `yaml:"dbPassword"`
		DBName         string `yaml:"dbName"`
		DBMaxOpenConns int    `yaml:"dbMaxOpenConns"`
		DBMaxIdleConns int    `yaml:"dbMaxIdleConns"`
		DBMaxLifeTime  int    `yaml:"dbMaxLifeTime"`
	}

	Alibabacloud struct {
		AccessKeyId                  string `yaml:"accessKeyId"`
		AccessKeySecret              string `yaml:"accessKeySecret"`
		SignName                     string `yaml:"SignName"`
		VerificationCodeTemplateCode string `yaml:"VerificationCodeTemplateCode"`
	}

	ImServer struct {
		IP        string `yaml:"ip"`
		Secret    string `yaml:"secret"`
		SuperCode string `yaml:"superCode"`
	}

	Log struct {
		StorageLocation       string   `yaml:"storageLocation"`
		RotationTime          int      `yaml:"rotationTime"`
		RemainRotationCount   uint     `yaml:"remainRotationCount"`
		ElasticSearchSwitch   bool     `yaml:"elasticSearchSwitch"`
		ElasticSearchAddr     []string `yaml:"elasticSearchAddr"`
		ElasticSearchUser     string   `yaml:"elasticSearchUser"`
		ElasticSearchPassword string   `yaml:"elasticSearchPassword"`
	}

	Api struct {
		Port string `yaml:"port"`
	}
}

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../")
)

func init() {
	path, _ := os.Getwd()
	fmt.Println(path + "/config/config.yaml")
	bytes, err := ioutil.ReadFile(path + "/config/config.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		panic(err)
	}
}
