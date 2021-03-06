package component

import (
	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/vmware-tanzu/octant/pkg/action"
)

type ExtensionTab struct {
	Tab          Component      `json:"tab"`
	ClosePayload action.Payload `json:"payload,omitempty"`
}

func (e *ExtensionTab) UnmarshalJSON(data []byte) error {
	x := struct {
		Tab          TypedObject            `json:"tab"`
		ClosePayload map[string]interface{} `json:"payload,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	tab, err := x.Tab.ToComponent()
	if err != nil {
		return err
	}

	e.Tab = tab
	e.ClosePayload = x.ClosePayload

	return nil
}

type ExtensionConfig struct {
	Tabs []ExtensionTab `json:"tabs"`
}

// +octant:component
type Extension struct {
	Base

	Config ExtensionConfig `json:"config"`
}

func NewExtension() *Extension {
	return &Extension{
		Base: newBase(TypeExtension, TitleFromString("Extension")),
	}
}

func (e *Extension) AddTab(tab ExtensionTab) {
	e.Config.Tabs = append(e.Config.Tabs, tab)
}

type extensionMarshal Extension

func (e *Extension) MarshalJSON() ([]byte, error) {
	m := extensionMarshal(*e)
	m.Metadata.Type = TypeExtension
	return json.Marshal(&m)
}

func (e *Extension) GetMetadata() Metadata {
	return e.Metadata
}
