package config

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const (
	FlagConfig = "config"

	Localhost         = "127.0.0.1"
	DefaultShard0Port = 26780
	DefaultShard1Port = 26781
	DefaultServerPort = 19021

	DefaultInitBlock       = "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"
	DefaultInitBlockParent = "000000000000000000925634d697d3dcd7a8f5aef312f043f4cb278fd9152baa"
	DefaultInitBlockHeight = 0
)

type Config struct {
	NodeHost string `mapstructure:"NODE_HOST"`

	InitBlock       string `mapstructure:"INIT_BLOCK"`
	InitBlockHeight uint   `mapstructure:"INIT_BLOCK_HEIGHT"`
	InitBlockParent string `mapstructure:"INIT_BLOCK_PARENT"`

	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort int    `mapstructure:"SERVER_PORT"`

	QueueShards []Shard `mapstructure:"QUEUE_SHARDS"`

	SaveMetrics bool `mapstructure:"SAVE_METRICS"`

	DataPrefix             string `mapstructure:"DATA_PREFIX"`
	OpenFilesCacheCapacity int    `mapstructure:"OPEN_FILES_CACHE_CAPACITY"`

	ProcessLimit struct {
		Utxos int `mapstructure:"UTXOS"`
	} `mapstructure:"PROCESS_LIMIT"`
}

var _config Config

var DefaultConfig = Config{
	NodeHost:        GetHost(8333),
	InitBlock:       DefaultInitBlock,
	InitBlockHeight: DefaultInitBlockHeight,
	InitBlockParent: DefaultInitBlockParent,
	ServerHost:      Localhost,
	ServerPort:      DefaultServerPort,
	QueueShards: []Shard{{
		Total: 2,
		Host:  Localhost,
		Port:  DefaultShard0Port,
	}, {
		Total: 2,
		Host:  Localhost,
		Port:  DefaultShard1Port,
	}},
}

const (
	NotFoundErrorMessage = "Config file not found"
)

func IsConfigNotFoundError(err error) bool {
	return jerr.HasError(err, NotFoundErrorMessage)
}

func Init(cmd *cobra.Command) error {
	config, _ := cmd.Flags().GetString(FlagConfig)
	if config != "" && !strings.HasPrefix(config, "config-") {
		config = "config-" + config
	} else if config == "" {
		config = "config"
	}
	viper.SetConfigName(config)
	viper.AddConfigPath("$HOME/.memo-server")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".config/memo")
	if err := viper.ReadInConfig(); err != nil {
		// Config not found, use default
		_config = DefaultConfig
		return nil
	}
	if err := viper.Unmarshal(&_config); err != nil {
		return jerr.Get("error unmarshalling config", err)
	}
	return nil
}

func GetNodeHost() string {
	return _config.NodeHost
}

func GetInitBlock() string {
	return _config.InitBlock
}

func GetInitBlockHeight() uint {
	return _config.InitBlockHeight
}

func GetInitBlockParent() string {
	return _config.InitBlockParent
}

func GetQueueShards() []Shard {
	return _config.QueueShards
}

func GetSaveMetrics() bool {
	return _config.SaveMetrics
}

func GetServerPort() int {
	return _config.ServerPort
}

func GetProcessLimitUtxos() int {
	return _config.ProcessLimit.Utxos
}

func GetSelfRpc() RpcConfig {
	var host = _config.ServerHost
	if host == "" {
		host = Localhost
	}
	return RpcConfig{
		Host: host,
		Port: _config.ServerPort,
	}
}

func GetDataPrefix() string {
	return _config.DataPrefix
}

func GetOpenFilesCacheCapacity() int {
	return _config.OpenFilesCacheCapacity
}

func GetHost(port uint) string {
	return fmt.Sprintf("%s:%d", Localhost, port)
}
