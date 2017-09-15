package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// controllerCmd represents the controller command
var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Generate controller files (c)",
	Long:  `This command helps generating controller structs.`,
	Example: `
  # generete users controller file. It is located at ./controllers directory by default
  gopress make controller users

  # use c for short
  gopress make c users

  # generate multiple controllerss in a single command
  gopress make controller users posts

  # generate controller files in directory other than ./controllers
  gopress make controller users --dir=actions`,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		dir := viper.Sub("controller").GetString("dir")
		if moduleRootPath != "" {
			dir = moduleRootPath
		}

		modules := make([]commonTemplateData, len(args))

		for i, m := range args {
			m = validateModuleFileName(m)
			mt := validateModuleTypeName(m)
			packageName, fileName := parseModulePath(dir, m, fileExtensionGolang)
			if err := checkFileName(fileName); err != nil {
				er(err)
			}
			modules[i] = commonTemplateData{packageName, fileName, m, mt}
		}

		for _, data := range modules {
			err := createControllerFile(
				data.packageName,
				data.fileName,
				data.moduleName,
				data.moduleTypeName,
			)
			if err != nil {
				er(err)
			}
			fmt.Printf("Controller file created: %s\n", data.fileName)
		}
	},
}

func init() {
	makeCmd.AddCommand(controllerCmd)

	viper.SetDefault("controller.dir", "controllers")
}

func createControllerFile(packageName, fileName, moduleName, moduleTypeName string) error {
	template := `package {{.packageName}}

import (
	"net/http"

	"github.com/fpay/gopress"
)

// {{.moduleTypeName}}
type {{.moduleTypeName}} struct {
	// Uncomment this line if you want to use services in the app
	// app *gopress.App
}

// New{{.moduleTypeName}} returns {{.moduleName}} controller instance.
func New{{.moduleTypeName}}() *{{.moduleTypeName}} {
	return new({{.moduleTypeName}})
}

// RegisterRoutes registes routes to app
// It is used to implements gopress.Controller.
func (c *{{.moduleTypeName}}) RegisterRoutes(app *gopress.App) {
	// Uncomment this line if you want to use services in the app
	// c.app = app

	app.GET("/{{.sampleRoute}}", c.{{.sampleAction}})
	// app.POST("/{{.sampleRoute}}", c.SamplePostAction)
	// app.PUT("/{{.sampleRoute}}", c.SamplePutAction)
	// app.DELETE("/{{.sampleRoute}}", c.SampleDeleteAction)
}

// {{.sampleAction}} Action
// Parameter gopress.Context is just alias of echo.Context
func (c *{{.moduleTypeName}}) {{.sampleAction}}(ctx gopress.Context) error {
	// Or you can get app from request context
	// app := gopress.AppFromContext(ctx)
	data := map[string]interface{}{}
	return ctx.Render(http.StatusOK, "{{.moduleName}}/sample", data)
}
`

	data := map[string]interface{}{
		"packageName":    packageName,
		"moduleName":     moduleName,
		"moduleTypeName": moduleTypeName + "Controller",
		"sampleAction":   "SampleGetAction",
		"sampleRoute":    moduleName + "/sample",
	}

	controllerScript, err := render(template, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(fileName, controllerScript)
	if err != nil {
		er(err)
	}

	return nil
}
