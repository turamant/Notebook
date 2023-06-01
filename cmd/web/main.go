package main

import (
	"flag"
	"log"
	"net/http"
)

func main(){
	addr := flag.String("addr", ":4000", "HTTP network address")
    flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/resume", downloadHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Запускаем сервер на %s порту", *addr)
	err := http.ListenAndServe(*addr,mux)
	log.Fatal(err)
}