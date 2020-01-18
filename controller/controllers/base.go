package controllers

import (
	networkv1 "github.com/firemiles/bifrost/controller/api/v1"
)

var (
	ownerKey = ".metadata.controller"
	apiGVStr = networkv1.GroupVersion.String()
)
