package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/veqryn/go-email/email"
)

func main() {

	emailBytes, err := ioutil.ReadFile(`C:\Users\cduncan\Downloads\temp\phantomjs\economist_email_raw.txt`)
	if err != nil {
		panic(err)
	}

	msg, err := email.ParseMessage(bytes.NewBuffer(emailBytes))
	if err != nil {
		panic(err)
	}

	for i, part := range msg.PartsContentTypePrefix("text/html") {
		ioutil.WriteFile(
			fmt.Sprintf(`C:\Users\cduncan\Downloads\temp\phantomjs\economist_email_text_html_%d.txt`, i),
			part.Body,
			0644,
		)
	}
}
