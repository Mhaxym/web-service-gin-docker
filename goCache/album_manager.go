package goCache

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"web-service-gin-docker/redis"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type AlbumManager struct {
	Manager
	albums       []Album
	albums_by_id map[string]Album
}

func (am *AlbumManager) Load() {
	err := am.loadFromRedis()
	if err != nil {
		am.loadFromLocalSeed()
	}
	am.indexData()
}

func (am *AlbumManager) loadFromLocalSeed() {
	am.albums = []Album{}
	for i := 0; i < 100000; i++ {
		album := Album{ID: strconv.Itoa(i), Title: "Blue Submarine " + strconv.Itoa(i), Artist: "John Doe", Price: 10.0}
		am.albums = append(am.albums, album)
	}
}

func (am *AlbumManager) loadFromRedis() error {
	redisService := redis.GetService()
	value, err := redisService.Get(GetPageKey(am.Manager.CacheName))
	if err != nil {
		log.Fatal(err)
		return err
	}
	// We check if the page exists and load the data
	dataPage := DataPage{Name: am.Manager.CacheName}
	dataPage.FromJSON(value)
	// We load the related keys data
	data, err := redisService.MGet(dataPage.PageKeys)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// We deserialize the data
	for _, value := range data {
		album := &Album{}
		json.Unmarshal([]byte(value.(string)), album)
		am.albums = append(am.albums, *album)
	}

	return nil
}

func (am *AlbumManager) indexData() {
	if len(am.albums) > 0 {
		am.albums_by_id = make(map[string]Album)
		for _, album := range am.albums {
			am.albums_by_id[album.ID] = album
		}
	}
}

func (am *AlbumManager) GetManagerName() string {
	return am.Manager.CacheName
}

func (am *AlbumManager) GetAlbums() []Album {
	return am.albums
}

func (am *AlbumManager) AddAlbum(album *Album) {
	am.albums = append(am.albums, *album)
	am.albums_by_id[(*album).ID] = *album
}

func (am *AlbumManager) GetAlbum(id string) (Album, error) {
	if album, ok := am.albums_by_id[id]; ok {
		return album, nil
	} else {
		return Album{}, errors.New("Album not found")
	}
}

func GetAlbumManager() *AlbumManager {
	goCacheManager := GetGoCacheManager()
	albumManager, ok := (*goCacheManager).Get("AlbumManager").(*AlbumManager)
	if !ok {
		albumManager = new(AlbumManager)
		albumManager.Manager.CacheName = "AlbumManager"
		(*goCacheManager).Set(albumManager)
	}
	return albumManager
}
