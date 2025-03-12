package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	CACHE_CONFIG_PATH = "./../cache/config"
	CACHE_PATH        = "./../cache/"
)

type RespMetadata struct {
	LastMod  string
	ETag     string
	Filename string
}

type Cache struct {
	cache map[string]*RespMetadata
}

func InitCache() *Cache {
	config, err := os.Open(CACHE_CONFIG_PATH)
	if err != nil {
		dir := filepath.Dir(CACHE_CONFIG_PATH)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Println("Failed to create cache")
			return nil
		}

		file, err := os.Create(CACHE_CONFIG_PATH)
		if err != nil {
			fmt.Println("Failed to create cache config")
			return nil
		}
		defer file.Close()
		return &Cache{cache: make(map[string]*RespMetadata)}
	}
	defer config.Close()

	scanner := bufio.NewScanner(config)
	c_map := make(map[string]*RespMetadata)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		c_map[line[0]] = &RespMetadata{Filename: line[1], LastMod: line[2], ETag: line[3]}
	}
	return &Cache{cache: c_map}
}

func (c *Cache) ReadCachedResponse(path string, w http.ResponseWriter) {
	file, err := os.Open(CACHE_PATH + c.cache[path].Filename)
	if err != nil {
		fmt.Println("Failed to read cached response")
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

func (c *Cache) SaveResponse(path, lm, etag string, respBody io.Reader) {
	targetURL := path

	path = strings.TrimPrefix(path, "http://")
	path = strings.TrimPrefix(path, "https://")
	path = strings.ReplaceAll(path, ".", "_")
	path = strings.ReplaceAll(path, "-", "_")
	path = strings.ReplaceAll(path, "/", "_")
	path += "_cache.txt"

	file, err := os.Create(CACHE_PATH + path)
	if err != nil {
		log.Printf("Failed to create cache file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, respBody)
	if err != nil {
		log.Printf("Failed to write to cache file: %v", err)
	}
	respMetadata := RespMetadata{Filename: path, LastMod: lm, ETag: etag}
	c.cache[targetURL] = &respMetadata

	config, err := os.OpenFile(CACHE_CONFIG_PATH, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Failed to open cache config")
		return
	}
	defer config.Close()
	if _, err := config.Write([]byte(targetURL + " " + respMetadata.Filename + " " + respMetadata.LastMod + " " + respMetadata.ETag + "\n")); err != nil {
		fmt.Println("Failed to write to cache config")
		log.Fatal(err)
	}
}

func (c *Cache) Contains(url string) bool {
	_, ok := c.cache[url]
	return ok
}

func (c *Cache) GetRespMetadata(url string) (string, string) {
	return c.cache[url].LastMod, c.cache[url].ETag
}
