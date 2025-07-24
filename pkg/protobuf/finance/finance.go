package finance

import (
	"strings"

	"liveJob/pkg/tools/errs"
)

func (req *GetZoneByPayTypeCodeReq) Check() error {
	if len(strings.TrimSpace(req.GetPayTypeCode())) == 0 {
		return errs.ErrArgs.WithDetail("para_err")
	}

	return nil
}

func (req *MyWalletReq) Check() error {
	if len(strings.TrimSpace(req.TimeZone)) == 0 {
		return errs.ErrArgs.WithDetail("para_err")
	}
	if len(strings.TrimSpace(req.Date)) == 0 {
		return errs.ErrArgs.WithDetail("para_err")
	}

	return nil
}
