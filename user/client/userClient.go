package main

import (
	pb "cdcon21builddriver/user/proto"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/metadata"
	"log"
	"time"
)

// serviceNameUser service identifier for user service
const serviceNameUser = "cdcon21builddriver.api.user"

// dateLayoutISO Default time format for dates entered as strings
const dateLayoutISO = "2006-01-02"

// CreateUser  calls the user service and create a new user
func CreateUser(ctx context.Context, srvClient pb.UserSrvService) (*pb.User, error) {
	// var outUser *pb.User
	// var err error

	_, validThru := timeStringToTimestamp("2021-05-24")

	newUser := pb.User{
		Firstname: "Huge",
		Lastname:  "Microbe",
		ValidFrom: ptypes.TimestampNow(),
		ValidThru: validThru,
		Active:    true,
		Pwd:       "1234",
		Email:     "microbes@tiny.com",
		Company:   "Tiny",
	}

	// if serverAddress != "" {
	// 	outUser, err = srvClient.CreateUser(context.Background(), &newUser, client.WithAddress(serverAddress))
	// } else {
	resp, err := srvClient.CreateUser(ctx, &newUser)
	// }

	if err != nil {
		log.Printf("Unable to create user. Error: %v", err)
		return nil, err
	}
	fmt.Printf("Created user %v\n", resp.GetUser())

	if len(resp.GetValidationErr().GetFailureDesc()) > 0 {
		fmt.Printf("Created user validations %v\n", resp.ValidationErr.FailureDesc)
	}
	return resp.GetUser(), nil
}

// UpdateUser calls the user service and update a user
func UpdateUser(ctx context.Context, srvClient pb.UserSrvService, user *pb.User) (*pb.User, error) {
	// var outUser *pb.User
	// var err error

	_, validThru := timeStringToTimestamp("2021-06-26")

	user.Firstname = "Incredible"
	user.Lastname = "Green Guy"
	user.ValidFrom = ptypes.TimestampNow()
	user.ValidThru = validThru
	user.Active = false
	user.Pwd = "5678"
	user.Email = "microbes@tiny.com"
	// user.Email = "cow@mymail.com"
	user.Company = "Tiny"

	resp, err := srvClient.UpdateUser(ctx, user)

	if err != nil {
		log.Printf("Unable to update user. Error: %v", err)
		return nil, err
	}
	fmt.Printf("Updated user %v\n", resp.GetUser())
	if len(resp.GetValidationErr().GetFailureDesc()) > 0 {
		fmt.Printf("Update user validations %v\n", resp.GetValidationErr().GetFailureDesc())
	}
	return resp.GetUser(), nil
}

// GetUserById  calls the user service and retrieve the user identified by a particular id
func GetUserById(ctx context.Context, srvClient pb.UserSrvService, searchId *pb.SearchId) (*pb.User, error) {
	var outUser *pb.User
	var err error

	outUser, err = srvClient.GetUserById(ctx, searchId)

	if err != nil {
		log.Printf("Unable to find user by Id. Error: %v", err)
		return nil, err
	}

	if outUser.Id == 0 {
		log.Printf("No user found for id %d\n", searchId.Id)
		return nil, fmt.Errorf("No user found for id %d\n", searchId.Id)
	}

	fmt.Printf("Pulled user by id %v\n", outUser)
	return outUser, nil
}

// DeleteUser calls the user service and DeleteUser the user identified by a given id
func DeleteUser(ctx context.Context, srvClient pb.UserSrvService, searchId *pb.SearchId) (int64, error) {

	resp, err := srvClient.DeleteUser(ctx, searchId)

	if err != nil {
		log.Printf("Unable to find user by Id. Error: %v", err)
		return 0, err
	}
	fmt.Printf("Count of users deleted %d\n", resp.GetAffectedCount())

	if len(resp.GetValidationErr().GetFailureDesc()) > 0 {
		fmt.Printf("Delete user validations %v\n", resp.GetValidationErr().GetFailureDesc())
	}

	return resp.GetAffectedCount(), nil
}

// GetUsers contact the user service and retrieve users based on a search criteria
func GetUsers(ctx context.Context, srvClient pb.UserSrvService) (*pb.Users, error) {
	_, searchDate := timeStringToTimestamp("2020-10-24")

	searchParms := pb.SearchParams{
		Fisrtname: "Super",
		Lastname:  "Duck",
		ValidDate: searchDate,
		Email:     "duck@mymail.com",
	}

	var outUsers *pb.Users
	var err error

	outUsers, err = srvClient.GetUsers(ctx, &searchParms)

	if err != nil {
		log.Fatalf("Unable to find users. Error: %v", err)
		return nil, err
	}
	if len(outUsers.GetUser()) == 0 {
		fmt.Printf("Users not found for parameters %v\n", &searchParms)
		return nil, fmt.Errorf("Users not found for parameters %v\n", &searchParms)
	}
	fmt.Printf("Pulled users %v\n", outUsers)
	return outUsers, nil

}

// authUser calls the user service and authenticate a user. receive a jwt token if successful
func authUser(srvClient pb.UserSrvService, user *pb.User) (*pb.Token, error) {
	token, err := srvClient.Auth(context.Background(), &pb.User{
		Email: user.Email,
		Pwd:   user.Pwd,
	})
	if err != nil {
		log.Printf("Unable to find token. Error: %v\n", err)
		return nil, err
	}
	fmt.Printf("Got token: %v\n", token)
	return token, err
}

// timeStringToTimestamp converts time string to gRPC timestamp
func timeStringToTimestamp(priceVTstr string) (error, *timestamp.Timestamp) {
	priceVTtime, err := time.Parse(dateLayoutISO, priceVTstr)
	if err != nil {
		log.Fatalf("Unable to Format date %v", priceVTstr)
	}
	priceVT, err := ptypes.TimestampProto(priceVTtime)
	if err != nil {
		log.Fatalf("Unable to convert time to timestamp %v", priceVTtime)
	}
	return err, priceVT
}

// loginUser calls authUser to get an authentication token and store it in the context for use on other tasks
func loginUser(srvClient pb.UserSrvService) (context.Context, error) {
	myUser := &pb.User{
		Pwd:   "1234",
		Email: "duck@mymail.com"}

	authToken, err := authUser(srvClient, myUser)
	if err != nil {
		return nil, err
	}

	ctx := metadata.NewContext(context.Background(), map[string]string{
		// "token": authToken.Token,
		"Authorization": "Bearer " + authToken.Token,
	})
	return ctx, nil
}

func main() {

	// define service
	service := micro.NewService(
		micro.Name("user.client"),
	)
	service.Init()
	fmt.Println("Client Running")
	srvClient := pb.NewUserSrvService(serviceNameUser, service.Client())

	//  send requests
	ctx, err := loginUser(srvClient)
	if err != nil {
		return
	}

	createdUser, err := CreateUser(ctx, srvClient)
	if err != nil {
		return
	}

	_, _ = UpdateUser(ctx, srvClient, createdUser)

	searchId := pb.SearchId{
		Id: createdUser.Id,
	}

	_, _ = GetUserById(ctx, srvClient, &searchId)
	_, _ = DeleteUser(ctx, srvClient, &searchId)
	_, _ = GetUserById(ctx, srvClient, &searchId)
	_, _ = GetUsers(ctx, srvClient)
}
