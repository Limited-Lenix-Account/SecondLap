package entry

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"secondlap.go/discord"
)

type Session struct {
	Cancel context.CancelFunc
	Ctx    context.Context
	Mu     sync.Mutex

	Log    *logrus.Entry
	Client *resty.Client
	ID     int

	UserConfig *Config
	ProxyList  []string
	EmailList  []string
	Proxy      string
	IPStr      string
	Webhook    discord.DiscordData

	state state
}

type state struct {
	FirstName  string
	LastName   string
	Email      string
	EntryCount int
	xsrf       string
	ReCap      string
}

type Profile struct {
	FirstName  string
	LastName   string
	Email      string
	EntryCount int
}

func NewSession(log *logrus.Logger, ctx context.Context, cancel context.CancelFunc, userConfig *Config, taskID int, profile Profile) (*Session, error) {

	return &Session{
		Cancel: cancel,
		Ctx:    ctx,
		// Log:        log.WithFields("task_id", taskID),
		Log: log.WithFields(logrus.Fields{
			"task_id": taskID,
			"email":   profile.Email,
		}),
		UserConfig: userConfig,
		Client:     resty.New(),
		state: state{
			FirstName:  profile.FirstName,
			LastName:   profile.LastName,
			Email:      profile.Email,
			EntryCount: profile.EntryCount,
		},
	}, nil

}

func NewProxyTest(log *logrus.Logger, ctx context.Context, cancel context.CancelFunc, userConfig *Config, taskID int) (*Session, error) {

	return &Session{
		ID:     taskID,
		Cancel: cancel,
		Ctx:    ctx,
		Log:    log.WithField("task_id", taskID),
		// Log: log.WithFields(logrus.Fields{
		// 	"task_id": taskID,
		// }),
		UserConfig: userConfig,
		Client:     resty.New(),
	}, nil
}

func (s *Session) setProxy() error {
	s.Log.Info("setting proxy...")

	if err := s.setRandomProxy(); err != nil {
		return fmt.Errorf("error setting index proxy: %w", err)
	}

	res, err := s.Client.R().Get("https://ipecho.net/plain")
	if err != nil {
		s.Log.Errorf("busy proxy, rotating proxy...")
		s.setProxy()
	}
	s.IPStr = res.String()

	return nil
}

func (s *Session) setRandomProxy() error {

	if len(s.ProxyList) == 0 {
		return nil
	}

	proxy := s.ProxyList[rand.Intn(len(s.ProxyList))]
	parts := strings.Split(proxy, ":")
	if len(parts) != 4 {
		return fmt.Errorf("wrong proxy format")
	}
	proxyURL, _ := url.Parse(("http://" + parts[2] + ":" + parts[3] + "@" + parts[0] + ":" + parts[1]))

	s.Client.SetTransport(&http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	})

	s.Client.SetProxy("http://" + parts[2] + ":" + parts[3] + "@" + parts[0] + ":" + parts[1])
	return nil

}

func (s *Session) setIndexProxy(i int) error {

	if len(s.ProxyList) == 0 {
		return nil
	}

	proxy := s.ProxyList[i]
	s.Proxy = proxy
	parts := strings.Split(proxy, ":")
	if len(parts) != 4 {
		return fmt.Errorf("wrong proxy format")
	}
	proxyURL, _ := url.Parse(("http://" + parts[2] + ":" + parts[3] + "@" + parts[0] + ":" + parts[1]))

	s.Client.SetTransport(&http.Transport{
		Proxy:             http.ProxyURL(proxyURL),
		DisableKeepAlives: true,
	})

	s.Client.SetProxy("http://" + parts[2] + ":" + parts[3] + "@" + parts[0] + ":" + parts[1])
	return nil
}
