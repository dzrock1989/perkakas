package configs

import (
	"context"
	"encoding/json"

	vault "github.com/hashicorp/vault/api"
	"github.com/tigapilarmandiri/perkakas/common/util"
)

type Vault struct {
	Address     string `json:"-"`
	Token       string `json:"-"`
	ServiceName string `json:"-"`
	SecretName  string `json:"-"`
}

func (vc Vault) load() (confs Configs) {
	config := vault.DefaultConfig()
	config.Address = vc.Address

	client, err := vault.NewClient(config)
	if err != nil {
		util.Log.Fatal().Msg(err.Error())
	}

	client.SetToken(vc.Token)

	ctx := context.Background()

	secret, err := client.KVv2(vc.ServiceName).Get(ctx, vc.SecretName)
	if err != nil {
		util.Log.Fatal().Msg(err.Error())
	}

	b, err := json.Marshal(secret)
	if err != nil {
		util.Log.Fatal().Msg(err.Error())
	}

	if err = json.Unmarshal(b, &confs); err != nil {
		util.Log.Fatal().Msg(err.Error())
	}

	return
}
