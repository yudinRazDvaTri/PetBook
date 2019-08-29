package controllers

import (
	"fmt"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/gorilla/context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func UploadFile() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Println("File Upload Endpoint Hit")
		id := context.Get(r, "id").(int)

		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			logger.Error(err)
			return
		}
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		path := "./web/static/usermedia/" + strconv.Itoa(id) + "/logo"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			_ = os.Mkdir("./web/static/usermedia/", os.ModeAppend)
			_ = os.Mkdir("./web/static/usermedia/"+strconv.Itoa(id), os.ModeAppend)
			_ = os.Mkdir(path, os.ModeAppend)
		}
		tempFile, err := ioutil.TempFile(path, "9*.png")
		if err != nil {
			logger.Error(err)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Error(err)
			return
		}
		tempFile.Write(fileBytes)
		//fmt.Fprintf(w, "Successfully Uploaded File\n")
		http.Redirect(w, r, "/edit", 301)
	}
}

func GetImgLogo(id int) []string {
	var files []string

	root := "./web/static/usermedia/"+strconv.Itoa(id)+"/logo"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".jpg" && filepath.Ext(path) != ".png" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		logger.Error(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
	for i,file:=range files{
		file=files[i]
		v:=strings.Replace(file,"\\","/",100)
		files[i]=v
	}
	for i,file:=range files{
		file=files[i]
		v:=strings.Replace(file,"web","..",100)
		files[i]=v
	}
	return files
}
//func GetImg() []string {
//	var files []string
//
//	root := "./web/static/usermedia"
//	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
//		if filepath.Ext(path) != ".jpg" && filepath.Ext(path) != ".png" {
//			return nil
//		}
//		files = append(files, path)
//		return nil
//	})
//	if err != nil {
//		panic(err)
//	}
//	for _, file := range files {
//		fmt.Println(file)
//	}
//	for i,file:=range files{
//		file=files[i]
//		v:=strings.Replace(file,"\\","/",100)
//		files[i]=v
//	}
//	for i,file:=range files{
//		file=files[i]
//		v:=strings.Replace(file,"web","..",100)
//		files[i]=v
//	}
//	return files
//}