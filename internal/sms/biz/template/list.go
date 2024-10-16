package template

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/rosas99/streaming/internal/pkg/meta"
	v1 "github.com/rosas99/streaming/pkg/api/sms/v1"
	"github.com/rosas99/streaming/pkg/log"
	"golang.org/x/sync/errgroup"
	"sync"
)

// List retrieves a list of all templates from the database.
func (t *templateBiz) List(ctx context.Context, rq *v1.ListTemplateRequest) (*v1.ListTemplateResponse, error) {

	count, list, err := t.ds.Templates().List(ctx, rq.TemplateCode, meta.WithOffset(rq.Offset), meta.WithLimit(rq.Limit))
	if err != nil {
		log.C(ctx).Errorw(err, "Failed to list templates from storage")
		return nil, err
	}

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)
	for _, item := range list {
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				var templateReply v1.TemplateReply
				_ = copier.Copy(templateReply, item)
				templateReply.CreatedAt = item.CreatedAt.Format("2006-01-02 15:04:05")
				templateReply.UpdatedAt = item.UpdatedAt.Format("2006-01-02 15:04:05")
				m.Store(item.ID, templateReply)
				return nil
			}

		})
	}

	if err := eg.Wait(); err != nil {
		log.C(ctx).Errorw(err, "Failed to wait all function calls returned")
		return nil, err
	}

	// The following code block is used to maintain the consistency of query order.
	templates := make([]*v1.TemplateReply, 0, len(list))
	for _, item := range list {
		template, _ := m.Load(item.ID)
		templates = append(templates, template.(*v1.TemplateReply))
	}

	log.C(ctx).Debugw("Get templates from backend storage", "count", len(templates))

	return &v1.ListTemplateResponse{TotalCount: count, Templates: templates}, nil
}

func (t *templateBiz) ListWithBadPerformance(ctx context.Context, rq *v1.ListTemplateRequest) (*v1.ListTemplateResponse, error) {

	count, list, err := t.ds.Templates().List(ctx, rq.TemplateCode, meta.WithOffset(rq.Offset), meta.WithLimit(rq.Limit))
	if err != nil {
		log.C(ctx).Errorw(err, "Failed to list orders from storage")
		return nil, err
	}

	templates := make([]*v1.TemplateReply, 0, len(list))

	for _, item := range list {

		var templateReply v1.TemplateReply
		_ = copier.Copy(templateReply, item)
		templateReply.CreatedAt = item.CreatedAt.Format("2006-01-02 15:04:05")
		templateReply.UpdatedAt = item.UpdatedAt.Format("2006-01-02 15:04:05")

		templates = append(templates, &templateReply)
	}

	log.C(ctx).Debugw("Get templates from backend storage", "count", len(templates))

	return &v1.ListTemplateResponse{TotalCount: count, Templates: templates}, nil
}
