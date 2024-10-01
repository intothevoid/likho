package generator

import (
	"fmt"
	"html/template"
	"os"

	"github.com/intothevoid/likho/pkg/utils"
	"go.uber.org/zap"
)

func executeTemplate(tmpl *template.Template, name, outputPath string, data interface{}) error {
	logger := utils.GetLogger()

	logger.Debug("executing template",
		zap.String("name", name),
		zap.String("outputPath", outputPath),
		zap.Any("data", data))

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", outputPath, err)
	}
	defer file.Close()

	err = tmpl.ExecuteTemplate(file, "base.html", data)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", name, err)
	}

	// Add this log after successful template execution
	logger.Info("html file generated", zap.String("path", outputPath))

	return nil
}
