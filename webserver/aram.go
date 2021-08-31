package webserver

import (
	"les-randoms/database"
	"les-randoms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAramRoute(c *gin.Context) {
	session := getSession(c)

	if isNotAuthentified(session) {
		redirectToAuth(c)
	}

	data := aramData{}

	data.LayoutData.NavData = newNavData(session)

	data.LayoutData.SubnavData.Title = "Aram Gaming"
	data.LayoutData.SubnavData.SubnavItems = append(data.LayoutData.SubnavData.SubnavItems, subnavItem{Name: "Golden List"})
	data.LayoutData.SubnavData.SubnavItems = append(data.LayoutData.SubnavData.SubnavItems, subnavItem{Name: "Black List"})
	data.LayoutData.SubnavData.SubnavItems = append(data.LayoutData.SubnavData.SubnavItems, subnavItem{Name: "Bot List"})
	//data.LayoutData.SubnavData.SubnavItems = append(data.LayoutData.SubnavData.SubnavItems, subnavItem{Name: "Tier List"})
	data.LayoutData.SubnavData.SelectedSubnavItemIndex = 0
	if c.Request.Method == "POST" {
		selectedItemName := c.PostForm("subnavSelectedItem")
		for i := 0; i < len(data.LayoutData.SubnavData.SubnavItems); i++ {
			if selectedItemName == data.LayoutData.SubnavData.SubnavItems[i].Name {
				data.LayoutData.SubnavData.SelectedSubnavItemIndex = i
				break
			}
		}

		data.ListTableData = customTableData{}
		List, err := database.List_SelectFirst("WHERE name = " + utils.Esc(selectedItemName))
		if err != nil {
			redirectToIndex(c)
			utils.LogError("Error while creating customTableData item with a DB List (Selected list : " + selectedItemName + ") : " + err.Error())
			return
		}
		data.ListTableData.HeaderList = List.Headers
	}

	c.HTML(http.StatusOK, "aram.tmpl.html", data)
}
