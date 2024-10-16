package service

import (
	"context"
	v1 "github.com/rosas99/streaming/pkg/api/sms/v1"
	"github.com/rosas99/streaming/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateTemplate is a method for creating a new template.
// It takes a CreateTemplateRequest as input and returns CreateTemplateResponse or an error.
func (s *SmsServerService) CreateTemplate(ctx context.Context, rq *v1.CreateTemplateRequest) (*v1.CreateTemplateResponse, error) {
	log.C(ctx).Infow("CreateTemplate function called")
	return s.biz.Templates().Create(ctx, rq)
}

// ListTemplate is a method for listing template.
// It takes a ListTemplateRequest as input and returns a ListTemplateResponse with the template or an error.
func (s *SmsServerService) ListTemplate(ctx context.Context, rq *v1.ListTemplateRequest) (*v1.ListTemplateResponse, error) {
	log.C(ctx).Infow("ListTemplate function called")
	return s.biz.Templates().List(ctx, rq)
}

// GetTemplate is a method for retrieving a specific template.
// It takes a GetTemplateRequest as input and returns a TemplateReply with the secret or an error.
func (s *SmsServerService) GetTemplate(ctx context.Context, id int64) (*v1.TemplateReply, error) {
	log.C(ctx).Infow("GetTemplate function called")
	return s.biz.Templates().Get(ctx, id)
}

// UpdateTemplate is a method for updating a template.
// It takes an UpdateTemplateRequest as input and returns an Empty message or an error.
func (s *SmsServerService) UpdateTemplate(ctx context.Context, id int64, rq *v1.UpdateTemplateRequest) (*emptypb.Empty, error) {
	log.C(ctx).Infow("UpdateTemplate function called")
	return &emptypb.Empty{}, s.biz.Templates().Update(ctx, id, rq)
}

// DeleteTemplate is a method for deleting a template.
// It takes a DeleteTemplateRequest as input and returns an Empty message or an error.
func (s *SmsServerService) DeleteTemplate(ctx context.Context, id int64) (*emptypb.Empty, error) {
	log.C(ctx).Infow("DeleteTemplate function called")
	return &emptypb.Empty{}, s.biz.Templates().Delete(ctx, id)
}
