package activity

import (
	"fmt"
	"log"

	"github.com/bxcodec/faker/v4"
)

func (p *ProfRep) ImportEmails(emails []string) {

	for i := range emails {
		if p.doesExist(emails[i]) {
			fmt.Println(emails[i], "exists")
			continue
		} else {
			fmt.Println(emails[i], "doesn't exist")
			_, err := p.insert(emails[i], faker.FirstName(), faker.LastName())
			if err != nil {
				log.Fatalf("cannot insert email: %s", err.Error())
			}
		}

	}

}
