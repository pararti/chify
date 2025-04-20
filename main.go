package main

import (
	"embed"
	"log"
	"pararti/chify/internal/registry"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

//go:embed translation
var translations embed.FS

//go:embed resources/chify.png
var icon []byte

func main() {
	a := app.New()
	resourceIcon := fyne.NewStaticResource("chify.png", icon)
	a.SetIcon(resourceIcon)

	window := a.NewWindow("Chify")
	window.Resize(fyne.NewSize(800, 640))
	window.SetFixedSize(false)

	// Load langs
	err := lang.AddTranslationsFS(translations, "translation")
	if err != nil {
		log.Fatal(err)
	}

	// Use DocTabs for closable tabs
	tabs := container.NewDocTabs()

	// Add initial tab
	initialTab := container.NewTabItem("aes", container.NewScroll(registry.DefaultService.BuildForm()))
	tabs.Append(initialTab)
	tabs.Select(initialTab)

	tabs.SetTabLocation(container.TabLocationTop)

	accordingItems := make([]*widget.AccordionItem, 0, len(registry.LeftServiceMenu))
	for _, menuEl := range registry.LeftServiceMenu {
		accordingItems = append(accordingItems, widget.NewAccordionItem(menuEl.Category,
			func() fyne.CanvasObject {
				buttonsBox := container.NewVBox()
				for _, subMenuEl := range menuEl.Elements {
					button := widget.NewButton(subMenuEl.Name, func() {
						form := subMenuEl.Service.BuildForm()
						// Update the *selected* tab's content and title
						if selectedTab := tabs.Selected(); selectedTab != nil {
							selectedTab.Text = subMenuEl.Name
							selectedTab.Content = container.NewScroll(form)
							tabs.Refresh()
						}
					})
					buttonsBox.Add(button)
				}
				return buttonsBox
			}()))
	}

	accordion := widget.NewAccordion(accordingItems...)
	accordion.Open(0)
	accordionContainer := container.NewVBox(accordion)

	// Add new tab button again with default aes mode
	newTabButton := widget.NewButton("+", func() {
		tabContent := container.NewScroll(registry.DefaultService.BuildForm())
		newTab := container.NewTabItem("aes", tabContent)
		tabs.Append(newTab)
		tabs.Select(newTab)
	})

	tabHeader := container.NewBorder(nil, nil, nil, newTabButton, tabs)

	split := container.NewHSplit(accordionContainer, tabHeader)
	split.Offset = 0.1 // Set left panel to take up 10% of space

	window.SetContent(split)
	window.ShowAndRun()
}
