package bdlutils

import (
	"fmt"
	"errors"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"strings"
	"time"
	"io/ioutil"
	"os"
	"path/filepath"
)

type IdImages struct {
	Id int `json:"idMediaServer"`
}

type Document struct {
	Id        string
	Id_images []IdImages
	Title     string
}

const (
	Url           = "https://www.bdl.servizirl.it/bdl/public/rest/json/item/{id}/bookreader/pages"
	CantaloupeUrl = "https://www.bdl.servizirl.it/cantaloupe/iiif/2/{id}/full/full/0/default.jpg"
	InfoUrl       = "https://www.bdl.servizirl.it/bdl/public/rest/json/item/{id}"
	nMaxGorutine  = 10
	hImage        = 210
	wImage        = 297
	posX          = 0
	posY          = 0
)

func (d *Document) GetImagesId(u string) error {
	var i []IdImages

	resp := request(u, d.Id)

	err := json.Unmarshal(resp, &i)
	if err != nil {
		return errors.New("Images ID not Found")
	}

	d.Id_images = i

	return nil
}

func (d *Document) GetTitle(u string) error {
	resp := request(u, d.Id)

	m2 := make(map[string]interface{})
	err := json.Unmarshal(resp, &m2)
	if err != nil {
		return errors.New("<error> Title not found! [id wrong]")
	}

	d.Title = fmt.Sprintf("%v", m2["title"])
	d.Title = strings.Replace(d.Title, ",", "", -1)
	d.Title = strings.Replace(d.Title, ".", "", -1)
	
	return nil
}

func (d *Document) GetImages(u string, dest string) error {
	var wg sync.WaitGroup
	nmg := nMaxGorutine

	err := CreateDir(dest + "/images/")
	if err != nil {
		return errors.New("Failed to create directory \"images\"")
	}

	for _, imageId := range d.Id_images {
		idString := strconv.Itoa(imageId.Id)

		if nmg != 0 {
			nmg = nmg - 1
			wg.Add(1)
			//log.Println("Acquisizione..." + idString)
			go getImage(idString, u, dest,  &wg)
		} else {
			wg.Wait()
			nmg = nMaxGorutine - 1
			wg.Add(1)
			//log.Println("Acquisizione..." + idString)
			go getImage(idString, u, dest, &wg)
		}
	}
	wg.Wait()
	return nil
}

func (d *Document) CreatePdf(dest string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 11)

	for _, imageId := range d.Id_images {
		nameImage := dest + "/images/" + strconv.Itoa(imageId.Id) + ".jpg"
		pdf.AddPage()
		pdf.Image(nameImage, posX, posY, hImage, wImage, false, "", 0, "")
		_ = os.Remove(nameImage)
	}
	pdf.OutputFileAndClose(dest + "/" + d.Title + ".pdf")

	if pdf.Error() != nil {
		return pdf.Error()
	}

	return nil
}

func request(u string, id string) []byte {
	var resp *http.Response
	var err error

	link := strings.Replace(u, "{id}", id, -1)

	for retries := 3; retries > 0; retries-- {
		resp, err = http.Get(link)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalln(err)
	}

	content, _ := ioutil.ReadAll(resp.Body)

	return content
}

func getImage(imageId, u string, dest string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp := request(u, imageId)
	nameImage := dest + "/images/" + imageId + ".jpg"
	err := ioutil.WriteFile(nameImage, resp, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	//log.Println(imageId + "...finito")
}

func CreateDir(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return err
	}
	return nil
}
