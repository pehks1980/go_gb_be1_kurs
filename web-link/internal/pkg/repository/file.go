package repository

import (
	"encoding/json"
	"fmt"
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
	List() ([]string, error)
}

type FileRepo struct {
	sync.RWMutex
	fileName string
	fileData map[string]model.DataEl
}

func (fr *FileRepo) New(filename string) RepoIf {
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
		fileRepo.FileRepoUnpackToStruct()
	}

	return fileRepo
}

func NewFileRepo(fileName string) *FileRepo {
	return &FileRepo{fileName: fileName}
}

// DumpMapToFile - no lock, as its has been done in upper level
func (fr *FileRepo) DumpMapToFile() error {
	// to do dump map to file.json
	// make slice of active links and write it to file
	var fileDataSlice model.Data

	for _, value := range fr.fileData {
		// stripe all not Active when dumping
		if value.Active == 1 {
			fileDataSlice.Data = append(fileDataSlice.Data, value)
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
	err = json.Unmarshal(byteValue, &fileDataSlice)
	if err != nil {
		return err
	}
	// quickly populate our file map

	// we iterate through array and make map key shortlink:filedata struct
	for _, datael := range fileDataSlice.Data {
		fr.fileData[datael.Shorturl] = datael
	}

	return nil
}

func (fr *FileRepo) Get(key string) (model.DataEl, error) {
	fr.RWMutex.RLock()
	defer fr.RWMutex.RUnlock()
	// get data needed
	// retrieve dat string

	if datael, ok := fr.fileData[key]; ok {
		if datael.Active == 0 {
			// deleted already
			err := fmt.Errorf("link deleted already")
			return model.DataEl{}, err
		}

		return datael, nil
		//no such key

	}
	err := fmt.Errorf("No such link")
	return model.DataEl{}, err
}

func (fr *FileRepo) Put(key string, value model.DataEl) error {
	fr.RWMutex.Lock()
	defer fr.RWMutex.Unlock()
/*	if _, ok := fr.fileData[key]; !ok {
		// key already exists
		err := fmt.Errorf("link %s dont exist", key)
		return err
	}*/
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
	if datael, ok := fr.fileData[key]; ok {
		datael.Active = 0
		fr.fileData[key] = datael
		// dump data to file straight away
		err := fr.DumpMapToFile()
		if err != nil {
			return err
		}
		return nil
	}
	err := fmt.Errorf("delete error key %s don't exist", key)
	return err
}

func (fr *FileRepo) List() ([]string, error) {
	fr.RWMutex.RLock()
	defer fr.RWMutex.RUnlock()
	// get data needed
	// retrieve list of keys as []string
	var keys []string
	for k, val := range fr.fileData {
		if val.Active == 1 {
			keys = append(keys, k)
		}
	}
	return keys, nil
}
