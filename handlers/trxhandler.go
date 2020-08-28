package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zdarovich/win-lose-api/database"
	"github.com/zdarovich/win-lose-api/models"
	"log"
	"strconv"
	"sync"

	"net/http"
)

type (
	TrxHandler struct {
		lock sync.RWMutex
	}
	TrxRequest struct {
		State    string `json:"state"`
		Amount    string `json:"amount"`
		TransactionId string `json:"transactionId"`
	}
	TrxResponse struct {
		Data interface{}
	}
)

// POST method handler. It saves transaction and updates user balance
func (r *TrxHandler) PostTransaction(c *gin.Context) {
	var req TrxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.State) == 0 {
		log.Println("state is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "state is empty"})
		return
	} else if req.State != "win"  && req.State != "lost" {
		log.Println("state is invalid")
		c.JSON(http.StatusBadRequest, gin.H{"error": "state is invalid"})
		return
	}else if req.Amount == "" {
		log.Println("amount is invalid")
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount is invalid"})
		return
	} else if req.TransactionId == "" {
		log.Println("transactionId is invalid")
		c.JSON(http.StatusBadRequest, gin.H{"error": "transactionId is invalid"})
		return
	}

	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.GetDB().Begin()
	var trx models.Transaction

	if !tx.First(&trx, "transaction_id = ?", req.TransactionId).RecordNotFound() {
		log.Println("duplicate transaction")
		c.JSON(http.StatusBadRequest, gin.H{"error": "duplicate transaction"})
		return
	}
	trx = models.Transaction{
		State:         req.State,
		Amount:        req.Amount,
		TransactionId: req.TransactionId,
		Canceled:      false,
	}
	if err := tx.Save(&trx).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	tx = database.GetDB().Begin()
	var user models.User
	if err := tx.First(&user).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedBalance float64
	if req.State == "lost" {
		updatedBalance = user.Balance - amount
		if updatedBalance < 0 {
			log.Println("user balance cannot be negative")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user balance cannot be negative"})
			return
		}
	} else {
		updatedBalance = user.Balance + amount
	}


	if err := tx.Model(&user).Update("balance", updatedBalance).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	log.Println("user balance is updated")
	log.Printf("%+v \n", user)

	tx = database.GetDB().Begin()
	if tx.First(&trx,  "transaction_id = ?", req.TransactionId).RecordNotFound() {
		log.Println("transaction not found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction not found"})
		return
	}

	if err := tx.Model(trx).Update("user", &user).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	log.Println("transaction is saved")
	log.Printf("%+v \n", trx)

	resp := new(TrxResponse)
	resp.Data = trx
	c.JSON(http.StatusCreated, resp)
}

