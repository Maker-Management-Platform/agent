package projects

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/eduardooliveira/stLib/core/data/database"
	"github.com/eduardooliveira/stLib/core/entities"
	"github.com/labstack/echo/v4"
	"github.com/morkid/paginate"
)

func index(c echo.Context) error {
	config := paginate.Config{
		FieldSelectorEnabled: true,
	}
	pg := paginate.New(config)

	q := database.DB.Model(&entities.Project{}).Preload("Tags")

	if c.QueryParams().Has("name") {
		q.Where("name LIKE ?", fmt.Sprintf("%%%s%%", c.QueryParam("name")))
	}
	if c.QueryParams().Has("tags") {
		for i, t := range strings.Split(c.QueryParam("tags"), ",") {
			q.Joins(fmt.Sprintf("LEFT JOIN project_tags as project_tags%d on project_tags%d.project_uuid = projects.uuid", i, i)).
				Where(fmt.Sprintf("project_tags%d.tag_value = ?", i), t)
		}

	}
	q.Order("name ASC")
	page := pg.With(q).Request(c.Request()).Response(&[]entities.Project{})
	if page.RawError != nil {
		log.Println(page.RawError)
		return echo.NewHTTPError(http.StatusInternalServerError, page.RawError.Error())
	}
	return c.JSON(http.StatusOK, page)
}
