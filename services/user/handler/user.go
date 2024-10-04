package handler

import (
	"context"
	"github.com/scg130/tools"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/google/uuid"

	"errors"
	user "user/proto/user"
	"user/repo"
)

type UserSrv struct {
	U repo.User
}

func (e *UserSrv) Login(ctx context.Context, req *user.Request, rsp *user.Response) error {
	userInfo, err := e.U.FindByPhone(req.Phone)
	if err != nil {
		rsp.Code = -1
		return err
	}
	if userInfo.Id == 0 {
		rsp.Code = -2
		rsp.Msg = "phone not exist"
		return nil
	}
	if tools.CompareHashAndPasswd(req.Passwd, userInfo.Password) {
		rsp.Code = 0
		rsp.Msg = "ok"
		rsp.Data = &user.UserInfo{
			UserId: userInfo.Id,
			Phone:  userInfo.Phone,
			Avatar: userInfo.Avatar,
		}
		return nil
	}
	rsp.Code = -1
	rsp.Msg = "login failure"
	return nil
}

func (e *UserSrv) Register(ctx context.Context, req *user.Request, rsp *user.Response) error {
	userInfo, err := e.U.FindByPhone(req.Phone)
	if err != nil {
		rsp.Code = -1
		return err
	}
	if userInfo.Id > 0 {
		rsp.Code = -2
		rsp.Msg = "user exist"
		return nil
	}

	passwd, err := tools.GeneratePasswd(req.Passwd)
	if err != nil {
		return err
	}
	u := &repo.User{
		Phone:      req.Phone,
		Password:   passwd,
		Avatar:     req.Avatar,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Version:    uuid.New().ID(),
	}
	_, err = e.U.Create(u)
	if err != nil {
		rsp.Code = -1
		rsp.Msg = "create user fail"
		return errors.New("fail")
	}
	rsp.Code = 0
	rsp.Msg = "ok"
	rsp.Data = &user.UserInfo{
		Phone:  req.Phone,
		UserId: u.Id,
	}
	return nil
}

func (e *UserSrv) Find(ctx context.Context, req *user.Request, rsp *user.Response) error {
	userInfo, err := e.U.FindByPhone(req.Phone)
	if err != nil {
		logrus.Error(err)
		rsp.Code = -1
		return err
	}

	if userInfo.Id == 0 {
		rsp.Code = -2
		rsp.Msg = "phone not exist"
		rsp.Data = &user.UserInfo{
			UserId: 0,
		}
		return nil
	}
	rsp.Code = 0
	rsp.Msg = "success"
	rsp.Data = &user.UserInfo{
		UserId: userInfo.Id,
		Phone:  userInfo.Phone,
		Avatar: userInfo.Avatar,
	}
	return nil
}

func (e *UserSrv) GetById(ctx context.Context, req *user.Request, rsp *user.Response) error {
	userInfo, err := e.U.GetByUserId(req.UserId)
	if err != nil {
		return err
	}
	if userInfo.Id == 0 {
		rsp.Code = -1
		rsp.Msg = err.Error()
		return nil
	}
	rsp.Code = 0
	rsp.Data = &user.UserInfo{
		UserId: userInfo.Id,
		Phone:  userInfo.Phone,
		Avatar: userInfo.Avatar,
	}
	return nil
}

func (e *UserSrv) ResetPasswd(ctx context.Context, req *user.Request, rsp *user.Response) error {
	passwd, err := tools.GeneratePasswd(req.Passwd)
	if err != nil {
		return err
	}

	err = e.U.UpdatePasswdById(req.UserId, passwd)
	if err != nil {
		return err
	}
	return nil
}
