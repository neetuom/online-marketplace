package models

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func SessionSetUp() *gocql.Session {

	//open a db connection
	var err error

	//cluster := gocql.NewCluster("127.0.0.1")
	cluster := gocql.NewCluster("cassandra-node1")
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}
	cluster.Keyspace = "go_db"
	cluster.Consistency = gocql.Quorum
	session, err = cluster.CreateSession()

	if err != nil {
		log.Println("%s\n", err)
		panic("Failed to connect to database!")
	}

	fmt.Println("cassandra init done")

	return session
}
