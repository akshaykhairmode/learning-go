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

	WithDeadline(ctx)
	WithTimeout(ctx)
	WithCancel(ctx)
	WithValue(ctx)

	//With Custom context
	realCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	log.Println(fileReadExample(realCtx))

	//With nested context
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	nestedFunction(ctxWithTimeout, 0)

	//With http server, caller cancelling request
	startHttpServer()

}

func nestedFunction(ctx context.Context, counter int) {

	for {
		select {
		case <-ctx.Done():
			log.Println("returning at counter", counter)
			return
		default:
			counter++
			log.Println("Sleeping for 1 second at counter", counter)
			time.Sleep(time.Second)
			nestedFunction(ctx, counter)
		}
	}

}

func startHttpServer() {
	s := http.Server{
		Addr: ":8000", //Start server on this port
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //On base path execute below code.
			data, err := fileReadExample(r.Context()) //Pass the http request context
			if err != nil {                           //Context cancels or any other erros will come here.
				log.Println(err)
				w.WriteHeader(500) //Give 500 status code if 500 error
				w.Write([]byte(data))
				return
			}

			w.WriteHeader(200) //200 status code if success
			w.Write([]byte(data))
			log.Println("Request Complete")
		}),
	}

	log.Println("Starting server on port 8000")
	log.Println(s.ListenAndServe()) //Start the http server
}

func fileReadExample(ctx context.Context) (string, error) {

	f, err := os.Open("data.txt")
	if err != nil {
		return "", err
	}

	s := bufio.NewScanner(f)
	out := strings.Builder{}

t:

	for {

		select {
		//If context is done then we wull return from this function with error.
		//We can choose if we want to return the data read till now or return empty string.
		case <-ctx.Done():
			return out.String() + "\n", fmt.Errorf("ctx cancelled")
		default:
			//Till my context is done default will be called to scan every line an append to our string builder.
			if s.Scan() {
				log.Printf("Appending %s", s.Text())
				out.WriteString(s.Text())
				time.Sleep(500 * time.Millisecond)
			} else {
				//When scan returns false, it means that the file has reached the end
				//so we break the loop with a label.
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
		select { //Default would be last priority here.
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
		select {
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

	go func() { //in 2 seconds call the cancel function.
		time.Sleep(2 * time.Second)
		cancel()
	}()

	for {
		select {
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

	//Add ID key to passed context. Replaces the old ID if existed.
	ctx = context.WithValue(ctx, "ID", "SOME-STRING")

	//Check comments in demo file.
	privateCtx := demo.Private()
	publicCtx := demo.Public()
	exportedPublicCtx := demo.PublicWithExportedType()

	//Trying modification on different context types.
	modifiedPrivateCtx := context.WithValue(privateCtx, "ID", "MY-MODIFIED-STRING")
	modifiedPublicCtx := context.WithValue(publicCtx, "ID", "MY-MODIFIED-STRING")
	modifiedExportedCtx := context.WithValue(exportedPublicCtx, demo.MyExportedKey("ID"), "MY-MODIFIED-STRING")

	//Maintain to print easily all of them.
	print := map[context.Context]string{
		ctx:                 "ctx",
		privateCtx:          "privateCtx",
		publicCtx:           "publicCtx",
		exportedPublicCtx:   "exportedPublicCtx",
		modifiedPrivateCtx:  "modifiedPrivateCtx",
		modifiedPublicCtx:   "modifiedPublicCtx",
		modifiedExportedCtx: "modifiedExportedCtx",
	}

	for ctx, what := range print {
		printContext(ctx, what)
	}

}

//Prints the context values
func printContext(ctx context.Context, which string) {
	log.Printf("%s ::: ID is %v, NAME is %v, exported ID is %v", which, ctx.Value("ID"), ctx.Value("NAME"), ctx.Value(demo.MyExportedKey("ID")))
}
