package spider

import (
	"crypto/tls"
	"errors"
	"fmt"
	"kanshu/util"
	"log"
	"net/http"
	"strings"
	"time"
	"net/url"
	"io"
	"encoding/json"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)
type Book struct {
	URLList     string `json:"url_list"`
	URLImg      string `json:"url_img"`
	ArticleName string `json:"articlename"`
	Author      string `json:"author"`
	Intro       string `json:"intro"`
}

func Xiangshu(novelName string) error {
	names := strings.Split(novelName, "|")
	URL := "https://www.53122c.cfd/user/search.html?q=%s&so=undefined"
	encodedKeyword := url.QueryEscape(names[0])
	URL = fmt.Sprintf(URL, encodedKeyword)
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", randomUA())
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != 200 {
		log.Println("request failed1111")
		return errors.New("request failed1111")
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("read body failed:", err)
		return err
	}
	var books []Book
	err = json.Unmarshal(bodyBytes, &books)
	if err != nil {
		log.Println("解析 JSON 失败:", err)
		return err
	}
	flag := true
	// 遍历打印
	for _, book := range books {
		detailUrl := "https://www.53122c.cfd" + book.URLList
		name := book.ArticleName
		author := book.Author
		author = strings.Trim(author,"作者：")
		if name != names[0] {
			continue
		}
		if len(names) == 2 && author != names[1] {
			continue
		}
		flag = false
		xiangshuNovel(detailUrl)
	}

	if flag {
		return errors.New("not found")
	}
	return nil
}

func xiangshuNovel(url string) {
	res, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return
	}
	if res.StatusCode != 200 {
		logrus.Error("request fail")
		return
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	img, _ := doc.Find(".book").Find("cover").Attr("src")
	title, _ := doc.Find(".book").Find("h1").Html()
	author, _ := doc.Find(".book").Find("small").Find("span").First().Html()
	intro := doc.Find(".book").Find(".intro").Find("dd").Last().Text()
	intro = util.CutChineseString(intro, 40)
	sql := "select * from novel.novel where name = ?"
	result := make(map[string]interface{})
	exist, err := x.Table("novel.novel").Where("name = ?", title).Get(&result)
	if err != nil {
		logrus.Error(err)
		return
	}
	var novelId, chapterCurent int64
	if exist {
		chapterCurent = int64(result["chapter_current"].(int32))
		novelId = int64(result["id"].(int32))
	} else {
		sql = "insert into novel(name,author,img,intro,cate_id) values (?, ?, ?, ?, 5)"
		insertRes, err := x.Exec(sql, title, author, img, intro)
		if err == nil && insertRes != nil {
			novelId, _ = insertRes.LastInsertId()
		}
	}
	doc.Find(".listmain").Find("dd").Each(func(i int, s *goquery.Selection) {
		if i < int(chapterCurent) {
			return
		}
		chapterTitle := s.Find("a").Text()
		href, _ := s.Find("a").Attr("href")
		href = "https://www.53122c.cfd" + href


		xiangshuChapter(href, chapterTitle, i, int(novelId))
	})
}

var (
	replaceContext1 = `请收藏本站：https://www.53122c.cfd。笔趣阁手机版：https://m.53122c.cfd <br/><br/>`
)

func xiangshuChapter(url, title string, num, novelId int) {
	var i = 0
loop:
	if i >= 3 {
		return
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 15,
	}
	rsp, err := client.Get(url)
	if err != nil {
		i++
		time.Sleep(time.Second)
		logrus.Error(err)
		goto loop
	}

	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		i++
		time.Sleep(time.Second)
		logrus.Error("request fail")
		goto loop
	}

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		i++
		time.Sleep(time.Second)
		logrus.Error(err)
		goto loop
	}
	content, _ := doc.Find("#chaptercontent").Html()
	content = strings.ReplaceAll(content, replaceContext1, "")
	re := regexp.MustCompile(`<p class="readinline">.*?</p>`)

    content = re.ReplaceAllString(content, "")
	if len(content) < 1000 {
		return
	}
	sql := "insert into chapter(title,content,novel_id,num,words) values (?, ?, ?, ?, ?)"
	_, err = x.Exec(sql, title, content, novelId, num+1, len(content))
	if err != nil {
		logrus.Error(err)
		return
	}
	sql = "update novel set chapter_total = chapter_total+1,chapter_current=?,words = words+? where id = ?"
	_, err = x.Exec(sql, num+1, len(content), novelId)
	if err != nil {
		logrus.Error(err)
	}
}
