package controller

import (
	"expensez/pkg/utils"
	v "expensez/pkg/validator"
	"expensez/src/models"
	"expensez/src/repository"
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/labstack/echo"
)

var (
	TagCtrl Controller
)

type TagController struct {
	repo      repository.Repository
	validator v.CustomValidator
}

func TagLoad() error {
	TagCtrl = &TagController{repo: repository.Repo, validator: *v.Validator}
	return nil
}

func (t *TagController) Create(c echo.Context) error {
	modal := &models.Tag{}
	if err := c.Bind(modal); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := t.repo.Create(c, modal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(modal)
}

func (t *TagController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = t.repo.Delete(c, models.Tag{}, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "Deleted")
}

func (t *TagController) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := t.repo.Get(c, models.Tag{ID: uint(id)})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

func (t *TagController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	modl := &models.Tag{}
	if err := c.Bind(modl); err != nil {
		return err
	}
	modl.ID = uint(id)
	err = t.repo.Update(c, modl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, modl)
}

func (t *TagController) List(c echo.Context) error {
	parentTag := c.QueryParam("parentTag")
	list := &[]models.FullTag{}
	where := ""
	if utils.IsValueNonZero(parentTag) {
		where = parentTag + " IN (t.id, t.tag_id)"
	}
	table := "tags t"
	silect := `t.id, t.name, t.tag_id, p.name AS parent, t.transaction_type_id, 
	tt.name AS type, COUNT(DISTINCT(m.id)) AS tag_count`
	joins := `LEFT JOIN tags p ON t.tag_id = p.id 
	LEFT JOIN tags c ON t.id = c.tag_id 
	LEFT JOIN transaction_types tt ON t.transaction_type_id = tt.id 
	LEFT JOIN activities m ON t.id in (m.tag_id, m.sub_tag_id) AND m.event_date > DATE_SUB(now(), INTERVAL 6 MONTH)`
	groupBy := "t.id, t.name, t.tag_id, p.name, t.transaction_type_id, tt.name"
	orderBy := "COUNT(DISTINCT(c.id)) DESC, t.name ASC"
	err := t.repo.FetchWithFullQuery(c, list, table, silect, joins, where, groupBy, orderBy, 0, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}
