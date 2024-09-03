package all

import (
	// activate init-config command
	_ "github.com/ncarlier/kcusers/cmd/init-config"
	// activate get-user command
	_ "github.com/ncarlier/kcusers/cmd/users/get"
	// activate count-users command
	_ "github.com/ncarlier/kcusers/cmd/users/count"
	// activate delete-users command
	_ "github.com/ncarlier/kcusers/cmd/users/delete"
	// activate count-sessions command
	_ "github.com/ncarlier/kcusers/cmd/sessions/count"
	// activate version command
	_ "github.com/ncarlier/kcusers/cmd/version"
)
