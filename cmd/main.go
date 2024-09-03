package main

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"sync"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"secondlap.go/activity"
	"secondlap.go/discord"
	"secondlap.go/entry"
)

var (
	db      *activity.ProfRep
	log     *logrus.Logger
	proxies []string
	emails  []string
)

func init() {

	cmd := exec.Command(`echo -en "\033]0;Second Lap\a`)
	cmd.Run()

	proxies, _ = entry.LoadTextFile("data/proxies.txt")
	emails, _ = entry.LoadTextFile("data/emails.txt")
	db, _ = activity.InitDB()

	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
}

func main() {

	//init stuff

	//get profile table size and update if email list is longer than
	//emails already entered - doesn't do anything if emails are removed from emails.txt
	size, err := db.GetTableSize()
	if err != nil {
		log.Fatal("error getting table size", err)
	}

	if size < len(emails) {
		log.Info("updating Emails")
		log.Info("current table size: ", size)
		db.ImportEmails(emails)
		size, _ = db.GetTableSize()
		log.Info("updated table size: ", size)
	}

	log.Info("current emails loaded: ", size)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := entry.LoadConfig()
	if err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var count int
	var profs []entry.Profile
	var choice int

	prompt := &survey.Select{
		Message: "What would you like to do",
		Options: []string{"View Profiles", "Place Entries", "Test Proxies", "Test Webhook"},
	}

	err = survey.AskOne(prompt, &choice)
	if err == terminal.InterruptErr {
		log.Fatal("interrupted")
	}

	switch choice {
	case 0:
		log.Info("Viewing Profiles...")
		db.PrintAllProfiles()

	case 1:

		for {
			log.Info("How many tasks would you like to run")
			if _, err := fmt.Scan(&count); err != nil {
				log.Fatalf("invalid input: %s", err)
			}

			if count > size {
				log.Warn("count cannot be greater than amount of profiles")
			} else {
				break
			}

		}

		profs = db.SelectProfiles(count)

		webhookChan := make(chan *discord.DiscordData)
		go func() {

			var wg sync.WaitGroup
			sem := semaphore.NewWeighted(int64(count))
			for i := 0; i < count; i++ {
				wg.Add(1)
				sem.Acquire(ctx, 1)
				go func(i int) {
					defer wg.Done()
					defer sem.Release(1)

					s, err := entry.NewSession(log, ctx, cancel, config, i+1, profs[i])
					if err != nil {
						log.Fatalf("error creating session: %s", err)
					}

					s.ProxyList = proxies
					s.EmailList = emails

					if success, hook, err := s.Ignite(); err != nil {
						log.Errorf("error creating entry: %s", err)
					} else {
						if success {

							if err := db.UpdateProfile(profs[i].Email); err != nil {
								log.Errorf("failed to update profile: %s", err.Error())
							} else {
								log.Info("profile updated")
								webhookChan <- hook
							}

						}
					}

				}(i)
			}

			wg.Wait()
			close(webhookChan)
		}()

		for i := range webhookChan {
			s := rand.Float64() * 3
			log.Infof("delaying webhook %f seconds", s)
			time.Sleep(time.Duration(s) * time.Second)
			i.SendEmbed()
		}

	case 2:

		// l := len(proxies)
		l := 50

		var wg sync.WaitGroup
		sem := semaphore.NewWeighted(int64(l))
		for i := 0; i < l; i++ {
			wg.Add(1)
			sem.Acquire(ctx, 1)
			time.Sleep(time.Duration(rand.Intn(15)))
			go func(i int) {
				defer wg.Done()
				defer sem.Release(1)

				log.Info("Starting Proxy Test...")

				s, err := entry.NewProxyTest(log, ctx, cancel, config, i+1)
				if err != nil {
					log.Fatalf("error creating session: %s", err)
				}
				s.ProxyList = proxies

				time.Sleep(time.Duration(rand.Intn(10)))

				ok, err := s.ProxyTest()
				if err != nil {
					s.Log.Error(err)
					// fmt.Println(s.Proxy)

					// entry.AppendToFile("data/deadproxy.txt", s.Proxy)
				}
				if ok {
					s.Log.Info("proxy works")
				}

			}(i)
		}
		wg.Wait()
	case 3:

		log.Info("Testing Webhooks")
		l := 50

		var wg sync.WaitGroup
		sem := semaphore.NewWeighted(int64(l))
		for i := 0; i < l; i++ {
			wg.Add(1)
			sem.Acquire(ctx, 1)
			go func(i int) {
				d := discord.DiscordData{
					Email:      "example@catchall.com",
					EntryCount: i,
				}
				d.SendEmbed()

			}(i)

		}
		wg.Wait()

	}
}
