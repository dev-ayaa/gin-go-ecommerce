package database

import "errors"

//Defining custom errors
var(
	ErrCantFindProduct = errors.New("can't find the product")
	ErrCantDecodeProducts  = errors.New("can't find or decode the product")
	ErrUserIdIsNotValid = errors.New("this user is not valid")
	ErrCantUpdateUser = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem = errors.New("unable to get the item from the cart")
	ErrCantBuyCartItem = errors.New("cannot update the item purchase presently")
)


func AddProductToCart(){

}


func RemoveCartItem()  {
	
}

func GetItemFromCart(){

}

func BuyItemFromCart()  {
	
}


func InstantBuy(){

}