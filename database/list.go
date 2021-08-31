package database

import (
	"database/sql"
	"fmt"
	"les-randoms/utils"
	"reflect"
)

type List struct {
	Id      int
	Name    string
	Headers []string
}

func (list List) ToStringSlice() []string {
	return []string{fmt.Sprint(list.Id), list.Name, utils.UnsliceStrings(list.Headers, " ")}
}

func Lists_ToDBStructSlice(lists []List) []DBStruct {
	var r []DBStruct
	for _, list := range lists {
		r = append(r, list)
	}
	return r
}

func List_GetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(&List{})).Type()
}

func List_SelectAll() ([]List, error) {
	rows, err := SelectDatabase("id, name, headers FROM List")
	if err != nil {
		return nil, err
	}
	lists := make([]List, 0)
	for rows.Next() {
		var id int
		var name string
		var headers string
		err = rows.Scan(&id, &name, &headers)
		if err != nil {
			return nil, err
		}
		lists = append(lists, List{Id: id, Name: name, Headers: utils.ParseDatabaseStringList(headers)})
	}
	return lists, nil
}

func List_Select(queryPart string) (List, error) {
	rows, err := SelectDatabase("id, name, headers FROM List " + queryPart)
	if err != nil {
		return List{}, err
	}
	rows.Next()
	var id int
	var name string
	var headers string
	err = rows.Scan(&id, &name, &headers)
	if err != nil {
		return List{}, err
	}
	return List{Id: id, Name: name, Headers: utils.ParseDatabaseStringList(headers)}, nil
}

func List_CreateNew(name string, headers string) (List, sql.Result, error) {
	result, err := InsertDatabase("List(name, headers) VALUES(" + utils.Esc(name) + ", " + utils.Esc(headers) + ")")
	if err != nil {
		return List{}, nil, err
	}
	newId, err := result.LastInsertId()
	if err != nil {
		return List{}, nil, err
	}
	return List{Id: int(newId), Name: name, Headers: utils.ParseDatabaseStringList(headers)}, result, err
}
