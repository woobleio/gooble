package models

// AS EXEMPLE MODEL BINDING FOR BODY

type DOM struct {
  Message string `json:"message" binding:"required"`
}