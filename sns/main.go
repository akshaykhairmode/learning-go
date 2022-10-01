package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var sess = session.Must(session.NewSession(
	&aws.Config{
		Region: aws.String("ap-south-1"),
	},
))

var snsService = sns.New(sess)

var publicHttpEndpoint string

type SNSEvent struct {
	RequestType  *string `json:"Type"`
	SubscribeURL *string `json:"SubscribeURL"`
	Message      *string `json:"Message"`
}

func main() {

	flag.StringVar(&publicHttpEndpoint, "endpoint", "", "-endpoint")
	flag.Parse()

	if publicHttpEndpoint == "" {
		log.Fatal("Public Endpoint is required")
		flag.PrintDefaults()
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startHttpServer(wg, ":8000")
	subscribe(publicHttpEndpoint)
	go publish()

	wg.Wait()
}

func startHttpServer(wg *sync.WaitGroup, port string) {
	defer wg.Done()
	server := http.Server{
		Addr:    port,
		Handler: handler(),
	}
	log.Println("starting server on port 8000")
	err := server.ListenAndServe()
	log.Fatal(err)
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}

		req := SNSEvent{}

		if err := json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		//confirmation request
		if *req.RequestType == "SubscriptionConfirmation" {
			confirm(req)
			return
		}

		//messages
		log.Println("Got Message", *req.Message)
		w.WriteHeader(200)
	}
}

//subscribe subscribes to the topic. This can be done from the AWS console also.
func subscribe(httpURL string) {
	inp := &sns.SubscribeInput{
		Endpoint: aws.String(httpURL),
		Protocol: aws.String("http"),
		TopicArn: aws.String("arn:aws:sns:ap-south-1:XXXXXXXXXX:test-topic"),
	}

	output, err := snsService.Subscribe(inp)
	if err != nil {
		log.Println(err)
	}

	log.Println("Output", output.String())
}

func confirm(request SNSEvent) {

	vals, err := url.ParseQuery(*request.SubscribeURL)
	if err != nil {
		log.Println(err)
	}

	token := vals.Get("Token")
	topicARN := vals.Get("TopicArn")

	confirm := &sns.ConfirmSubscriptionInput{
		Token:    aws.String(token),
		TopicArn: aws.String(topicARN),
	}
	output, err := snsService.ConfirmSubscription(confirm)
	if err != nil {
		log.Println(err)
	}

	log.Println("confirm output", output.String())
}

func publish() {
	for {
		input := &sns.PublishInput{
			Message:  aws.String("Current time is " + time.Now().String()),
			TopicArn: aws.String("arn:aws:sns:ap-south-1:XXXXXXXXXX:test-topic"),
		}
		//We are ignoring the output here. Its the messageid.
		_, err := snsService.Publish(input)
		if err != nil {
			log.Println("publish error", err)
		}

		time.Sleep(2 * time.Second)
	}
}
