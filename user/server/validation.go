package main

import (
	"cdcon21builddriver/globalUtils"
	"cdcon21builddriver/globalerrors"
	pb "cdcon21builddriver/user/proto"
	"context"
	"strconv"
)

//mb: Broker instance to send/receive message from pub/sub system
var mb globalUtils.MyBroker

//checkMandatoryFields: Ensure that all mandatory fields are populated properly
func checkMandatoryFields(user *pb.User) ([]string, error) {
	var failureDesc []string
	if user.GetFirstname() == "" {
		failureDesc = append(failureDesc, glErr.MissingField("first name"))
	}
	if user.GetLastname() == "" {
		failureDesc = append(failureDesc, glErr.MissingField("last name"))
	}
	if user.GetPwd() == "" {
		failureDesc = append(failureDesc, glErr.MissingField("password"))
	}
	if user.GetEmail() == "" {
		failureDesc = append(failureDesc, glErr.MissingField("email"))
	}
	dateValidation, err := globalUtils.CheckValidityDates(user.ValidFrom, user.ValidThru)
	if err != nil {
		return nil, err
	}
	failureDesc = append(failureDesc, dateValidation...)

	return failureDesc, nil
}

//checkEmail: Search the DB for users with a given email and call checkEmailUnique
func checkEmail(ctx context.Context, user *pb.User, u *User, isInsert bool) (string, error) {
	var usersWithSameEmail pb.Users
	err := u.GetUsersByEmail(ctx, &pb.SearchString{Value: user.Email}, &usersWithSameEmail)
	if err != nil {
		return "", err
	}
	return checkEmailUnique(user, &usersWithSameEmail, isInsert), nil
}

//checkEmailUnique: Check if an email is unique in the database. If this is an insert operation, there should
//be no users in the db with that user. If it is an update, then only one user should have the address and that
//user should be the current user
func checkEmailUnique(user *pb.User, users *pb.Users, isInsert bool) string {

	if users == nil {
		return ""
	}
	usersCount := len(users.User)
	if isInsert && usersCount > 0 {
		return userErr.InsertDupEmail()
	}
	if !isInsert && usersCount > 0 {
		if usersCount > 1 || (usersCount == 1 && users.User[0].Id != user.Id) {
			return userErr.UpdateDupEmail()
		}
	}
	return ""

}

//func SetMandatoryFields(user *pb.User){
//	validThru, _ := globalUtils.TimeToTimeStampPPB(time.Now().AddDate(1,0,0))
//	user.Active = false
//	user.ValidFrom = ptypes.TimestampNow()
//	user.ValidThru = validThru[0]
//}

//BeforeCreateUser calls data validations before creating a user
func (u *User) BeforeCreateUser(ctx context.Context, user *pb.User, validationErr *pb.ValidationErr) error {
	_ = ctx
	validation, err := checkMandatoryFields(user)
	if err != nil {
		return err
	}
	validationErr.FailureDesc = append(validationErr.FailureDesc, validation...)

	dupEmail, err2 := checkEmail(ctx, user, u, true)
	if err2 != nil {
		return err2
	}
	if dupEmail != "" {
		validationErr.FailureDesc = append(validationErr.FailureDesc, dupEmail)
	}

	if len(validationErr.FailureDesc) > 0 {
		return &globalerrors.ValidationError{Source: "BeforeCreateUser", FailureDesc: validationErr.FailureDesc}
	}
	return nil
}

//BeforeUpdateUser calls data validations before updating a user
func (u *User) BeforeUpdateUser(ctx context.Context, user *pb.User, validationErr *pb.ValidationErr) error {
	_ = ctx

	validation, err := checkMandatoryFields(user)
	if err != nil {
		return err
	}
	validationErr.FailureDesc = append(validationErr.FailureDesc, validation...)

	dupEmail, err2 := checkEmail(ctx, user, u, false)
	if err2 != nil {
		return err2
	}
	if dupEmail != "" {
		validationErr.FailureDesc = append(validationErr.FailureDesc, dupEmail)
	}

	if len(validationErr.FailureDesc) > 0 {
		return &globalerrors.ValidationError{Source: "BeforeUpdateUser", FailureDesc: validationErr.FailureDesc}
	}
	return nil
}

//BeforeDeleteUser calls data validations before deleting a user
func (u *User) BeforeDeleteUser(ctx context.Context, user *pb.User, validationErr *pb.ValidationErr) error {
	_ = ctx
	if user.GetActive() {
		validationErr.FailureDesc = append(validationErr.FailureDesc, userErr.DelUserActive())
	}
	if len(validationErr.FailureDesc) > 0 {
		return &globalerrors.ValidationError{Source: "BeforeDeleteUser", FailureDesc: validationErr.FailureDesc}
	}
	return nil
}

//AfterCreateUser calls processes to be run after user create
func (u *User) AfterCreateUser(ctx context.Context, user *pb.User, afterFuncErr *pb.AfterFuncErr) error {
	_ = ctx

	failureDesc := u.sendUserAudit(ctx, serviceName, "AfterCreateUser", "insert", "user", user.GetId(), user)
	if len(failureDesc) > 0 {
		afterFuncErr.FailureDesc = append(afterFuncErr.FailureDesc, failureDesc)
	}

	//if len(afterFuncErr.FailureDesc) > 0 {
	//	return &globalerrors.ValidationError{Source: "AfterCreateUser", FailureDesc: afterFuncErr.FailureDesc}
	//}
	return nil
}

//AfterUpdateUser calls processes to be run after user update
func (u *User) AfterUpdateUser(ctx context.Context, user *pb.User, afterFuncErr *pb.AfterFuncErr) error {
	_ = ctx

	failureDesc := u.sendUserAudit(ctx, serviceName, "AfterUpdateUser", "update", "user", user.GetId(), user)
	if len(failureDesc) > 0 {
		afterFuncErr.FailureDesc = append(afterFuncErr.FailureDesc, failureDesc)
	}

	//if len(afterFuncErr.FailureDesc) > 0 {
	//	return &globalerrors.ValidationError{Source: "AfterCreatePromotion"}
	//}
	return nil
}

//AfterDeleteUser calls processes to be run after user delete
func (u *User) AfterDeleteUser(ctx context.Context, user *pb.User, afterFuncErr *pb.AfterFuncErr) error {
	_ = ctx

	failureDesc := u.sendUserAudit(ctx, serviceName, "AfterDeleteUser", "Delete", "user", user.GetId(), user)
	if len(failureDesc) > 0 {
		afterFuncErr.FailureDesc = append(afterFuncErr.FailureDesc, failureDesc)
	}

	//if len(afterFuncErr.FailureDesc) > 0 {
	//	return &globalerrors.ValidationError{Source: "AfterCreatePromotion"}
	//}
	return nil
}

//sendUserAudit converts a user to a byte array, and call AuditUtil to send message with updated promotion to audit service
func (u *User) sendUserAudit(ctx context.Context, serviceName, actionFunc, actionType string, objectName string, iObjectId int64, user *pb.User) string {

	if !glDisableAuditRecords {

		objectId := strconv.FormatInt(iObjectId, 10)
		byteUser, err := mb.ProtoToByte(user)
		if err != nil {
			return glErr.AudFailureSending(actionType, objectId, err)
		}

		return globalUtils.AuditSend(ctx, mb, serviceName, actionFunc, actionType, objectName, objectId, byteUser)
	}
	return ""

}
