/**FileHeader
 * @Author: Liangkang Zhang
 * @Date: 2026/1/24 00:35:39
 * @LastEditors: Liangkang Zhang
 * @LastEditTime: 2026/1/25 00:24:57
 * @Description:
 * @Copyright: Copyright (©)}) 2026 Liangkang Zhang<lkzhang98@gmail.com>. All rights reserved. Use of this source code is governed by a MIT style license that can be found in the LICENSE file.. All rights reserved.
 * @Email: lkzhang98@gmail.com
 * @Repository: https://github.com/geminik12/autostack
 */
// Package log provides a simple logger.

package log

import (
	"fmt"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

// Options contains configuration options for logging.
type Options struct {
	// DisableCaller specifies whether to include caller information in the log.
	DisableCaller bool `json:"disable-caller,omitempty" mapstructure:"disable-caller"`
	// DisableStacktrace specifies whether to record a stack trace for all messages at or above panic level.
	DisableStacktrace bool `json:"disable-stacktrace,omitempty" mapstructure:"disable-stacktrace"`
	// EnableColor specifies whether to output colored logs.
	EnableColor bool `json:"enable-color"       mapstructure:"enable-color"`
	// Level specifies the minimum log level. Valid values are: debug, info, warn, error, dpanic, panic, and fatal.
	Level string `json:"level,omitempty" mapstructure:"level"`
	// Format specifies the log output format. Valid values are: console and json.
	Format string `json:"format,omitempty" mapstructure:"format"`
	// OutputPaths specifies the output paths for the logs.
	OutputPaths []string `json:"output-paths,omitempty" mapstructure:"output-paths"`
	// EnableFile specifies whether to enable file logging.
	EnableFile bool `json:"enable-file,omitempty" mapstructure:"enable-file"`
	// LogDir specifies the directory to store the logs.
	LogDir string `json:"log-dir,omitempty" mapstructure:"log-dir"`
	// MaxSize is the maximum size in megabytes of the log file before it gets rotated.
	MaxSize int `json:"max-size,omitempty" mapstructure:"max-size"`
	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int `json:"max-backups,omitempty" mapstructure:"max-backups"`
	// MaxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename.
	MaxAge int `json:"max-age,omitempty" mapstructure:"max-age"`
	// Compress determines if the rotated log files should be compressed using gzip.
	Compress bool `json:"compress,omitempty" mapstructure:"compress"`
}

// NewOptions creates a new Options object with default values.
func NewOptions() *Options {
	return &Options{
		Level:       zapcore.InfoLevel.String(),
		Format:      "console",
		OutputPaths: []string{"stdout"},
		EnableFile:  false,
		LogDir:      "./logs",
		MaxSize:     100,
		MaxBackups:  5,
		MaxAge:      30,
		Compress:    false,
	}
}

// Validate verifies flags passed to LogsOptions.
func (o *Options) Validate() []error {
	errs := []error{}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	if o.Format != "console" && o.Format != "json" {
		errs = append(errs, fmt.Errorf("invalid log format: %s", o.Format))
	}

	return errs
}

// AddFlags adds command line flags for the configuration.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, "log.level", o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, "log.disable-caller", o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, "log.disable-stacktrace", o.DisableStacktrace, ""+
		"Disable the log to record a stack trace for all messages at or above panic level.")
	fs.BoolVar(&o.EnableColor, "log.enable-color", o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringVar(&o.Format, "log.format", o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.StringSliceVar(&o.OutputPaths, "log.output-paths", o.OutputPaths, "Output paths of log.")
	fs.BoolVar(&o.EnableFile, "log.enable-file", o.EnableFile, "Enable file logging.")
	fs.StringVar(&o.LogDir, "log.dir", o.LogDir, "Directory to store the logs.")
	fs.IntVar(&o.MaxSize, "log.max-size", o.MaxSize, "Maximum size in megabytes of the log file before it gets rotated.")
	fs.IntVar(&o.MaxBackups, "log.max-backups", o.MaxBackups, "Maximum number of old log files to retain.")
	fs.IntVar(&o.MaxAge, "log.max-age", o.MaxAge, "Maximum number of days to retain old log files.")
	fs.BoolVar(&o.Compress, "log.compress", o.Compress, "Compress rotated log files.")
}
