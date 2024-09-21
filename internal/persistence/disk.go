package persistence

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

const fileName = ".cache"

type DiskPersistence struct {
	basePath  string
	cacheFile string
	Cache     *Cache
}

type Cache struct {
	Model   string `json:"model"`
	Profile string `json:"profile"`
}

func NewDiskPersistence(basePath string) *DiskPersistence {
	filePath := path.Join(basePath, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		p := &DiskPersistence{
			basePath:  basePath,
			cacheFile: ".cache",
			Cache: &Cache{
				Model:   "default",
				Profile: "default",
			},
		}

		p.Save()
		return p
	}

	p := &DiskPersistence{
		basePath:  basePath,
		cacheFile: fileName,
	}

	p.loadCache()

	return p
}

func (p *DiskPersistence) filePath() string {
	return path.Join(p.basePath, p.cacheFile)
}

func (p *DiskPersistence) loadCache() {
	content, err := p.readFrom(p.cacheFile)
	if err != nil {
		panic(err)
	}

	parsed := &Cache{}

	err = json.Unmarshal(content, parsed)

	if err != nil {
		panic(err)
	}

	p.Cache = parsed
}

func (p *DiskPersistence) Save() error {
	content, err := json.Marshal(p.Cache)

	if err != nil {
		return err
	}

	return p.writeTo(p.cacheFile, content)
}

func (p *DiskPersistence) readFrom(filename string) ([]byte, error) {

	file, err := os.Open(p.filePath())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return content, nil
}

func (p *DiskPersistence) writeTo(filename string, content []byte) error {
	file, err := os.Create(p.filePath())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)

	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) SetModel(modelName string) {
	c.Model = modelName
}

func (c *Cache) SetProfile(profileName string) {
	c.Profile = profileName
}
