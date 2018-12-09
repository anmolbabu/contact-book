package handlers

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/anmolbabu/contact-book/models"
	"github.com/anmolbabu/contact-book/utils"
)

type ListReq struct {
	Contact   models.Contact
	PageLimit int
	PageNo    int
}

func GetDefaultListReq() ListReq {
	return ListReq{
		Contact: models.Contact{
			Name:    utils.INVALID_STRING,
			EmailID: utils.INVALID_STRING,
		},
		PageLimit: utils.INVALID_INDEX,
		PageNo:    utils.INVALID_INDEX,
	}
}

func (lr ListReq) Serialise(params url.Values) (ListReq, error) {
	if params == nil {
		return lr, nil
	}
	if name, ok := params["name"]; ok && len(name) == 1 {
		lr.Contact.Name = name[0]
	}
	if emailid, ok := params["emailid"]; ok && len(emailid) == 1 {
		lr.Contact.EmailID = emailid[0]
	}
	if pagelimit, ok := params["pagelimit"]; ok && len(pagelimit) == 1 {
		pLimit, err := strconv.Atoi(pagelimit[0])
		if err != nil {
			err = fmt.Errorf("Invalid pagelimit %+v passed", pagelimit[0])
			return lr, err
		}
		lr.PageLimit = pLimit
	}
	if pageno, ok := params["page"]; ok && len(pageno) == 1 {
		pNo, err := strconv.Atoi(pageno[0])
		if err != nil {
			err = fmt.Errorf("Invalid pagelimit %+v passed", pageno[0])
			return lr, err
		}
		lr.PageNo = pNo
	}
	if lr.PageLimit == utils.INVALID_INDEX && lr.PageNo != utils.INVALID_INDEX {
		lr.PageLimit = utils.DEFAULT_PAGE_LIMIT
	}
	return lr, nil
}
