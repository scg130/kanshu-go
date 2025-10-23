package spider

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"regexp"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)
type Book struct {
	DetailURL     string `json:"detail_url"`
	URLImg      string `json:"url_img"`
	ArticleName string `json:"articlename"`
	Author      string `json:"author"`
	Intro       string `json:"intro"`
}
var (
	replaceContext = `<p>喜欢.*?更新速度全网最快。</p>`
	cookie = `sex=boy; server_name_session=18baaff35356a21a8b00b52f2cc8a14f; Hm_lvt_56f6cec7fc20e74e50f246aadba47e1c=1761185633; HMACCOUNT=9971CCFA9B807CB6; novel_81963=459767%7C1761185651; novel_72936=0%7C1761186111; novel_36889=0%7C1761190289; 9shuyuan_user=%7B%22id%22%3A%22106172%22%2C%22name%22%3A%22scg130%22%2C%22pass%22%3A%2261f72d8d58f264aaec33f8820bfad6a7%22%2C%22time%22%3A1761193787%7D; vv=1761193802; qd_vt=1761193802; Hm_lpvt_56f6cec7fc20e74e50f246aadba47e1c=1761193803; RT="z=1&dm=9shuyuan.org&si=a7714d31-7221-47b4-8e37-2969068ef0d7&ss=mh2wf3qe&sl=c&tt=n2d&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&nu=3p7jcs99&cl=welm&r=23u5n0p9&ul=weln"`
)

// 每15秒请求一次
func searchNovel(novelName string) (string) {
	// 1. 表单参数
	formData := url.Values{}
	formData.Set("searchtype", "all")
	formData.Set("searchkey", novelName) // URL解码后的中文

	// 2. 构造请求
	req, err := http.NewRequest("POST", "https://www.9shuyuan.org/search.html", strings.NewReader(formData.Encode()))
	if err != nil {
		panic(err)
	}

	// 3. 设置 Header 模拟浏览器
	now := time.Now().Unix()
	// 匹配 9~10 位数字
	re := regexp.MustCompile(`\b\d{9,10}\b`)
	// 替换所有匹配的时间戳
	newCookie := re.ReplaceAllString(cookie, fmt.Sprintf("%d", now))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Origin", "https://www.9shuyuan.org")
	req.Header.Set("Referer", "https://www.9shuyuan.org/search/920/1.html")
	req.Header.Set("User-Agent", randomUA())
	req.Header.Set("Cookie", newCookie)

	// 4. 禁止自动跟随重定向，这样可以拿到 Location
	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// 5. 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 6. 获取 Location
	location := resp.Header.Get("Location")
	if location == "" {
		fmt.Println("没有 Location 头")
		return ""
	}

	// 7. 拼成完整 URL
	base, _ := url.Parse(req.URL.String())
	redirectURL, _ := url.Parse(location)
	fullURL := base.ResolveReference(redirectURL)
	fmt.Println("重定向后的完整 URL:", fullURL.String())
	return fullURL.String()
}

func JiuShuYuan(novelName string) error {
	names := strings.Split(novelName, "|")
	URL := searchNovel(names[0])
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

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err)
		return err
	}
	flag := true
	doc.Find("#sitembox").Find("dl").Each(func(i int, s *goquery.Selection) {
		book := Book{
		}
		detailURL,_ := s.Find("dt").Find("a").Attr("href")
		book.URLImg,_ = s.Find("dt").Find("a").Find("img").Attr("src")
		book.ArticleName,_ = s.Find("dt").Find("a").Attr("title")
		book.Author = s.Find(".book_other").First().Find("span").First().Text()
		book.Intro = s.Find(".book_des").Text()
		book.DetailURL = "https://www.9shuyuan.org" + detailURL
	
		if book.ArticleName != names[0] {
			return
		}
		if len(names) == 2 && book.Author != names[1] {
			return
		}
		flag = false
		jiuShuYuanNovel(book)
	})

	if flag {
		return errors.New("not found")
	}
	return nil
}

func jiuShuYuanNovel(book Book) {
	listURL := strings.TrimRight(book.DetailURL,"/")+".html"

	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		panic(err)
	}

	// 3. 设置 Header 模拟浏览器
	now := time.Now().Unix()
	// 匹配 9~10 位数字
	re := regexp.MustCompile(`\b\d{9,10}\b`)
	// 替换所有匹配的时间戳
	newCookie := re.ReplaceAllString(cookie, fmt.Sprintf("%d", now))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Origin", "https://www.9shuyuan.org")
	req.Header.Set("Referer", "https://www.9shuyuan.org/search/920/1.html")
	req.Header.Set("User-Agent", randomUA())
	req.Header.Set("Cookie", newCookie)

	// 4. 禁止自动跟随重定向，这样可以拿到 Location
	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// 5. 发送请求
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	img := book.URLImg
	title := book.ArticleName
	author := book.Author
	intro := book.Intro

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
	num := 1
	doc.Find("#list").Find("dl").Each(func(i int, dl *goquery.Selection) {
		dtCount := 0
		capture := false

		dl.Children().Each(func(j int, node *goquery.Selection) {
			if goquery.NodeName(node) == "dt" {
				dtCount++
				if dtCount == 2 {
					// 找到第二个 dt
					capture = true
				} else if dtCount > 2 {
					// 遇到第三个 dt，就停止捕获
					capture = false
				}
			}

			if capture && goquery.NodeName(node) == "dd" {
				if num <= int(chapterCurent) {
					num++
					return
				}
				chapterTitle := node.Find("a").Text()
				href, _ := node.Find("a").Attr("href")
				href = "https://www.9shuyuan.org" + href
				num++
				jiuShuYuanChapter(href, chapterTitle, num, int(novelId))
			}
		})
	})
}

func jiuShuYuanChapter(url, title string, num, novelId int) {
	var i = 0
loop:
	if i >= 3 {
		return
	}
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	// 3. 设置 Header 模拟浏览器
	now := time.Now().Unix()
	// 匹配 9~10 位数字
	re := regexp.MustCompile(`\b\d{9,10}\b`)
	// 替换所有匹配的时间戳
	newCookie := re.ReplaceAllString(cookie, fmt.Sprintf("%d", now))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Origin", "https://www.9shuyuan.org")
	req.Header.Set("Referer", "https://www.9shuyuan.org/search/920/1.html")
	req.Header.Set("User-Agent", randomUA())
	req.Header.Set("Cookie", newCookie)

	// 4. 禁止自动跟随重定向，这样可以拿到 Location
	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// 5. 发送请求
	rsp, err := client.Do(req)
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
	re = regexp.MustCompile(replaceContext)

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
	_, err = x.Exec(sql, num, len(content), novelId)
	if err != nil {
		logrus.Error(err)
	}
}
