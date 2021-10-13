package dbconn

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

	if err == nil {
		log.Println("Session connection was successful!")
	} else {
		log.Println("%s\n", err)
		panic("Failed to connect to database!")
	}

	// Check if the table already exists. Create if table does not exist
	keySpaceMeta, _ := session.KeyspaceMetadata("go_db")

	if _, exists := keySpaceMeta.Tables["payment"]; exists != true {
		/* Create payment table */
		session.Query("CREATE TABLE IF NOT EXISTS payment (payment_id text,order_id text,order_date text,invoice_number text," +
			"invoice_date text,payment_method text,tracking_id text,grand_total double,PRIMARY KEY(payment_id))").Exec()
	}

	fmt.Println("cassandra init done")

	return session
}
