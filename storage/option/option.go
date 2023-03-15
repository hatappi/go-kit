package option

type SaveOption struct {
	ContentType        *string
	ContentDisposition *string
}

type SaveOptionFunc func(opt *SaveOption)

func SaveOptionWithContentType(ct string) SaveOptionFunc {
	return func(opt *SaveOption) {
		opt.ContentType = &ct
	}
}

func SaveOptionWithContentDisposition(cd string) SaveOptionFunc {
	return func(opt *SaveOption) {
		opt.ContentDisposition = &cd
	}
}
