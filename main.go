package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

type Event struct {
	Records []Records `json:"Records"`
}

type Records struct {
	EventVersion string `json:"eventVersion"`
	EventSource  string `json:"eventSource"`
	S3           struct {
		Bucket struct {
			Name string `json:"name"`
		} `json:"bucket"`
		Object struct {
			Key string `json:"key"`
		} `json:"object"`
	} `json:"s3"`
}

func callback(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, err)
	} else {
		var event Event
		err = json.Unmarshal(body, &event)
		if err != nil {
			sendError(w, err)
		} else {
			log.Infof("successfully parsed the body")
			bucketName := event.Records[0].S3.Bucket.Name
			log.Infof("bucket %s", bucketName)
			key := event.Records[0].S3.Object.Key
			log.Infof("key %s", key)
			if err := createStorageConfig(bucketName, key); err != nil {
				sendError(w, err)
			} else {
				fmt.Fprintf(w, "ack")
			}
		}
	}
}

func createStorageConfig(bucketName, key string) error {
	access_key_id := os.Getenv("AWS_ACCESS_KEY_ID")
	secret_access_key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	secure, err := strconv.ParseBool(os.Getenv("SECURE_ENDPOINT"))
	if err != nil {
		return err
	}
	s3Client, err := minio.New(endpoint, access_key_id, secret_access_key, secure)
	if err != nil {
		return err
	}
	storageconfig_key := strings.Replace(key, "kopia.repository", ".storageconfig", 1)
	if _, err := s3Client.FPutObject(bucketName, storageconfig_key, "/etc/ceph-callback/.storageconfig", minio.PutObjectOptions{
		ContentType: "application/json",
	}); err != nil {
		return err
	}
	log.Infof("Successfully uploaded %s", storageconfig_key)
	return nil
}

func sendError(w http.ResponseWriter, err error) {
	log.WithFields(log.Fields{
		"error": err,
	}).Errorf("An error happened")
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = fmt.Sprintf("%v", err)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func main() {
	http.HandleFunc("/", callback)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
