package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/horizon/core/services/infra/api-gateway/internal/errors"
)

func CreateUser(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "User service not wired") }
func ListUsers(c *gin.Context)    { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func GetUser(c *gin.Context)      { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ActivateUser(c *gin.Context) { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func SuspendUser(c *gin.Context)  { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ArchiveUser(c *gin.Context)  { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }

func CreateInstitution(c *gin.Context)    { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ListInstitutions(c *gin.Context)     { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func GetInstitution(c *gin.Context)       { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }

func CreateEvent(c *gin.Context)    { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ListEvents(c *gin.Context)     { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func GetEvent(c *gin.Context)       { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ConfirmEvent(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func PostEvent(c *gin.Context)      { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ReverseEvent(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ArchiveEvent(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }

func CreateGoal(c *gin.Context)    { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ListGoals(c *gin.Context)     { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func GetGoal(c *gin.Context)       { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ActivateGoal(c *gin.Context)  { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func CompleteGoal(c *gin.Context)  { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ArchiveGoal(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }

func CreateAccount(c *gin.Context)    { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ListAccounts(c *gin.Context)     { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func GetAccount(c *gin.Context)       { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ActivateAccount(c *gin.Context)  { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func FreezeAccount(c *gin.Context)    { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func CloseAccount(c *gin.Context)     { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }

func CreateAllocation(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ListAllocations(c *gin.Context)    { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func GetAllocation(c *gin.Context)      { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }

func CreateAsset(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ListAssets(c *gin.Context)    { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func GetAsset(c *gin.Context)      { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }

func CreateLiability(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ListLiabilities(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func GetLiability(c *gin.Context)      { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func SettleLiability(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }

func CreatePortfolio(c *gin.Context)   { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func ListPortfolios(c *gin.Context)    { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
func GetPortfolio(c *gin.Context)      { errors.Respond(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "", "") }
