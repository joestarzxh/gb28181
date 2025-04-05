// Code generated by godddx, DO AVOID EDIT.
package mediadb

import (
	"context"

	"github.com/gowvp/gb28181/internal/core/media"
	"github.com/ixugo/goddd/pkg/orm"
	"gorm.io/gorm"
)

var _ media.StreamPushStorer = StreamPush{}

// StreamPush Related business namespaces
type StreamPush DB

// NewStreamPush instance object
func NewStreamPush(db *gorm.DB) StreamPush {
	return StreamPush{db: db}
}

// Find implements media.StreamPushStorer.
func (d StreamPush) Find(ctx context.Context, bs *[]*media.StreamPush, page orm.Pager, opts ...orm.QueryOption) (int64, error) {
	return orm.FindWithContext(ctx, d.db, bs, page, opts...)
}

// Get implements media.StreamPushStorer.
func (d StreamPush) Get(ctx context.Context, model *media.StreamPush, opts ...orm.QueryOption) error {
	return orm.FirstWithContext(ctx, d.db, model, opts...)
}

// Add implements media.StreamPushStorer.
func (d StreamPush) Add(ctx context.Context, model *media.StreamPush) error {
	return d.db.WithContext(ctx).Create(model).Error
}

// Edit implements media.StreamPushStorer.
func (d StreamPush) Edit(ctx context.Context, model *media.StreamPush, changeFn func(*media.StreamPush), opts ...orm.QueryOption) error {
	return orm.UpdateWithContext(ctx, d.db, model, changeFn, opts...)
}

// Del implements media.StreamPushStorer.
func (d StreamPush) Del(ctx context.Context, model *media.StreamPush, opts ...orm.QueryOption) error {
	return orm.DeleteWithContext(ctx, d.db, model, opts...)
}

func (d StreamPush) Session(ctx context.Context, changeFns ...func(*gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, fn := range changeFns {
			if err := fn(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (d StreamPush) EditWithSession(tx *gorm.DB, model *media.StreamPush, changeFn func(b *media.StreamPush) error, opts ...orm.QueryOption) error {
	return orm.UpdateWithSession(tx, model, changeFn, opts...)
}
