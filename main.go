package main

import (
	"github.com/kaifei-bianjie/mock/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	rootCmd.AddCommand(
		cmd.FaucetInitCmd(),
		cmd.GenSignedTxDataCmd(),
	)

	err := initFlag(rootCmd)
	if err != nil {
		panic(err)
	}

	err = rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func initFlag(rootCmd *cobra.Command) error {
	// viper bind flag
	viper.BindPFlags(rootCmd.Flags())
	for _, c := range rootCmd.Commands() {
		viper.BindPFlags(c.Flags())
	}

	//homeDir := viper.GetString(cmd.FlagHome)
	//viper.Set(cmd.FlagHome, homeDir)
	//viper.SetConfigName("config")                         // name of config file (without extension)
	//viper.AddConfigPath(homeDir)                          // search root directory
	//viper.AddConfigPath(filepath.Join(homeDir, "config")) // search root directory /config
	//
	//// If a config file is found, read it in.
	//if err := viper.ReadInConfig(); err == nil {
	//	// stderr, so if we redirect output to json file, this doesn't appear
	//	// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	//} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
	//	// ignore not found error, return other errors
	//	return err
	//}
	return nil
}