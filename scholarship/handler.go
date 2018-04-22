package scholarship

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/harkce/beasiswakita"
	"github.com/harkce/beasiswakita/errors"
	"github.com/harkce/beasiswakita/server/response"
	"github.com/harkce/beasiswakita/utils"

	"github.com/harkce/beasiswakita/authentication"
	"github.com/julienschmidt/httprouter"
)

type ScholarshipHandler struct{}

func (h *ScholarshipHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	owner, err := authentication.Token(r.Header.Get("Authorization"), []string{"organization"})
	if err != nil {
		response.Error(w, errors.Unauthorized)
		log.Println(err)
		return
	}

	var s beasiswakita.Scholarship
	err = utils.Decode(r, &s)
	if err != nil {
		response.Error(w, errors.InternalServerError)
		log.Println(err)
		return
	}

	s.OrganizationID = owner.ProfileID

	beasiswakita.Transaction, err = beasiswakita.DbMap.Begin()
	if err != nil {
		beasiswakita.Transaction.Rollback()
		response.Error(w, errors.InternalServerError)
		log.Println(err)
		return
	}

	newScholarship, err := CreateScholarship(s)
	if err != nil {
		beasiswakita.Transaction.Rollback()
		response.Error(w, errors.InternalServerError)
		log.Println(err)
		return
	}

	err = beasiswakita.Transaction.Commit()
	if err != nil {
		beasiswakita.Transaction.Rollback()
		response.Error(w, errors.InternalServerError)
		log.Println(err)
		return
	}

	response.Created(w, newScholarship)
	return
}

func (h *ScholarshipHandler) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filter := make(map[string]string)
	if sort, ok := r.URL.Query()["sort"]; ok {
		filter["sort"] = sort[0]
	} else {
		filter["sort"] = ""
	}
	if keywords, ok := r.URL.Query()["keywords"]; ok {
		filter["keywords"] = keywords[0]
	} else {
		filter["keywords"] = ""
	}
	if startDate, ok := r.URL.Query()["start_date"]; ok {
		filter["start_date"] = startDate[0]
	} else {
		filter["start_date"] = ""
	}
	if endDate, ok := r.URL.Query()["end_date"]; ok {
		filter["end_date"] = endDate[0]
	} else {
		filter["end_date"] = ""
	}
	if organizationID, ok := r.URL.Query()["organization_id"]; ok {
		filter["organization_id"] = organizationID[0]
	} else {
		filter["organization_id"] = ""
	}

	if limit, ok := r.URL.Query()["limit"]; ok {
		filter["limit"] = limit[0]
	} else {
		filter["limit"] = "10"
	}

	if filter["limit"] == "" {
		filter["limit"] = "10"
	}
	limit, err := strconv.Atoi(filter["limit"])
	if err != nil {
		response.Error(w, errors.UnprocessableEntity)
		log.Println(err)
		return
	}

	if offset, ok := r.URL.Query()["offset"]; ok {
		filter["offset"] = offset[0]
	} else {
		filter["offset"] = "0"
	}

	if filter["offset"] == "" {
		filter["offset"] = "0"
	}
	offset, err := strconv.Atoi(filter["offset"])
	if err != nil {
		response.Error(w, errors.UnprocessableEntity)
		log.Println(err)
		return
	}

	s, total, err := GetScholarships(filter, limit, offset)
	if err != nil {
		response.Error(w, errors.InternalServerError)
		log.Println(err)
		return
	}
	fmt.Println()
	meta := response.Meta{
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
	response.OKMeta(w, s, meta)
	return
}
