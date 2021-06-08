package repository

import (
	"encoding/json"
	"github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type RepoIf interface {
	New(filename string) RepoIf
	Get(key string) (model.DataEl, error)
	Put(key string, value model.DataEl) error
	Del(key string) error
}

type FileRepo struct {
	sync.RWMutex
	fileName string
	fileData map[string]model.DataEl
}

func (fr *FileRepo) New (filename string) RepoIf{
	// todo init

	fileRepo := &FileRepo{
		fileName: filename,
		fileData: make(map[string]model.DataEl),
	}
	//check if file exists
	// if yes load from disk and populate repo structs
	// so Image of file is held in map and it gets flushed every time change occurs
	// todo need to flush file when shutdown server
	if _, err := os.Stat(filename); err == nil {
		// path/to/whatever exists
		fr.FileRepoUnpackToStruct()
	}

	return fileRepo
}

func NewFileRepo(fileName string) *FileRepo {
	return &FileRepo{fileName: fileName}
}

// DumpMapToFile
func (fr *FileRepo) DumpMapToFile() error{
	fr.RWMutex.Lock()
	defer fr.RWMutex.Unlock()
	// to do dump map to file.json
	// make slice of active links and write it to file
	var fileDataSlice model.Data

	for _, value := range fr.fileData {
		// stripe all not Active when dumping
		if value.Active == 1 {
			fileDataSlice.Data = append (fileDataSlice.Data, value)
		}
	}

	filedata, _ := json.MarshalIndent(fileDataSlice, "", " ")

	err := ioutil.WriteFile(fr.fileName, filedata, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// FileRepoUnpackToStruct
func (fr *FileRepo) FileRepoUnpackToStruct() error {
	fr.RWMutex.Lock()
	defer fr.RWMutex.Unlock()
	// по ссылке извлекаем строку файлового хранилища
	// читаем все в мапу и делаем поиск
	jsonFile, err := os.Open(fr.fileName)
	if err != nil {
		return err
	}

	// Не забываем закрыть файл при выходе из функции
	defer func() {
		var ferr = jsonFile.Close()
		if ferr != nil {
			log.Printf("can't close file: %v", ferr)
		}
	}()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// we initialize our data array
	var fileDataSlice model.Data
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'fileDataSlice' which we defined above
	json.Unmarshal(byteValue, fileDataSlice)
	// quickly populate our file map

	// we iterate through array and make map key shortlink:filedata struct
	for _, datael := range fileDataSlice.Data{
		fr.fileData[datael.Shorturl] = datael
	}

	return nil
}

func (fr *FileRepo) Get(key string) (model.DataEl, error) {
	fr.RWMutex.RLock()
	defer fr.RWMutex.RUnlock()
	// get data needed
	// retrieve dat string
	datael := fr.fileData[key]
	// if its deleted return {}
	// todo no err!!!!
	if datael.Active == 0 {
		return model.DataEl{}, nil
	}
	return datael, nil
}

func (fr *FileRepo) Put(key string, value model.DataEl) error {
	fr.RWMutex.Lock()
	defer fr.RWMutex.Unlock()
	fr.fileData[key] = value
	// changes needs to be flushed to file
	err := fr.DumpMapToFile()
	if err != nil {
		return err
	}
	return nil
}

// Del - mark Active = 0 to 'delete'
func (fr *FileRepo) Del(key string) error {
	// TODO: impl
	fr.RWMutex.Lock()
	defer fr.RWMutex.Unlock()
	datael := fr.fileData[key]
	datael.Active = 0
	fr.fileData[key] = datael
	return nil
}
