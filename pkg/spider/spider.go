package spider

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"kanshu/util"
	"xorm.io/xorm"
)

const URL = "http://m.biquge5200.cc/modules/article/waps.php?keyword=%s"

var (
	lockMap *sync.Map
	enc     mahonia.Decoder
	x       *xorm.Engine
	client  *http.Client // ✅ 全局 HTTP Client
)

func init() {
	enc = mahonia.NewDecoder("gbk")
	initDB()
	lockMap = &sync.Map{}
	client = &http.Client{
		Timeout: 10 * time.Second, // ✅ 设置全局超时
	}
}

func Run(pyName, zhName string) error {
	key := fmt.Sprintf("%s:%s", pyName, zhName)
	if _, ok := lockMap.Load(key); ok {
		logrus.Error("task exist!")
		return errors.New("task exist!")
	}
	lockMap.Store(key, 1)
	defer lockMap.Delete(key)
	err := JiuShuYuan(zhName)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func find(name string) error {
	names := strings.Split(name, "|")
	encodedKeyword := url.QueryEscape(names[0])
	url := fmt.Sprintf(URL, encodedKeyword)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", randomUA())
	req.Header.Set("Referer", "http://m.biquge5200.cc/")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logrus.Errorf("请求失败，状态码: %d", res.StatusCode)
		return errors.New("fail")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err)
		return err
	}

	flag := true
	doc.Find(".cover").Find("p").Each(func(i int, s *goquery.Selection) {
		author := enc.ConvertString(s.Find("a").Last().Text())
		novelName := enc.ConvertString(s.Find(".blue").Text())
		detailUrl, _ := s.Find(".blue").Attr("href")
		if author == names[1] && novelName == names[0] {
			flag = false
			novel(fmt.Sprintf("http://m.abcbiquge.com%s", detailUrl), author)
		}
	})
	if flag {
		return errors.New("not found")
	}
	return nil
}

func novel(url string, author string) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", randomUA())
	res, err := client.Do(req) // ✅ 替换 http.Get
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logrus.Errorf("请求失败 novel 状态码: %d", res.StatusCode)
		return errors.New("fail")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err)
		return err
	}

	img, _ := doc.Find(".block_img2").Find("img").Attr("src")
	title, _ := doc.Find("#block_txt2").Find("h2").Find("a").Html()
	title = enc.ConvertString(title)
	intro, _ := doc.Find("#intro_info").Html()
	intro = util.CutChineseString(intro, 40)

	sql := "select * from novel.novel where name = ?"
	result := make(map[string]interface{})
	exist, err := x.Table("novel.novel").Where("name = ?", title).Get(&result)
	if err != nil {
		logrus.Error(err)
		return err
	}

	var novelId, chapterCurent int64
	if exist {
		chapterCurent = result["chapter_current"].(int64)
		novelId = result["id"].(int64)
	} else {
		sql = "insert into novel(name,author,img,intro,cate_id) values (?, ?, ?, ?, 1)"
		insertRes, err := x.Exec(sql, title, author, img, intro)
		if err == nil && insertRes != nil {
			novelId, _ = insertRes.LastInsertId()
		}
	}

	listUrl, _ := doc.Find(".ablum_read").First().Find("span").Last().Find("a").Attr("href")
	url1 := "http://m.abcbiquge.com" + listUrl[:len(listUrl)-1] + "_1/"

	req, _ = http.NewRequest("GET", url1, nil)
	req.Header.Set("User-Agent", randomUA())
	res, err = client.Do(req) // ✅ 替换 http.Get
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logrus.Error("request fail 3333")
		return errors.New("fail")
	}

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err)
		return err
	}

	page := enc.ConvertString(doc.Find(".page").Last().Text())
	pages, _ := strconv.Atoi(page[strings.Index(page, "/")+1 : strings.Index(page, ")")-3])

	for i := 1; i <= pages; i++ {
		chaptersUrl := fmt.Sprintf("http://m.abcbiquge.com/"+listUrl[:len(listUrl)-1]+"_%d/", i)
		chapters(chaptersUrl, chapterCurent, novelId, i)
	}
	return nil
}

func chapters(url string, chapterCurent int64, novelId int64, page int) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", randomUA())
	res, err := client.Do(req) // ✅ 替换 http.Get
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logrus.Error("request fail 5555")
		return errors.New("fail")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err)
		return err
	}

	doc.Find(".chapter").Find("li").Each(func(i int, s *goquery.Selection) {
		if i+(20*page-1) < int(chapterCurent) {
			return
		}
		href, _ := s.Find("a").Attr("href")
		chapter(fmt.Sprintf("http://m.abcbiquge.com%s", href), enc.ConvertString(s.Text()), i, int(novelId))
	})
	return nil
}

func chapter(url, title string, num, novelId int) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", randomUA())
	res, err := client.Do(req) // ✅ 替换 http.Get
	if err != nil {
		logrus.Error(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logrus.Error("request fail chapter 6666")
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	content, _ := doc.Find(".text").Html()
	content = enc.ConvertString(content)
	sql := "insert into chapter(title,content,novel_id,num,words) values (?, ?, ?, ?, ?)"
	var t string
	if len(strings.Split(title, " ")) == 2 {
		t = strings.Split(title, " ")[1]
	}
	_, err = x.Exec(sql, t, content, novelId, num+1, len(content))
	if err != nil {
		logrus.Error(err)
		return
	}
	sql = "update novel set chapter_total = chapter_total+1,chapter_current=?,words = words+? where id = ?"
	_, err = x.Exec(sql, num+1, len(content), novelId)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func randomUA() string {
	uaList := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
	}
	return uaList[rand.Intn(len(uaList))]
}

// const (
// 	user = "root"
// 	pass = "smd123456"
// 	db   = "novel"
// 	port = "3306"
// 	// host = "60.205.191.37"
// 	host = "127.0.0.1"
// )

func initDB() {
	host := GetEnv("MYSQL_HOST","127.0.0.1")
	user := GetEnv("MYSQL_USER","root")
	pass := GetEnv("MYSQL_PASSWD","smd013012")
	db := GetEnv("MYSQL_NOVEL_DB","novel")
	mysql_log := GetEnv("MYSQL_LOG","true") == "true"
	port := GetEnv("MYSQL_PORT","3306")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, pass, host, port, db)
	var err error
	x, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatalf("mysql connect err:%s", err.Error())
	}
	if err := x.Ping(); err != nil {
		log.Fatalf("mysql ping err:%s", err.Error())
	}
	x.ShowSQL(mysql_log)

	x.SetMaxIdleConns(5)

	x.SetMaxOpenConns(5)
}

func GetEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}