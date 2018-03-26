package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sourcegraph/go-webkit2/webkit2"
	"github.com/sqs/gojs"
)

func main() {
	fmt.Println(os.Args)

	StartGTK()

	view := NewView()

	view.Open(os.Args[1])
	waitErr := view.Wait()
	if waitErr != nil {
		panic(waitErr)
	}

	// Wait for render
	time.Sleep(5 * time.Second)

	wg := sync.WaitGroup{}
	wg.Add(1)

	fmt.Println("Getting Snapshot")
	view.GetSnapshot(func(result *image.RGBA, err error) {
		defer wg.Done()
		if err != nil || result == nil {
			panic(err)
		}

		fmt.Println(result.Rect)

		f, err := os.OpenFile(os.Args[2], os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer func() {
			syncErr := f.Sync()
			if syncErr != nil {
				panic(syncErr)
			}
			closeErr := f.Close()
			if closeErr != nil {
				panic(closeErr)
			}
			fmt.Println("Flushed")
		}()

		err = png.Encode(f, result)
		if err != nil {
			panic(err)
		}
	})

	fmt.Println("Waiting")
	wg.Wait()

	fmt.Println("Finished")

	//ready, waitErr := view.EvaluateJavaScript("window.$renderStaticReady")
	//if waitErr != nil {
	//	panic(waitErr)
	//}
	//if ready, _ := ready.(bool); !ready {
	//	fmt.Println("not ready\n")
	//}
	//result, waitErr := view.EvaluateJavaScript("document.documentElement.outerHTML")
	//if waitErr != nil {
	//	panic(waitErr)
	//}
	//fmt.Println(result.(string))
}

var startGTKOnce sync.Once
var ErrLoadFailed = errors.New("load failed")

// StartGTK ensures that the GTK+ main loop has started. If it has already been
// started by StartGTK, it will not start it again. If another goroutine is
// already running the GTK+ main loop, StartGTK's behavior is undefined.
func StartGTK() {
	startGTKOnce.Do(func() {
		gtk.Init(nil)
		go func() {
			runtime.LockOSThread()
			gtk.Main()
		}()
	})
}

// View represents a WebKit view that can load resources at a given URL and
// query information about them.
type View struct {
	*webkit2.WebView

	load        chan struct{}
	lastLoadErr error

	destroyed bool
}

// NewView creates a new View in the context.
func NewView() *View {
	view := make(chan *View, 1)
	glib.IdleAdd(func() bool {
		webView := webkit2.NewWebView()
		settings := webView.Settings()
		settings.SetEnableWriteConsoleMessagesToStdout(true)
		settings.SetUserAgentWithApplicationDetails("WebLoop", "v1")
		v := &View{WebView: webView}
		loadChangedHandler, _ := webView.Connect("load-changed", func(_ *glib.Object, loadEvent webkit2.LoadEvent) {
			switch loadEvent {
			case webkit2.LoadFinished:
				// If we're here, then the load must not have failed, because
				// otherwise we would've disconnected this handler in the
				// load-failed signal handler.
				v.load <- struct{}{}
			}
		})
		webView.Connect("load-failed", func() {
			v.lastLoadErr = ErrLoadFailed
			webView.HandlerDisconnect(loadChangedHandler)
		})
		view <- v
		return false
	})
	return <-view
}

// Open starts loading the resource at the specified URL.
func (v *View) Open(url string) {
	v.load = make(chan struct{}, 1)
	v.lastLoadErr = nil
	glib.IdleAdd(func() bool {
		if !v.destroyed {
			v.WebView.LoadURI(url)
		}
		return false
	})
}

// Wait waits for the current page to finish loading.
func (v *View) Wait() error {
	<-v.load
	return v.lastLoadErr
}

// EvaluateJavaScript runs the JavaScript in script in the view's context and
// returns the script's result as a Go value.
func (v *View) EvaluateJavaScript(script string) (result interface{}, err error) {
	resultChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	glib.IdleAdd(func() bool {
		v.WebView.RunJavaScript(script, func(result *gojs.Value, err error) {
			glib.IdleAdd(func() bool {
				if err == nil {
					goval, err := result.GoValue()
					if err != nil {
						errChan <- err
						return false
					}
					resultChan <- goval
				} else {
					errChan <- err
				}
				return false
			})
		})
		return false
	})

	select {
	case result = <-resultChan:
		return result, nil
	case err = <-errChan:
		return nil, err
	}
}
