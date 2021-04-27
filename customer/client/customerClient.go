package main

import (
	"cdcon21builddriver/customer/proto"
	"cdcon21builddriver/globalProtos"
	userSrv "cdcon21builddriver/user/proto"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"log"
	"time"
)

const (
	// serviceName: service identifier
	serviceName = "cdcon21builddriver.api.customer"
	// serviceNameUser: service identifier for user service
	serviceNameUser = "cdcon21builddriver.api.user"
)

// dateLayoutISO default time format for dates entered as strings
const dateLayoutISO = "2006-01-02"

// GetCustomerById calls the customer service and retrieve the customer identified by a particular id
func GetCustomerById(ctx context.Context, customerClient proto.CustomerSrvService, custId *proto.SearchId) (*proto.Customer, error) {

	customer, err := customerClient.GetCustomerById(ctx, custId)

	if err != nil {
		log.Printf("Unable to find customer by Id. Error: %v", err)
		return nil, err
	}

	if customer.XKey == "" {
		fmt.Printf("No Customer found for id %s\n", custId.XKey)
		return nil, err
	}
	fmt.Printf("Pulled customer by id %v\n", customer)
	return customer, nil

}

// GetCustomers contacts the customer service and retrieve customer based on a search criteria
func GetCustomers(ctx context.Context, customerClient proto.CustomerSrvService) (*proto.Customers, error) {
	_, searchDate := timeStringToTimestamp("2020-10-24")

	searchParms := proto.SearchParams{
		XKey:      "ducksrus",
		Name:      "Ducks R Us",
		ValidDate: searchDate,
	}

	customers, err := customerClient.GetCustomers(ctx, &searchParms)

	if err != nil {
		log.Printf("Unable to find customers. Error: %v", err)
		return nil, err
	}
	if len(customers.GetCustomer()) == 0 {
		fmt.Printf("Customers not found for parameters %v\n", &searchParms)
		return nil, fmt.Errorf("Customers not found for parameters %v\n", &searchParms)
	}

	fmt.Printf("Pulled customers %v\n", customers)
	return customers, nil

}

// CreateCustomer calls the customer service and create a new customer
func CreateCustomer(ctx context.Context, customerClient proto.CustomerSrvService) (*proto.Customer, error) {

	// var cust *proto.Customer
	// var err error

	_, validThru := timeStringToTimestamp("2021-05-24")

	newCust := proto.Customer{
		XKey: "6308345766332077057",
		Name: "Awesome Customer",
		ValidityDates: &globalProtos.GlValidityDate{
			ValidFrom: ptypes.TimestampNow(),
			ValidThru: validThru,
		},
		Modifications: &globalProtos.GlModification{
			CreateDate: ptypes.TimestampNow(),
			UpdateDate: ptypes.TimestampNow(),
			ModifiedBy: "123456789",
		},
	}

	resp, err := customerClient.CreateCustomer(ctx, &newCust)

	if err != nil {
		log.Printf("Unable to create customer. Error: %v", err)
		return nil, err
	}
	fmt.Printf("Created customer %v\n", resp.GetCustomer())

	if len(resp.GetValidationErr().GetFailureDesc()) > 0 {
		fmt.Printf("Created customer validations %v\n", resp.ValidationErr.FailureDesc)
	}
	return resp.GetCustomer(), nil
}

// UpdateCustomer calls the customer service and update a customer
func UpdateCustomer(ctx context.Context, customerClient proto.CustomerSrvService, cust *proto.Customer) (*proto.Customer, error) {
	_, validThru := timeStringToTimestamp("2021-06-26")

	cust.Name = "Just Ok Customer"
	cust.ValidityDates.ValidFrom = ptypes.TimestampNow()
	cust.ValidityDates.ValidThru = validThru
	cust.Modifications.UpdateDate = ptypes.TimestampNow()
	cust.Modifications.ModifiedBy = "3308341401806443521"

	resp, err := customerClient.UpdateCustomer(ctx, cust)

	if err != nil {
		log.Printf("Unable to update customer. Error: %v", err)
		return nil, err
	}
	fmt.Printf("Updated customer %v\n", resp.GetCustomer())

	if len(resp.GetValidationErr().GetFailureDesc()) > 0 {
		fmt.Printf("Update customer validations %v\n", resp.GetValidationErr().GetFailureDesc())
	}

	return resp.GetCustomer(), nil
}

// DeleteCustomer calls the customer service and delete the customer identified by a given id
func DeleteCustomer(ctx context.Context, customerClient proto.CustomerSrvService, custId *proto.SearchId) (int64, error) {

	resp, err := customerClient.DeleteCustomer(ctx, custId)

	if err != nil {
		log.Printf("Unable to find customer by Id. Error: %v", err)
		return 0, err
	}
	fmt.Printf("Count of customers deleted %d\n", resp.AffectedCount)

	if len(resp.GetValidationErr().GetFailureDesc()) > 0 {
		fmt.Printf("Delete customer validations %v\n", resp.GetValidationErr().GetFailureDesc())
	}
	return resp.GetAffectedCount(), nil
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

// authUser calls the user service and authenticate a user. receive a jwt token if successful
func authUser(srvClient userSrv.UserSrvService, user *userSrv.User) (*userSrv.Token, error) {
	token, err := srvClient.Auth(context.Background(), &userSrv.User{
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

// loginUser calls authUser to get an authentication token and store it in the context for use on other tasks
func loginUser(srvClient userSrv.UserSrvService) (context.Context, error) {
	myUser := &userSrv.User{
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
	service := micro.NewService(
		micro.Name("customer.client"),
	)
	service.Init()
	fmt.Println("Client Running")

	//  send requests
	ctx, err := loginUser(userSrv.NewUserSrvService(serviceNameUser, client.DefaultClient))
	if err != nil {
		return
	}

	customerClient := proto.NewCustomerSrvService(serviceName, service.Client())

	createdPromo, err := CreateCustomer(ctx, customerClient)
	if err != nil {
		return
	}

	_, _ = UpdateCustomer(ctx, customerClient, createdPromo)

	searchId := proto.SearchId{
		XKey: createdPromo.GetXKey(),
	}
	_, _ = GetCustomerById(ctx, customerClient, &searchId)
	_, _ = DeleteCustomer(ctx, customerClient, &searchId)
	_, _ = GetCustomerById(ctx, customerClient, &searchId)
	_, _ = GetCustomers(ctx, customerClient)
}
