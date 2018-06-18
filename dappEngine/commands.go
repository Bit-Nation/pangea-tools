package dappEngine

import (
	"context"
	"fmt"

	dapp "github.com/Bit-Nation/panthalassa/dapp"
	reg "github.com/Bit-Nation/panthalassa/dapp/registry"
	cli "github.com/urfave/cli"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	lp2p "gx/ipfs/QmWsV6kzPaYGBDVyuUfWBvyQygEc9Qrv9vzo8vZ7X4mdLN/go-libp2p"
)

var devServerAddr string

func config(cfg *lp2p.Config) error {
	addr, err := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		return err
	}
	cfg.ListenAddrs = []ma.Multiaddr{
		addr,
	}
	return lp2p.Defaults(cfg)
}

type DAppClient struct {
	startDApp func(dApp dapp.JsonRepresentation) error
}

func (c *DAppClient) HandleReceivedDApp(dApp dapp.JsonRepresentation) error {
	return c.startDApp(dApp)
}

var StartEngine = cli.Command{
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "dev-server",
			Destination: &devServerAddr,
		},
	},
	Name: "engine:start",
	Action: func(c *cli.Context) error {

		h, err := lp2p.New(context.Background(), lp2p.Defaults, config)
		if err != nil {
			panic(err)
		}

		addr, err := ma.NewMultiaddr(devServerAddr)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		client := &DAppClient{}

		registry := reg.NewDAppRegistry(h, client)

		// start the DApp when we receive it
		client.startDApp = func(dApp dapp.JsonRepresentation) error {
			return registry.StartDApp(&dApp)
		}

		if err := registry.ConnectDevelopmentServer(addr); err != nil {
			panic(err)
		}

		fmt.Println("started engine")

		select {}

		return nil
	},
}
