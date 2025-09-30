// Мини приложение для последовательного выполнения миграций
package main

import (
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type sqlFile struct {
	prefix int
	name   string
}

func main() {
	godotenv.Load()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	pathToMigrations := os.Getenv("MIGRATIONS_PATH")
	databaseUrl := os.Getenv("DATABASE_URL")

	files, err := os.ReadDir(pathToMigrations)
	if err != nil {
		log.Fatal().Msg("Unable to read the migrations directory. Please ensure that the MIGRATIONS_PATH variable is set.")
		return
	}

	re := regexp.MustCompile(`^(\d{4})_.*\.sql$`)

	var sqlFiles []sqlFile

	for _, file := range files {
		if file.IsDir() {
			log.Error().Msg(file.Name() + " is directory")
			continue
		}
		regexMatch := re.FindStringSubmatch(file.Name())
		if len(regexMatch) == 2 {
			num, err := strconv.Atoi(regexMatch[1])
			if err != nil {
				log.Error().Err(err)
				continue
			}
			sqlFiles = append(sqlFiles, sqlFile{
				prefix: num,
				name:   regexMatch[0],
			})
		} else {
			log.Error().Msg(file.Name() + " does not fit the format")
		}
	}
	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].prefix < sqlFiles[j].prefix
	})

	for _, f := range sqlFiles {
		log.Info().Str("filename", f.name).Msg("Executing")
		cmd := exec.Command("psql", databaseUrl, "-f", pathToMigrations+"/"+f.name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Error().Str("filename", f.name).Err(err)
			os.Exit(1)
		}
	}
}
