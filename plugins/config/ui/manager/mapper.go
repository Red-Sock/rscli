package manager

import (
	rscliuitkit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/composit-items/multiselect"
	config "github.com/Red-Sock/rscli/plugins/config/processor"
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
)

func mapToConfig(menuItem string) (args []string) {
	switch menuItem {
	case restHttpType:
		return []string{"-" + config.RESTHTTPServer}
	case grpcType:
		return []string{"-" + config.GRPCServer}

	case pgCon:
		return []string{"-" + config.SourceNamePg}
	case redisCon:
		return []string{"-" + config.SourceNameRds}
	default:
		return nil
	}
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
		multiselect.Header("DataSources"),
		multiselect.Items(c.Items...),
		multiselect.SeparatorChecked([]rune{'x'}),
		multiselect.Checked(checked),
	)
}

func (c *ConfigMenuSubItem) handleResponse(args []string) rscliuitkit.UIElement {
	c.Flags = args
	return c.PrevMenu()
}

func (c *ConfigMenuSubItem) BuildFlagsForConfig() (res []string) {
	for _, item := range c.Flags {
		res = append(res, mapToConfig(item)...)
	}
	return res
}
