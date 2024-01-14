package goelp

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/johnbalvin/goelp/trace"
)

func TestSaveOnDisk() {
	test1()
}
func TestNoImages() ([]Data, error) {
	return test2(false)
}

func TestImages() ([]Data, error) {
	return test2(true)
}
func test3() {
	client := DefaulClient()
	client.Language = "en_US"
	searchValue := "law"
	datas, err := client.Search(30, searchValue)
	if err != nil {
		log.Println("err: ", err)
		return
	}
	f, _ := os.Create("./datas.json")
	json.NewEncoder(f).Encode(datas)
}

func test1() {
	//Make sure you have write permissions
	if err := os.MkdirAll("./test/0", 0644); err != nil {
		log.Println("test2:1 -> err: ", err)
		return
	}
	client := DefaulClient()
	client.Location = "San Francisco, CA"
	client.Language = "es_MX"
	searchValue := "law"
	searchResults, err := client.Search(30, searchValue)
	if err != nil {
		log.Println("test2:2 -> err: ", err)
		return
	}
	f, _ := os.Create("./test/searchResults.json")
	json.NewEncoder(f).Encode(searchResults)
	f.Close()
	var datas []Data
	for i, item := range searchResults.Items {
		folderPath := fmt.Sprintf("./test/%d/images", i)
		os.MkdirAll(folderPath, 0644)
		data, err := GetFromYelpBizURL(item.YelpURL, nil)
		if err != nil {
			log.Println("test2:3 -> err: ", err)
			continue
		}
		if err := data.SetImages(client.ProxyURL); err != nil {
			log.Println("test2:4 -> err: ", err)
			continue
		}
		for j, img := range data.Images {
			fname3 := fmt.Sprintf("./test/%d/images/%d%s", i, j, img.Extension)
			os.WriteFile(fname3, img.Content, 0644)
		}
		datas = append(datas, data)
		log.Printf("Progress getting data: %d/%d\n", i+1, len(datas))
		filePath := fmt.Sprintf("./test/%d/data.json", i)
		f, err := os.Create(filePath)
		if err != nil {
			log.Println("test2:5 -> err: ", err)
			continue
		}
		json.NewEncoder(f).Encode(data)
		f.Close()
	}
	f2, _ := os.Create("./test/datas.json")
	json.NewEncoder(f2).Encode(datas)
}

func test2(images bool) ([]Data, error) {
	//Make sure you have write permissions
	if err := os.MkdirAll("./test/0", 0644); err != nil {
		return nil, trace.NewOrAdd(1, "main", "test2", err, "")
	}
	//location := "San Francisco, CA"
	client := DefaulClient()
	client.Language = "es_MX"
	searchValue := "law"
	searchResults, err := client.Search(30, searchValue)
	if err != nil {
		return nil, trace.NewOrAdd(2, "main", "test2", err, "")
	}
	f, _ := os.Create("./test/searchResults.json")
	json.NewEncoder(f).Encode(searchResults)
	f.Close()
	var datas []Data
	for i, item := range searchResults.Items {
		folderPath := fmt.Sprintf("./test/%d/images", i)
		os.MkdirAll(folderPath, 0644)
		data, err := GetFromYelpBizURL(item.YelpURL, nil)
		if err != nil {
			errData := trace.NewOrAdd(3, "main", "test2", err, "")
			log.Println(errData)
			continue
		}
		if images {
			if err := data.SetImages(client.ProxyURL); err != nil {
				errData := trace.NewOrAdd(4, "main", "test2", err, "")
				log.Println(errData)
			}
		}
		datas = append(datas, data)
		log.Printf("Progress getting data: %d/%d id: %s\n", i+1, len(searchResults.Items), item.YelpURL)
	}
	return datas, nil
}
