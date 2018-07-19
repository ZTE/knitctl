package config

import (
	"fmt"
	"github.com/ZTE/knitctl/api"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

type configOptions struct {
	yamlFile   string
	knitterurl string

	Out    io.Writer
	ErrOut io.Writer
}

var (
	configShort = ("config knitter server from a file or option.")

	configLong = (`  Config knitter server ip and port.
YAML format are accepted.`)

	configExample = (`  # Config knitter server ip and port with command.
  knitctl config --knitterurl=ip:port

  # Config knitter server ip and port with yaml file.
  knitctl config -f xx.yaml`)
)

func NewCommandConfig(out io.Writer, errOut io.Writer) *cobra.Command {
	o := configOptions{
		Out:    out,
		ErrOut: errOut,
	}

	cmd := &cobra.Command{
		Use: "config [flags]",
		DisableFlagsInUseLine: true,
		Short:   configShort,
		Long:    configLong,
		Example: configExample,
		Run: func(cmd *cobra.Command, args []string) {
			api.ExcuteError(o.RunConfig(cmd, args), errOut)
		},
	}

	cmd.Flags().StringVarP(&o.yamlFile, "filename", "f", o.yamlFile, "If present, config knitter server from this yaml file")
	cmd.Flags().StringVarP(&o.knitterurl, "knitterurl", "", o.knitterurl, "If present, config knitter server from this option.")

	return cmd
}

func (o *configOptions) RunConfig(cmd *cobra.Command, args []string) error {
	if api.ContainsElem(args, "options") {
		cmd.HelpFunc()(cmd, args)
		return nil
	}
	if len(args) > 0 {
		return fmt.Errorf("invalid argument to config knitter server.\n " +
			`Use "knitctl config --help" for more information about a given command.`)
	}
	configFile := o.yamlFile
	knitterurl := o.knitterurl
	if configFile == "" && knitterurl == "" {
		return fmt.Errorf("no found options filename or knitterurl.\n" +
			"run 'knitctl config --filename=xx.yaml or knitctl config --knitterurl=ip:port'")
	}

	knitterurlValue, err := o.getKnitterurlFromUrl()
	if err != nil {
		return err
	}
	if knitterurlValue == "" {
		knitterurlValue, _ = o.getKnitterurlFromOption()
	}
	if knitterurlValue == "" {
		return fmt.Errorf("no found knitterurl from config yaml file or knitterurl option.\n" +
			"run 'knitctl config --filename=xx.yaml or knitctl config --knitterurl=ip:port' to config knitter server")
	}

	mapConfig := make(map[string]string)
	if !strings.Contains(knitterurlValue, "http") {
		knitterurlValue = "http://" + knitterurlValue
	}
	mapConfig["knitterurl"] = knitterurlValue
	if err := api.ModifyConfigFile(mapConfig); err != nil {
		return err
	}
	strRe := fmt.Sprintf("config knitter server success: %s.\n", mapConfig)
	o.Out.Write([]byte(strRe))

	return nil
}

func (o *configOptions) getKnitterurlFromUrl() (string, error) {
	configFile := o.yamlFile
	if configFile != "" {
		if find, _ := api.FileExists(configFile); !find {
			return "", fmt.Errorf("no such file or directory: %s", configFile)
		} else {
			mapResult, err := api.GetFromYaml(configFile)
			if err != nil {
				return "", err
			}
			if len(mapResult) > 0 {
				if value, found := mapResult[0]["knitterurl"]; found {
					return value.(string), nil
				}
			}
		}
		return "", fmt.Errorf("no found knitterurl from file %s", configFile)
	}
	return "", nil
}

func (o *configOptions) getKnitterurlFromOption() (string, error) {
	return o.knitterurl, nil
}
