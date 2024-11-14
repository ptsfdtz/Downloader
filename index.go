package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	OnlineToken string `json:"online_token"`
}

type VideoData struct {
	Seasons []struct {
		ID string `json:"id"`
	} `json:"seasons"`
}

type SeasonData struct {
	Episodes []struct {
		URL string `json:"url"`
	} `json:"episodes"`
}

func encodeUrlPartially(u string) string {
	u = strings.ReplaceAll(u, " ", "%20")
	u = strings.ReplaceAll(u, "(", "%28")
	u = strings.ReplaceAll(u, ")", "%29")
	return u
}

func getData(apiURL, token string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "online_token="+token)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	if req.Method != "GET" {
		fmt.Println("请求失败")
	}
	data, _ := io.ReadAll(resp.Body)
	return data
}

func fetchAndDownloadFiles(videoId string) {
	configFile, _ := os.ReadFile("config.json")
	var config Config
	json.Unmarshal(configFile, &config)

	url := fmt.Sprintf("https://online.njtech.edu.cn/api/v2/videos/%s?id=%s", videoId, videoId)
	videoDataJson := getData(url, config.OnlineToken)
	var videoData VideoData
	json.Unmarshal(videoDataJson, &videoData)

	seasonId := videoData.Seasons[0].ID
	apiUrl := fmt.Sprintf("https://online.njtech.edu.cn/api/v2/video_seasons/%s?order=asc&orderBy=index", seasonId)
	seasonDataJson := getData(apiUrl, config.OnlineToken)
	var seasonData SeasonData
	json.Unmarshal(seasonDataJson, &seasonData)

	dirPath := filepath.Join("media", videoId)
	os.MkdirAll(dirPath, os.ModePerm)

	for _, episode := range seasonData.Episodes {
		url := encodeUrlPartially(episode.URL)
		fmt.Printf("正在下载: %s\n", url)

		fileName := filepath.Base(url)
		fileName = strings.ReplaceAll(fileName, "%20", " ")
		filePath := filepath.Join(dirPath, fileName)

		out, _ := os.Create(filePath)
		defer out.Close()
		resp, _ := http.Get(url)
		if resp.StatusCode != 200 {
			fmt.Println("下载失败")
		}
		defer resp.Body.Close()
		io.Copy(out, resp.Body)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入视频 ID: ")
	videoId, _ := reader.ReadString('\n')
	videoId = strings.TrimSpace(videoId)

	if videoId != "" {
		fetchAndDownloadFiles(videoId)
		fmt.Println("下载完成。按回车键退出")
		reader.ReadString('\n')
	} else {
		fmt.Println("视频 ID 不能为空")
	}
}
