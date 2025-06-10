package infa

import "kahoot_bsu/internal/domain/models"

type StateProvider struct {	
}

func (s StateProvider) WaitLogin() models.State {
	return models.StateAwaitingLogin
}

func (s StateProvider) WaitOTP() models.State {
	return models.StateAwaitingOTP
}

func (s StateProvider) Registered() models.State {
	return models.StateRegistered
}