package handlers

import (
	"net/http"

	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/Alan-Luc/VertiLog/backend/services"
	"github.com/Alan-Luc/VertiLog/backend/utils/apiErrors"
	"github.com/Alan-Luc/VertiLog/backend/utils/auth"
	"github.com/gin-gonic/gin"
)

func UserRegisterHandler(ctx *gin.Context) {
	var user models.User
	var err error

	err = ctx.ShouldBindJSON(&user)
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid input. Please check the submitted data and try again.",
		err,
		http.StatusBadRequest,
	) {
		return
	}

	err = services.PrepareUser(&user)
	if apiErrors.HandleAPIError(
		ctx,
		"An error occurred while processing your request. Please try again later.",
		err,
		http.StatusInternalServerError,
	) {
		return
	}

	err = services.CreateUser(&user)
	if apiErrors.HandleAPIError(
		ctx,
		"We encountered an issue while registering your account. Please try again later.",
		err,
		http.StatusInternalServerError,
	) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user_id": user,
	})
}

func UserLoginHandler(ctx *gin.Context) {
	var loginInput models.User
	var err error

	err = ctx.ShouldBindJSON(&loginInput)
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid input. Please check the submitted data and try again.",
		err,
		http.StatusBadRequest,
	) {
		return
	}

	jwt, err := services.VerifyUser(loginInput.Username, loginInput.Password)
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid username or password. Please check your credentials and try again.",
		err,
		http.StatusUnauthorized,
	) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": jwt})
}

func UserProfileHandler(ctx *gin.Context) {
	var updateInput models.PasswordUpdateInput
	var userID int
	var err error

	err = ctx.ShouldBindJSON(&updateInput)
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid input. Please check the submitted data and try again.",
		err,
		http.StatusBadRequest,
	) {
		return
	}

	userID, err = auth.ExtractUserIdFromJWT(ctx)
	if apiErrors.HandleAPIError(
		ctx,
		"Authorization token is invalid or missing. Please log in and try again.",
		err,
		http.StatusUnauthorized,
	) {
		return
	}

	err = services.UpdateUserPassword(userID, updateInput.CurrentPassword, updateInput.NewPassword)
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid password. Please check your credentials and try again.",
		err,
		http.StatusUnauthorized,
	) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}
