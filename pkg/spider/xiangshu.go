package spider

import (
	"crypto/tls"
	"errors"
	"fmt"
	"kanshu/util"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func Xiangshu(novelName string) error {
	names := strings.Split(novelName, "|")
	url := "https://www.xbiquge.la/modules/article/waps.php?" + "searchkey=" + names[0]
	res, err := http.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != 200 {
		log.Println("request failed")
		return errors.New("request failed")
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	flag := true
	doc.Find(".grid").Find("tr").Each(func(i int, s *goquery.Selection) {
		detailUrl, _ := s.Find(".even").Find("a").Attr("href")
		name := s.Find(".even").Find("a").Text()
		if name != names[0] {
			return
		}
		flag = false
		xiangshuNovel(detailUrl)
	})
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
	img, _ := doc.Find("#fmimg").Find("img").Attr("src")
	title, _ := doc.Find("#info").Find("h1").Html()
	author, _ := doc.Find("#info").Find("p").First().Html()
	intro, _ := doc.Find("#intro").Find("p").Last().Html()
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
		chapterCurent = result["chapter_current"].(int64)
		novelId = result["id"].(int64)
	} else {
		sql = "insert into novel(name,author,img,intro,cate_id,new_chapter) values (?, ?, ?, ?, 5,'')"
		insertRes, err := x.Exec(sql, title, author, img, intro)
		if err == nil && insertRes != nil {
			novelId, _ = insertRes.LastInsertId()
		}
	}
	doc.Find("#list").Find("dd").Each(func(i int, s *goquery.Selection) {
		if i < int(chapterCurent) {
			return
		}
		chapterTitle := s.Find("a").Text()
		href, _ := s.Find("a").Attr("href")
		href = "https://www.ibiquge.la" + href
		xiangshuChapter(href, chapterTitle, i, int(novelId))
	})
}

func xiangshuChapter(url, title string, num, novelId int) {
	var i = 0
loop:
	if i >= 3 {
		return
	}
	fmt.Println(url)
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
	content, _ := doc.Find("#content").Html()
	bytes := regexp.MustCompile("<p><a.*</p>").ReplaceAll([]byte(content), []byte(""))
	content = string(bytes)
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
