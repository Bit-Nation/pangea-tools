package main

import (
	"github.com/Bit-Nation/pangea-tools/dappDevelopment"
	"github.com/Bit-Nation/pangea-tools/dappEngine"
	"github.com/Bit-Nation/pangea-tools/signingKey"
	"github.com/urfave/cli"
	"os"
)

func main() {

	app := cli.NewApp()

	app.Commands = []cli.Command{
		signingKey.KeyNew,
		dappEngine.StartEngine,
		dappDevelopment.DAppStream,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}

}
