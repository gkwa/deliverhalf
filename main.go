/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/taylormonacelli/deliverhalf/cmd"
	_ "github.com/taylormonacelli/deliverhalf/cmd/common"
	_ "github.com/taylormonacelli/deliverhalf/cmd/config"
	_ "github.com/taylormonacelli/deliverhalf/cmd/ec2"
	_ "github.com/taylormonacelli/deliverhalf/cmd/ec2/volume"
	_ "github.com/taylormonacelli/deliverhalf/cmd/meta"
	_ "github.com/taylormonacelli/deliverhalf/cmd/sns"
)

func main() {
	cmd.Execute()
}
