package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin2/pkg/logging"
)

func MakeErrors(errors []*validation.Error) {

	for _, error := range errors {

		logging.Info(error.Key,error.Message)
	}
	return
}
