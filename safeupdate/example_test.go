package safeupdate_test

import (
	"fmt"

	"github.com/crast/dynatools/safeupdate"
	"gopkg.in/underarmour/dynago.v1"
)

var executor dynago.Executor
var personKeys = []string{"Id"} // the primary keys of the "Person" type

func Example() {
	client := dynago.NewClient(executor)

	// Create the person we want in the database.
	person := Person{Id: 5, Name: "Bob", Age: 40, Address: "123 fake street"}
	client.PutItem("Person", person.AsDocument()).Execute()

	// Bob's address was fake, let's blank it out:
	person.Address = ""
	update := safeupdate.Build(personKeys, person.AsDocument())

	// output: SET Age=:a, #b=:b REMOVE Address
	fmt.Println(update.Expression)

	// Okay, we proved our point, let's execute it:
	update.Apply(client, "Person").Execute()
}

// Just an example of a normal "ORM" model
type Person struct {
	Id      int
	Name    string
	Age     int
	Address string
}

func (p Person) AsDocument() dynago.Document {
	return dynago.Document{
		"Id":      p.Id,
		"Name":    p.Name,
		"Age":     p.Age,
		"Address": p.Address,
	}
}

func init() {
	// We'd want a real executor if this was a real application
	executor = &dynago.MockExecutor{}
}
