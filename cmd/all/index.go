package all

import (
	// activate init-config command
	_ "github.com/ncarlier/kcusers/cmd/init-config"
	// activate get-user command
	_ "github.com/ncarlier/kcusers/cmd/users/get"
	// activate count-users command
	_ "github.com/ncarlier/kcusers/cmd/users/count"
	// activate delet-users command
	_ "github.com/ncarlier/kcusers/cmd/users/delete"
)
