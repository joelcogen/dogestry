package config

import (
  "code.google.com/p/gcfg"
)

type RemoteConfig struct {
  Url string
}

type S3Config struct {
  Access_Key_Id string
  Secret_Key string
}

type CompressorConfig struct {
  Lz4 string
}

type DockerConfig struct {
  Connection string
}

type Config struct {
  Remote map[string]*RemoteConfig
  S3 S3Config
  Compressor CompressorConfig
  Docker DockerConfig
}


func ParseConfig(configFilePath string) (config Config, err error) {
  err = gcfg.ReadFileInto(&config, configFilePath)
  return
}
