package main

import (
	"encoding/json"
	"fmt"
	"github.com/leandroveronezi/go-recognizer"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"snippetbox/pkg/models"
	"strings"
)

const fotosDir = "fotos"
const dataDir = "models"
const testDir = "test"

func (app *application) recognition(w http.ResponseWriter, r *http.Request) {
	//file uploading part starts here
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("./test/test.jpeg", fileBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("file uploaded")
	//file uploading part ends here

	faces, err := app.rec.Classify(filepath.Join(testDir, "test.jpeg"))

	if err != nil {
		fmt.Println(err)
		w.Write([]byte("photo might contain multiple faces or face may not be in database !!"))
		return
	}
	var dob, name, place1, place2, branch string
	for _, face := range faces {
		id := face.Id
		id = id[:8]
		fmt.Println(id)
		if student, ok := app.studentMap[strings.ToUpper(id)]; ok {
			dob = student["dob"].(string)
			name = student["stu_name"].(string)
			branch = student["dname"].(string)
			place1 = id[:4]
			place2 = student["line2"].(string)
		} else {
			fmt.Println(ok)
		}
		mod := &models.Stranger{
			Name:   name,
			Place1: place1,
			Place2: place2,
			Branch: branch,
			Dob:    dob,
		}
		app.render(w, r, "showface.page.tmpl", &templateData{
			Stranger: mod,
		})
		return
	}
}
func addFile(rec *recognizer.Recognizer, Path, Id string) {

	err := rec.AddImageToDataset(Path, Id)

	if err != nil {
		fmt.Println(err)
		return
	}
}
func (app *application) dataTrain() {
	err := app.rec.Init(dataDir)
	if err != nil {
		fmt.Println(err)
		//app.serverError(w, err)
		return
	}
	app.rec.Tolerance = 0.35
	app.rec.UseGray = false
	app.rec.UseCNN = false
	//defer app.rec.Close()
	items, _ := ioutil.ReadDir("./fotos")
	for _, foto := range items {
		fmt.Println(foto.Name()[:8])
		addFile(&app.rec, filepath.Join(fotosDir, string(foto.Name())), string(foto.Name()))
		break
	}
	app.rec.SetSamples()
}
func (app *application) jsonToMap() {
	data, err := ioutil.ReadFile("./faceISM/stu_data/info.json")
	if err != nil {
		panic(err)
	}
	var studentData []map[string]interface{}
	if err := json.Unmarshal(data, &studentData); err != nil {
		panic(err)
	}
	//studentMap := make(map[string]map[string]interface{})
	for _, student := range studentData {
		id := student["stu_details"].(map[string]interface{})["id"].(string)
		app.studentMap[id] = student["stu_details"].(map[string]interface{})
	}
}
