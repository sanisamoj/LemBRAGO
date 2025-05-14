package config

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var iconDomains = []string{
	// Social
	"facebook.com", "instagram.com", "x.com", "linkedin.com", "youtube.com",
	"twitch.tv", "tiktok.com", "pinterest.com", "reddit.com", "whatsapp.com",
	"discord.com", "snapchat.com", "telegram.org",

	// Trabalho
	"slack.com", "figma.com", "trello.com", "notion.so", "github.com",
	"docker.com", "asana.com", "atlassian.com", "zoom.us", "miro.com",
	"drive.google.com", "teams.microsoft.com", "gitlab.com", "hubspot.com",
	"office.com", "salesforce.com",

	// Serviços
	"gmail.com", "paypal.com", "stripe.com", "amazon.com", "mercadolivre.com.br",
	"ebay.com", "booking.com", "airbnb.com", "ifood.com.br", "outlook.com",
	"1password.com", "lastpass.com", "bitwarden.com", "uber.com", "99app.com",
	"rappi.com.br", "deezer.com", "music.apple.com", "nubank.com.br",

	// Tech
	"cloudflare.com", "digitalocean.com", "google.com", "chrome.com", "firefox.com",
	"apple.com", "dropbox.com", "aws.amazon.com", "azure.microsoft.com",
	"microsoft.com", "stackoverflow.com", "vercel.com", "netlify.com", "theverge.com",
	"godaddy.com", "hostinger.com", "mongodb.com", "oracle.com",

	// Geral
	"wikipedia.org", "netflix.com", "spotify.com", "globo.com", "uol.com.br",
	"nytimes.com", "imdb.com", "maps.google.com", "weather.com", "tripadvisor.com",
	"canva.com", "bing.com", "disneyplus.com", "max.com", "primevideo.com",
	"tv.apple.com", "coursera.org", "udemy.com", "duolingo.com", "quora.com",
	"imgur.com",
}

func downloadIcon(domain string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	// Concurrency limiter
	sem <- struct{}{}
	defer func() { <-sem }()

	url := fmt.Sprintf("https://icon.horse/icon/%s", domain)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("[ERROR] Failed to download %s: %v\n", domain, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[ERROR] %s returned status %d\n", domain, resp.StatusCode)
		return
	}

	// Detect file extension from Content-Type
	contentType := resp.Header.Get("Content-Type")
	exts, _ := mime.ExtensionsByType(contentType)

	ext := ".webp" // default/fallback
	if len(exts) > 0 && strings.HasPrefix(contentType, "image/") {
		ext = exts[0]
	}

	// Clean domain name (remove everything after first dot)
	re := regexp.MustCompile(`^[^.]+`)
	base := re.FindString(domain)
	filename := base + ext
	path := filepath.Join("uploads", filename)

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("[INFO] The icon %s already exists, skipping download.\n", filename)
		return
	}

	out, err := os.Create(path)
	if err != nil {
		fmt.Printf("[ERROR] Failed to create file %s: %v\n", path, err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] Failed to save %s: %v\n", path, err)
		return
	}

	fmt.Printf("[OK] Downloaded: %s (%s)\n", filename, contentType)
}



func IconPopulate() {
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		fmt.Println("Error creating uploads directory:", err)
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for _, domain := range iconDomains {
		wg.Add(1)
		go downloadIcon(domain, &wg, sem)
	}

	wg.Wait()
	fmt.Println("All icons downloaded.")
}
