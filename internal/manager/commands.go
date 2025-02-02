package manager

import (
	"os"

	"github.com/trustwallet/assets/internal/config"
	"github.com/trustwallet/assets/internal/file"
	"github.com/trustwallet/assets/internal/processor"
	"github.com/trustwallet/assets/internal/report"
	"github.com/trustwallet/assets/internal/service"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configPath, root string

func InitCommands() {
	rootCmd.Flags().StringVar(&configPath, "config", ".github/assets.config.yaml",
		"config file (default is $HOME/.github/assets.config.yaml)")
	rootCmd.Flags().StringVar(&root, "root", ".", "root path to files")

	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(fixCmd)
	rootCmd.AddCommand(updateAutoCmd)
	rootCmd.AddCommand(updateManualCmd)
	rootCmd.AddCommand(addTokenCmd)
}

var (
	rootCmd = &cobra.Command{
		Use:   "assets",
		Short: "",
		Long:  "",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: " Execute validation checks",
		Run: func(cmd *cobra.Command, args []string) {
			assetsService := InitAssetsService()
			assetsService.RunJob(assetsService.Check)
		},
	}
	fixCmd = &cobra.Command{
		Use:   "fix",
		Short: "Perform automatic fixes where possible",
		Run: func(cmd *cobra.Command, args []string) {
			assetsService := InitAssetsService()
			assetsService.RunJob(assetsService.Fix)
		},
	}
	updateAutoCmd = &cobra.Command{
		Use:   "update-auto",
		Short: "Run automatic updates from external sources",
		Run: func(cmd *cobra.Command, args []string) {
			assetsService := InitAssetsService()
			assetsService.RunUpdateAuto()
		},
	}
	updateManualCmd = &cobra.Command{
		Use:   "update-manual",
		Short: "Run manual updates from external sources",
		Run: func(cmd *cobra.Command, args []string) {
			assetsService := InitAssetsService()
			assetsService.RunUpdateManual()
		},
	}

	addTokenCmd = &cobra.Command{
		Use:   "add-token",
		Short: "Creates info.json template for the asset",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("1 argument was expected")
			}

			err := CreateAssetInfoJSONTemplate(args[0])
			if err != nil {
				log.Fatalf("Can't create asset info json template: %v", err)
			}
		},
	}
)

func InitAssetsService() *service.Service {
	setup()

	paths, err := file.ReadLocalFileStructure(root, config.Default.ValidatorsSettings.RootFolder.SkipFiles)
	if err != nil {
		log.WithError(err).Fatal("Failed to load file structure.")
	}

	fileService := file.NewService(paths...)
	validatorsService := processor.NewService(fileService)
	reportService := report.NewService()

	return service.NewService(fileService, validatorsService, reportService, paths)
}

func setup() {
	if err := config.SetConfig(configPath); err != nil {
		log.WithError(err).Fatal("Failed to set config.")
	}

	logLevel, err := log.ParseLevel(config.Default.App.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("Failed to parse log level.")
	}

	log.SetLevel(logLevel)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
