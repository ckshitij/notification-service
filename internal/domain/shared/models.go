package shared

type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSlack Channel = "slack"
	ChannelInApp Channel = "in_app"
)

type TemplateType string

const (
	SystemTemplate TemplateType = "system"
	UserTemplate   TemplateType = "user"
)
