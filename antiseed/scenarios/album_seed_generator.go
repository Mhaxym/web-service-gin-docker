package scenarios

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"web-service-gin-docker/goCache"
	"web-service-gin-docker/redis"
)

const N_ALBUM = 100000

var NAMES = [10]string{"John Coltrane", "Miles Davis", "Bill Evans", "Thelonious Monk", "Charlie Parker", "Ornette Coleman", "Duke Ellington", "Sonny Rollins", "Art Blakey", "Herbie Hancock"}
var ALBUM_NAMES = [10]string{"Blue Train", "Kind of Blue", "Portrait in Jazz", "Monk's Dream", "KoKo", "Something Else", "Money Jungle", "Saxophone Colossus", "A Night at Birdland", "Empyrean Isles"}

func CreateAlbums() {

	var data = map[string]interface{}{}
	var dataPage = goCache.DataPage{Name: "AlbumManager"}

	for i := 0; i < N_ALBUM; i++ {
		item := goCache.Album{
			ID:     strconv.Itoa(i),
			Title:  fmt.Sprintf("%s (%d)", ALBUM_NAMES[rand.Intn(len(ALBUM_NAMES))], i),
			Artist: NAMES[rand.Intn(len(NAMES))],
			Price:  math.Round(rand.Float64()*10000) / 100,
		}

		key := fmt.Sprintf("{%s}/%s", dataPage.Name, item.ID)

		data[key], _ = json.Marshal(&item)
		dataPage.AddPageKey(key)
	}

	var service redis.Service = *redis.GetService()
	// First we save the data
	service.MSet(data)
	// Then we save the page
	service.Set(goCache.GetPageKey("AlbumManager"), dataPage.ToJSON())
}
