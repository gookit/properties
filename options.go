package properties

import "github.com/mitchellh/mapstructure"

// DefaultTagName for mapping data to struct.
var DefaultTagName = "properties"

// OpFunc custom option func
type OpFunc func(opts *Options)

// Options for config
type Options struct {
	// Debug open debug mode
	Debug bool
	// ParseEnv parse ENV var name, default True. eg: "${SHELL}"
	ParseEnv bool
	// ParseVar reference. eg: "${other.var.name}". default: true
	ParseVar bool
	// ParseTime string on binding struct. eg: 3s -> 3*time.Second
	ParseTime bool
	// TagName for binding data to struct. default: properties
	TagName string
	// TrimValue trim "\n" for value string. default: false
	TrimValue bool

	// InlineComment support split inline comments. default: false
	//
	// allow chars: #, //
	InlineComment bool
	// InlineSlice support parse the inline slice. eg: [23, 34]. default: false
	InlineSlice bool
	// MapStructConfig for binding data to struct.
	MapStructConfig mapstructure.DecoderConfig
	// BeforeCollect value handle func, you can return a new value.
	BeforeCollect func(name string, val any) any
}

func (opts *Options) makeDecoderConfig() *mapstructure.DecoderConfig {
	decConf := opts.MapStructConfig
	// compatible with settings on opts.TagName
	if decConf.TagName == "" {
		decConf.TagName = opts.TagName
	}

	// parse time string on binding to struct
	if opts.ParseTime || decConf.DecodeHook == nil {
		decConf.DecodeHook = ValDecodeHookFunc()
	}

	return &decConf
}

func newDefaultOption() *Options {
	return &Options{
		ParseVar: true,
		TagName:  DefaultTagName,
		// map struct config
		MapStructConfig: mapstructure.DecoderConfig{
			TagName: DefaultTagName,
			// will auto convert string to int/uint
			WeaklyTypedInput: true,
		},
	}
}

// WithDebug open debug mode
func WithDebug(opts *Options) {
	opts.Debug = true
}

// ParseEnv open parse ENV var string.
func ParseEnv(opts *Options) {
	opts.ParseEnv = true
}

// ParseTime open parse time string.
func ParseTime(opts *Options) {
	opts.ParseTime = true
}

// ParseInlineSlice open parse inline slice
func ParseInlineSlice(opts *Options) {
	opts.InlineSlice = true
}

// WithTagName custom tag name on binding struct
func WithTagName(tagName string) OpFunc {
	return func(opts *Options) {
		opts.TagName = tagName
		opts.MapStructConfig.TagName = tagName
	}
}
