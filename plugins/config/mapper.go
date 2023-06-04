package config

import (
	rscliuitkit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/composit-items/multiselect"
	"github.com/Red-Sock/rscli-uikit/utils/common"

	shared_ui "github.com/Red-Sock/rscli/internal/shared-ui"
	config "github.com/Red-Sock/rscli/plugins/config/pkg/const"
)

const (
	// main config menu Items
	TransportTypeMenu = "transport type"
	DataSourceMenu    = "data source"
	CommitConfig      = "done"

	// destination menu Items
	pgCon    = "pg connection"
	redisCon = "redis connection"

	// transport layer menu Items
	restHttpType = "HTTP/rest"
	grpcType     = "grpc"
	telegramType = "telegram bot"
)

var mapToConfig = map[string][]string{
	restHttpType: {"-" + config.RESTHTTPServer},
	grpcType:     {"-" + config.GRPCServer},
	telegramType: {"-" + config.TelegramServer},
	pgCon:        {"-" + config.SourceNamePostgres},
	redisCon:     {"-" + config.SourceNameRedis},
}

func MainMenuItems() []string {
	return []string{
		TransportTypeMenu,
		DataSourceMenu,
		CommitConfig,
	}
}

func TransportTypeItems() []string {
	return []string{
		restHttpType,
		grpcType,
		telegramType,
	}

}

func DataSourcesItems() []string {
	return []string{
		pgCon,
		redisCon,
	}
}

type ConfigMenuSubItem struct {
	Items []string
	Flags []string

	PrevMenu func() rscliuitkit.UIElement
}

func NewConfigMenuSubItem(items []string, prevMenu func() rscliuitkit.UIElement) *ConfigMenuSubItem {
	return &ConfigMenuSubItem{
		Items:    items,
		PrevMenu: prevMenu,
	}
}

func (c *ConfigMenuSubItem) UiElement() rscliuitkit.UIElement {
	checked := make([]int, 0, 1)

	for idx, item := range c.Items {
		for _, selectedItem := range c.Flags {
			if item == selectedItem {
				checked = append(checked, idx)
				break
			}
		}
	}

	return multiselect.New(
		c.handleResponse,
		multiselect.HeaderLabel(shared_ui.GetHeaderFromText("Data sources")),
		multiselect.Items(c.Items...),
		multiselect.SeparatorChecked([]rune{'x'}),
		multiselect.Checked(checked),
		multiselect.Position(common.NewRelativePositioning(common.NewFillSpacePositioning(), common.NewFillSpacePositioning(), 0.4, 0.4)),
	)
}

func (c *ConfigMenuSubItem) handleResponse(args []string) rscliuitkit.UIElement {
	c.Flags = args
	return c.PrevMenu()
}

func (c *ConfigMenuSubItem) BuildFlagsForConfig() (res []string) {
	for _, item := range c.Flags {
		res = append(res, mapToConfig[item]...)
	}
	return res
}
