package controller

import "github.com/gin-gonic/gin"

type EventControolerType struct{}

var EventController *EventControolerType

func init() {
	EventController = new(EventControolerType)
}

func (c EventControolerType) Create(ctx *gin.Context) {

}
