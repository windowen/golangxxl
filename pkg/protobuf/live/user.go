package live

import "liveJob/pkg/tools/errs"

func (req *BillsListReq) Check() error {
	if req.GetBillType() < 0 {
		return errs.ErrArgs.WithDetail("invalid_request")
	}
	if req.GetLastId() < 0 {
		return errs.ErrArgs.WithDetail("invalid_request")
	}
	if req.GetPageSize() <= 0 || req.GetPageSize() > 999 {
		return errs.ErrArgs.WithDetail("invalid_request")
	}
	return nil
}
