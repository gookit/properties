package properties

import "github.com/mitchellh/mapstructure"

// DefaultTagName for mapping data to struct.
var DefaultTagName = "properties"

// OpFunc custom option func
type OpFunc func(opts *Options)

// Options for config
type Options struct {
	// ParseEnv parse ENV var name, default True. eg: "${SHELL}"
	ParseEnv bool
	// ParseVar reference. eg: "${other.var.name}". default: true
	ParseVar bool
	// ParseTime string on binding struct. eg: 3s -> 3*time.Second
	ParseTime bool
	// TagName for binding data to struct. default: properties
	TagName string

	// MapStructConfig for binding data to struct.
	MapStructConfig mapstructure.DecoderConfig

	// InlineComment support split inline comments. default: false
	//
	// allow chars: #, //
	InlineComment bool
	// InlineSlice support parse the inline slice. eg: [23, 34]. default: false
	InlineSlice bool
	// TrimMultiLine trim "\n" for multi line value. default: false
	TrimMultiLine bool
	// BeforeCollect value handle func.
	BeforeCollect func(name, value string) (val interface{}, ok bool)
}

func (opts *Options) shouldAddHookFunc() bool {
	if opts.MapStructConfig.DecodeHook == nil {
		return opts.ParseTime || opts.ParseEnv
	}
	return false
}

func (opts *Options) makeDecoderConfig() *mapstructure.DecoderConfig {
	decConf := opts.MapStructConfig
	// compatible with settings on opts.TagName
	if decConf.TagName == "" {
		decConf.TagName = opts.TagName
	}

	// add hook on decode value to struct
	if opts.shouldAddHookFunc() {
		decConf.DecodeHook = ValDecodeHookFunc(opts.ParseEnv, opts.ParseTime)
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

// WithTagName custom tag name on binding struct
func WithTagName(tagName string) OpFunc {
	return func(opts *Options) {
		opts.TagName = tagName
		opts.MapStructConfig.TagName = tagName
	}
}
