package link_migrate

import (
	clickModel "url_shortener_pro/src/services/link_service/click/model"
	linkModel "url_shortener_pro/src/services/link_service/link/model"
	analysticsModel "url_shortener_pro/src/services/link_service/analystics/model"
)

func Models() []any {
	return []any{
		linkModel.Link{},
		clickModel.Click{},
		analysticsModel.LinkStats{},
	}
}
