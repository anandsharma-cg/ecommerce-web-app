package cart

import (
	"fmt"
	"net/http"
)

// CartService handles business logic for product-related operations.
type CartService struct {
	Repo *CartRepository
}

// NewCartService creates a new CartService.
func NewCartService(repo *CartRepository) *CartService {
	return &CartService{
		Repo: repo,
	}
}

// AddOrUpdateCartItem adds a product to the cart or updates the quantity if it already exists in the cart.
func (s *CartService) addOrUpdateCartItemService(cartID, productID, quantity int) (int, error) {
	// Ensure the quantity is greater than zero
	if quantity <= 0 {
		return http.StatusBadRequest, fmt.Errorf("invalid quantity: must be greater than zero")
	}

	// Call the repository to perform the upsert
	err := s.Repo.addOrUpdateCartItem(cartID, productID, quantity)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to add or update cart item: %v", err)
	}

	return http.StatusOK, nil
}
