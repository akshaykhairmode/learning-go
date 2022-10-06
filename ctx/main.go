package main

import (
	"bufio"
	"context"
	"ctx/demo"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	ctx := context.Background() //Empty context
	log.Println("empty context", ctx)

	// WithDeadline(ctx)
	// WithTimeout(ctx)
	// WithCancel(ctx)
	// WithValue(ctx)

	//With Custom context
	realCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	log.Println(someRealWorldExample(realCtx))

	//With http server, caller cancelling request
	startHttpServer()
}

func startHttpServer() {
	s := http.Server{
		Addr: ":8000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := someRealWorldExample(r.Context()) //Pass the http request context
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				w.Write([]byte(data))
				return
			}

			w.WriteHeader(200)
			w.Write([]byte(data))
			log.Println("Request Complete")
		}),
	}

	log.Println("Starting server on port 8000")
	log.Println(s.ListenAndServe())
}

func someRealWorldExample(ctx context.Context) (string, error) {

	f, err := os.Open("data.txt")
	if err != nil {
		return "", err
	}

	s := bufio.NewScanner(f)

	out := strings.Builder{}

t:

	for {

		select {
		case <-ctx.Done():
			return out.String() + "\n", fmt.Errorf("ctx cancelled")
		default:
			if s.Scan() {
				log.Printf("Appending %s", s.Text())
				out.WriteString(s.Text())
				time.Sleep(500 * time.Millisecond)
			} else {
				break t
			}
		}
	}

	return out.String() + "\n", nil
}

func WithDeadline(ctx context.Context) {

	log.Println("WithDeadline")

	now := time.Now()
	log.Println("Start time is", now)

	after5seconds := now.Add(5 * time.Second)
	log.Println("Deadline is", after5seconds)

	ctx, cancel := context.WithDeadline(ctx, after5seconds)
	defer cancel() //Release the context

	for {
		select { //Random selection if both cases are ready.
		case <-ctx.Done():
			log.Println("Exiting function WithDeadline")
			return
		default:
			time.Sleep(500 * time.Millisecond)
			log.Printf("This line will keep printing every 500ms")
		}
	}

}

func WithTimeout(ctx context.Context) {

	log.Println("WithTimeout")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) //Will wait for 5 seconds
	defer cancel()                                         //Release the context

	for {
		select { //Random selection if both cases are ready.
		case <-ctx.Done():
			log.Println("Exiting function WithTimeout")
			return
		default:
			time.Sleep(500 * time.Millisecond)
			log.Printf("This line will keep printing every 500ms")
		}
	}
}

func WithCancel(ctx context.Context) {

	log.Println("WithCancel")

	ctx, cancel := context.WithCancel(ctx) //returns a cancel func which can be called separately to send signal to done channel
	defer cancel()                         //Release the context

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	for {
		select { //Random selection if both cases are ready.
		case <-ctx.Done():
			log.Println("Exiting function WithCancel")
			return
		default:
			time.Sleep(500 * time.Millisecond)
			log.Printf("This line will keep printing every 500ms till 2 second passes")
		}
	}
}

func WithValue(ctx context.Context) {

	ctx = context.WithValue(ctx, "ID", "SOME-STRING") //returns a new context which contains the key value pair

	printContext(ctx, "ctx")

	privateCtx := demo.Private()

	printContext(privateCtx, "privateCtx")

	publicCtx := demo.Public()

	printContext(publicCtx, "publicCtx")

	exportedPublic := demo.PublicWithExportedType()

	printContext(exportedPublic, "exportedPublic")

	//Modifications below

	modifiedPrivateCtx := context.WithValue(privateCtx, "ID", "MY-MODIFIED-STRING")

	printContext(modifiedPrivateCtx, "modifiedPrivateCtx")

	modifiedPublicCtx := context.WithValue(publicCtx, "ID", "MY-MODIFIED-STRING")

	printContext(modifiedPublicCtx, "modifiedPublicCtx")

	modifiedExportedCtx := context.WithValue(exportedPublic, demo.MyExportedKey("ID"), "MY-MODIFIED-STRING")

	printContext(modifiedExportedCtx, "modifiedExportedCtx")

}

func printContext(ctx context.Context, which string) {
	log.Printf("%s ::: ID is %v, NAME is %v, exported ID is %v", which, ctx.Value("ID"), ctx.Value("NAME"), ctx.Value(demo.MyExportedKey("ID")))
}
