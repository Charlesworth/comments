package main

import (
	json "encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	port := flag.Int("port", 8000, "the port to be used by the service")
	flag.Parse()
	return ":" + string(*port)
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

	w.WriteHeader(http.StatusOK)
	printComments(w, encodedCommentSlice)
	return
}

func printComments(w http.ResponseWriter, byteSlice [][]byte) {
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rlog.Error("Unable to read request body, error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newComment := inputComment{}
	err = json.Unmarshal(body, &newComment)
	if err != nil {
		rlog.Error("Unable to unmarshal body json with error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	commentTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	storedComment := comment{
		Poster:   newComment.Poster,
		Page:     page,
		Msg:      newComment.Msg,
		TimeUnix: commentTime,
	}
	encodedComment, err := json.Marshal(storedComment)

	if err != nil {
		rlog.Error("unable to PUT comment into byte store:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bs.Put(storedComment.Page, commentTime, encodedComment)
	rlog.Info("comment added to byteStore: ", string(encodedComment))

	w.WriteHeader(http.StatusOK)
	return
}
