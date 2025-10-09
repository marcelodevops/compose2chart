package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"helm-compose2chart-plugin/internal/convert"
)

var (
	composeFile string
	outDir      string
	chartName   string
	appVersion  string
	version     string
)

var rootCmd = &cobra.Command{
	Use:   "helm-compose2chart",
	Short: "Generate a Helm chart from docker-compose.yml",
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&composeFile, "file", "f", "docker-compose.yml", "path to docker-compose file")
	rootCmd.Flags().StringVarP(&outDir, "out", "o", "./chart", "output directory for generated Helm chart")
	rootCmd.Flags().StringVarP(&chartName, "name", "n", "generated-chart", "chart name")
	rootCmd.Flags().StringVar(&appVersion, "app-version", "0.1.0", "Chart appVersion")
	rootCmd.Flags().StringVar(&version, "version", "0.1.0", "Chart version")
}

func run() error {
	opts := convert.Options{
		ComposeFile: composeFile,
		OutDir:      outDir,
		ChartName:   chartName,
		AppVersion:  appVersion,
		Version:     version,
	}
	return convert.GenerateChart(opts)
}