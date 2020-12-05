package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/themakers/gphoto2"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	gpc := gphoto2.NewContext()
	defer gpc.Cancel()

	cam := gphoto2.NewCamera(gpc)

	if err := cam.Init(); err != nil {
		panic(err)
	} else {
		defer cam.Exit()
	}

	files, err := cam.ListFilesRec("")
	if err != nil {
		panic(err)
	}

	var download string
	for _, file := range files {
		log.Print(file)
		if !strings.HasSuffix(strings.ToLower(file), ".jpg") {
			download = file
		}
	}

	file := cam.File(download)

	if info, err := file.GetInfo(); err != nil {
		panic(err)
	} else {
		log.Println(info.File)
	}

	if fr, err := file.ReadSeeker(); err != nil {
		panic(err)
	} else {
		t := time.Now()
		buf := bytes.NewBuffer([]byte{})
		if n, err := io.Copy(buf, fr); err != nil {
			panic(err)
		} else {
			d := time.Now().Sub(t)
			log.Printf("downloaded %.2fMB %s %.2fMB/s", float64(n)/(1024*1024), d, (float64(n)/(1024*1024))/(float64(d)/float64(time.Second)))

			if err := ioutil.WriteFile(file.Name(), buf.Bytes(), 0777); err != nil {
				panic(err)
			}
		}
	}

	if w, writeConfig, err := cam.RootWidget(); err != nil {
		panic(err)
	} else if err := w.Walk(func(depth int, pw, w *gphoto2.CameraWidget) error {
		name, err := w.Name()
		if err != nil {
			panic(err)
		}

		label, err := w.Label()
		if err != nil {
			panic(err)
		}

		ro, err := w.Readonly()
		if err != nil {
			panic(err)
		}

		wt, wvt, err := w.Type()
		if err != nil {
			panic(err)
		}

		value, err := w.Value()
		if err != nil {
			panic(err)
		}

		log.Println(
			strings.Repeat("•", (depth+1)*2), depth,
			"name[", name, "]",
			"label[", label, "]",
			"ro[", ro, "]",
			"type[", wt, "/", wvt, "]",
			"value[", value, "]",
		)

		return nil
	}, true); err != nil {
		panic(err)
	} else if w, err := w.ChildByName("capturetarget"); err != nil {
		panic(err)
	} else {
		log.Println(w.Label())
		log.Println(w.Type())
		log.Println(w.Value())
		log.Println(w.Choices())

		if err := w.SetValueString("Internal RAM"); err != nil {
			panic(err)
		} else if err := writeConfig(); err != nil {
			panic(err)
		} else if err := cam.TriggerCapture(); err != nil {
			panic(err)
		}
	}
}
