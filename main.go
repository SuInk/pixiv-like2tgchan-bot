package main

import (
	"bookmark2tgchan-bot/config"
	"database/sql"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type Like struct {
	title       string
	description string
	url         string
	pubDate     string
	link        string
	src         []string
}

func init() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func Client() (http.Client, error) {
	// 设置clash代理
	uri, err := url.Parse(config.ProxyURL)
	if err != nil {
		log.Fatal(err)
		return http.Client{}, nil
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}
	return client, nil
}

// GetLikes 获取用户收藏
func GetLikes() ([]Like, error) {
	client, err := Client()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Get(config.RssURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// Load the XML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	likes := make([]Like, 0)
	doc.Find("item").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the content
		title := regexp.MustCompile("CDATA\\[(.*)\\]\\]>").FindStringSubmatch(s.Find("title").Text())[1]
		description := s.Find("description").Text()
		pubDate := s.Find("pubDate").Text()
		link := s.Find("guid").Text()
		srcs := make([]string, 0)
		s.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			srcs = append(srcs, src)
		})
		like := Like{
			title:       title,
			description: description,
			pubDate:     pubDate,
			link:        link,
			src:         srcs,
		}
		likes = append(likes, like)
		//fmt.Println(i, title, pubDate, link, description)
	})
	// log.Println(likes)
	return likes, nil
}
func GetDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func InitDb(db *sql.DB) *sql.DB {
	table := `CREATE TABLE IF NOT EXISTS likes (
        title TEXT NOT NULL,
    	description TEXT NOT NULL,
        pubDate DATE NULL,
    	link TEXT NOT NULL UNIQUE
    );`
	_, err := db.Exec(table)
	//defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func SaveLike(db *sql.DB, like Like) {
	sqlStr := `INSERT INTO likes (title, description, pubDate, link) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(sqlStr, like.title, like.description, like.pubDate, like.link)
	//defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
func Check(db *sql.DB, link string) bool {
	sqlStr := `SELECT * FROM likes WHERE link = ?`
	row := db.QueryRow(sqlStr, link)
	err := row.Scan(&link)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}
func SendMessage(like Like) {
	client, err := Client()
	if err != nil {
		log.Fatal(err)
	}
	formData := url.Values{
		"chat_id": {config.ChatID},
		"caption": {like.title + "\n" + like.link},
		"photo":   {like.src[0]},
	}
	tgBotUrl := "https://api.telegram.org/bot" + config.TgBotToken + "/sendPhoto"
	resp, err := client.PostForm(tgBotUrl, formData)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//body, err := io.ReadAll(resp.Body)
	//log.Println(string(body))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(like)
}
func main() {
	cursor := GetDb()
	cursor = InitDb(cursor)
	//t := time.Tick(time.Minute * 60)
	// 定时运行
	log.Println("start task...")
	//for _ = range t {
	likes, err := GetLikes()
	log.Println(likes)
	if err != nil {
		log.Fatal(err)
	}
	for _, like := range likes {
		// 如果数据库中存在该条记录，不处理
		if Check(cursor, like.link) {
			// SendMessage(like)
			log.Println("已存在", like.link)
		} else {
			SendMessage(like)
			SaveLike(cursor, like)
			log.Println("保存成功", like.link)
		}
	}
	//}
}
