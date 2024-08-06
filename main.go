package main

import (
	"github.com/sirupsen/logrus"
	"lolcrawler/config"
	"lolcrawler/crawler"
	"lolcrawler/store"
	"os"
	"strconv"
)

func panicOnError(err error, logger *logrus.Logger) {
	if err != nil {
		logger.Panic(err)
	}
}

func main() {
	logger := logrus.New()

	if len(os.Args) < 2 {
		logger.Panic("match count required as first argument")
	}

	crawlCount, err := strconv.Atoi(os.Args[1])

	cfg, err := config.Load()
	panicOnError(err, logger)

	file, err := os.Create(cfg.LogFile)
	panicOnError(err, logger)

	logger.SetOutput(file)

	db, err := store.NewDatabase(cfg.Crawler.DataSource)
	panicOnError(err, logger)

	crwlr, err := crawler.New(store.New(db), cfg.RiotApi.Key, cfg.RiotApi.Region, cfg.Crawler.StartName, cfg.Crawler.StartTag, logger)
	panicOnError(err, logger)

	panicOnError(crwlr.Crawl(crawlCount), logger)
	panicOnError(crwlr.FillMissingData(), logger)
}
