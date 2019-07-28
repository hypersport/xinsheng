package xinsheng

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type WebPage struct {
	DayPosts   int
	WeekPosts  int
	MonthPosts int
	YearPosts  int
	TodayPosts []Posts
}

type Posts struct {
	Url         string
	Title       string
	Description string
}

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	http.Handle("/css/", http.FileServer(http.Dir("web/static")))
	http.Handle("/js/", http.FileServer(http.Dir("web/static")))
	http.Handle("/image/", http.FileServer(http.Dir("web/static")))

	logFile, err := os.OpenFile("log/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0766)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}

	Info = log.New(os.Stdout, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, logFile), "Error:", log.Ldate|log.Ltime|log.Lshortfile)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	article := WebPage{
		DayPosts:   getPostsNum("day"),
		WeekPosts:  getPostsNum("week"),
		MonthPosts: getPostsNum("month"),
		YearPosts:  getPostsNum("year"),
		TodayPosts: getTodayPosts(),
	}
	t, err := template.ParseFiles("web/template/index.html")
	if err != nil {
		Error.Fatal(err)
		fmt.Fprintf(w, "页面没有准备好，请稍后再访问 ...")
		return
	}

	err = t.Execute(w, article)
	if err != nil {
		Error.Fatal(err)
		fmt.Fprintf(w, "页面没有准备好，请稍后再访问 ...")
	}
}

func parseWebPage(postType string) (*goquery.Document, error) {
	format := "http://xinsheng.huawei.com/cn/index.php?app=search&mod=Isearch&act=index&key=惯例&type=&filter_type=topic&sort=createtime&filter_nodes=&filter_date=%s&filter_location=&ipage="

	url := fmt.Sprintf(format, postType)
	res, err := http.Get(url)
	if err != nil {
		Error.Println("Open URL Error")
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		Warning.Printf("Status Code is %d\n", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		Error.Println("Open Web Page Error")
		return nil, err
	}

	return doc, nil
}

func parseResultNum(doc *goquery.Document) int {
	var result string
	doc.Find("div.search-header-filter").Each(func(_ int, selection *goquery.Selection) {
		result = selection.Text()
	})

	sum, tmp := 0, 0
	for _, char := range result {
		if char == 39033 {
			break
		}
		if char >= 48 && char <= 57 {
			tmp, _ = strconv.Atoi(string(char))
			sum = sum*10 + tmp
		}
	}
	return sum
}

func getPostsNum(postType string) int {
	doc, err := parseWebPage(postType)
	if err != nil {
		Warning.Println(err)
		return -1
	}
	return parseResultNum(doc)
}

func getTodayPosts() []Posts {
	result := make([]Posts, 10)
	todayPost := Posts{}
	doc, err := parseWebPage("day")
	if err != nil {
		Error.Println(err)
		return nil
	}
	doc.Find("div.itemDiv").Each(func(_ int, s *goquery.Selection) {
		todayPost.Url, _ = s.Find("a").Attr("href")
		todayPost.Title, _ = s.Find("a").Attr("title")
		todayPost.Description = s.Find(".discription").Text()
		result = append(result, todayPost)
	})
	return result
}
