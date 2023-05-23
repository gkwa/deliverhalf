package main

import (
	"github.com/taylormonacelli/deliverhalf/cmd"
	_ "github.com/taylormonacelli/deliverhalf/cmd/client"
	_ "github.com/taylormonacelli/deliverhalf/cmd/client/send"
	_ "github.com/taylormonacelli/deliverhalf/cmd/common"
	_ "github.com/taylormonacelli/deliverhalf/cmd/config"
	_ "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	_ "github.com/taylormonacelli/deliverhalf/cmd/ec2/ami"
	_ "github.com/taylormonacelli/deliverhalf/cmd/ec2/instance"
	_ "github.com/taylormonacelli/deliverhalf/cmd/ec2/launchtemplate"
	_ "github.com/taylormonacelli/deliverhalf/cmd/ec2/launchtemplate/test"
	_ "github.com/taylormonacelli/deliverhalf/cmd/ec2/volume"
	_ "github.com/taylormonacelli/deliverhalf/cmd/ec2/volume/test"
	_ "github.com/taylormonacelli/deliverhalf/cmd/logs"
	_ "github.com/taylormonacelli/deliverhalf/cmd/meta"
	_ "github.com/taylormonacelli/deliverhalf/cmd/sns"
	_ "github.com/taylormonacelli/deliverhalf/cmd/update"
	_ "github.com/taylormonacelli/deliverhalf/cmd/watchdog"
)

func main() {
	cmd.Execute()
}
