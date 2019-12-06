package cmd

import (
	"fmt"
	"github.com/microsoft/abstrakt/internal/chartservice"
	"github.com/microsoft/abstrakt/internal/composeservice"
	"github.com/microsoft/abstrakt/internal/tools/logger"
	"github.com/spf13/cobra"
	"path"
	"strings"
)

var templateType string
var constellationFilePath string
var mapsFilePath string
var outputPath string
var zipChart *bool

var composeCmd = &cobra.Command{
	Use:   "compose [chart name]",
	Short: "Compose a package into requested template type",
	Long: `Compose is for composing a package based on mapsFilePath and constellationFilePath and template (default value is helm).

Example: abstrakt compose [chart name] -t [templateType] -f [constellationFilePath] -m [mapsFilePath] -o [outputPath] -z`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {

		chartName := args[0]

		if templateType != "helm" && templateType != "" {
			return fmt.Errorf("Template type: %v is not known", templateType)
		}

		service := composeservice.NewComposeService()
		_ = service.LoadFromFile(constellationFilePath, mapsFilePath)
		chart, err := service.Compose(chartName, outputPath)
		if err != nil {
			return fmt.Errorf("Could not compose: %v", err)
		}

		err = chartservice.SaveChartToDir(chart, outputPath)

		if err != nil {
			return fmt.Errorf("There was an error saving the chart: %v", err)
		}

		logger.Infof("Chart was saved to: %v", outputPath)

		out, err := chartservice.BuildChart(path.Join(outputPath, chartName))

		if err != nil {
			return fmt.Errorf("There was an error saving the chart: %v", err)
		}

		if *zipChart {
			_, err = chartservice.ZipChartToDir(chart, outputPath)
			if err != nil {
				return fmt.Errorf("There was an error zipping the chart: %v", err)
			}
		}

		logger.PrintBuffer(out, true)

		logger.Debugf("args: %v", strings.Join(args, " "))
		logger.Debugf("template: %v", templateType)
		logger.Debugf("constellationFilePath: %v", constellationFilePath)
		logger.Debugf("mapsFilePath: %v", mapsFilePath)
		logger.Debugf("outputPath: %v", outputPath)
		return nil
	},
}

func init() {
	composeCmd.Flags().StringVarP(&constellationFilePath, "constellationFilePath", "f", "", "constellation file path")
	_ = composeCmd.MarkFlagRequired("constellationFilePath")
	composeCmd.Flags().StringVarP(&mapsFilePath, "mapsFilePath", "m", "", "maps file path")
	_ = composeCmd.MarkFlagRequired("mapsFilePath")
	composeCmd.Flags().StringVarP(&outputPath, "outputPath", "o", "", "destination directory")
	_ = composeCmd.MarkFlagRequired("outputPath")
	composeCmd.Flags().StringVarP(&templateType, "templateType", "t", "helm", "output template type")
	zipChart = composeCmd.Flags().BoolP("zipChart", "z", false, "zips the chart")
}
