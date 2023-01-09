## Details
This is just a small project to play around with **Go** and its connection with **Redis**.

It also has some additional conceptual implementations I wanted to test like:
- In memory cache
- Redis pre-loading
- Small *antiseed* to generate mocked data.

### Execution
In theory you just have to run the docker compose
```
docker-compose up --build
```
And then insert your own `album` data or run the `seed` scenario manually inside its folder.
```
go run antiseed.go
```

Once you have all set up, you just have to send the `POST`/`GET` requests.

**Get Albums**: Returns a list of all the loaded albums.
```
GET http://localhost:8080/albums
```
**Get Album**: Returns a single instance of an album.
```
http://localhost:8080/albums/1
```
**Create Album**: Creates and saves an album instance.
```
POST http://localhost:8080/albums

{
    "id": "9999999",
    "title": "The Modern Sound of Betty Carter",
    "artist": "Betty Carter",
    "price": 49.99
}
```
