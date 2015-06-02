package robots

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Info struct {
	Pic string
	Bio string
}

const infoTempl = "http://base-generate.bijint.com/%s/"
const picTempl = "http://www.bijint.com/%s/tokei_images/"

var links = map[string]Info{
	"japan":     Info{Pic: fmt.Sprintf(picTempl, "jp")},
	"thailand":  Info{Pic: fmt.Sprintf(picTempl, "jp/img/clk")},
	"2012":      Info{Pic: fmt.Sprintf(picTempl, "2012jp")},
	"2011":      Info{Pic: fmt.Sprintf(picTempl, "2011jp")},
	"tokyo":     Info{Pic: fmt.Sprintf(picTempl, "tokyo")},
	"hokkaido":  Info{Pic: fmt.Sprintf(picTempl, "hokkaido")},
	"sendai":    Info{Pic: fmt.Sprintf(picTempl, "sendai")},
	"akita":     Info{Pic: fmt.Sprintf(picTempl, "akita")},
	"gunma":     Info{Pic: fmt.Sprintf(picTempl, "gunma")},
	"niigata":   Info{Pic: fmt.Sprintf(picTempl, "niigata")},
	"kanazawa":  Info{Pic: fmt.Sprintf(picTempl, "kanazawa")},
	"fukui":     Info{Pic: fmt.Sprintf(picTempl, "fukui")},
	"nagoya":    Info{Pic: fmt.Sprintf(picTempl, "nagoya")},
	"kyoto":     Info{Pic: fmt.Sprintf(picTempl, "kyoto")},
	"osaka":     Info{Pic: fmt.Sprintf(picTempl, "osaka")},
	"kobe":      Info{Pic: fmt.Sprintf(picTempl, "kobe")},
	"okayama":   Info{Pic: fmt.Sprintf(picTempl, "okayama")},
	"kagawa":    Info{Pic: fmt.Sprintf(picTempl, "kagawa")},
	"fukuoka":   Info{Pic: fmt.Sprintf(picTempl, "fukuoka")},
	"kagoshima": Info{Pic: fmt.Sprintf(picTempl, "kagoshima")},
	"okinawa":   Info{Pic: fmt.Sprintf(picTempl, "okinawa")},
	"kumamoto":  Info{Pic: fmt.Sprintf(picTempl, "kumamoto")},
	"saitama":   Info{Pic: fmt.Sprintf(picTempl, "saitama")},
	"hiroshima": Info{Pic: fmt.Sprintf(picTempl, "hiroshima")},
	"chiba":     Info{Pic: fmt.Sprintf(picTempl, "chiba")},
	"nara":      Info{Pic: fmt.Sprintf(picTempl, "nara")},
	"yamaguchi": Info{Pic: fmt.Sprintf(picTempl, "yamaguchi")},
	"nagano":    Info{Pic: fmt.Sprintf(picTempl, "nagano")},
	"shizuoka":  Info{Pic: fmt.Sprintf(picTempl, "shizuoka")},
	"miyazaki":  Info{Pic: fmt.Sprintf(picTempl, "miyazaki")},
	"tottori":   Info{Pic: fmt.Sprintf(picTempl, "tottori")},
	"iwate":     Info{Pic: fmt.Sprintf(picTempl, "iwate")},
	"ibaraki":   Info{Pic: fmt.Sprintf(picTempl, "ibaraki")},
	"tochigi":   Info{Pic: fmt.Sprintf(picTempl, "tochigi")},
	"taiwan":    Info{Pic: fmt.Sprintf(picTempl, "taiwa/")},
	"hawaii":    Info{Pic: fmt.Sprintf(picTempl, "hawaii")},
	"seifuku":   Info{Pic: fmt.Sprintf(picTempl, "seifuku")},
	"megane":    Info{Pic: fmt.Sprintf(picTempl, "megane")},
	"sara":      Info{Pic: fmt.Sprintf(picTempl, "sara")},
	"hairstyle": Info{Pic: fmt.Sprintf(picTempl, "hairstyle")},
	"asahi":     Info{Pic: fmt.Sprintf(picTempl, "tv-asahi")},
	"circuit":   Info{Pic: fmt.Sprintf(picTempl, "cc")},
	"hanayome":  Info{Pic: fmt.Sprintf(picTempl, "hanayome")},
	"waseda":    Info{Pic: fmt.Sprintf(picTempl, "wasedastyle")},
}

