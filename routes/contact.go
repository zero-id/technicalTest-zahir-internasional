package routes

import (
  "contact/handlers"
  "contact/pkg/mysql"
  "contact/repositories"

  "github.com/labstack/echo/v4"
)

func ContactRoutes(e *echo.Group) {
  contactRepository := repositories.RepositoryContact(mysql.DB)
  h := handlers.HandlerContact(contactRepository)

  e.GET("/contacts", h.FindContacts)
  e.GET("/contact/:id", h.GetContact)
  e.POST("/contact", h.CreateContact)
  e.PATCH("/contact/:id", h.UpdateContact)
  e.DELETE("/contact/:id", h.DeleteContact)
}