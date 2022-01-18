package explorer

func getMenus() []Menus {
	return []Menus{
		{
			Identifier: "home",
			Name:       "home",
			Title:      "This is the home section",
			URL:        "/",
		},
		{
			Identifier: "login",
			Name:       "login",
			Title:      "this is the login section",
			URL:        "/login",
		},
	}
}

func getMenusAuth() []Menus {
	return []Menus{
		{
			Identifier: "home",
			Name:       "home",
			Title:      "This is the home section",
			URL:        "/",
		},
		{
			Identifier: "wallet",
			Name:       "wallet",
			Title:      "this is the wallet section",
			URL:        "/wallets",
		},
	}
}
