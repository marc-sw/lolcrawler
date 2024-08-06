package crawler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/account"
	"github.com/sirupsen/logrus"
	"lolcrawler/store"
)

type Crawler struct {
	client     *golio.Client
	store      *store.Store
	logger     *logrus.Logger
	accountRow int
	matchIdRow int
}

func New(s *store.Store, riotApiKey, riotApiRegion, startName, startTag string, logger *logrus.Logger) (*Crawler, error) {
	client := golio.NewClient(riotApiKey, golio.WithLogger(logger), golio.WithRegion(api.Region(riotApiRegion)))
	count, err := s.GetAccountsCount()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		acc, err := client.Riot.Account.GetByRiotID(startName, startTag)
		if err != nil {
			return nil, err
		}
		err = s.AddOrIgnoreAccount(*acc)
		if err != nil {
			return nil, err
		}
	}

	accountRow, err := s.GetCount(store.CountAccountRow)
	if err != nil {
		return nil, err
	}
	matchIdRow, err := s.GetCount(store.CountMatchIdRow)
	if err != nil {
		return nil, err
	}
	return &Crawler{
		client:     client,
		store:      s,
		logger:     logger,
		accountRow: accountRow,
		matchIdRow: matchIdRow,
	}, nil
}

func (crawler *Crawler) Logger() *logrus.Logger {
	return crawler.logger
}

func (crawler *Crawler) crawlNextMatch() error {
	matchId, err := crawler.store.GetMatchIdAtRow(crawler.matchIdRow)
	if err != nil {
		return err
	}
	match, err := crawler.client.Riot.LoL.Match.Get(matchId)
	if err != nil {
		if errors.Is(err, api.ErrNotFound) {
			crawler.logger.Error(matchId + " " + err.Error())
			err = crawler.store.DeleteMatchId(matchId)
			if err != nil {
				return err
			}
			crawler.logger.Info("deleted " + matchId)
			return crawler.crawlNextMatch()
		}
		return err
	}
	var acc account.Account
	for _, participant := range match.Info.Participants {
		if participant.PUUID == "BOT" {
			continue
		}
		acc = account.Account{
			Puuid:    participant.PUUID,
			GameName: participant.RiotIDGameName,
			TagLine:  participant.RiotIDTagline,
		}
		err = crawler.store.AddOrIgnoreAccount(acc)
		if err != nil {
			return err
		}
	}
	crawler.matchIdRow++
	return crawler.store.SetCount(store.CountMatchIdRow, crawler.matchIdRow)
}

func (crawler *Crawler) crawlNextAccount() error {
	acc, err := crawler.store.GetAccountAtRow(crawler.accountRow)
	if err != nil {
		return err
	}
	matchIds, err := crawler.client.Riot.LoL.Match.List(acc.Puuid, 0, 100)
	if err != nil {
		return err
	}
	err = crawler.store.AddOrIgnoreMatchIds(matchIds)
	if err != nil {
		return err
	}
	crawler.accountRow++
	return crawler.store.SetCount(store.CountAccountRow, crawler.accountRow)
}

func (crawler *Crawler) CrawlNext() error {
	err := crawler.crawlNextMatch()
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	err = crawler.crawlNextAccount()
	if err == nil {
		return crawler.CrawlNext()
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNothingToCrawl
	}
	return err
}

func (crawler *Crawler) Crawl(total int) error {
	fmt.Println(fmt.Sprintf("crawling %d matches", total))
	progress := NewProgress(total)
	var err error
	for !progress.Done() {
		if err = crawler.CrawlNext(); err != nil {
			crawler.logger.Error(err)
			continue
		}
		progress.Increase()
	}
	fmt.Println()
	return nil
}

func (crawler *Crawler) FillMissingData() error {

	puuids, err := crawler.store.GetEmptyAccountsPuuid()
	if err != nil {
		return err
	}
	total := len(puuids)
	fmt.Println(fmt.Sprintf("filling %d accounts", total))

	progress := NewProgress(total)
	for _, puuid := range puuids {
		acc, err := crawler.client.Riot.Account.GetByPUUID(puuid)
		if err != nil {
			if errors.Is(err, api.ErrNotFound) {
				if err = crawler.store.DeleteAccount(puuid); err != nil {
					return err
				}
				crawler.logger.Error(fmt.Sprintf("no account found with puuid '%s'", puuid))
				crawler.logger.Info(fmt.Sprintf("deleting account with puuid '%s'", puuid))
				continue
			}
			return err
		}
		if err = crawler.store.UpdateAccount(*acc); err != nil {
			return err
		}
		progress.Increase()
	}
	fmt.Println()
	return nil
}
