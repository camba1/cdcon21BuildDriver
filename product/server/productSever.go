package main

import (
	"context"
	"fmt"
	adb "github.com/arangodb/go-driver"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"goTemp/globalUtils"
	"goTemp/globalerrors"
	"goTemp/product/proto"
	"goTemp/product/server/statements"
	pb "goTemp/user/proto"
	userSrv "goTemp/user/proto"
	"log"
	"os"
)

//serviceName: service identifier
const serviceName = "product"

const (
	//dbName: Name of the DB hosting the data
	dbName = "product"
	//dbAddressEnvVarName: Name of Environment variable that contains the address to the DB
	dbAddressEnvVarName = "DB_ADDRESS"
	//dbUserEnvVarName: Name of Environment variable that contains the database username
	dbUserEnvVarName = "DB_USER"
	//dbPassEnvVarName: Name of Environment variable that contains the database password
	dbPassEnvVarName = "DB_PASS"
)

// conn: Database connection
var conn *adb.Database

//prodErr: Holds service specific errors
var prodErr statements.ProdErr

//glErr: Holds the service global errors that are shared cross services
var glErr globalerrors.SrvError

//mb: Broker instance to send/receive message from pub/sub system
var mb globalUtils.MyBroker

//product: Main entry point for product related services
type product struct{}

//AuthWrapper: Authentication middleware
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return fmt.Errorf(glErr.AuthNoMetaData(req.Endpoint()))
		}

		token := meta["Token"]
		log.Printf("endpoint: %v", req.Endpoint())

		userClient := userSrv.NewUserSrvService("user", client.DefaultClient)
		outToken, err := userClient.ValidateToken(ctx, &pb.Token{Token: token})
		if err != nil {
			return err
		}
		if outToken.Valid != true {
			return fmt.Errorf(glErr.AuthInvalidToken())
		}

		if outToken.EUid == "" {
			return fmt.Errorf("unable to get user id from token for endpoint %v\n", req.Endpoint())
		}
		ctx2 := metadata.Set(ctx, "userid", outToken.EUid)

		return fn(ctx2, req, resp)
	}
}

//getDBConnString: Get the connection string to the DB
func getDBConnString() *globalUtils.DbConnParams {
	addressString := os.Getenv(dbAddressEnvVarName)
	if addressString == "" {
		log.Fatalf(glErr.DbNoConnectionString(dbAddressEnvVarName))
	}
	userNameString := os.Getenv(dbUserEnvVarName)
	if userNameString == "" {
		log.Fatalf(glErr.DbNoConnectionString(dbUserEnvVarName))
	}
	passString := os.Getenv(dbPassEnvVarName)
	if passString == "" {
		log.Fatalf(glErr.DbNoConnectionString(dbPassEnvVarName))
	}
	return &globalUtils.DbConnParams{
		Address:  addressString,
		Username: userNameString,
		Password: passString,
	}
}

//connectToDB: Call the Util pgxDBConnect to connect to the database. Service will try to connect a few times
//before giving up and throwing an error
func connectToDB() *adb.Database {
	var dbConnect globalUtils.ArangoConnect
	db, err := dbConnect.ConnectToDBWithRetry(dbName, getDBConnString())
	if err != nil {
		log.Fatalf(glErr.DbNoConnection(dbName, err))
	}
	return db
}

func main() {
	//instantiate service
	service := micro.NewService(
		micro.Name(serviceName),
		micro.WrapHandler(AuthWrapper),
	)

	service.Init()
	err := proto.RegisterProductSrvHandler(service.Server(), new(product))
	if err != nil {
		log.Fatalf(glErr.SrvNoHandler(err))
	}

	//Connect to DB
	conn = connectToDB()

	//defer conn.Close(context.Background())

	//setup the nats broker
	mb.Br = service.Options().Broker

	// Run Service
	err = service.Run()
	if err != nil {
		log.Fatalf(glErr.SrvNoStart(serviceName, err))
	}
}