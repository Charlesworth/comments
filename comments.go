package main

import (
	json "encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/charlesworth/byteStore"
	"github.com/julienschmidt/httprouter"
)

type comment struct {
	Poster   string
	Page     string
	Msg      string
	TimeUnix string
}

var bs byteStore.ByteStore

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	startByteStore()
}

func main() {
	http.Handle("/", newRouter())
	port := getPort()
	log.Println("Comment service started on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func newRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/:page", getComments)
	router.POST("/:page", postComment)
	router.DELETE("/:page/:time", deleteComment)
	router.GET("/", getCmds)
	return router
}

func getPort() string {
	portPtr := flag.Int("port", 8000, "the port to be used by the service")
	flag.Parse()
	return fmt.Sprintf(":%v", *portPtr)
}

func startByteStore() {
	var err error
	bs, err = byteStore.New("comments.byteStore")
	if err != nil {
		log.Fatalln(err)
	}
}

func getCmds(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.WithFields(log.Fields{
		"IP": r.RemoteAddr,
	}).Info("GET /")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Comment Service: \nGET /:page \nPOST /:page \nDELETE /:page/:time")
	return
}

func deleteComment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	page := params.ByName("page")
	time := params.ByName("time")
	log.WithFields(log.Fields{
		"IP":   r.RemoteAddr,
		"page": page,
	}).Info("DELETE /:page/:time")

	bs.Delete(page, time)
	return
}

func getComments(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	page := params.ByName("page")
	log.WithFields(log.Fields{
		"IP":   r.RemoteAddr,
		"page": page,
	}).Info("GET /:page")

	encodedCommentSlice := bs.GetBucketValues(page)
	if encodedCommentSlice == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// enable CORS header and set as origin URL
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	if callback, exists := queryStringCallback(r); exists {
		printJSONPComments(w, callback, encodedCommentSlice)
		return
	}
	printJSONComments(w, encodedCommentSlice)
	return
}

func queryStringCallback(r *http.Request) (callbackName string, ok bool) {
	qsValues := r.URL.Query()
	callback := qsValues.Get("callback")
	if callback == "" {
		return "", false
	}
	return callback, true
}

func printJSONPComments(w http.ResponseWriter, callbackName string, byteSlice [][]byte) {
	fmt.Fprint(w, callbackName+"({\"Comments\":")
	printJSONComments(w, byteSlice)
	fmt.Fprint(w, "});")
}

func printJSONComments(w http.ResponseWriter, byteSlice [][]byte) {
	//if single element, print and return
	if len(byteSlice) == 1 {
		fmt.Fprint(w, string(byteSlice[0]))
		return
	}

	//else print as array
	fmt.Fprint(w, "[")
	len := len(byteSlice)
	for i, val := range byteSlice {
		fmt.Fprint(w, string(val))
		if i == len-1 {
			fmt.Fprint(w, "]")
			return
		}
		fmt.Fprint(w, ",")
	}
}

type inputComment struct {
	Msg    string
	Poster string
}

func postComment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	page := params.ByName("page")
	rlog := log.WithFields(log.Fields{
		"IP":   r.RemoteAddr,
		"page": page,
	})
	rlog.Info("POST /:page")

	r.ParseForm()

	poster := r.FormValue("poster")
	msg := r.FormValue("msg")

	commentTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	storedComment := comment{
		Poster:   poster,
		Page:     page,
		Msg:      msg,
		TimeUnix: commentTime,
	}
	encodedComment, err := json.Marshal(storedComment)

	if err != nil {
		rlog.Error("unable to POST comment into byte store:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs.Put(storedComment.Page, commentTime, encodedComment)
	rlog.Info("comment added to byteStore: ", string(encodedComment))

	w.WriteHeader(http.StatusOK)
	return
}
