package main

import (
	"cdcon21builddriver/globalProtos"
	"cdcon21builddriver/globalUtils"
	"cdcon21builddriver/globalerrors"
	"cdcon21builddriver/product/proto"
	"context"
	"log"
	"strconv"
	"time"
)

// checkMandatoryFields ensures that all mandatory fields are populated properly
func checkMandatoryFields(product *proto.Product) ([]string, error) {
	var failureDesc []string
	if product.GetName() == "" {
		failureDesc = append(failureDesc, glErr.MissingField("name"))
	}
	if product.GetHierarchyLevel() == "" {
		failureDesc = append(failureDesc, glErr.MissingField("Hierarchy Level"))
	}
	dateValidation, err := globalUtils.CheckValidityDates(product.GetValidityDates().GetValidFrom(), product.GetValidityDates().GetValidThru())
	if err != nil {
		return nil, err
	}
	failureDesc = append(failureDesc, dateValidation...)

	return failureDesc, nil
}

// SetMandatoryFields presets the mandatory fields that need to be populated before insert,delete or update
func SetMandatoryFields(ctx context.Context, product *proto.Product, isInsert bool) error {

	log.Println("Start Set Mandatory Fields")

	tempDates, _ := globalUtils.TimeToTimeStampPPB(time.Now(), time.Now().AddDate(1, 0, 0))
	if product.GetValidityDates() == nil {
		product.ValidityDates = &globalProtos.GlValidityDate{}
	}

	if product.GetModifications() == nil {
		product.Modifications = &globalProtos.GlModification{}
	}

	if isInsert {
		if product.GetValidityDates().GetValidFrom() == nil {
			product.GetValidityDates().ValidFrom = tempDates[0]
			product.GetValidityDates().ValidThru = tempDates[1]
		}
		product.Modifications.CreateDate = tempDates[0]
	}

	log.Println("Set Mandatory Fields - Set Mod Date")

	product.GetModifications().UpdateDate = tempDates[0]

	log.Println("Set Mandatory Fields - Set user")

	currentUser, err := getCurrentUser(ctx)
	if err != nil {
		return err
	}
	product.Modifications.ModifiedBy = currentUser

	log.Println("End Set Mandatory Fields")

	return nil
}

// getCurrentUser gets the user from the context. Notice that the authorization service returns a int64 and we convert to string
func getCurrentUser(ctx context.Context) (string, error) {
	var auth globalUtils.AuthUtils
	currentUser, err := auth.GetCurrentUserFromContext(ctx)
	if err != nil {
		log.Printf(glErr.AuthNoUserInToken(err))
		return "", err
	}
	return strconv.FormatInt(currentUser, 10), nil
}

// BeforeCreateProduct calls data validations before creating a product
func (p *Product) BeforeCreateProduct(ctx context.Context, product *proto.Product, validationErr *proto.ValidationErr) error {
	_ = ctx

	err := SetMandatoryFields(ctx, product, true)
	if err != nil {
		return err
	}

	validation, err := checkMandatoryFields(product)
	if err != nil {
		return err
	}
	validationErr.FailureDesc = append(validationErr.FailureDesc, validation...)

	if len(validationErr.FailureDesc) > 0 {
		return &globalerrors.ValidationError{Source: "BeforeCreateUser", FailureDesc: validationErr.FailureDesc}
	}
	return nil
}

// BeforeUpdateProduct calls data validations before updating a product
func (p *Product) BeforeUpdateProduct(ctx context.Context, product *proto.Product, validationErr *proto.ValidationErr) error {
	_ = ctx

	err := SetMandatoryFields(ctx, product, false)
	if err != nil {
		return err
	}

	validation, err := checkMandatoryFields(product)
	if err != nil {
		return err
	}
	validationErr.FailureDesc = append(validationErr.FailureDesc, validation...)

	if len(validationErr.FailureDesc) > 0 {
		return &globalerrors.ValidationError{Source: "BeforeCreateUser", FailureDesc: validationErr.FailureDesc}
	}
	return nil
}

// BeforeDeleteProduct calls data validations before deleting a product
func (p *Product) BeforeDeleteProduct(ctx context.Context, product *proto.Product, validationErr *proto.ValidationErr) error {
	_ = ctx

	err := SetMandatoryFields(ctx, product, false)
	if err != nil {
		return err
	}

	if len(validationErr.FailureDesc) > 0 {
		return &globalerrors.ValidationError{Source: "BeforeDeleteProduct", FailureDesc: validationErr.FailureDesc}
	}
	return nil
}

// AfterCreateProduct calls processes to be run after product create
func (p *Product) AfterCreateProduct(ctx context.Context, product *proto.Product, afterFuncErr *proto.AfterFuncErr) error {
	_ = ctx

	failureDesc := p.sendUserAudit(ctx, serviceName, "AfterCreateProduct", "insert", "product", product.GetXKey(), product)
	if len(failureDesc) > 0 {
		afterFuncErr.FailureDesc = append(afterFuncErr.FailureDesc, failureDesc)
	}

	// if len(afterFuncErr.FailureDesc) > 0 {
	// 	return &globalerrors.ValidationError{Source: "AfterCreateUser", FailureDesc: afterFuncErr.FailureDesc}
	// }
	return nil
}

// AfterUpdateProduct calls processes to be run after product update
func (p *Product) AfterUpdateProduct(ctx context.Context, product *proto.Product, afterFuncErr *proto.AfterFuncErr) error {
	_ = ctx

	failureDesc := p.sendUserAudit(ctx, serviceName, "AfterUpdateProduct", "update", "product", product.GetXKey(), product)
	if len(failureDesc) > 0 {
		afterFuncErr.FailureDesc = append(afterFuncErr.FailureDesc, failureDesc)
	}

	// if len(afterFuncErr.FailureDesc) > 0 {
	// 	return &globalerrors.ValidationError{Source: "AfterCreatePromotion"}
	// }
	return nil
}

// AfterDeleteProduct calls processes to be run after product delete
func (p *Product) AfterDeleteProduct(ctx context.Context, product *proto.Product, afterFuncErr *proto.AfterFuncErr) error {
	_ = ctx

	failureDesc := p.sendUserAudit(ctx, serviceName, "AfterDeleteProduct", "Delete", "product", product.GetXKey(), product)
	if len(failureDesc) > 0 {
		afterFuncErr.FailureDesc = append(afterFuncErr.FailureDesc, failureDesc)
	}

	// if len(afterFuncErr.FailureDesc) > 0 {
	// 	return &globalerrors.ValidationError{Source: "AfterCreatePromotion"}
	// }
	return nil
}

// sendUserAudit Converts a user to a byte array, and call AuditUtil to send message with updated record to audit service
func (p *Product) sendUserAudit(ctx context.Context, serviceName, actionFunc, actionType string, objectName string, objectId string, product *proto.Product) string {
	if !glDisableAuditRecords {
		byteMessage, err := mb.ProtoToByte(product)
		if err != nil {
			return glErr.AudFailureSending(actionType, objectId, err)
		}

		return globalUtils.AuditSend(ctx, mb, serviceName, actionFunc, actionType, objectName, objectId, byteMessage)
	}
	return ""

}
