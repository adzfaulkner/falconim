package recaptcha

import "github.com/dpapathanasiou/go-recaptcha"

type Checker func(remoteIp, response string) bool

func Check(key string) Checker {
	recaptcha.Init(key)
	return func(remoteIp, response string) bool {
		result, err := recaptcha.Confirm(remoteIp, response)

		if err != nil {
			return false
		}

		return result
	}
}
