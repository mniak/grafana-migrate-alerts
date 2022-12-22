package main

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "migalert",
	Run: func(cmd *cobra.Command, args []string) {
		grafanaBaseURL := lo.Must1(cmd.Flags().GetString("base-url"))
		grafanaSession := lo.Must1(cmd.Flags().GetString("session"))
		grafanaIncludeAll := lo.Must1(cmd.Flags().GetBool("include-all"))
		groupSuffix := lo.Must1(cmd.Flags().GetString("group-suffix"))
		ruleSuffix := lo.Must1(cmd.Flags().GetString("rule-suffix"))

		grafanaFolders := lo.Must1(GetGrafanaRules(grafanaBaseURL, grafanaSession))
		gen := TerraformGenerator{
			IncludeAPIGroups:    grafanaIncludeAll,
			RuleGroupNameSuffix: groupSuffix,
			RuleTitleSuffix:     ruleSuffix,
		}

		tfFiles := lo.Must1(gen.GenerateFiles(grafanaFolders))
		for fileName, tfFile := range tfFiles {
			file := lo.Must1(os.Create(fmt.Sprintf("%s.tf", fileName)))
			defer file.Close()
			lo.Must1(tfFile.WriteTo(file))
		}
	},
}

func main() {
	rootCmd.Flags().String("base-url", "", "The grafana base URL (i.e. https://my-experiment.grafana.net).")
	rootCmd.MarkFlagRequired("base-url")
	rootCmd.Flags().String("session", "", "The grafana session id. It can be found on the Cookie 'grafana_session'.")
	rootCmd.MarkFlagRequired("session")

	rootCmd.Flags().Bool("include-all", false, "Include definitions even for Alert Rules that were created using the API")
	rootCmd.Flags().String("group-suffix", "", "Suffix to add to Rule Group names")
	rootCmd.Flags().String("rule-suffix", "", "Suffix to add to Rule titles")

	rootCmd.Execute()
}
