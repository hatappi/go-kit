package option

type SaveOption struct {
	ContentType *string
}

type SaveOptionFunc func(opt *SaveOption)

func SaveOptionWithContentType(ct string) SaveOptionFunc {
	return func(opt *SaveOption) {
		opt.ContentType = &ct
	}
}
