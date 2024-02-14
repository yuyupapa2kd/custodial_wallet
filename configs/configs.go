package configs

import (
	"math/big"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var RuntimeConf = RuntimeConfig{}

type RuntimeConfig struct {
	RpcEndpoint string      `yaml:"rpcEndpoint"`
	Server      Server      `yaml:"server"`
	EvmNetIDStr EvmNetIDStr `yaml:"evmNetIDStr"`
	EvmNetID    EvmNetID    `yaml:"evmNetID"`
	LogPath     LogPath     `yaml:"logPath"`
	WhiteIP     string      `yaml:"whiteIP"`
}

type Server struct {
	Ip   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type EvmNetIDStr struct {
	Eth string `yaml:"eth"`
	Gnd string `yaml:"gnd"`
}

type EvmNetID struct {
	Eth *big.Int `yaml:"eth"`
	Gnd *big.Int `yaml:"gnd"`
}

type LogPath struct {
	Info  string
	Error string
	Debug string
}

func SetRuntimeConfig(profile string) error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName(profile)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "viper read in config failed")
	}

	err = viper.Unmarshal(&RuntimeConf)
	if err != nil {
		return errors.Wrap(err, "viper unmarshal to RuntimeConf failed")
	}

	ethBigInt := new(big.Int)
	ethBigInt, ok := ethBigInt.SetString(RuntimeConf.EvmNetIDStr.Eth, 10)
	if !ok {
		return errors.Wrap(errors.New("convert failed"), "convert string to bigint failed for eth")
	}
	RuntimeConf.EvmNetID.Eth = ethBigInt

	gndBigInt := new(big.Int)
	gndBigInt, ok = ethBigInt.SetString(RuntimeConf.EvmNetIDStr.Gnd, 10)
	if !ok {
		return errors.Wrap(errors.New("convert failed"), "convert string to bigint failed for gnd")
	}
	RuntimeConf.EvmNetID.Gnd = gndBigInt

	return nil
}
