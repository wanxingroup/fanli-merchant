package merchant

import (
	"fmt"
	"time"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/errors"
	rpclog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	idcreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/validate"
)

func (controller *Controller) RecordCode(ctx context.Context, req *protos.RecordCodeRequest) (reply *protos.RecordCodeReply, err error) {
	logger := rpclog.WithRequestId(ctx, log.GetLogger()).WithField("requestData", req)

	if req == nil {
		logger.Error("get request data is nil")
		return nil, fmt.Errorf(constant.ErrorMessageGetRequestDataIsNil)
	}

	if err := controller.validateRecordCodeRequestData(req); err != nil {

		logger.WithError(err).Info("validate parameter failed")
		return &protos.RecordCodeReply{
			Err: convertError(err, int64(errors.RequestParamError.ErrorCode)),
		}, nil
	}

	if !validate.CheckMobile(req.GetMobile()) {
		return &protos.RecordCodeReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeMobileFormatInvalid,
				Message: constant.ErrorMessageMobileFormatInvalid,
				Stack:   nil,
			},
		}, nil
	}

	var expiredTime time.Time
	expiredTimeString := req.GetExpiredAt()
	if validation.IsEmpty(expiredTimeString) {
		expiredTime = time.Now().Add(DefaultExpireTime * time.Minute)
	} else {
		expiredTime, err = time.Parse("2006-01-02 15:04:05", expiredTimeString)
		if err != nil {
			logger.WithError(err).Error("expired time is error")
			expiredTime = time.Now().Add(DefaultExpireTime * time.Minute)
		}
	}

	record := &model.MerchantCode{
		Id:        idcreator.NextID(),
		Mobile:    req.GetMobile(),
		Code:      req.GetCode(),
		ExpiredAt: expiredTime,
		CodeType:  model.CodeTypeFindPassword,
	}

	// 写入数据库
	err = database.GetDB(constant.DatabaseConfigKey).Create(record).Error
	if err != nil {

		logger.WithError(err).Error("record code error")
		return &protos.RecordCodeReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeRecordCodeError,
				Message: constant.ErrorMessageRecordCodeError,
			},
		}, nil
	}

	return &protos.RecordCodeReply{
		RecordId: record.Id,
	}, nil
}

func (controller *Controller) validateRecordCodeRequestData(req *protos.RecordCodeRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Code, CodeRules...),
	)
}
