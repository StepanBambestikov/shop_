package baseapp

import (
	"context"
)

type IApp interface {
	Start(ctx context.Context) error
}
