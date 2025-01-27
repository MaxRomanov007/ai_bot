package paginator

type Option func(p *Paginator)

// PerPage sets the number of items to be displayed per page.
func PerPage(perPage int) Option {
	return func(p *Paginator) {
		p.perPage = perPage
	}
}

// WithText sets the text to be displayed
func WithText(text string) Option {
	return func(p *Paginator) {
		p.Text = text
	}
}

func NoDeleteBeforeHandler() Option {
	return func(p *Paginator) {
		p.isDeleteBeforeHandler = false
	}
}

// WithCloseButtonText sets the close button text to be displayed
func WithCloseButtonText(text string) Option {
	return func(p *Paginator) {
		p.closeButton = text
	}
}

// OnError sets the error handler
func OnError(f OnErrorHandler) Option {
	return func(p *Paginator) {
		p.onError = f
	}
}

// WithPrefix is a keyboard option that sets a prefix for the widget
func WithPrefix(s string) Option {
	return func(w *Paginator) {
		w.prefix = s
	}
}
