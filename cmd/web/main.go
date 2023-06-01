package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main(){
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
    
    f, e := os.OpenFile("/tmp/info2.log", os.O_RDWR|os.O_CREATE, 0666)
	if e != nil {
		log.Fatal(e)
	}
	defer f.Close()

	infoLogFile := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	errorLogFile := log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    infoLogTerminal := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/resume", downloadHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLogFile,
		Handler: mux,
	}

	infoLogFile.Printf("Starting server on %s", *addr)
	infoLogTerminal.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLogFile.Fatal(err)
	
}