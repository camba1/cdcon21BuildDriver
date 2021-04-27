package main

import (
	"cdcon21builddriver/customer/proto"
	"cdcon21builddriver/globalProtos"
	"cdcon21builddriver/globalUtils"
	"cdcon21builddriver/globalerrors"
	"context"
	"log"
	"strconv"
	"time"
)

// checkMandatoryFields: Ensure that all mandatory fields are populated properly
func checkMandatoryFields(customer *proto.Customer) ([]string, error) {
	var failureDesc []string
	if customer.GetName() == "" {
		failureDesc = append(failureDesc, glErr.MissingField("name"))
	}
	dateValidation, err := globalUtils.CheckValidityDates(customer.GetValidityDates().GetValidFrom(), customer.GetValidityDates().GetValidThru())
	if err != nil {
		return nil, err
	}
	failureDesc = append(failureDesc, dateValidation...)

	return failureDesc, nil
}

// SetMandatoryFields presets the mandatory fields that need to be populated before insert,delete or update
func SetMandatoryFields(ctx context.Context, customer *proto.Customer, isInsert bool) error {

	log.Println("Start Set Mandatory Fields")

	tempDates, _ := globalUtils.TimeToTimeStampPPB(time.Now(), time.Now().AddDate(1, 0, 0))
	if customer.GetValidityDates() == nil {
		customer.ValidityDates = &globalProtos.GlValidityDate{}
	}

	if customer.GetModifications() == nil {
		customer.Modifications = &globalProtos.GlModification{}
	}
	if isInsert {
		if customer.GetValidityDates().GetValidFrom() == nil {
			customer.GetValidityDates().ValidFrom = tempDates[0]
			customer.GetValidityDates().ValidThru = tempDates[1]
		}
		customer.Modifications.CreateDate = tempDates[0]
	}

	log.Println("Set Mandatory Fields - Set Mod Date")

	customer.GetModifications().UpdateDate = tempDates[0]

	log.Println("Set Mandatory Fields - Set user")

	currentUser, err := getCurrentUser(ctx)
	if err != nil {
		return err
	}
	customer.Modifications.ModifiedBy = currentUser

	log.Println("End Set Mandatory Fields")

	return nil
}

// getCurrentUser: Get the user from the context. Notice that the authorization service returns a int64 and we convert to string
func getCurrentUser(ctx context.Context) (string, error) {
	var auth globalUtils.AuthUtils
	currentUser, err := auth.GetCurrentUserFromContext(ctx)
	if err != nil {
		log.Printf(glErr.AuthNoUserInToken(err))
		return "", err
	}
	return strconv.FormatInt(currentUser, 10), nil
}

// BeforeCreateCustomer calls data validations before creating a customer
func (c *customer) BeforeCreateCustomer(ctx context.Context, customer *proto.Customer, validationErr *proto.ValidationErr) error {
	_ = ctx

	err := SetMandatoryFields(ctx, customer, true)
	if err != nil {
		return err
	}

	validation, err := checkMandatoryFields(customer)
	if err != nil {
		return err
	}
	validationErr.FailureDesc = append(validationErr.FailureDesc, validation...)

	if len(validationErr.FailureDesc) > 0 {
		return &globalerrors.ValidationError{Source: "BeforeCreateUser", FailureDesc: validationErr.FailureDesc}
	}
	return nil
}

// BeforeUpdateCustomer calls data validations before updating a customer
func (c *customer) BeforeUpdateCustomer(ctx context.Context, customer *proto.Customer, validationErr *proto.ValidationErr) error {
	_ = ctx

	err := SetMandatoryFields(ctx, customer, false)
	if err != nil {
		return err
	}

	validation, err := checkMandatoryFields(customer)
	if err != nil {
		return err
	}
	validationErr.FailureDesc = append(validationErr.FailureDesc, validation...)

	if len(validationErr.FailureDesc) > 0 {
		return &globalerrors.ValidationError{Source: "BeforeCreateCustomer", FailureDesc: validationErr.FailureDesc}
	}
	return nil
}

// BeforeDeleteUser calls data validations before deleting a customer
func (c *customer) BeforeDeleteCustomer(ctx context.Context, customer *proto.Customer, validationErr *proto.ValidationErr) error {
	_ = ctx

	err := SetMandatoryFields(ctx, customer, false)
	if err != nil {
		return err
	}

	if len(validationErr.FailureDesc) > 0 {
		return &globalerrors.ValidationError{Source: "BeforeDeleteUser", FailureDesc: validationErr.FailureDesc}
	}
	return nil
}

// AfterCreateCustomer calls processes to be run after customer creation
func (c *customer) AfterCreateCustomer(ctx context.Context, customer *proto.Customer, afterFuncErr *proto.AfterFuncErr) error {
	_ = ctx

	failureDesc := c.sendUserAudit(ctx, serviceName, "AfterCreateCustomer", "insert", "customer", customer.GetXKey(), customer)
	if len(failureDesc) > 0 {
		afterFuncErr.FailureDesc = append(afterFuncErr.FailureDesc, failureDesc)
	}

	// if len(afterFuncErr.FailureDesc) > 0 {
	// 	return &globalerrors.ValidationError{Source: "AfterCreateUser", FailureDesc: afterFuncErr.FailureDesc}
	// }
	return nil
}

// AfterUpdateCustomer calls processes to be run after customer update
func (c *customer) AfterUpdateCustomer(ctx context.Context, customer *proto.Customer, afterFuncErr *proto.AfterFuncErr) error {
	_ = ctx

	failureDesc := c.sendUserAudit(ctx, serviceName, "AfterUpdateCustomer", "update", "customer", customer.GetXKey(), customer)
	if len(failureDesc) > 0 {
		afterFuncErr.FailureDesc = append(afterFuncErr.FailureDesc, failureDesc)
	}

	// if len(afterFuncErr.FailureDesc) > 0 {
	// 	return &globalerrors.ValidationError{Source: "AfterCreatePromotion"}
	// }
	return nil
}

// AfterDeleteCustomer calls processes to be run after customer delete
func (c *customer) AfterDeleteCustomer(ctx context.Context, customer *proto.Customer, afterFuncErr *proto.AfterFuncErr) error {
	_ = ctx

	failureDesc := c.sendUserAudit(ctx, serviceName, "AfterDeleteCustomer", "Delete", "customer", customer.GetXKey(), customer)
	if len(failureDesc) > 0 {
		afterFuncErr.FailureDesc = append(afterFuncErr.FailureDesc, failureDesc)
	}

	// if len(afterFuncErr.FailureDesc) > 0 {
	// 	return &globalerrors.ValidationError{Source: "AfterCreatePromotion"}
	// }
	return nil
}

// sendUserAudit converts a user to a byte array, and call AuditUtil to send message with updated promotion to audit service
func (c *customer) sendUserAudit(ctx context.Context, serviceName, actionFunc, actionType string, objectName string, objectId string, customer *proto.Customer) string {
	if !glDisableAuditRecords {
		byteMessage, err := mb.ProtoToByte(customer)
		if err != nil {
			return glErr.AudFailureSending(actionType, objectId, err)
		}

		return globalUtils.AuditSend(ctx, mb, serviceName, actionFunc, actionType, objectName, objectId, byteMessage)
	}
	return ""

}
