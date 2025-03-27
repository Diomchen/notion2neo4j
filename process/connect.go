package process

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/viper"
)

var NEO4J_DIVER neo4j.DriverWithContext

func init() {
	viper.SetConfigFile("conf/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	uri := viper.GetString("db.NEO4J_URI")
	username := viper.GetString("db.NEO4J_USERNAME")
	password := viper.GetString("db.NEO4J_PASSWORD")
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		panic(err)
	}
	context.Background()
	NEO4J_DIVER = driver
}
