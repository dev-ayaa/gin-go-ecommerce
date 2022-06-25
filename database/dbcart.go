package database

import (
	_"context"
	"errors"

	_"github.com/yusuf/gin-gonic-ecommerce/controller"
)

//Defining custom errors
var(
	ErrCantFindProduct = errors.New("cannot find the product")
	ErrCantDecodeProducts  = errors.New("cannot find or decode the product")
	ErrUserIdIsNotValid = errors.New("this user is not valid")
	ErrCantUpdateUser = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem = errors.New("unable to get the item from the cart")
	ErrCantBuyCartItem = errors.New("cannot update the item purchase presently")
)
// var app controller.Application
// var ctx context

func AddProductToCart() error{

}


func RemoveCartItem()  {
	
}

func GetItemFromCart(){

}

func BuyItemFromCart()  {
	
}


func InstantBuy(){

}