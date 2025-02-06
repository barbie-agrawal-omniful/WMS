package appinit

import (
	"context"
	"time"

	postgres "wms/database"

	"github.com/omniful/go_commons/config"
	opostgres "github.com/omniful/go_commons/db/sql/postgres"
	"github.com/omniful/go_commons/log"
)

var db *opostgres.DbCluster

func GetClient() *opostgres.DbCluster {
	return db
}

func setClient(databaseInstance *opostgres.DbCluster) {
	db = databaseInstance
}

func Initialize(ctx context.Context) {
	databaseInstance := initializeDB(ctx)
	setClient(databaseInstance)
}

func initializeDB(ctx context.Context) *opostgres.DbCluster {
	maxOpenConnections := config.GetInt(ctx, "postgresql.maxOpenConns")
	maxIdleConnections := config.GetInt(ctx, "postgresql.maxIdleConns")

	database := config.GetString(ctx, "postgresql.database")
	connIdleTimeout := 10 * time.Minute

	// Read Write endpoint config
	mysqlWriteServer := config.GetString(ctx, "postgresql.master.host")
	mysqlWritePort := config.GetString(ctx, "postgresql.master.port")
	mysqlWritePassword := config.GetString(ctx, "postgresql.master.password")
	mysqlWriterUsername := config.GetString(ctx, "postgresql.master.username")

	// Fetch Read endpoint config
	//mysqlReadServers := config.GetString(ctx, "postgresql.slaves.hosts")
	//mysqlReadPort := config.GetString(ctx, "postgresql.slaves.port")
	//mysqlReadPassword := config.GetString(ctx, "postgresql.slaves.password")
	//mysqlReadUsername := config.GetString(ctx, "postgresql.slaves.username")

	debugMode := config.GetBool(ctx, "postgresql.debugMode")

	// Master config i.e. - Write endpoint
	masterConfig := opostgres.DBConfig{
		Host:               mysqlWriteServer,
		Port:               mysqlWritePort,
		Username:           mysqlWriterUsername,
		Password:           mysqlWritePassword,
		Dbname:             database,
		MaxOpenConnections: maxOpenConnections,
		MaxIdleConnections: maxIdleConnections,
		ConnMaxLifetime:    connIdleTimeout,
		DebugMode:          debugMode,
	}

	// Slave config i.e. - array with read endpoints
	slavesConfig := make([]opostgres.DBConfig, 0)
	//for _, host := range strings.Split(mysqlReadServers, ",") {
	//	slaveConfig := opostgres.DBConfig{
	//		Host:               host,
	//		Port:               mysqlReadPort,
	//		Username:           mysqlReadUsername,
	//		Password:           mysqlReadPassword,
	//		Dbname:             database,
	//		MaxOpenConnections: maxOpenConnections,
	//		MaxIdleConnections: maxIdleConnections,
	//		ConnMaxLifetime:    connIdleTimeout,
	//		DebugMode:          debugMode,
	//	}
	//	slavesConfig = append(slavesConfig, slaveConfig)
	//}

	db := opostgres.InitializeDBInstance(masterConfig, &slavesConfig)
	log.InfofWithContext(ctx, "Initialized Postgres DB client")
	postgres.SetCluster(db)
	return db
}
