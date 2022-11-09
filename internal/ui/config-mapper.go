package ui

import (
	rscliuitkit "github.com/Red-Sock/rscli-uikit"
	"github.com/Red-Sock/rscli-uikit/composit-items/multiselect"
	"github.com/Red-Sock/rscli/pkg/service/config"
	"github.com/Red-Sock/rscli/pkg/service/help"
)

const (
	transportTypeMenu = "transport type"
	dataSourceMenu    = "data source"
	commitConfig      = "done"

	pgCon    = "pg connection"
	redisCon = "redis connection"

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

func mainMenuItems() []string {
	return []string{
		transportTypeMenu,
		dataSourceMenu,
		commitConfig,
	}
}

func transportTypeItems() []string {
	return []string{
		restHttpType,
		grpcType,
	}

}

func dataSourcesItems() []string {
	return []string{
		pgCon,
		redisCon,
	}
}

type configMenuSubItem struct {
	items []string
	flags []string

	prevMenu func() rscliuitkit.UIElement
}

func newConfigMenuSubItem(items []string, prevMenu func() rscliuitkit.UIElement) *configMenuSubItem {
	return &configMenuSubItem{
		items:    items,
		prevMenu: prevMenu,
	}
}

func (c *configMenuSubItem) uiElement() rscliuitkit.UIElement {
	checked := make([]int, 0, 1)

	for idx, item := range c.items {
		for _, selectedItem := range c.flags {
			if item == selectedItem {
				checked = append(checked, idx)
				break
			}
		}
	}

	return multiselect.New(
		c.handleResponse,
		multiselect.Header(help.Header+"DataSources"),
		multiselect.Items(c.items...),
		multiselect.SeparatorChecked([]rune{'x'}),
		multiselect.Checked(checked),
	)
}

func (c *configMenuSubItem) handleResponse(args []string) rscliuitkit.UIElement {
	c.flags = args
	return c.prevMenu()
}

func (c *configMenuSubItem) buildFlagsForConfig() (res []string) {
	for _, item := range c.flags {
		res = append(res, mapToConfig(item)...)
	}
	return res
}
