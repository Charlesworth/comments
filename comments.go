package main

import (
	json "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	byteStore "github.com/charlesworth/byteStore"
	"github.com/julienschmidt/httprouter"
)

type comment struct {
	Poster   string
	Page     string
	Msg      string
	TimeUnix string
}

//storage should be a bucket per pages
//in each page bucket have each comment as a time / comment storage
// get can walk through the bucket
// put can append the bucket

func main() {
	//set logging
	log.SetFormatter(&log.JSONFormatter{})

	//init bolt store
	byteStore.Get("gdfg", "adf")

	//set port
	port := ":3000"
	// if len(os.Args) == 1 {
	// 	port = ":" + os.Args[1]
	// }

	//start the server and listen for requests
	http.Handle("/", newRouter())
	log.Println("Comment service started: listening on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func newRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/:page", getComments)
	router.PUT("/:page", putComment)
	// router.DELETE("/:page/:time", putComment)
	router.GET("/", getCmds)
	return router
}

func getCmds(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.WithFields(log.Fields{
		"IP": r.RemoteAddr,
	}).Info("GET / request")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Comment Service: \nGET /:page \nPUT /:page \nDELETE /:page/:time")
	return
}

// //should just walk through the bucket and return all of the comments in itself
// func deleteComment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	pageName := params.ByName("page")
// 	time := params.ByName("time")
// 	log.WithFields(log.Fields{
// 		"IP":   r.RemoteAddr,
// 		"page": pageName,
// 	}).Info("DELETE comment request")
//
// 	byteStore.Delete()
// 	return
// }

//should just walk through the bucket and return all of the comments in itself
func getComments(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pageName := params.ByName("page")
	log.WithFields(log.Fields{
		"IP":   r.RemoteAddr,
		"page": pageName,
	}).Info("GET comments request")

	encodedCommentSlice := byteStore.GetBucket(pageName)
	if encodedCommentSlice == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	arrayPrint(w, encodedCommentSlice)
	return
}

func arrayPrint(w http.ResponseWriter, byteSlice [][]byte) {
	if len(byteSlice) == 1 {
		fmt.Fprint(w, string(byteSlice[0]))
		return
	}

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
	return
}

type inputComment struct {
	Msg    string
	Poster string
}

//should require a msg and name, will set time itself
func putComment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	pageName := params.ByName("page")
	rlog := log.WithFields(log.Fields{
		"IP":   r.RemoteAddr,
		"page": pageName,
	})
	rlog.Info("PUT comment request")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rlog.Error("Unable to read request body")
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
		Page:     pageName,
		Msg:      newComment.Msg,
		TimeUnix: commentTime,
	}
	encodedComment, err := json.Marshal(storedComment)

	if err != nil {
		rlog.Error("unable to PUT comment into byte store:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	byteStore.Put(storedComment.Page, commentTime, encodedComment)
	rlog.Info("comment added to byteStore: ", string(encodedComment))

	w.WriteHeader(http.StatusOK)
	return
}
