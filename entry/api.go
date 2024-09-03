package entry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/disgoorg/log"
	"secondlap.go/discord"
)

type entryRes struct {
	Name    string `json:"name"`
	Success string `json:"success"`
	Message string `json:"message"`
	Fail    string `json:"fail"`
}

func (s *Session) getHomepage() error {

	s.Log.Info("getting homepage...")

	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"Host":            "www.motosport.com",
			"Connection":      "keep-alive",
			"Cache-Control":   "max-age=0",
			"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"Accept-Language": "en-US,en;q=0.9",
		}).Get("https://www.motosport.com/win?utm_medium=referral&utm_source=event&utm_campaign=2024-MXtour-Sweeps")

	if err != nil {
		// s.setProxy()
		s.Log.Error("busy proxy is busy, retrying")
		time.Sleep(3 * time.Second)
		s.getHomepage()
		// return err
	}
	if res.StatusCode() != 200 {
		fmt.Println(res.String())
		return fmt.Errorf("invalid status code: %d", res.StatusCode())
	} else if res.StatusCode() == 403 {
		s.setProxy()
		log.Errorf("403 error, changing proxy")
		time.Sleep(3 * time.Second)
		s.getHomepage()
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		return err
	}

	str, err := getXSRFToken(doc)
	if err != nil {
		return fmt.Errorf("error parsing xsrf token")
	}

	s.Log.Info("session made: ", str)
	s.state.xsrf = str
	s.Client.SetCookies(res.Cookies())

	return nil
}

func (s *Session) submitEntry() error {

	s.Log.Info("submitting entry...")

	res, err := s.Client.R().
		SetHeaders(map[string]string{
			"Host":            "www.motosport.com",
			"Connection":      "keep-alive",
			"Cache-Control":   "max-age=0",
			"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"Accept-Language": "en-US,en;q=0.9",
			"X-CSRF-TOKEN":    s.state.xsrf,
			"Referer":         "https://www.motosport.com/win?utm_medium=referral&utm_source=event&utm_campaign=2024-MXtour-Sweeps",
			"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
		}).SetFormData(map[string]string{
		"url":                  "/win",
		"first_name":           s.state.FirstName,
		"last_name":            s.state.LastName,
		"email":                s.state.Email,
		"newsletter":           "newsletter",
		"sweepstakes_submit":   "submit",
		"_csrf_token":          s.state.xsrf,
		"g-recaptcha-response": s.state.ReCap,
	}).Post("https://www.motosport.com/win/submit?utm_medium=referral&utm_source=event&utm_campaign=2024-MXtour-Sweeps")

	if err != nil {
		return err
	}

	if res.StatusCode() != 200 {
		fmt.Println(res.String())
		return fmt.Errorf("invalid status code: %w", err)
	}

	var entry entryRes
	if err := json.Unmarshal(res.Body(), &entry); err != nil {
		fmt.Println(res.String())
		return err
	}

	if entry.Success != "" {
		s.Log.Info("entry successful!")
		s.state.EntryCount++
		d := discord.DiscordData{
			Email:      s.state.Email,
			EntryCount: s.state.EntryCount,
		}
		s.Webhook = d
		// d.SendEmbed()

	} else {

		if strings.Contains(res.String(), "You are already entered") {
			return ErrEmailAlreadyEntered
		}
		// fmt.Println(res.String())
		// return fmt.Errorf("entry already submitted: %s", entry.Message)
	}

	return nil

}
