package gofrcontext

import (
	"context"

	"github.com/devum25/cloudnativego/transactionlog"
)

type GofrContext struct{
	context.Context
	logger transactionlog.TransactionLogger
}
