package blockeditor

import (
	"github.com/gouniverse/hb"
	"github.com/gouniverse/ui"
	"github.com/samber/lo"
)

// blockToCard creates a card for a block
func (b *editor) blockToCard(block ui.BlockInterface) *hb.Tag {
	buttonMoveUp := b.cardButtonMoveUp(block.ID())
	buttonMoveDown := b.cardButtonMoveDown(block.ID())
	buttonEdit := b.cardButtonSettings(block.ID())
	buttonDelete := b.cardButtonDelete(block.ID())
	buttonDropdown := b.cardButtonDropdown(block)

	definition := b.findDefinitionByType(block.Type())

	hasRenderer := false

	if definition != nil {
		hasRenderer = definition.ToHTML != nil
	}

	render := lo.IfF(hasRenderer, func() string {
		return definition.ToHTML(block)
	}).ElseF(func() string {
		return hb.NewTag("center").
			Child(definition.Icon).
			Style("font-size: 40px;").ToHTML()
	})

	card := hb.Div().
		Class(`BlockCard card`).
		Child(
			hb.Div().
				Class(`card-header bg-info`).
				Style(`--bs-bg-opacity: 0.2;`).
				Style(`padding: 2px 10px;font-size: 11px;`).
				Child(buttonDropdown).
				Text(block.Type()).
				Child(buttonDelete).
				Child(buttonEdit).
				Child(buttonMoveUp).
				Child(buttonMoveDown),
		).
		Child(hb.Div().
			Class(`card-body bg-info`).
			ClassIf(block.Type() == "row", `row`).
			Style(`--bs-bg-opacity: 0.1;`).
			// ChildIf(len(block.Children()) < 1, b.blockDivider().Child(b.buttonBlockInsert(blockExt.ID, 0, false))).
			ChildrenIfF(len(block.Children()) > 0, func() []hb.TagInterface {
				return lo.Map(block.Children(), func(child ui.BlockInterface, position int) hb.TagInterface {
					return hb.Wrap().
						// Child(b.blockDivider().Child(b.buttonBlockInsert(blockExt.ID, position, false))).
						Child(b.blockToCard(child))
				})
			}).
			// ChildIf(len(block.Children()) > 0, b.blockDivider().Child(b.buttonBlockInsert(blockExt.ID, len(block.Children()), false))).
			HTMLIf(len(block.Children()) < 1, render))

	if block.Type() == "column" {
		width := block.Parameter("width")

		if width == "" {
			width = "12"
		}

		return hb.Div().
			Class("col-" + width).
			Child(card)
	}

	return card

}
