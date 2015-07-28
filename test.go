package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"text/template"
)

type User struct {
	Person    string
	Email     string
	Password  string
	Receivers []string
}
type TemplateMask struct {
	U   *User
	Rec string
}

func newUser() *User {
	u := new(User)
	r := bufio.NewReader(os.Stdin)
	for i := 0; i < 4; i++ {
		str, _ := r.ReadString('\n')
		switch i {
		case 0:
			u.Person = str
		case 1:
			u.Email = str
		case 2:
			u.Password = str
		case 3:
			u.Receivers = strings.Split(str, ",")
		}
	}
	return u
}

var B bytes.Buffer

func (U *User) send(R []string) {
	auth := smtp.PlainAuth(U.Person, U.Email, U.Password, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, U.Email, R, B.Bytes())
	if err != nil {
		log.Println(err)
		fmt.Println("Error in sending the mail.")
	} else {
		fmt.Println("Email sent successfully")
	}
}

const multMessage = `{{define "multMessage"}}
      Dear All
      Sorry I am too busy and will try to reply later this week!

      From,
      {{.U.Person}}
      {{end}}
`
const singleMessage = `{{define "singleMessage"}}
      Dear {{.Rec}}
      Sorry I am too busy and will try to reply later this week!

      From,
      {{.U.Person}}
      {{end}}
`

func (T *TemplateMask) prepareTemplate(choice string, tt string) {
	t := template.Must(template.New(choice).Parse(tt))
	err := t.Execute(&B, T)
	if err != nil {
		fmt.Println("Error in trying to execute the template.")
	}
}

func main() {
	fmt.Println("Enter your name, email address, password, and recipients each on a new line.\n(separate each recipient by a comma. Hit enter to start a new line.)")
	obj := newUser()
	if len(obj.Receivers) > 1 {
		fmt.Println("Would you like to send a separate email to each user or one to address all users? Type yes or no")
		r := bufio.NewReader(os.Stdin)
		str, _ := r.ReadString('\n')
		if str == "yes" {
			count := 0
			for range obj.Receivers {
				//set the receiver in the template
				t := TemplateMask{obj, obj.Receivers[count]}
				//make a slice of the recipients array of each current recipient and send
				X := obj.Receivers[count : count+1]
				t.prepareTemplate("singleMessage", singleMessage)
				obj.send(X)
				count++
			}
		} else {
			t := TemplateMask{obj, ""}
			t.prepareTemplate("multMessage", multMessage)
			obj.send(obj.Receivers)
		}
	} else {
		t := TemplateMask{obj, obj.Receivers[0]}
		t.prepareTemplate("singleMessage", singleMessage)
		obj.send(obj.Receivers)
	}
}
