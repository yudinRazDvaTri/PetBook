package controllers

import (
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MediaStore struct {
	DB *sqlx.DB
}

func (c *Controller) UploadMedia() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		id := context.Get(r, "id").(int)

		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("myMedia")
		if err != nil {
			logger.Error(err)
			return
		}
		defer file.Close()
		path := "./web/static/usermedia/" + strconv.Itoa(id) + "/gallery"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			_ = os.Mkdir("./web/static/usermedia/", os.ModeAppend)
			_ = os.Mkdir("./web/static/usermedia/"+strconv.Itoa(id), os.ModeAppend)
			_ = os.Mkdir(path, os.ModeAppend)
		}
		tempFile, err := ioutil.TempFile(path, "*.png")
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
		renamedFiles:= folderMediaPath(id)
		for _, element := range renamedFiles {
			c.MediaStore.AddMediaPathDb(element,id)
		}
		http.Redirect(w, r, "/edit", 301)
	}
}

func folderMediaPath(id int) []string {
		var files []string
		root := "./web/static/usermedia/"+strconv.Itoa(id)+"/gallery"
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
		return changePath(files)
	}


func (c *Controller) UploadLogo() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		id := context.Get(r, "id").(int)

		r.ParseMultipartForm(10 << 20)
		file, _, err := r.FormFile("myFile")
		if err != nil {
			logger.Error(err)
			return
		}
		defer file.Close()
		path := "./web/static/usermedia/" + strconv.Itoa(id) + "/logo"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			_ = os.Mkdir("./web/static/usermedia/", os.ModeAppend)
			_ = os.Mkdir("./web/static/usermedia/"+strconv.Itoa(id), os.ModeAppend)
			_ = os.Mkdir(path, os.ModeAppend)
		}
		tempFile, err := ioutil.TempFile(path, "*.png")
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
		renamedFiles:= folderLogoPath(id)
		for _, element := range renamedFiles {
			c.MediaStore.AddLogoPathDb(element,id)
		}
		http.Redirect(w, r, "/edit", 301)
	}
}

func folderLogoPath(id int) []string {
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
	return changePath(files)
}

func changePath(files [] string) []string{
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
func (c *Controller) DeleteImgHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		path := r.FormValue("path")
		p:=strings.Replace(path,"%2f","/",100)

		c.MediaStore.DeleteFile(".."+p)
		http.Redirect(w, r, "/", 301)
	}
}



