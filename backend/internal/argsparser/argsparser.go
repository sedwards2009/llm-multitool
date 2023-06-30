package argsparser

import "flag"

type CommandLineArguments struct {
	ConfigFilePath *string
}

func Parse(args *[]string) (parsed *CommandLineArguments, errString string) {
	result := &CommandLineArguments{ConfigFilePath: nil}

	var configFlag = flag.String("config", "", "Path to the configuration file.")
	flag.Parse()

	result.ConfigFilePath = configFlag

	return result, ""
}
