

package errors

import "errors"

var (
	ERR_ASSET_NAME_INVALID        = errors.New("asset name invalid, too long")
	ERR_ASSET_PRECISION_INVALID   = errors.New("asset precision invalid")
	ERR_ASSET_AMOUNT_INVALID      = errors.New("asset amount invalid")
	ERR_ASSET_CHECK_OWNER_INVALID = errors.New("asset owner invalid")
)
