package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/tidwall/gjson"
)

func dirviewer(w http.ResponseWriter, req *http.Request, c chan string) bool {
	fmt.Println("> dirviewerHandleFunc")
	fmt.Println(req.Method + " :" + req.Header.Get("X-Forwarded-For") + " " + req.RemoteAddr + " -> " + req.RequestURI)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("content-type", "application/x-www-form-urlencoded")
	// w.Header().Set("content-type", "multipart/form-data")
	// w.Header().Set("content-type", "text/plain")
	w.Header().Set("content-type", "application/json")

	reqURL := strings.Split(req.RequestURI, "?")
	if reqURL[0] == (suburl + dirviewerpath) {
		if req.Method == http.MethodOptions {
			w.WriteHeader(200)
			c <- http.MethodOptions
			return false
		}

		req.ParseMultipartForm(32 << 20)
		formcID, fid := req.Form["id"]
		formcSecret, fs := req.Form["secret"]
		formcPath, fp := req.Form["path"]
		missingParameter := ""
		if !fid {
			missingParameter += "id"
		}
		if !fs {
			if missingParameter != "" {
				missingParameter += ","
			}
			missingParameter += "secret"
		}
		if missingParameter != "" {
			w.WriteHeader(400)
			c <- code(201, "")
			return false
		}

		dirs := gjson.Get(conf, "dirs")
		if len(dirs.Array()) > 0 {
			dirinfo := gjson.Get(dirs.String(), formcID[0])
			if !dirinfo.Exists() {
				w.WriteHeader(400)
				c <- code(202, "")
				return false
			}
			password := gjson.Get(dirinfo.String(), "secret")
			if password.String() != formcSecret[0] {
				w.WriteHeader(400)
				c <- code(203, "")
				return false
			}
			bucketName := gjson.Get(dirinfo.String(), "bucket").String()
			dir := gjson.Get(dirinfo.String(), "dir").String() + "/"

			// 获取存储空间。
			bucket, err := client.Bucket(bucketName)
			if err != nil {
				fmt.Println("Error:", err)
			}
			// 列举文件。
			marker := ""
			if fp {
				dir += formcPath[0]
			}
			var files []map[string]interface{}
			for {
				lsRes, err := bucket.ListObjects(oss.MaxKeys(1000), oss.Marker(marker), oss.Prefix(dir), oss.Delimiter("/"))
				if err != nil {
					fmt.Println("Error:", err)
				}
				for _, dirName := range lsRes.CommonPrefixes {
					fileName := dirName[len(dir):]
					files = append(files, map[string]interface{}{
						"path": fileName,
						"size": 0,
						"type": "folder",
					})
				}
				// 打印列举文件，默认情况下一次返回100条记录。
				for _, object := range lsRes.Objects {
					if object.Key == dir {
						continue
					}
					fileName := object.Key[len(dir):]
					files = append(files, map[string]interface{}{
						"path": fileName,
						"size": object.Size,
						"type": "file",
					})
				}
				if lsRes.IsTruncated {
					marker = lsRes.NextMarker
				} else {
					break
				}
			}

			w.WriteHeader(200)
			c <- code(100, files)
			return false

		} else {
			w.WriteHeader(400)
			c <- code(300, "")
			return false
		}
	}
	return false
}

func dirviewerHandleFunc(w http.ResponseWriter, req *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	c := make(chan string)
	go dirviewer(w, req, c)
	re := <-c
	wg.Done()
	w.Write([]byte(re))
	// fmt.Fprint(w, re)
	wg.Wait()
}
