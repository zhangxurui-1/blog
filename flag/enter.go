package flag

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"os"
	"server/global"
)

// 定义 CLI 标志，用于不同操作的命令行选项
var (
	sqlFlag = &cli.BoolFlag{
		Name:  "sql",
		Usage: "Initialize the structure of MySQL database table.",
	}
	sqlExportFlag = &cli.BoolFlag{
		Name:  "sql-export",
		Usage: "Export SQL data to a specified file.",
	}
	sqlImportFlag = &cli.StringFlag{
		Name:  "sql-import",
		Usage: "Import SQL data from a specified file.",
	}
	esFlag = &cli.BoolFlag{
		Name:  "es",
		Usage: "Initialize the Elasticsearch index.",
	}
	esExportFlag = &cli.BoolFlag{
		Name:  "es-export",
		Usage: "Export ES data to a specified file.",
	}
	esImportFlag = &cli.StringFlag{
		Name:  "es-import",
		Usage: "Import data into ES from a specified file.",
	}
	adminFlag = &cli.BoolFlag{
		Name:  "admin",
		Usage: "Create an administrator using the name, email and address specified in the config.yaml.",
	}
)

// Run 执行基于命令行标志的相应操作
// 它处理不同的标志，执行相应操作，并记录成功或错误的消息
func Run(c *cli.Context) {
	// 检查是否设置了多个标志
	if c.NumFlags() > 1 {
		err := cli.NewExitError("Only one flag can be specified!", 1)
		global.Log.Error("Invalid flag usage:", zap.Error(err))
		os.Exit(1)
	}

	// 根据不同的标志选择执行的操作
	switch {
	case c.Bool(sqlFlag.Name):
		if err := SQL(); err != nil {
			global.Log.Error("Failed to initialize table structure:", zap.Error(err))
			return
		}
		global.Log.Info("Successfully initialized table structure")

	case c.Bool(sqlExportFlag.Name):
		if err := SQLExport(); err != nil {
			global.Log.Error("Failed to export SQL data:", zap.Error(err))
		} else {
			global.Log.Info("Successfully exported SQL data")
		}
	case c.IsSet(sqlImportFlag.Name):
		if errs := SQLImport(c.String(sqlImportFlag.Name)); len(errs) > 0 {
			var combinedErrs string
			for _, err := range errs {
				combinedErrs += err.Error() + "\n"
			}
			err := errors.New(combinedErrs)
			global.Log.Error("Failed to import SQL data:", zap.Error(err))
		} else {
			global.Log.Info("Successfully import SQL data")
		}
	case c.Bool(esFlag.Name):
		if err := Elasticsearch(); err != nil {
			global.Log.Error("Failed to create ES index:", zap.Error(err))
		} else {
			global.Log.Info("Successfully created ES index")
		}
	case c.Bool(esExportFlag.Name):
		if err := ElasticsearchExport(); err != nil {
			global.Log.Error("Failed to export ES data:", zap.Error(err))
		} else {
			global.Log.Info("Successfully exported ES data")
		}
	case c.IsSet(esImportFlag.Name):
		if num, err := ElasticsearchImport(c.String(esImportFlag.Name)); err != nil {
			global.Log.Error("Failed to import ES data:", zap.Error(err))
		} else {
			global.Log.Info(fmt.Sprintf("Successfully imported ES data, totaling %v records", num))
		}
	case c.Bool(adminFlag.Name):
		if err := Admin(); err != nil {
			global.Log.Error("Failed to Create an administrator:", zap.Error(err))
		} else {
			global.Log.Info("Successfully created an administrator")
		}
	default:
		err := cli.NewExitError("unknown flag", 1)
		global.Log.Error(err.Error(), zap.Error(err))
	}

}

// NewApp 创建并配置一个新的 CLI 应用程序，设置标志和默认操作
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Go blog"
	app.Flags = []cli.Flag{
		sqlFlag,
		sqlExportFlag,
		sqlImportFlag,
		esFlag,
		esExportFlag,
		esImportFlag,
		adminFlag,
	}

	app.Action = Run
	return app
}

// InitFlag 初始化并运行 CLI 应用程序
func InitFlag() {
	if len(os.Args) > 1 {
		app := NewApp()
		if err := app.Run(os.Args); err != nil {
			global.Log.Error("Application execution encountered an error:", zap.Error(err))
			os.Exit(1)
		}
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Println("Displaying help message...")
		}
		os.Exit(0)
	}
}
