package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/tomatobobot/dict/db"
)

const (
	// BingDictURL bing dict url
	BingDictURL = "https://cn.bing.com/dict/search?q="
	// DBPath Data file path
	DBPath = "dict.db"
	// NONE none
	NONE = "\033[00m"
	// BOLD bold
	BOLD = "\033[1m"
)

// request 到bing在线词典查询单词释意
func request(word string) (map[string]string, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodGet, BingDictURL+word, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	dictMap := make(map[string]string)
	doc.Find(".qdef ul li").Each(func(i int, s *goquery.Selection) {
		pos := s.Find("span").First().Text()
		def := s.Find("span").Last().Text()
		dictMap[pos] = def
	})
	return dictMap, nil
}
func newDB(path string) (*db.DB, error) {
	db := &db.DB{}

	if err := db.Open(path); err != nil {
		return nil, err
	}
	return db, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// printDictBytes 打印dict字典 先将bytes转换成map[string]string
func printDictBytes(b []byte) {
	var dict map[string]string
	json.Unmarshal(b, &dict)
	printDict(dict)
}

// printDict 打印dict字典
func printDict(dict map[string]string) {
	for k, v := range dict {
		fmt.Printf("%4s:%s\n", color.RedString(k), v)
	}
}
func printWord(word string) {
	fmt.Printf("%s%s %s\n", BOLD, word, NONE)
}
func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("请输入要查询的单词")
		return
	}
	word := args[0]
	data, err := newDB(DBPath)
	checkErr(err)

	defer data.Close()
	var result []byte
	err = data.View(func(tx *db.Tx) error {
		d, err := tx.Dict([]byte(word))
		if err != nil {
			return err
		}
		result = d.Result
		return nil
	})
	// 没有错误说明已数据库里已存在此单词，故直接打印保存的单词信息并退出程序
	if err == nil {
		printWord(word)
		printDictBytes(result)
		return
	}
	paraphraseMap, err := request(word)
	checkErr(err)
	// 将查询到的单词释意存入数据库
	data.Update(func(tx *db.Tx) error {
		result, err := json.Marshal(paraphraseMap)
		if err != nil {
			return err
		}
		d := &db.Dict{
			Tx:     tx,
			Word:   []byte(word),
			Result: result,
		}
		d.Save()
		return nil
	})
	printWord(word)
	printDict(paraphraseMap)
}
