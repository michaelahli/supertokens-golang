package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/constants"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

func MakeAPIImplementation() models.APIImplementation {
	return models.APIImplementation{
		EmailExistsGET: func(email string, options models.APIOptions) models.EmailExistsGETResponse {
			user := options.RecipeImplementation.GetUserByEmail(email)
			return models.EmailExistsGETResponse{
				Ok:    true,
				Exist: user != nil,
			}
		},

		GeneratePasswordResetTokenPOST: func(formFields []models.FormFieldValue, options models.APIOptions) models.GeneratePasswordResetTokenPOSTResponse {
			var email string
			for _, formField := range formFields {
				if formField.ID == constants.FormFieldEmailID {
					email = formField.Value
				}
			}
			user := options.RecipeImplementation.GetUserByEmail(email)
			if user == nil {
				return models.GeneratePasswordResetTokenPOSTResponse{
					Ok: true,
				}
			}
			response := options.RecipeImplementation.CreateResetPasswordToken(user.ID)
			if response.Status == "UNKNOWN_USER_ID" {
				return models.GeneratePasswordResetTokenPOSTResponse{
					Ok: true,
				}
			}
			passwordResetLink := options.Config.ResetPasswordUsingTokenFeature.GetResetPasswordURL(*user) + "?token=" + response.Token + "&rid=" + options.RecipeID
			options.Config.ResetPasswordUsingTokenFeature.CreateAndSendCustomEmail(*user, passwordResetLink)

			return models.GeneratePasswordResetTokenPOSTResponse{
				Ok: true,
			}
		},

		PasswordResetPOST: func(formFields []models.FormFieldValue, token string, options models.APIOptions) models.PasswordResetPOSTResponse {
			var newPassword string
			for _, formField := range formFields {
				if formField.ID == constants.FormFieldPasswordID {
					newPassword = formField.Value
				}
			}
			response := options.RecipeImplementation.ResetPasswordUsingToken(token, newPassword)

			return models.PasswordResetPOSTResponse{
				Status: response.Status,
			}
		},

		SignInPOST: func(formFields []models.FormFieldValue, options models.APIOptions) models.SignInUpResponse {
			var email string
			for _, formField := range formFields {
				if formField.ID == constants.FormFieldEmailID {
					email = formField.Value
				}
			}
			var newPassword string
			for _, formField := range formFields {
				if formField.ID == constants.FormFieldPasswordID {
					newPassword = formField.Value
				}
			}

			response := options.RecipeImplementation.SignIn(email, newPassword)
			if response.Status == "WRONG_CREDENTIALS_ERROR" {
				return response
			}

			user := response.User
			jwtPayload := options.Config.SessionFeature.SetJwtPayload(user, formFields, "signin")
			sessionData := options.Config.SessionFeature.SetSessionData(user, formFields, "signin")

			_, err := session.CreateNewSession(options.Res, user.ID, jwtPayload, sessionData)
			if err != nil {
				return models.SignInUpResponse{
					Status: "WRONG_CREDENTIALS_ERROR",
				}
			}

			return models.SignInUpResponse{
				User:   user,
				Status: "OK",
			}
		},

		SignUpPOST: func(formFields []models.FormFieldValue, options models.APIOptions) models.SignInUpResponse {
			var email string
			for _, formField := range formFields {
				if formField.ID == constants.FormFieldEmailID {
					email = formField.Value
				}
			}
			var newPassword string
			for _, formField := range formFields {
				if formField.ID == constants.FormFieldPasswordID {
					newPassword = formField.Value
				}
			}

			response := options.RecipeImplementation.SignIn(email, newPassword)
			if response.Status == "EMAIL_ALREADY_EXISTS_ERROR" {
				return response
			}

			user := response.User
			jwtPayload := options.Config.SessionFeature.SetJwtPayload(user, formFields, "signin")
			sessionData := options.Config.SessionFeature.SetSessionData(user, formFields, "signin")

			_, err := session.CreateNewSession(options.Res, user.ID, jwtPayload, sessionData)
			if err != nil {
				return models.SignInUpResponse{
					Status: "EMAIL_ALREADY_EXISTS_ERROR",
				}
			}

			return models.SignInUpResponse{
				User:   user,
				Status: "OK",
			}
		},
	}
}
