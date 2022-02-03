package database

import (
	"errors"
	"les-randoms/utils"
)

func VerifyDatabase() {
	utils.LogInfo("Starting to verify database validity..")

	testing := utils.CreateTesting("DATABASE VALIDITY TEST")

	testDatabasePing(&testing)

	tableStates := testDatabaseTables(&testing, []string{"User", "List", "ListItem", "AccessRight"})

	defaultAccessRightsValidity := testDefaultAccessRights(&testing)

	// Display tests conclusion
	err := testing.Conclusion()
	if err != nil {
		utils.LogError(err.Error() + " failed tests but no fatal tests failed")
		utils.LogInfo("Trying to fix errors..")
		fixTables(tableStates)
		fixDefaultAccessRightsValidity(defaultAccessRightsValidity)

		utils.LogSuccess("All fixes trials done")
		VerifyDatabase()
		return
	}
	utils.LogInfo("Program will continue..")
	return
}

func fixTables(tableStates map[string]int) {
	for key, value := range tableStates {
		if value != 0 {
			utils.LogInfo("Fixing " + key + " table..")
			if value == 2 { // table does not exists
				createTable(getSpecificTableCreationQuery(key))
			} else if value == 1 { // table exists but structure is not valid
				utils.HandlePanicError(errors.New("Fixing 'table exists but structure is not valid' error is not implemented yet. Sorry bro, GLHF !"))
			}
		}
	}
}

func fixDefaultAccessRightsValidity(b bool) {
	if !b {
		utils.LogInfo("Fixing Default AccessRights Validity..")
		DeleteDatabase("FROM AccessRight WHERE userId=1 AND path='/database'")
		AccessRight_CreateNew(1, "/database", RightTypes.Authorized)
	}
}

func testDatabasePing(testing *utils.Testing) {
	// Test to ping database
	testing.TestError(Database.Ping(), "Successful ping to database", "Failed to ping database", true)
}

func testDatabaseTables(testing *utils.Testing, tables []string) map[string]int {
	// Tests on tables validity

	tableStates := make(map[string]int)
	// 0 Means exists and valid
	// 1 Means exists but structure is not valid
	// 2 Means does not exists

	for _, t := range tables {
		tableStates[t] = 2
	}

	rows, _ := SelectDatabase("name, sql FROM sqlite_schema WHERE type IN ('table','view') AND name NOT LIKE 'sqlite_%' ORDER BY 1")
	var name string
	var sql string
	for rows.Next() {
		rows.Scan(&name, &sql)
		if testing.TestStringEqual(getSpecificTableCreationQuery(name), sql, name+" table exists and is valid", name+" table exists but is not valid", false) == nil {
			tableStates[name] = 0
		} else {
			tableStates[name] = 1
		}
	}
	for key, value := range tableStates {
		if value == 2 {
			testing.TestBool(true, false, "if you see this, there is a problem", key+" table does not exists", false)
		}
	}

	return tableStates
}

func testDefaultAccessRights(testing *utils.Testing) bool {
	right, err := AccessRight_SelectFirst("WHERE userId=1 AND path='/database'")
	b := false
	if err == nil {
		b = right.RightType == RightTypes.Authorized
	}

	testing.TestBool(true, b, "User of id 1 is authorized to access /database", "User of id 1 is forbidden from accessing /database", false)

	return b
}
