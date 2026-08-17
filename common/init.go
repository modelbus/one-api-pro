package common

import (
	"flag"
	"fmt"
	"path/filepath"

	"log"
	"os"

	"github.com/Leon-PanPan/one-api-pro/common/config"
	"github.com/Leon-PanPan/one-api-pro/common/logger"
	"github.com/joho/godotenv"
)

var (
	Port         = flag.Int("port", 3000, "the listening port")
	PrintVersion = flag.Bool("version", false, "print version and exit")
	PrintHelp    = flag.Bool("help", false, "print help and exit")
	LogDir       = flag.String("log-dir", "./logs", "specify the log directory")
	EnvFile      = flag.String("env", "", "specify the .env file path (supports relative path)")
)

func printHelp() {
	fmt.Println("One Api Pro " + Version + " - All in one API service for OpenAI API.")
	fmt.Println("Copyright (C) 2026 Leon PanPan. All rights reserved.")
	fmt.Println("GitHub: https://github.com/Leon-PanPan/one-api-pro")
	fmt.Println("Usage: one-api-pro [--port <port>] [--log-dir <log directory>] [--env <env file path>] [--version] [--help]")
}

func Init() {
	loadVersionFromFile()

	flag.Parse()

	if *PrintVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *PrintHelp {
		printHelp()
		os.Exit(0)
	}

	loadEnvFile()

	if os.Getenv("SESSION_SECRET") != "" {
		if os.Getenv("SESSION_SECRET") == "random_string" {
			logger.SysError("SESSION_SECRET is set to an example value, please change it to a random string.")
		} else {
			config.SessionSecret = os.Getenv("SESSION_SECRET")
		}
	}
	if os.Getenv("SQLITE_PATH") != "" {
		SQLitePath = os.Getenv("SQLITE_PATH")
	}
	if *LogDir != "" {
		var err error
		*LogDir, err = filepath.Abs(*LogDir)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := os.Stat(*LogDir); os.IsNotExist(err) {
			err = os.Mkdir(*LogDir, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}
		logger.LogDir = *LogDir
	}
}

func loadEnvFile() {
	var envPath string
	if *EnvFile != "" {
		absPath, err := filepath.Abs(*EnvFile)
		if err != nil {
			log.Fatalf("failed to resolve env file path '%s': %s", *EnvFile, err.Error())
		}
		envPath = absPath
	} else {
		envPath = ".env"
	}

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		if *EnvFile != "" {
			log.Fatalf("specified env file not found: %s", envPath)
		}
		return
	}

	if err := godotenv.Load(envPath); err != nil {
		if *EnvFile != "" {
			log.Fatalf("failed to load env file '%s': %s", envPath, err.Error())
		}
		logger.SysError(fmt.Sprintf("failed to load .env file: %s", err.Error()))
	}
}
