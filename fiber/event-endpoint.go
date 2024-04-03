package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func EventEndpointGet(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		option := c.Query("option", "")
		if option == "default" {
			skip, err := strconv.ParseInt(c.Query("skip", "0"), 10, 64)
			if err != nil {
				log.Println(err.Error())
				return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
			}
			take, err := strconv.Atoi(c.Query("take", "10"))
			if err != nil {
				log.Println(err.Error())
				return c.Status(400).JSON(fiber.Map{"message": "参数错误"})
			}
			option := RetrieveOption{
				Skip: skip,
				Take: take,
			}
			var filter RetrieveFilter
			query := c.Queries()
			if query["equal"] != "" {
				filter.Equal = strings.Split(query["equal"], ",")
			}
			if query["object-contain"] != "" {
				filter.ObjectContain = strings.Split(query["object-contain"], ",")
			}
			if query["array-contain"] != "" {
				filter.ArrayContain = strings.Split(query["array-contain"], ",")
			}
			if query["like"] != "" {
				filter.Like = strings.Split(query["like"], ",")
			}
			if query["object-like"] != "" {
				filter.ObjectLike = strings.Split(query["object-like"], ",")
			}
			if query["in"] != "" {
				filter.In = strings.Split(query["in"], ",")
			}
			if query["lesser"] != "" {
				filter.Lesser = strings.Split(query["lesser"], ",")
			}
			if query["greater"] != "" {
				filter.Greater = strings.Split(query["greater"], ",")
			}
			result, err := EventDefaultFilter(option, filter)
			if err != nil {
				slogger.Error(err.Error())
				return c.Status(500).JSON(fiber.Map{"message": "服务器错误"})
			}
			if len(result) == 0 {
				return c.SendString("[]")
			}
			var response []EventExtended
			for _, it := range result {
				extendedEvent := EventExtended{
					Id:           it.Id,
					RelationId:   it.RelationId,
					ReferenceId:  it.ReferenceId,
					Tags:         it.Tags,
					Detail:       it.Detail,
					Time:         it.Time,
					Id_:          strconv.FormatInt(it.Id, 10),
					RelationId_:  strconv.FormatInt(it.RelationId, 10),
					ReferenceId_: strconv.FormatInt(it.ReferenceId, 10),
				}
				response = append(response, extendedEvent)
			}
			c.Set(HEADER_API_VERSION, "2024-01-06")
			return c.JSON(response)
		}
	}
	if id > 0 {
		return c.SendStatus(200)
	}
	return c.Status(200).SendString("")
}
