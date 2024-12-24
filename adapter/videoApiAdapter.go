package adapter

import (
	"encoding/json"
	"fmt"
	"instaLinkCollector/dto"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strings"
)

type VideoApiAdapter struct{}

func NewVideoService() *VideoApiAdapter { return &VideoApiAdapter{} }

func (v *VideoApiAdapter) transformURL(dynamicUrl string) (string, error) {
	checkUrl := strings.ReplaceAll(dynamicUrl, " ", "")

	if checkUrl == "" {
		return "", fmt.Errorf("error: empty URL")
	}

	var dynamicURLBuilder strings.Builder
	slashCount := 0

	for i := 0; i < len(checkUrl); i++ {
		if checkUrl[i] == '/' {
			slashCount++
		}
		dynamicURLBuilder.WriteByte(checkUrl[i])

		if slashCount == 5 {
			break
		}
	}

	finalURL := dynamicURLBuilder.String() + "?__a=1&__d=dis"
	return finalURL, nil
}

func (v *VideoApiAdapter) FetchInstagramWithCookies(urlInsta string) (r []dto.VideoResponse, err error) {
	//Inisialisasi cookiejar untuk cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error creating cookie jar: %v", err)
		return nil, err
	}

	//Karena perlu cookies kombinasi http client dengan cookie jar
	client := &http.Client{
		Transport: &http.Transport{},
		Jar:       jar,
	}

	//sementara isi url ini dulu
	instagramURL, err := v.transformURL(urlInsta)
	if err != nil {
		return nil, err
	}

	//buat request HTTP dan tambahkan cookies (session_id yang valid)
	req, err := http.NewRequest("GET", instagramURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
		return nil, err
	}

	//menggunakan cookies (session_id) agar bisa akses link
	req.AddCookie(&http.Cookie{
		Name:  "sessionid",
		Value: os.Getenv("INSTA_SESSIONID"), // Ganti dengan session_id yang valid
	})

	//kirim request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	//memastikan berhasil/200
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to load page, status code: %d", resp.StatusCode)
		return nil, err
	}

	//baca response body ke dalam strings.Builder
	var htmlContent strings.Builder
	_, err = io.Copy(&htmlContent, resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return nil, err
	}

	//hapus 'amp;' dari html kalo ada (terutama di url)
	htmlContentStr := htmlContent.String()
	htmlContentStr = strings.ReplaceAll(htmlContentStr, "amp;", "")

	//regex mencari json yang berisi video
	re := regexp.MustCompile(`"video_versions":(\[.*?\])`)
	matches := re.FindStringSubmatch(htmlContentStr)
	if len(matches) < 2 {
		log.Fatalf("No video_versions found")
		return nil, err
	}

	//parse json untuk mendapatkan data yang dibutuhkan
	var videoVersions []dto.VideoVersion
	err = json.Unmarshal([]byte(matches[1]), &videoVersions)
	if err != nil {
		log.Fatalf("Error parsing video_versions: %v", err)
		return nil, err
	}

	//isi videoReponse
	var result []dto.VideoResponse
	for _, v := range videoVersions {
		result = append(result, dto.VideoResponse{
			Resolution: fmt.Sprintf("%dx%d", v.Width, v.Height),
			URL:        v.URL,
		})
	}

	return result, nil
}
