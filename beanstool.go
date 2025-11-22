package main

import (
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/src-d/beanstool/cli"

	"github.com/agtorre/gocolorize"
	"github.com/jessevdk/go-flags"
	"github.com/lmittmann/tint"
	"github.com/rafaelespinoza/logg"
	"golang.org/x/term"
)

func main() {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		gocolorize.SetPlain(true)
	}

	type GlobalFlags struct {
		LogLevel  string `short:"l" long:"log-level" description:"configure logging level" choice:"DEBUG" choice:"INFO" choice:"WARN" choice:"ERROR" default:"INFO"`
		LogFormat string `short:"f" long:"log-format" description:"configure logging format" choice:"JSON" choice:"TEXT" choice:"TINT" default:"TINT"`
	}
	var opts GlobalFlags

	parser := flags.NewParser(&opts, flags.Default)
	parser.CommandHandler = func(cmdr flags.Commander, args []string) error {
		handler := makeSlogHandler(os.Stderr, opts.LogLevel, opts.LogFormat)
		logg.SetDefaults(handler, nil)
		return cmdr.Execute(args)
	}
	parser.AddCommand("stats", "print stats on all tubes", "", &cli.StatsCommand{})
	parser.AddCommand("tail", "tails a tube and prints his content", "", &cli.TailCommand{})
	parser.AddCommand("peek", "peeks a job from a queue", "", &cli.PeekCommand{})
	parser.AddCommand("delete", "delete a job from a queue", "", &cli.DeleteCommand{})
	parser.AddCommand("kick", "kicks jobs from buried back into ready", "", &cli.KickCommand{})
	parser.AddCommand("put", "put a job into a tube", "", &cli.PutCommand{})
	parser.AddCommand("bury", "bury existing jobs from ready state", "", &cli.BuryCommand{})

	_, err := parser.Parse()
	if err != nil {
		if flagErr, ok := err.(*flags.Error); ok && flagErr.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}

var loggingLevels = []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}

func makeSlogHandler(w io.Writer, logLevel, logFormat string) slog.Handler {
	// The implied default logging level is INFO because it's the zero value in
	// the log/slog package.
	var lvl slog.Level
	levels := make([]string, len(loggingLevels))
	for i, validLevel := range loggingLevels {
		levels[i] = validLevel.String()
	}
	if ind := slices.Index(levels, strings.ToUpper(strings.TrimSpace(logLevel))); ind >= 0 {
		lvl = loggingLevels[ind]
	} else {
		slog.Debug("unknown logging level, setting to default",
			slog.String("input_level", logLevel),
			slog.String("default_level", lvl.String()),
		)
	}

	var handler slog.Handler
	opts := slog.HandlerOptions{Level: lvl}
	switch strings.ToUpper(strings.TrimSpace(logFormat)) {
	case "JSON":
		handler = slog.NewJSONHandler(w, &opts)
	case "TEXT":
		handler = slog.NewTextHandler(w, &opts)
	case "TINT":
		opts := tint.Options{Level: lvl, TimeFormat: time.RFC3339}
		handler = tint.NewHandler(w, &opts)
	default:
		slog.Debug("unknown logging format, setting to default",
			slog.String("input_format", logFormat),
			slog.String("default_format", "JSON"),
		)
		handler = slog.NewJSONHandler(w, &opts)
	}

	return handler
}
