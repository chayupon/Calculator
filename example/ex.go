package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/", ServeHTTP)

?
	server := http.Server{Addr: ":8080", Handler: mux}
	log.Fatal(server.ListenAndServe())

}

func ServeHTTP(rw http.ResponseWriter, r *http.Request) {


	rw.Write([]byte("Welcome to the \"Just Enough Go\" blog series!!"))
}

func createPost(w http.ResponseWriter, r *http.Request) {
  
	
  // input struct
  var post Post

  // convert r.Body to struct input
  _ = json.NewDecoder(r.Body).Decode(post)


  post.ID = strconv.Itoa(rand.Intn(1000000))
  posts = append(posts, post)

  // set http header
  w.Header().Set("Content-Type", "application/json")

  // set status 200
  //{"result":1000,"time_now": timeformat.RFC33xx}
  // combile header, status code and body
  json.NewEncoder(w).Encode(&post)
}