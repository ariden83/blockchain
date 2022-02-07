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
			Identifier: "about",
			Name:       "about",
			Title:      "this is the about-us page section",
			URL:        "/about",
		},
		{
			Identifier: "contact",
			Name:       "contact",
			Title:      "this is the contact page section",
			URL:        "/contact",
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
		{
			Identifier: "contact",
			Name:       "contact",
			Title:      "this is the contact page section",
			URL:        "/contact",
		},
	}
}
