package webserver

import (
	"fmt"
	"les-randoms/database"
	"les-randoms/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAramRoute(c *gin.Context) {
	session := getSession(c)

	data := aramData{}

	setupNavData(&data.LayoutData.NavData, session)

	selectedItemName := setupSubnavData(&data.LayoutData.SubnavData, c, "Aram Gaming", []string{"Golden List", "Black List", "Bot List" /*, "Tier List"*/})

	setupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	data.ListTableData = customTableData{}
	list, err := database.List_SelectFirst("WHERE name = " + utils.Esc(selectedItemName))
	if err != nil {
		redirectToIndex(c)
		utils.LogError("Error while creating customTableData item with a DB List (Selected list : " + selectedItemName + ") : " + err.Error())
		return
	}
	data.ListTableData.HeaderList = list.Headers

	listItems, err := database.ListItem_SelectAll("WHERE listId = " + fmt.Sprint(list.Id) + " ORDER BY date")
	data.ListTableData.ItemList = make([]tableItemData, 0)
	for _, listItem := range listItems {
		data.ListTableData.ItemList = append(data.ListTableData.ItemList, tableItemData{FieldList: append([]string{listItem.Date.Local().Format(utils.DateFormat)}, utils.ParseDatabaseStringList(listItem.Value)...)})
	}

	c.HTML(http.StatusOK, "aram.tmpl.html", data)
}
