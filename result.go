package wechatpay

import "encoding/json"

// Result is a result after call client.Do
type Result struct {
	Body      []byte
	Timestamp int64
	Nonce     string
	Signature string
	SerialNo  string
	Err       error
}

// Scan data from the response into the dest object.
func (r *Result) Scan(dest interface{}) error {
	if r.Err != nil {
		return r.Err
	}

	if len(r.Body) == 0 {
		return nil
	}

	if err := json.Unmarshal(r.Body, dest); err != nil {
		return err
	}

	return nil
}

// Error return the error.
func (r *Result) Error() error {
	return r.Err
}