type BijinBot struct {
	Location *time.Location
}

type BijinConfiguration struct {
	Timezone string `json:"timezone"`
}

var BijinConfig = new(BijinConfiguration)

func init() {
	flag.Parse()
	configFile := filepath.Join(*ConfigDirectory, "bijin.json")
	if _, err := os.Stat(configFile); err == nil {
		config, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Printf("ERROR: Error opening bijin config: %s", err)
			return
		}
		err = json.Unmarshal(config, BijinConfig)
		if err != nil {
			log.Printf("ERROR: Error parsing bijin config: %s", err)
			return
		}
	} else {
		log.Printf("WARNING: Could not find configuration file bijin.json in %s", *ConfigDirectory)
	}
	b := &BijinBot{}
	if BijinConfig.Timezone != "" {
		loc, err := time.LoadLocation(BijinConfig.Timezone)
		if err != nil {
			log.Printf("Error loading timezone %q, falling back to UTC (%s)", BijinConfig.Timezone, err.Error())
			b.Location = time.UTC
		} else {
			b.Location = loc
		}
	}

	RegisterRobot("bijin", b)
}

func (r BijinBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	go r.DeferredAction(p)
	return ""
}

func (r BijinBot) DeferredAction(p *Payload) {
	response := &IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "Bijin Bot",
		IconEmoji:   ":ghost:",
		UnfurlLinks: true,
	}

	rand.Seed(time.Now().UTC().UnixNano())
	t := time.Now().In(r.Location)
	hours := fmt.Sprintf("%02d", t.Hour())
	minutes := fmt.Sprintf("%02d", t.Minute())
	region, link, _ := GetLink(strings.ToLower(strings.TrimSpace(p.Text)))
	// b = GetInfo(info, hours, minutes)
	response.Text = fmt.Sprintf("<@%s|%s> Here's your <%s%s%s.jpg|%s:%s 美人 (%s)> ", p.UserID, p.UserName, link, hours, minutes, hours, minutes, strings.ToTitle(region))
	response.Send()
}

// type Bijin struct {
// 	Birthday time.Time
// }

// func GetInfo(info, hours, minutes string) Bijin {
// 	const u = "%scache/%s%s.html"
// 	b := &Bijin{}
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", fmt.Sprintf(u, info, hours, minutes), nil)
// 	if err != nil {
// 		return b
// 	}
// 	req.Header.Add("Host", "http://www.bijint.com")
// 	req.Header.Add("Referer", info)
// 	resp, err := client.Do(req)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		return b
// 	}
// 	t := html.NewTokenizerFragment(resp.Body, "div")
// 	p := false
// 	for {
// 		if t.Next() == html.ErrorToken {
// 			// Returning io.EOF indicates success.
// 			break
// 		}
// 		token := t.Token()
// 		if !p && token.Type == html.CommentToken && strings.Contains(token.String(), "profile") {
// 			p = true
// 		}
// 		log.Println(token.String())
// 	}
// 	return b
// }

func GetLink(region string) (string, string, string) {
	var r string
	var i Info
	if i, ok := links[region]; ok {
		return region, i.Pic, i.Bio
	}
	for r, i = range links {
		break
	}
	return r, i.Pic, i.Bio
}

func (r BijinBot) Description() (description string) {
	return "Displays the current time's 美人 (hope bijint.com doesn't get mad me)\n\tUsage: /bijin\n\tExpected Response: (beautiful woman)"
}
