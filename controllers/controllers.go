// Package controllers R.H.T.(Rent Household Things)
//
// A simple application which allows users to rent household
// things easily at affordable prices !
//
//     Schemes: http
//     BasePath: /api
//     Version: 1.0.0
//     Host: localhost:8080
//     License: MIT license
//     Contact: ckar1604212@gmail.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"Newton/db"
	"Newton/helpers"
	"Newton/models"
	model "Newton/models"
	"Newton/query"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	otp     string
	testobj primitive.ObjectID
	teststr string
)

// import "Newton/models"

func Check(url string, method string, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/"+url {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != method {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
}

// swagger:route POST /account User Account
// Checks if the phone number exists in database or not.
// responses:
//   200: accountResponse

//
// swagger:response accountResponse
type accountResponseWrapper struct {
	// The generated response
	// in:body
	Body models.User
}

// consumes:
// - application/json
// swagger:parameters Account
type accountParamsWrapper struct {
	// in:body
	// type:application/json
	Body models.Account
}

// AccountHandler : Checks if the phone number/User exists in database or not
func AccountHandler(w http.ResponseWriter, r *http.Request) {

	Check("account", "POST", w, r)

	w.Header().Set("Content-Type", "application/json")
	var data model.Account
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &data)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)

	}
	var user model.User

	_ = query.FindoneID("user", data.ID, "_id").Decode(&user)
	match, err := regexp.MatchString("[0-9]{10}", user.Phone)
	//fmt.Println(match)
	if data.Exist == false && user.Phone == "" {
		res.Result = "Not registered"
		json.NewEncoder(w).Encode(res)
	} else if data.Exist == false && match {
		res.Result = "Login required"
		json.NewEncoder(w).Encode(res)
	} else if data.Exist == true {
		json.NewEncoder(w).Encode(user)

	}
}

//carousel

func Carousel(w http.ResponseWriter, r *http.Request) {

	Check("carousel", "GET", w, r)

	var picture model.Carousel

	picture.Carousel = [7]string{"https://rht007.s3.amazonaws.com/carousel/1.jpg", "https://rht007.s3.amazonaws.com/carousel/2.jpg", "https://rht007.s3.amazonaws.com/carousel/3.jpg", "https://rht007.s3.amazonaws.com/carousel/4.jpg", "https://rht007.s3.amazonaws.com/carousel/5.jpg", "https://rht007.s3.amazonaws.com/carousel/6.jpg", "https://rht007.s3.amazonaws.com/carousel/7.jpg"}

	json.NewEncoder(w).Encode(picture)

}

// swagger:operation GET /productslist/{location}/{sub}/{page} User Productlist
// ---
// summary: List of products according to location subcategory and page number
// description: All products in the given location Subcategory and page number
// parameters:
// - name: location
//   in: path
//   description: location id string
//   type: string
//   required: true
// - name: sub
//   in: path
//   description: subcategoryid string
//   type: string
//   required: true
// - name: page
//   in: path
//   description: pagenumber
//   type: string
//   required: true
// responses:
//   200: ProductlistResponse

// This text will appear as description of your response body.
// swagger:response ProductlistResponse
type ProductsListResponseWrapper struct {
	// The generated response
	// in:body
	Body models.List
}

// ProductsList : handles the Productslist
func ProductsList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var items model.Items
	params := mux.Vars(r)
	Check("productslist/"+params["location"]+"/"+params["sub"]+"/"+params["page"], "GET", w, r)
	//var list []model.Items
	var res model.ResponseResult
	LocID, err := primitive.ObjectIDFromHex(params["location"])
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	SubID, err1 := primitive.ObjectIDFromHex(params["sub"])
	if err1 != nil {
		res.Error = err1.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	s := params["page"]
	page, _ := strconv.Atoi(s)
	filter := bson.M{"locationid": LocID, "subcategoryid": SubID}

	cursor, totaldocuments := query.FindAll("products", filter, page)
	totalpage := (totaldocuments / 6)
	extra := totaldocuments % 6
	var response model.List
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		items.Img = nil
		items.Itemsid = nil
		if err = cursor.Decode(&items); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		//list = append(list, items)
		response.Products = append(response.Products, items)

	}
	response.Page = s

	if int64(page) == totalpage-1 && extra == 0 {
		response.IsLastpage = "true"
	} else if extra != 0 && int64(page) == totalpage {
		response.IsLastpage = "true"
	} else {
		response.IsLastpage = "false"
	}

	json.NewEncoder(w).Encode(response)

}

// swagger:operation GET /usercreation User UserCreation
// ---
// summary: Creates Userid and wishlist and cart for that Userid
// description: Creates Userid and wishlist and cart for that Userid and returns the new userid
// responses:
//   200: UserCreationResponse

// This text will appear as description of your response body.
// swagger:response UserCreationResponse
type UserCreationResponseWrapper struct {
	// The generated response
	// in:body
	Body []model.Id
}

// UserCreationHandler :  Creates Userid and wishlist and cart for that Userid and returns the new userid
func UserCreationHandler(w http.ResponseWriter, r *http.Request) {

	Check("usercreation", "GET", w, r)

	var user model.User

	var id model.Id
	result := query.InsertOne("user", user)

	oid, _ := result.InsertedID.(primitive.ObjectID)
	id.ID1 = oid

	var wish model.Wishlist
	wish.Userid = oid

	result1 := query.InsertOne("wishlist", wish)
	oidw, _ := result1.InsertedID.(primitive.ObjectID)

	fmt.Println(oidw)

	result2 := query.InsertOne("cart", wish)
	oidc, _ := result2.InsertedID.(primitive.ObjectID)

	fmt.Println(oidc)
	json.NewEncoder(w).Encode(id)
}

// swagger:route POST /wishlist User Wishlisting
// Adding Products into User Wishlist and Removing Products from user Wishlist.
// responses:
//   200: wishlistResponse

// Returning Bool True for Adding in wishlist and False for Removing From Wishlist
// swagger:response wishlistResponse
type wishlistResponseWrapper struct {
	// The generated response
	// in:body
	Body bool
}

// consumes:
// - application/json
// swagger:parameters Wishlisting
type wishlistParamsWrapper struct {
	// Requires Userid, Productid and Bool value to add or remove from wishlist
	// in:body
	// type:application/json
	Body models.Cartproduct
}

// WishlistHandler : Adding Products into User Wishlist and Removing Products from user Wishlist
func WishlistHandler(w http.ResponseWriter, r *http.Request) {

	Check("wishlist", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")
	var wishlist model.Cartproduct
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &wishlist)

	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"userid": wishlist.Userid}

	if wishlist.Status == true {
		update := bson.M{"$push": bson.M{"itemsId": wishlist.Productid}}
		query.UpdateOne("wishlist", filter, update)
		response := true
		json.NewEncoder(w).Encode(response)

	} else if wishlist.Status == false {
		update := bson.M{"$pull": bson.M{"itemsId": wishlist.Productid}}
		query.UpdateOne("wishlist", filter, update)
		response := false
		json.NewEncoder(w).Encode(response)
	}
}

// swagger:route POST /wishlistproducts User wishlistproducts
// Returns the wishlisted Products details of a User.
// responses:
//   200: wishlistproductsResponse

// Array of Products in Wishlist Of a User
// swagger:response wishlistproductsResponse
type wishlistproductsResponseWrapper struct {
	// The generated response
	// in:body
	Body []models.Items
}

// consumes:
// - application/json
// swagger:parameters wishlistproducts
type WishlistProductsParamsWrapper struct {
	//
	// in:body
	// type:application/json
	Body models.Id
}

// WishlistProductsHandler : Returns the wishlisted Products details of a User
func WishlistProductsHandler(w http.ResponseWriter, r *http.Request) {

	Check("wishlistproducts", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")
	var res model.ResponseResult
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return

	}
	var product model.Wishlistarray
	var list []model.Items
	var item model.Items
	err = query.FindoneID("wishlist", id.ID1, "userid").Decode(&product)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(list)
		return

	}

	collection, client := query.Connection("products")
	for i := 0; i < len(product.Wisharr); i++ {
		item.Img = nil
		item.Itemsid = nil
		_ = collection.FindOne(context.TODO(), bson.M{"_id": product.Wisharr[i]}).Decode(&item)
		list = append(list, item)
	}
	defer query.Endconn(client)
	json.NewEncoder(w).Encode(list)

}

// swagger:operation GET /productdetails/{productid} User ProductDetails
// ---
// summary: Details of a Product
// description: Details of a particular Product
// parameters:
// - name: productid
//   in: path
//   description: Product Id String
//   type: string
//   required: true
// responses:
//   200: ProductDetailsResponse

// This text will appear as description of your response body.
// swagger:response ProductDetailsResponse
type ProductDetailsResponseWrapper struct {
	// The generated response
	// in:body
	Body models.Items
}

// ProductDetailsHandler : handles the Details of a Product
func ProductDetailsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	Check("productdetails/"+params["productid"], "GET", w, r)
	w.Header().Set("Content-Type", "application/json")

	ProID, err := primitive.ObjectIDFromHex(params["productid"])
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var item model.Items
	err = query.FindoneID("products", ProID, "_id").Decode(&item)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(item)

}

// swagger:route POST /checkout User Checkout
// checkout from cart moves products from cart to user intransit and decrease stock.
// responses:
//   200: checkoutResponse
// swagger:response checkoutResponse
type checkoutResponseWrapper struct {
	// The generated response
	// in:body
	Body models.ResponseResult
}

// consumes:
// - application/json
// swagger:parameters Checkout
type checkoutParamsWrapper struct {
	// Details of product to be cancelled
	// in:body
	// type:application/json
	Body models.Id
	// The auth token
	// in:header
	// type:string
	Token string
	// Privilege
	// in:header
	// type:string
	Admin string
}

// CheckoutHandler :checkout from cart moves products from cart to user intransit and decrease stock
// CheckoutHandler :checkout from cart moves products from cart to user intransit and decrease stock
func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	Check("checkout", "POST", w, r)
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)

	if err != nil {
		helpers.ErrHandler(err, w)
	}

	filter := bson.M{"userid": id.ID1}
	update := bson.M{"$set": bson.M{"product": bson.A{}}}

	var products model.Cart
	products.Product = nil

	err = query.FindoneID("cart", id.ID1, "userid").Decode(&products)
	if err != nil {
		helpers.ErrHandler(err, w)
	}

	query.UpdateOne("cart", filter, update)

	var order model.Order
	order.Items_count = nil
	order.PayDates = nil
	var orders []model.Order
	orders = nil
	var prod model.Items
	prod.Itemsid = nil
	prod.Img = nil
	collection, client := query.Connection("products")
	for i := 0; i < len(products.Product); i++ {
		prod.Itemsid = nil
		order.Items_count = nil
		filter := bson.M{"_id": products.Product[i].P_id}
		update := bson.M{"$inc": bson.M{"stock": -products.Product[i].Count, "demand": products.Product[i].Count}}
		err = collection.FindOne(context.TODO(), bson.M{"_id": products.Product[i].P_id}).Decode(&prod)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		for j := 0; j < products.Product[i].Count; j++ {
			update1 := bson.M{"$pull": bson.M{"itemsid": prod.Itemsid[j]}}
			_, err = collection.UpdateOne(context.TODO(), filter, update1)
			if err != nil {
				helpers.ErrHandler(err, w)
			}
			order.Items_count = append(order.Items_count, prod.Itemsid[j])
		}

		order.P_id = products.Product[i].P_id
		order.Rent = products.Product[i].Rent
		order.Count = products.Product[i].Count
		order.Duration = products.Product[i].Duration
		order.IsCancellled = false
		order.Due = 0
		order.Date = time.Now()
		orders = append(orders, order)
		_, err = collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
		if err != nil {
			helpers.ErrHandler(err, w)
		}

	}

	query.Endconn(client)
	var orderids []primitive.ObjectID
	orderids = nil
	collection1, client1 := query.Connection("orders")
	for i := 0; i < len(products.Product); i++ {
		result, err := collection1.InsertOne(context.TODO(), orders[i])
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		oid := result.InsertedID.(primitive.ObjectID)
		orderids = append(orderids, oid)
	}
	query.Endconn(client1)
	collection2, client2 := query.Connection("user")
	filter1 := bson.M{"_id": id.ID1}
	for i := 0; i < len(orderids); i++ {
		update1 := bson.M{"$push": bson.M{"intransit": orderids[i]}}
		_, err = collection2.UpdateOne(context.TODO(), filter1, update1)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
	}

	query.Endconn(client2)
}

//SearcEngine api
func SearchEngine(w http.ResponseWriter, r *http.Request) {

	Check("searchengine", "POST", w, r)

	w.Header().Set("Content-Type", "application/json")

	body, _ := ioutil.ReadAll(r.Body)

	var srch model.SearchProduct

	var res model.ResponseResult

	err := json.Unmarshal(body, &srch)
	if err != nil {
		log.Fatal(w, "error occured while unmarshling")
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
	}
	collection, client := query.Connection("products")
	search := bson.M{"$text": bson.M{"$search": srch.Search}}

	cursor, err := collection.Find(context.TODO(), search)
	if err != nil {
		log.Fatal(err)
	}

	var show []model.Items
	var product model.Items
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		product.Img = nil
		product.Itemsid = nil
		if err = cursor.Decode(&product); err != nil {
			log.Fatal(err)
		}
		show = append(show, product)
	}

	json.NewEncoder(w).Encode(show)
	query.Endconn(client)
}

//cart products api
func CartProducts(w http.ResponseWriter, r *http.Request) {

	Check("cartproducts", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var item model.Cart

	err = query.FindoneID("cart", id.ID1, "userid").Decode(&item)
	if err != nil {
		fmt.Println("Incart")
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var cartProdList []model.CartOrder
	cartProdList = nil
	var cartItem model.Cart
	cartItem.Order = nil
	cartItem.Product = nil
	for i := 0; i < len(item.Product); i++ {

		var prodCollectionVar model.Items

		err = query.FindoneID("products", item.Product[i].P_id, "_id").Decode(&prodCollectionVar)
		if err != nil {
			//fmt.Println("Inproducts")
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		cartItem.Order = make([]model.CartOrder, len(item.Product))

		cartItem.Order[i].P_id = item.Product[i].P_id
		cartItem.Order[i].Name = prodCollectionVar.Name
		cartItem.Order[i].Count = item.Product[i].Count
		cartItem.Order[i].Rent = item.Product[i].Rent
		cartItem.Order[i].Duration = item.Product[i].Duration
		cartItem.Order[i].Img = prodCollectionVar.Img[0]

		cartItem.Order[i].Deposit = prodCollectionVar.Deposit

		cartProdList = append(cartProdList, cartItem.Order[i])
		//fmt.Println("prod", cartProdList)

	}
	json.NewEncoder(w).Encode(cartProdList)
}

//Cart First Time
func CartFirstTime(w http.ResponseWriter, r *http.Request) {

	Check("cartfirsttime", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var ct model.CartInput
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &ct)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	collection, client := query.Connection("cart")
	ct.Order.Pstr = ct.Order.P_id.Hex() + strconv.Itoa(ct.Order.Duration)
	var doc model.Cart
	//fmt.Println(ct.Order)
	err = collection.FindOne(context.TODO(), bson.M{"userid": ct.Userid, "product.p_id": ct.Order.P_id, "product.duration": ct.Order.Duration}).Decode(&doc)

	if err != nil {

		_, err = collection.UpdateOne(context.TODO(), bson.M{"userid": ct.Userid}, bson.M{"$push": bson.M{"product": ct.Order}})
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		res1 := "New product added"
		json.NewEncoder(w).Encode(res1)

	} else {

		update := bson.M{"$set": bson.M{"product.$.count": ct.Order.Count + 1}} //ct.Product.Count + 1}}

		_, err = collection.UpdateOne(context.TODO(), bson.M{"userid": ct.Userid, "product.pstr": ct.Order.Pstr, "product.duration": ct.Order.Duration}, update)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		res2 := "Count of product increased"
		json.NewEncoder(w).Encode(res2)

	}
	query.Endconn(client)

}

/*
//CartInput
func CartInput(w http.ResponseWriter, r *http.Request) {

	Check("cartinput", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var ct model.CartInput
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &ct)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
	}
	collection, client := query.Connection("cart")

	var doc model.Cart

	err = collection.FindOne(context.TODO(), bson.M{"userid": ct.Userid, "product.p_id": ct.Product.P_id, "product.count": ct.Product.Count, "product.duration": ct.Product.Duration}).Decode(&doc)
	if err != nil {
		//didn't found any match
		_, err = collection.UpdateOne(context.TODO(), bson.M{"userid": ct.Userid}, bson.M{"$push": bson.M{"product": ct.Product}})

		respn := "New Product Added"
		json.NewEncoder(w).Encode(respn)

	} else {
		// if found the match
		_, err = collection.UpdateOne(context.TODO(), bson.M{"userid": ct.Userid, "product.p_id": ct.Product.P_id}, bson.M{"$set": bson.M{"product.count": ct.Product.Count, "product.duration": ct.Product.Duration, "product._rent": ct.Product.Rent}})

		respm := "Existing Product Updated"
		json.NewEncoder(w).Encode(respm)

	}
	query.Endconn(client)
}
*/
//Remove Cart Products
func RemoveCartProduct(w http.ResponseWriter, r *http.Request) {

	Check("removecartproduct", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var rem model.RemoveCartProduct

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &rem)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	fmt.Println(rem)
	filter := bson.M{"userid": rem.UserId}
	update := bson.M{"$pull": bson.M{"product": bson.M{"pstr": rem.ProductId.Hex() + strconv.Itoa(rem.Duration), "duration": rem.Duration}}}
	query.UpdateOne("cart", filter, update)

	if err != nil {
		helpers.ErrHandler(err, w)
	}

	respn := "Data Removed"
	json.NewEncoder(w).Encode(respn)
}

// swagger:operation GET /stock/{productid} User ProductStock
// ---
// summary: Returns the stock of the product
// description: Returns the stock of the desired Product
// parameters:
// - name: productid
//   in: path
//   description: Product Id String
//   type: string
//   required: true
// responses:
//   200: ProductStockResponse

// This text will appear as description of your response body.
// swagger:response ProductStockResponse
type ProductStockResponseWrapper struct {
	// The generated response
	// in:body
	Body int
}

// ProductStock : Returns the stock of the desired Product
func ProductStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	Check("stock/"+params["productid"], "GET", w, r)
	w.Header().Set("Content-Type", "application/json")

	var sd model.Items
	var res model.ResponseResult
	LocId, err := primitive.ObjectIDFromHex(params["productid"])
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	err1 := query.FindoneID("products", LocId, "_id").Decode(&sd)
	if err1 != nil {
		res.Error = err1.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	json.NewEncoder(w).Encode(sd.Stock)
}

func CartUpdate(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/cartupdate" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//	userid, value, productid,count,duration.

	var ct model.CartInput
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &ct)

	if err != nil {
		helpers.ErrHandler(err, w)
	}
	fmt.Println(ct)
	/*pipeline := bson.M{"$and": []bson.M{
	{"userid": ct.Userid},

	{"$and": []bson.M{
		{"product.duration": bson.M{"$eq": ct.Product.Duration}},
		{"product.pstr": bson.M{"$eq": ct.Product.P_id.Hex()}}}}}}*/
	//fmt.Println(ct.Product.P_id.Hex() + strconv.Itoa(ct.Product.Duration))
	filter := bson.M{"userid": ct.Userid, "product.pstr": ct.Product.P_id.Hex() + strconv.Itoa(ct.Product.Duration)}

	collection, client := query.Connection("cart")
	if ct.Status == true {

		update1 := bson.M{"$inc": bson.M{"product.$.count": ct.Value}}
		_, err := collection.UpdateOne(context.TODO(), filter, update1)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		res2 := "Count of product increased"
		json.NewEncoder(w).Encode(res2)

	} else if ct.Status == false {

		update2 := bson.M{"$inc": bson.M{"product.$.count": -ct.Value}}
		_, err := collection.UpdateOne(context.TODO(), filter, update2)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		res2 := "Count of product decreased"
		json.NewEncoder(w).Encode(res2)
	}
	query.Endconn(client)
}

// swagger:route POST /stockcheck User Stockcheck
// Checks if all the products added in cart is available or not.
// responses:
//   200: stockcheckResponse

//
// swagger:response stockcheckResponse
type stockcheckResponseWrapper struct {
	// The generated response
	// in:body
	Body string
}

// consumes:
// - application/json
// swagger:parameters Stockcheck
type stockcheckParamsWrapper struct {
	// in:body
	// type:application/json
	Body models.Id
}

// StockCheckHandler : Checks if all the products added in cart is available or not,returns string upon Success
func StockCheckHandler(w http.ResponseWriter, r *http.Request) {
	Check("stockcheck", "POST", w, r)
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var products model.Cart

	err1 := query.FindoneID("cart", id.ID1, "userid").Decode(&products)
	if err1 != nil {
		res.Error = err1.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var count []int
	var productid []primitive.ObjectID
	count = nil
	productid = nil
	dict1 := make(map[primitive.ObjectID]int)
	for i := 0; i < len(products.Product); i++ {
		count = append(count, products.Product[i].Count)
		productid = append(productid, products.Product[i].P_id)
		dict1[products.Product[i].P_id] = products.Product[i].Count

	}
	var ids []primitive.ObjectID
	ids = nil
	for key := range dict1 {
		ids = append(ids, key)

	}
	var finalcount []int
	finalcount = nil
	for i := 0; i < len(ids); i++ {
		sum := 0
		for j := 0; j < len(productid); j++ {
			if ids[i] == productid[j] {
				sum += count[j]
			}
		}
		finalcount = append(finalcount, sum)
	}
	var stock []int
	var item model.Items
	stock = nil
	collection, client := query.Connection("products")

	for i := 0; i < len(ids); i++ {
		err := collection.FindOne(context.TODO(), bson.M{"_id": ids[i]}).Decode(&item)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		stock = append(stock, item.Stock)

	}
	fmt.Println(ids, finalcount, stock)

	for i := 0; i < len(stock); i++ {

		if finalcount[i] > stock[i] {

			json.NewEncoder(w).Encode("fail")
			return
		}
	}
	if len(count) == 0 {
		json.NewEncoder(w).Encode("No Item in Cart")
	} else {
		json.NewEncoder(w).Encode("success")
	}

	query.Endconn(client)
}

// swagger:operation GET /userlist/{status} Admin Userlist
// ---
// summary: List the orders given by users
// description: All orders including intransit/current/past
// parameters:
// - name: status
//   in: path
//   description: status of the order
//   type: string
//   required: true
// - name: Token
//   in: header
//   description: auth token
//   type: string
// - name: Admin
//   in: header
//   description: privilege
//   type: string
// responses:
//   200: userlistResponse

// This text will appear as description of your response body.
// swagger:response userlistResponse
type userlistResponseWrapper struct {
	// The generated response
	// in:body
	Body []models.User
}

// UserListHandler : handles the userlist
func UserList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Groom")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	Check("userlist/"+params["status"], "GET", w, r)
	var projection bson.M
	var cursor *mongo.Cursor
	var user model.User
	var list []model.User
	list = nil
	if params["status"] == "intransit" {
		projection = bson.M{"currentorder": 0, "pastorder": 0, "isadmin": 0}
	} else if params["status"] == "currentorder" {
		projection = bson.M{"intransit": 0, "pastorder": 0, "isadmin": 0}
	} else if params["status"] == "pastorder" {
		projection = bson.M{"intransit": 0, "currentorder": 0, "isadmin": 0}
	} else if params["status"] == "cancelled" {
		projection = bson.M{"isadmin": 0, "intransit": 0, "currentorder": 0}
	}

	collection, client := query.Connection("user")

	var filter bson.M
	if params["status"] == "cancelled" {
		filter = bson.M{"pastorder": bson.M{"$exists": true, "$not": bson.M{"$size": 0}}}
	} else {
		filter = bson.M{params["status"]: bson.M{"$exists": true, "$not": bson.M{"$size": 0}}}
	}

	cursor, err := collection.Find(
		context.TODO(),
		filter,
		options.Find().SetProjection(projection),
	)
	if err != nil {
		helpers.ErrHandler(err, w)
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		user.CurrentOrder = nil
		user.InTransit = nil
		user.PastOrder = nil

		if err := cursor.Decode(&user); err != nil {
			fmt.Fprintln(w, "failed to decode")
		}
		list = append(list, user)
	}
	query.Endconn(client)

	Ordercollection, client1 := query.Connection("orders")
	var order model.Order
	order.Items_count = nil
	order.PayDates = nil

	var response []model.Result
	response = nil
	var result model.Result
	var product model.Items
	for i := 0; i < len(list); i++ {
		result.CurrentOrder = nil
		result.InTransit = nil
		result.PastOrder = nil
		result.ID = list[i].ID
		result.Address = list[i].Address
		result.Email = list[i].Email
		result.Name = list[i].Name
		result.Phone = list[i].Phone

		if params["status"] == "intransit" {
			for j := 0; j < len(list[i].InTransit); j++ {
				order.Items_count = nil
				order.PayDates = nil
				product.Itemsid = nil
				product.Img = nil
				err = Ordercollection.FindOne(context.TODO(), bson.M{"_id": list[i].InTransit[j]}).Decode(&order)
				if err != nil {
					helpers.ErrHandler(err, w)
				}
				err = query.FindoneID("products", order.P_id, "_id").Decode(&product)
				if err != nil {
					helpers.ErrHandler(err, w)
				}
				order.Img = product.Img[0]
				order.Name = product.Name
				order.Deposit = product.Deposit
				result.InTransit = append(result.InTransit, order)
			}
		} else if params["status"] == "currentorder" {
			order.Items_count = nil
			order.PayDates = nil
			product.Itemsid = nil
			product.Img = nil
			for j := 0; j < len(list[i].CurrentOrder); j++ {
				err = Ordercollection.FindOne(context.TODO(), bson.M{"_id": list[i].CurrentOrder[j]}).Decode(&order)
				if err != nil {
					helpers.ErrHandler(err, w)
				}
				err = query.FindoneID("products", order.P_id, "_id").Decode(&product)
				if err != nil {
					helpers.ErrHandler(err, w)
				}
				order.Img = product.Img[0]
				order.Name = product.Name
				order.Deposit = product.Deposit
				result.CurrentOrder = append(result.CurrentOrder, order)
			}
		} else if params["status"] == "pastorder" {
			order.Items_count = nil
			order.PayDates = nil
			product.Itemsid = nil
			product.Img = nil
			for j := 0; j < len(list[i].PastOrder); j++ {
				err = Ordercollection.FindOne(context.TODO(), bson.M{"_id": list[i].PastOrder[j]}).Decode(&order)
				if err != nil {
					helpers.ErrHandler(err, w)
				}
				if order.IsCancellled == true {
					continue
				}
				err = query.FindoneID("products", order.P_id, "_id").Decode(&product)
				if err != nil {
					helpers.ErrHandler(err, w)
				}
				order.Img = product.Img[0]
				order.Name = product.Name
				order.Deposit = product.Deposit
				result.PastOrder = append(result.PastOrder, order)
			}
		} else if params["status"] == "cancelled" {
			order.Items_count = nil
			order.PayDates = nil
			product.Itemsid = nil
			product.Img = nil
			for j := 0; j < len(list[i].PastOrder); j++ {
				err = Ordercollection.FindOne(context.TODO(), bson.M{"_id": list[i].PastOrder[j], "iscancelled": true}).Decode(&order)
				if err != nil {
					helpers.ErrHandler(err, w)
				}
				err = query.FindoneID("products", order.P_id, "_id").Decode(&product)
				if err != nil {
					helpers.ErrHandler(err, w)
				}
				order.Img = product.Img[0]
				order.Name = product.Name
				order.Deposit = product.Deposit
				result.PastOrder = append(result.PastOrder, order)
			}
		}
		if len(result.CurrentOrder) == 0 && len(result.PastOrder) == 0 && len(result.InTransit) == 0 {
			continue
		} else {
			response = append(response, result)
		}

	}
	json.NewEncoder(w).Encode(response)
	query.Endconn(client1)
	fmt.Println("Bride")
}

// swagger:route POST /editprofile User profilehandling
// Allows users to edit their profile.
// responses:
//   200: editprofileNewResponse

// Allows a user to edit his/her profile
// swagger:response editprofileResponse
type editprofileResponseWrapper struct {
	// The generated response
	// in:body
	Body models.ResponseResult
}

// consumes:
// - application/json
// swagger:parameters profilehandling
type editprofileParamsWrapper struct {
	// The edit details
	// in:body
	// type:application/json
	Body models.User
}

// NewSignupHandler : Handles the user signup
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	Check("editprofile", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")
	var data model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &data)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	filter := bson.M{"_id": data.ID}
	update := bson.M{"$set": data}
	query.UpdateOne("user", filter, update)

	res.Result = "Details successfully updated!"
	res.Error = ""
	json.NewEncoder(w).Encode(res)

}

// swagger:operation POST /values Admin productcreation
// ---
// summary: Allows admin to add a new product
// description: Allows admin to add a new product
// parameters:
// - name: Name
//   in: formData
//   description: name of the product
//   type: string
//   required: true
// - name: Price
//   in: formData
//   description: price of the product
//   type: string
// - name: Details
//   in: formData
//   description: details about the product
//   type: string
// - name: Rent
//   in: formData
//   description: rent of the product
//   type: string
// - name: Deposit
//   in: formData
//   description: deposit of the product
//   type: string
// - name: Stock
//   in: formData
//   description: stock of the product
//   type: string
// - name: Img0
//   in: formData
//   description: The first imagez
//   type: file
// - name: Img1
//   in: formFile
//   description: The second imagez
//   type: file

// responses:
//   200: valuesResponse

// This text will appear as description of your response body.
// swagger:response valuesResponse
type valuesResponseWrapper struct {
	// The generated response
	// in:body
	Body models.ResponseResult
}

// ValueHandler : Handles the creation of a product
func ValueHandler(w http.ResponseWriter, r *http.Request) {

	Check("values", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	name := r.FormValue("name")
	price := r.FormValue("price")
	details := r.FormValue("details")
	rent := r.FormValue("rent")
	deposit := r.FormValue("deposit")
	stock := r.FormValue("stock")
	subcategoryid := r.FormValue("subcategoryid")
	locationid := r.FormValue("locationid")
	subname := r.FormValue("subname")
	catname := r.FormValue("catname")
	indices := r.FormValue("indices")
	from := r.FormValue("from")
	fmt.Println(r)

	fmt.Println(subname, catname)

	jsonData := map[string]string{"name": name, "price": price, "details": details, "rent": rent, "deposit": deposit, "subcategoryid": subcategoryid, "locationid": locationid}
	jsonValue, _ := json.Marshal(jsonData)

	detailsResponse, err := http.Post("http://localhost:8080/api/details", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		insertedID, _ := ioutil.ReadAll(detailsResponse.Body)
		fmt.Println(string(insertedID))
		fmt.Println(insertedID)

	}

	jsonData1 := map[string]string{"quantity": stock, "productid": teststr}
	jsonValue1, _ := json.Marshal(jsonData1)

	stockResponse, err := http.Post("http://localhost:8080/api/stocker", "application/json", bytes.NewBuffer(jsonValue1))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		fmt.Println("hella")
		response, _ := ioutil.ReadAll(stockResponse.Body)
		fmt.Println(string(response))
	}

	helpers.ProductImageHandler(w, r, testobj, catname, subname, indices, from)

}

// DetailsHandler ...
func DetailsHandler(w http.ResponseWriter, r *http.Request) {

	Check("details", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var details model.ProductDetails
	var convDetails model.ProductUpload

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &details)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}

	convDetails.Subcategoryid, _ = primitive.ObjectIDFromHex(details.Subcategoryid)
	convDetails.LocationID, _ = primitive.ObjectIDFromHex(details.LocationID)
	convDetails.Name = details.Name
	convDetails.Details = details.Details
	convDetails.Price, _ = strconv.Atoi(details.Price)
	convDetails.Rent, _ = strconv.Atoi(details.Rent)
	convDetails.Deposit, _ = strconv.Atoi(details.Deposit)
	convDetails.Itemsid = []primitive.ObjectID{}
	convDetails.Img = []string{}
	convDetails.Stock = 0
	convDetails.Demand = 0
	convDetails.Createdat = time.Now()

	res := query.InsertOne("products", convDetails)
	json.NewEncoder(w).Encode(res.InsertedID)
	testobj = res.InsertedID.(primitive.ObjectID)
	teststr = testobj.Hex()
}

// StockHiandler ...
func StockHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hellao")

	Check("stocker", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var stock model.ProductStock
	var items model.ProductItems
	var itemarr []primitive.ObjectID

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &stock)
	if err != nil {
		json.NewEncoder(w).Encode("Error in Unmarshalling")
		return
	}

	quantity, _ := strconv.Atoi(stock.Quantity)
	fmt.Println("hello")
	items.Productid, _ = primitive.ObjectIDFromHex(stock.Productid)
	items.Createdat = time.Now()

	for i := 0; i < quantity; i++ {
		res := query.InsertOne("items", items)
		if err != nil {
			json.NewEncoder(w).Encode("Error in creating items")
			return
		}
		itemarr = append(itemarr, res.InsertedID.(primitive.ObjectID))
	}

	for i := 0; i < quantity; i++ {
		filter := bson.M{"_id": items.Productid}
		update := bson.M{"$push": bson.M{"itemsid": itemarr[i]}}
		query.UpdateOne("products", filter, update)
	}

	filter := bson.M{"_id": items.Productid}
	update := bson.M{"$inc": bson.M{"stock": quantity}}
	query.UpdateOne("products", filter, update)

}

// swagger:route POST /delete Admin deletehandling
// Allows admin to delete a product.
// responses:
//   200: deleteResponse

// Allows an admin to delete a product
// swagger:response deleteResponse
type deleteResponseWrapper struct {
	// The generated response
	// in:body
	Body models.ResponseResult
}

// consumes:
// - application/json
// swagger:parameters deletehandling
type deleteParamsWrapper struct {
	// The edit details
	// in:body
	// type:application/json
	Body models.Delete
	// The auth token
	// in:header
	// type:string
	Token string
	// Privilege
	// in:header
	// type:string
	Admin string
}

// DeleteHandler : Handles the deletion of a product
func DeleteHandler(w http.ResponseWriter, r *http.Request) {

	Check("delete", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var delete model.Delete
	var deleted model.Deleteditems

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &delete)
	if err != nil {
		json.NewEncoder(w).Encode("Error in Unmarshalling")
		return
	}

	var itemsarr []primitive.ObjectID

	result := query.FindoneID("products", delete.Productid, "_id")
	if err = result.Decode(&deleted); err != nil {
		log.Fatal(err)
	}
	itemsarr = deleted.Itemsid
	fmt.Println(itemsarr)

	collection1, client1 := query.Connection("items")

	for i := 0; i < len(itemsarr); i++ {
		_, err = collection1.DeleteOne(context.TODO(), bson.M{"_id": itemsarr[i]})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	defer query.Endconn(client1)

	collection2, client2 := query.Connection("products")
	_, err = collection2.DeleteOne(context.TODO(), bson.M{"_id": delete.Productid})
	if err != nil {
		fmt.Println(err)
		return
	}

	defer query.Endconn(client2)
	fmt.Println("here", delete.Productid)

	collection3, client3 := query.Connection("wishlist")
	filter := bson.M{}
	update := bson.M{"$pull": bson.M{"itemsId": delete.Productid}}
	_, err = collection3.UpdateMany(context.Background(), filter, update)
	defer query.Endconn(client3)

}

// swagger:route POST /adminstock Admin stockhandling
// Allows admin to handle stock.
// responses:
//   200: adminstockResponse

// Allows an admin to handle stock of a product
// swagger:response adminstockResponse
type adminstockResponseWrapper struct {
	// The generated response
	// in:body
	Body models.ResponseResult
}

// consumes:
// - application/json
// swagger:parameters stockhandling
type adminstockParamsWrapper struct {
	// The edit details
	// in:body
	// type:application/json
	Body models.ProductStock
	// The auth token
	// in:header
	// type:string
	Token string
	// Privilege
	// in:header
	// type:string
	Admin string
}

// AdminStockHandler : Updates the stock of a particular product
func AdminStockHandler(w http.ResponseWriter, r *http.Request) {

	Check("adminstock", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var stock model.ProductStock
	var items model.ProductItems
	var product model.Items
	var itemarr []primitive.ObjectID

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &stock)
	if err != nil {
		json.NewEncoder(w).Encode("Error in Unmarshalling")
		return
	}

	quantity, _ := strconv.Atoi(stock.Quantity)
	items.Productid, _ = primitive.ObjectIDFromHex(stock.Productid)
	items.Createdat = time.Now()

	collection1, client1, err := db.GetDBCollection("products")
	err = collection1.FindOne(context.TODO(), bson.D{{Key: "_id", Value: items.Productid}}).Decode(&product)

	if product.Stock < quantity {

		for i := 0; i < (quantity - product.Stock); i++ {
			res := query.InsertOne("items", items)
			if err != nil {
				json.NewEncoder(w).Encode("Error in creating items")
				return
			}
			itemarr = append(itemarr, res.InsertedID.(primitive.ObjectID))
		}

		for i := 0; i < (quantity - product.Stock); i++ {
			filter := bson.M{"_id": items.Productid}
			update := bson.M{"$push": bson.M{"itemsid": itemarr[i]}}
			query.UpdateOne("products", filter, update)
		}

		filter := bson.M{"_id": items.Productid}
		update := bson.M{"$inc": bson.M{"stock": (quantity - product.Stock)}}
		query.UpdateOne("products", filter, update)
	} else {

		for i := 0; i < (product.Stock - quantity); i++ {
			filter := bson.M{"_id": items.Productid}
			update := bson.M{"$pull": bson.M{"itemsid": product.Itemsid[i]}}
			query.UpdateOne("products", filter, update)
		}

		filter := bson.M{"_id": items.Productid}
		update := bson.M{"$inc": bson.M{"stock": (quantity - product.Stock)}}
		query.UpdateOne("products", filter, update)

	}
	err = client1.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

}

// swagger:route POST /adminupdate Admin productupdation
// Allows admin to create a product.
// responses:
//   200: adminupdateResponse

// Allows an admin to create a product
// swagger:response adminupdateResponse
type adminupdateResponseWrapper struct {
	// The generated response
	// in:body
	Body models.ResponseResult
}

// consumes:
// - multipart/form-data
// swagger:parameters productupdation
type adminupdateParamsWrapper struct {
	// The name of the product
	// in:formData
	// type:string
	Name string
	// The price of the product
	// in:formData
	// type:string
	Price string
	// The details of the product
	// in:formData
	// type:string
	Details string
	// The rent of the product
	// in:formData
	// type:string
	Rent string
	// The initial deposit of the product
	// in:formData
	// type:string
	Deposit string
	// The index of image file insertion
	// in:formData
	// type:string
	Indices string
	// The entry point creation/updation
	// in:formData
	// type:string
	From string
	// The first imagez
	// in:formData
	// type:file
	Img0 os.File
	// The second imagez
	// in:formData
	// type: file
	Img1 os.File
	// The third imagez
	// in:formData
	// type:file
	Img2 os.File
	// The auth token
	// in:header
	// type:string
	Token string
	// Privilege
	// in:header
	// type:string
	Admin string
}

// ValueHandler : Handles the creation of a product
func AdminUpdateHandler(w http.ResponseWriter, r *http.Request) {

	Check("adminupdate", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	id := r.FormValue("id")
	fmt.Println(id)
	name := r.FormValue("name")
	price := r.FormValue("price")
	details := r.FormValue("details")
	rent := r.FormValue("rent")
	deposit := r.FormValue("deposit")
	indices := r.FormValue("indices")
	from := r.FormValue("from")

	subname := "testing"
	catname := "admin"

	var produpdate model.ProductUpdate
	produpdate.Productid = query.DocId(id)
	produpdate.Name = name
	produpdate.Price, _ = strconv.Atoi(price)
	produpdate.Details = details
	produpdate.Rent, _ = strconv.Atoi(rent)
	produpdate.Deposit, _ = strconv.Atoi(deposit)

	filter := bson.M{"_id": query.DocId(id)}
	update := bson.M{"$set": produpdate}
	query.UpdateOne("products", filter, update)

	helpers.ProductImageHandler(w, r, query.DocId(id), catname, subname, indices, from)

}

// swagger:route POST /intransit User Intransit
// Checks if 72 hours is over from checkoutdate or not,After 72 hours moves product to currentorder and returns products in intransit.
// responses:
//   200: intransitResponse

// swagger:response intransitResponse
type intransitResponseWrapper struct {
	// The generated response
	// in:body
	Body []models.Product
}

// consumes:
// - application/json
// swagger:parameters Intransit
type intransitParamsWrapper struct {
	// Orders in transit
	// in:body
	// type:application/json
	Body models.Id
	// The auth token
	// in:header
	// type:string
	Token string
	// Privilege
	// in:header
	// type:string
	Admin string
}

// IntransitHandler : Checks if 72 hours is over from checkoutdate or not,After 72 hours moves product to currentorder and returns products in intransit
func IntransitHandler(w http.ResponseWriter, r *http.Request) {
	Check("intransit", "POST", w, r)
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)

	if err != nil {
		helpers.ErrHandler(err, w)
	}

	var user model.User
	collection, client := query.Connection("user")
	err = collection.FindOne(context.TODO(), bson.M{"_id": id.ID1}).Decode(&user)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	var order model.Order
	order.Items_count = nil
	order.PayDates = nil
	var response []model.Order
	response = nil
	for i := 0; i < len(user.InTransit); i++ {
		err = query.FindoneID("orders", user.InTransit[i], "_id").Decode(&order)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		response = append(response, order)

	}

	filter := bson.M{"_id": id.ID1}

	today := time.Now()

	for i := 0; i < len(response); i++ {
		target := response[i].Date.AddDate(0, 0, 3)
		if today.After(target) {
			filter1 := bson.M{"_id": response[i].ID}

			update1 := bson.M{"$push": bson.M{"paydates": time.Now()}}
			query.UpdateOne("orders", filter1, update1)

			update := bson.M{"$push": bson.M{"currentorder": response[i].ID}}
			_, err = collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				helpers.ErrHandler(err, w)
			}
			update2 := bson.M{"$pull": bson.M{"intransit": response[i].ID}}
			_, err = collection.UpdateOne(context.TODO(), filter, update2)
			if err != nil {
				helpers.ErrHandler(err, w)
			}

		}
	}
	err = collection.FindOne(context.TODO(), bson.M{"_id": id.ID1}).Decode(&user)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	response = nil
	var product model.Items
	order.Items_count = nil
	order.PayDates = nil

	for i := 0; i < len(user.InTransit); i++ {
		product.Itemsid = nil
		product.Img = nil
		err = query.FindoneID("orders", user.InTransit[i], "_id").Decode(&order)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		err = query.FindoneID("products", order.P_id, "_id").Decode(&product)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		order.Img = product.Img[0]
		order.Name = product.Name
		order.Deposit = product.Deposit
		response = append(response, order)

	}

	query.Endconn(client)
	json.NewEncoder(w).Encode(response)
}

// swagger:route POST /currentorder User CurrentOrder
// Checks if tenure is over from checkoutdate or not,After tenure is over moves product to pastorder and returns products in currentorder and updates stock in product collection.
// responses:
//   200: currentorderResponse

// swagger:response currentorderResponse
type currentorderResponseWrapper struct {
	// The generated response
	// in:body
	Body []models.Product
}

// consumes:
// - application/json
// swagger:parameters CurrentOrder
type currentorderParamsWrapper struct {
	// Products in Current Order
	// in:body
	// type:application/json
	Body models.Id
	// The auth token
	// in:header
	// type:string
	Token string
	// Privilege
	// in:header
	// type:string
	Admin string
}

// CurrentOrderHandler :Checks if tenure is over from checkoutdate or not,After tenure is over moves product to pastorder and returns products in currentorder and updates stock in product collection
func CurrentOrderHandler(w http.ResponseWriter, r *http.Request) {
	Check("currentorder", "POST", w, r)
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)

	if err != nil {
		helpers.ErrHandler(err, w)
	}

	var user model.User
	collection, client := query.Connection("user")
	err = collection.FindOne(context.TODO(), bson.M{"_id": id.ID1}).Decode(&user)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	var order model.Order
	order.Items_count = nil
	order.PayDates = nil
	var response []model.Order
	response = nil
	for i := 0; i < len(user.CurrentOrder); i++ {
		err = query.FindoneID("orders", user.CurrentOrder[i], "_id").Decode(&order)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		response = append(response, order)

	}
	var collected []model.Order
	collected = nil

	today := time.Now()

	for i := 0; i < len(response); i++ {
		if response[i].Duration == 12 {
			target := response[i].Date.AddDate(1, 0, 0)

			if today.After(target) {
				collected = append(collected, response[i])
				query.CurrentUpdate(response[i].ID, id.ID1, collection)
			}

		} else if response[i].Duration == 24 {

			target := response[i].Date.AddDate(2, 0, 0)
			if today.After(target) {
				collected = append(collected, response[i])
				query.CurrentUpdate(response[i].ID, id.ID1, collection)
			}
		} else if response[i].Duration == 6 {
			target := response[i].Date.AddDate(0, 6, 0)

			if today.After(target) {
				collected = append(collected, response[i])
				query.CurrentUpdate(response[i].ID, id.ID1, collection)
			}
		} else if response[i].Duration == 3 {
			target := response[i].Date.AddDate(0, 3, 0)

			if today.After(target) {
				collected = append(collected, response[i])
				query.CurrentUpdate(response[i].ID, id.ID1, collection)
			}
		}
	}
	err = collection.FindOne(context.TODO(), bson.M{"_id": id.ID1}).Decode(&user)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	var product model.Items
	response = nil
	order.Items_count = nil
	order.PayDates = nil
	for i := 0; i < len(user.CurrentOrder); i++ {
		product.Itemsid = nil
		product.Img = nil
		err = query.FindoneID("orders", user.CurrentOrder[i], "_id").Decode(&order)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		err = query.FindoneID("products", order.P_id, "_id").Decode(&product)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		order.Img = product.Img[0]
		order.Name = product.Name
		order.Deposit = product.Deposit
		response = append(response, order)

	}
	query.Endconn(client)
	collection1, client1 := query.Connection("products")
	for i := 0; i < len(collected); i++ {
		filter := bson.M{"_id": collected[i].P_id}
		update := bson.M{"$inc": bson.M{"collected": collected[i].Count * collected[i].Rent * collected[i].Duration, "stock": collected[i].Count}}
		_, err := collection1.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		for j := 0; j < len(collected[i].Items_count); j++ {
			update1 := bson.M{"$push": bson.M{"itemsid": collected[i].Items_count[j]}}
			_, err := collection1.UpdateOne(context.TODO(), filter, update1)
			if err != nil {
				helpers.ErrHandler(err, w)
			}
		}
	}
	query.Endconn(client1)

	json.NewEncoder(w).Encode(response)

}

// swagger:route POST /pastorder User PastOrder
// Returns all the cancelled and returned products of user.
// responses:
//   200: pastorderResponse

// swagger:response pastorderResponse
type pastorderResponseWrapper struct {
	// The generated response
	// in:body
	Body []models.Product
}

// consumes:
// - application/json
// swagger:parameters PastOrder
type pastorderParamsWrapper struct {
	// Products in Past Order
	// in:body
	// type:application/json
	Body models.Id
	// The auth token
	// in:header
	// type:string
	Token string
	// Privilege
	// in:header
	// type:string
	Admin string
}

// PastOrderHandler :Returns all the cancelled and returned products of user
func PastOrderHandler(w http.ResponseWriter, r *http.Request) {
	Check("pastorder", "POST", w, r)
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)

	if err != nil {
		helpers.ErrHandler(err, w)
	}

	var user model.User

	err = query.FindoneID("user", id.ID1, "_id").Decode(&user)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	var response []model.Order
	response = nil
	var order model.Order
	order.Items_count = nil
	order.PayDates = nil
	var product model.Items
	for i := 0; i < len(user.PastOrder); i++ {
		product.Itemsid = nil
		product.Img = nil
		err = query.FindoneID("orders", user.PastOrder[i], "_id").Decode(&order)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		err = query.FindoneID("products", order.P_id, "_id").Decode(&product)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
		order.Img = product.Img[0]
		order.Name = product.Name
		order.Deposit = product.Deposit
		response = append(response, order)

	}
	json.NewEncoder(w).Encode(response)

}

// swagger:route POST /payment Admin Payment
// Updates payment Report due,lastpaymentdate and monthspaid for a product in currentorder of a user.
// responses:
//   200: paymentResponse

//
// swagger:response paymentResponse
type paymentResponseWrapper struct {
	// The generated response
	// in:body
	Body models.UserReport
}

// consumes:
// - application/json
// swagger:parameters Payment
type paymentParamsWrapper struct {
	// in:body
	// type:application/json
	Body models.UserReportId
}

// PaymentHandler : Updates payment Report due,lastpaymentdate and monthspaid for a product in currentorder of a user
func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	Check("payment", "POST", w, r)
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	collection, client := query.Connection("orders")

	filter := bson.M{"_id": id.ID1}
	update := bson.M{"$set": bson.M{"due": 0}, "$push": bson.M{"paydates": time.Now()}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	query.Endconn(client)

}

// swagger:route POST /cancelorder User Cancelorder
// cancels products from intransit and current order and  returns bool.
// responses:
//   200: cancelorderResponse

// swagger:response cancelorderResponse
type cancelorderResponseWrapper struct {
	// The generated response
	// in:body
	Body bool
}

// consumes:
// - application/json
// swagger:parameters Cancelorder
type cancelorderParamsWrapper struct {
	// Details of product to be cancelled
	// in:body
	// type:application/json
	Body models.Id
	// The auth token
	// in:header
	// type:string
	Token string
	// Privilege
	// in:header
	// type:string
	Admin string
}

// CancelHandler :cancels products from intransit and current order and  returns bool
func CancelHandler(w http.ResponseWriter, r *http.Request) {
	Check("cancelorder", "POST", w, r)
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)

	if err != nil {
		helpers.ErrHandler(err, w)
	}
	var order model.Order
	order.Items_count = nil
	order.PayDates = nil
	collection, client := query.Connection("orders")

	err = collection.FindOne(context.TODO(), bson.M{"_id": id.ID1}).Decode(&order)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	update := bson.M{"$set": bson.M{"iscancelled": true}, "$push": bson.M{"paydates": time.Now()}}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id.ID1}, update)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	query.Endconn(client)
	collection1, client1 := query.Connection("user")

	filter := bson.M{"_id": id.UserId}
	update1 := bson.M{"$push": bson.M{"pastorder": id.ID1}, "$pull": bson.M{id.From: id.ID1}}
	_, err = collection1.UpdateOne(context.TODO(), filter, update1)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	query.Endconn(client1)
	collection2, client2 := query.Connection("products")
	update2 := bson.M{"$inc": bson.M{"stock": order.Count, "collected": order.Count * order.Rent * len(order.PayDates)}}
	_, err = collection2.UpdateOne(context.TODO(), filter, update2)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	for j := 0; j < len(order.Items_count); j++ {
		update3 := bson.M{"$push": bson.M{"itemsid": order.Items_count[j]}}
		_, err := collection2.UpdateOne(context.TODO(), filter, update3)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
	}
	query.Endconn(client2)
	status := true
	json.NewEncoder(w).Encode(status)
}

// swagger:route POST /userreport User Userreport
// Sends UserReport for a product in currentorder.
// responses:
//   200: userreportResponse

// swagger:response userreportResponse
type UserReportResponseWrapper struct {
	// The generated response
	// in:body
	Body models.UserReport
}

// consumes:
// - application/json
// swagger:parameters Userreport
type UserReportParamsWrapper struct {
	// in:body
	// type:application/json
	Body models.UserReportId
}

// UserReport : UserReport for a product in currentorder
func UserReport(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	Check("userreport", "POST", w, r)

	var res model.ResponseResult
	var id model.UserReportId
	var ur model.UserReport
	var sd model.StockData
	//var user model.User2
	var order model.Order

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)

	err = query.FindoneID("orders", id.Id, "_id").Decode(&order)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	err = query.FindoneID("products", order.P_id, "_id").Decode(&sd)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	price, rent, count, duration, totalPrice, profit, totalRent := 0, 0, 0, 0, 0, 0, 0
	price = sd.Price
	rent = order.Rent
	count = order.Count
	duration = order.Duration
	totalPrice = price * count
	totalRent = rent * count * duration
	profit = totalPrice - totalRent

	paid := len(order.PayDates) * rent * count
	payable := totalRent - paid

	due := order.Due

	ur.Name = sd.Name
	ur.Paid = paid
	ur.Payable = payable

	ur.Profit = profit
	ur.TotalRent = totalRent
	ur.TotalPrice = totalPrice
	ur.Due = due

	fmt.Println("len(order.PayDates)", len(order.PayDates))

	ur.LastPaymentDate = order.PayDates[len(order.PayDates)-1]
	ur.NextPaymentDate = order.PayDates[len(order.PayDates)-1].AddDate(0, 1, 0)

	fmt.Println("\npaid :", paid, "    payable : ", payable, " due  :", due)
	json.NewEncoder(w).Encode(ur)

	fmt.Println("UR is ", ur)

}

// swagger:route POST /signupNew User usersignup
// Allows users to signup to RHT.
// responses:
//   200: signupNewResponse

// Allows a user to signup
// swagger:response signupNewResponse
type signupNewResponseWrapper struct {
	// The generated response
	// in:body
	Body []models.ResponseResult
}

// consumes:
// - application/json
// swagger:parameters usersignup
type signupNewParamsWrapper struct {
	// The signup details
	// in:body
	// type:application/json
	Body models.User
}

// NewSignupHandler : Handles the user signup
func NewSignupHandler(w http.ResponseWriter, r *http.Request) {
	Check("signupNew", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var data model.User

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &data)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
	}
	fmt.Println(data)

	collection, client, err := db.GetDBCollection("user")
	err1 := collection.FindOne(context.TODO(), bson.M{"phone": data.Phone}).Err()

	if err1 != nil {
		res.Result = "Phone Authentication Required!"
		json.NewEncoder(w).Encode(res)
		fmt.Println(data.ID)

		otpauth(data.ID)
	} else {
		res.Result = "Phone Number already Registered!!"
		json.NewEncoder(w).Encode(res)
	}
	defer query.Endconn(client)

}

func otpauth(ID primitive.ObjectID) {

	helpers.LoadEnv()
	accountSid := helpers.GetEnvWithKey("ACCOUNT_SID")
	authToken := helpers.GetEnvWithKey("AUTH_TOKEN")
	urlStr := helpers.GetEnvWithKey("URL_STR")

	max := 9999
	min := 1000
	rand.Seed(time.Now().UnixNano())
	otp = strconv.Itoa(rand.Intn(max-min+1) + min)
	fmt.Println(otp)
	fmt.Println(ID)
	fmt.Println("otp getting set")

	filter := bson.M{"_id": ID}
	update := bson.M{"$set": bson.M{"otp": otp, "expiry": time.Now().Add(time.Minute * 5)}}
	query.UpdateOne("user", filter, update)

	msgData := url.Values{}
	msgData.Set("To", "+919438476609")
	msgData.Set("From", "+18589237950")
	msgData.Set("Body", "<#> RHT: Your code is "+otp+" "+helpers.GetEnvWithKey("APP_SIGNATURE"))
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)

		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
}

// swagger:route POST /resend User resendotp
// Allows users to resend otp.
// responses:
//   200: resendResponse

// Allows a user to resend his/her otp
// swagger:response resendResponse
type resendResponseWrapper struct {
	// The generated response
	// in:body
	Body []models.ResponseResult
}

// consumes:
// - application/json
// swagger:parameters resendotp
type resendParamsWrapper struct {
	// details for otp resend
	// in:body
	// type:application/json
	Body models.User
}

// ResendotpHandler : handles the otp resend
func Resendotp(w http.ResponseWriter, r *http.Request) {

	Check("resend", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var data model.User

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &data)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
	}
	fmt.Println("koli", data.ID)
	otpauth(data.ID)
}

// swagger:route POST /auth User authentication
// Allows users to authenticate details.
// responses:
//   200: authResponse

// Allows a user to authenticate his/her details
// swagger:response authResponse
type authResponseWrapper struct {
	// The generated response
	// in:body
	Body models.ResponseResult
}

// consumes:
// - application/json
// swagger:parameters authentication
type authParamsWrapper struct {
	// details for authentication
	// in:body
	// type:application/json
	Body models.OtpContainer
}

// AuthHandler : handles the authentication
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	Check("auth", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")
	var userOtp model.OtpContainer
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &userOtp)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var guestcart model.Cart
	var guestwishlist model.Wishlist
	var data models.User
	var code models.OtpCreds

	fmt.Println("hello")
	collection, client := query.Connection("user")
	_ = collection.FindOne(context.TODO(), bson.M{"_id": userOtp.ID}).Decode(&code)
	fmt.Println("the otp set is", code.Otp)
	t1 := time.Now()
	t2 := code.Expiry

	if userOtp.OtpEntered == code.Otp && t2.Sub(t1) > 0 && userOtp.From == "signup" {

		data.Address = userOtp.Address
		data.Name = userOtp.Name
		data.ID = userOtp.ID
		data.Phone = userOtp.Number
		data.Email = userOtp.Email
		data.Isadmin = false

		filter := bson.M{"_id": data.ID}
		update := bson.M{"$set": data}
		query.UpdateOne("user", filter, update)

		fmt.Println("The signUp authentication is successful!")
		res.Result = "The signUp authentication is successful!"
		var tokenacc string
		var tokenref string
		tokenacc, err = helpers.GenerateJWTAccess("false", data)
		tokenref, err = helpers.GenerateJWTRefresh("false", data)

		fmt.Println(tokenacc, tokenref)
		w.Header().Add("Token", tokenacc)
		w.Header().Add("Token1", tokenref)
		json.NewEncoder(w).Encode(res)

	} else if userOtp.OtpEntered != code.Otp && userOtp.From == "login" {
		res.Error = "OTP Did not Match!"
		json.NewEncoder(w).Encode(res)

	} else if userOtp.OtpEntered == code.Otp && t2.Sub(t1) > 0 && userOtp.From == "login" {

		check := "false"
		err = collection.FindOne(context.TODO(), bson.D{{Key: "phone", Value: userOtp.Number}}).Decode(&data)
		if data.Isadmin == true {
			check = "true"
		}

		if userOtp.ID != data.ID {
			collection1, client1, err := db.GetDBCollection("cart")
			err = collection1.FindOne(context.TODO(), bson.D{{Key: "userid", Value: userOtp.ID}}).Decode(&guestcart)

			for i := 0; i < len(guestcart.Product); i++ {
				filter := bson.M{"userid": data.ID}
				update := bson.M{"$push": bson.M{"product": guestcart.Product[i]}}
				query.UpdateOne("cart", filter, update)
				if err != nil {
					res.Error = err.Error()
					json.NewEncoder(w).Encode(res)
					return
				}
			}

			collection2, client2, err := db.GetDBCollection("wishlist")
			err = collection2.FindOne(context.TODO(), bson.D{{Key: "userid", Value: userOtp.ID}}).Decode(&guestwishlist)
			fmt.Println(guestwishlist)

			for i := 0; i < len(guestwishlist.ItemsId); i++ {
				err1 := collection2.FindOne(context.TODO(), bson.M{"userid": data.ID, "itemsId": guestwishlist.ItemsId[i]}).Err()
				if err1 != nil {
					filter := bson.M{"userid": data.ID}
					update := bson.M{"$push": bson.M{"itemsId": guestwishlist.ItemsId[i]}}
					query.UpdateOne("wishlist", filter, update)
					if err != nil {
						res.Error = err.Error()
						json.NewEncoder(w).Encode(res)
						return
					}
				}
			}
			_, err = collection2.DeleteOne(context.TODO(), bson.D{{Key: "userid", Value: userOtp.ID}})
			defer query.Endconn(client2)

			_, err = collection1.DeleteOne(context.TODO(), bson.D{{Key: "userid", Value: userOtp.ID}})
			defer query.Endconn(client1)

			_, err = collection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: userOtp.ID}})
			defer query.Endconn(client)
		}
		fmt.Println("The Login authentication is successful!")
		res.Result = "The Login authentication is successful!"
		fmt.Println("the data entering the jwt", data)

		var tokenacc string
		var tokenref string
		tokenacc, err = helpers.GenerateJWTAccess(check, data)
		tokenref, err = helpers.GenerateJWTRefresh(check, data)
		fmt.Println(tokenacc, tokenref)
		w.Header().Add("Token", tokenacc)
		w.Header().Add("Token1", tokenref)
		json.NewEncoder(w).Encode(res)

	} else if userOtp.OtpEntered != code.Otp && userOtp.From == "signup" {
		res.Error = "OTP Did not Match!"
		json.NewEncoder(w).Encode(res)
	} else if userOtp.OtpEntered != code.Otp && t2.Sub(t1) > 0 && userOtp.From == "numberchange" {
		res.Error = "OTP Did not Match!"
		json.NewEncoder(w).Encode(res)
	} else if userOtp.OtpEntered == code.Otp && t2.Sub(t1) > 0 && userOtp.From == "numberchange" {

		collection, client := query.Connection("user")

		err1 := collection.FindOne(context.TODO(), bson.M{"phone": userOtp.Number}).Err()
		if err1 != nil {
			filter := bson.M{"email": userOtp.Email}
			update := bson.M{"$set": bson.M{"phone": userOtp.Number}}
			query.UpdateOne("user", filter, update)
			defer query.Endconn(client)
			res.Result = "Number change successful!!"
			json.NewEncoder(w).Encode(res)
		} else {
			res.Result = "Number already taken!!"
			json.NewEncoder(w).Encode(res)

		}

	}
}

// swagger:route POST /loginNew User login
// Allows users to login into the app.
// responses:
//   200: loginNewResponse

// Allows a user to login into the app
// swagger:response loginNewResponse
type loginNewResponseWrapper struct {
	// hello2
	// The generated response
	// in:body
	Body models.ResponseResult
}

// consumes:
// - application/json
// swagger:parameters login
type loginNewParamsWrapper struct {
	// hello1
	// details for logging in
	// in:body
	// type:application/json
	Body models.Newlogin
}

// NewLoginHandler : handles the login
func NewLoginHandler(w http.ResponseWriter, r *http.Request) {

	Check("loginNew", "POST", w, r)
	var login model.Newlogin
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &login)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println(login.Contact)

	collection, client, err := db.GetDBCollection("user")
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{Key: "phone", Value: ("+91" + login.Contact)}}).Decode(&result)
	fmt.Println(result)
	fmt.Println(err)
	if err != nil {

		res.Result = "User Doesnot Exist!!"
		json.NewEncoder(w).Encode(res)
		return
	}
	if err == nil {
		fmt.Println("Otp getting called")
		otpauth(login.Userid)
		res.Result = "User Exists!!"
		json.NewEncoder(w).Encode(res)
		return
	}
	defer query.Endconn(client)
}

// swagger:operation GET /categorylist User CategoryList
// ---
// summary: Returns All the Subcategories Under all Categories
// description: Returns All the Subcategories Under all Categories with Icon Urls
// responses:
//   200: CategoryListResponse

// CategoryListResponseHandler : handles the category
func CategoryListResponseHandler(w http.ResponseWriter, r *http.Request) {
	Check("categorylist", "GET", w, r)
	var res model.ResponseResult
	collection, client := query.Connection("category")
	cursor, err := collection.Find(context.TODO(), bson.M{"archived": false}, options.Find().SetSort(bson.M{"categoryName": 1}))
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var category model.Category
	var catlist []model.Category
	catlist = nil

	for cursor.Next(context.TODO()) {

		if err = cursor.Decode(&category); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		catlist = append(catlist, category)

	}
	cursor.Close(context.TODO())
	query.Endconn(client)
	collection1, client1 := query.Connection("subcategory")
	var subcategory model.Subcategory

	var response model.CatResponse
	var final []model.CatResponse
	final = nil
	for i := 0; i < len(catlist); i++ {
		response.SubArray = nil
		response.CategoryName = catlist[i].Catname
		result, err := collection1.Find(context.TODO(), bson.M{"categoryid": catlist[i].Catid, "archived": false})
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		for result.Next(context.TODO()) {

			if err = result.Decode(&subcategory); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			response.SubArray = append(response.SubArray, subcategory)

		}

		final = append(final, response)
		result.Close(context.TODO())
	}
	json.NewEncoder(w).Encode(final)
	query.Endconn(client1)

}

// swagger:route POST /paymentstatus User paymentstatus
// Adds due after one month from lastpayment date and returns new Due,lastpaydate and next paydate.
// responses:
//   200: paymentstatusResponse

//
// swagger:response paymentstatusResponse
type paymentstatusResponseWrapper struct {
	// The generated response
	// in:body
	Body []models.PaymentStatus
}

// consumes:
// - application/json
// swagger:parameters paymentstatus
type paymentstatusParamsWrapper struct {
	//
	// in:body
	// type:application/json
	Body models.Id
}

// PaymentStatusHandler :Adds due after one month from lastpayment date and returns new Due,lastpaydate and next paydate
func PaymentStatusHandler(w http.ResponseWriter, r *http.Request) {
	Check("paymentstatus", "POST", w, r)
	var id model.Id
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &id)

	if err != nil {
		helpers.ErrHandler(err, w)
	}
	var order model.Order
	order.Items_count = nil
	order.PayDates = nil
	filter := bson.M{"_id": id.ID1}
	collection, client := query.Connection("orders")
	err = collection.FindOne(context.TODO(), filter).Decode(&order)
	if err != nil {
		helpers.ErrHandler(err, w)
	}

	flag := false
	if time.Now().After(order.PayDates[len(order.PayDates)-1].AddDate(0, 1, 0)) {
		order.Due += order.Rent
		flag = true
	}
	if flag == true {
		filter := bson.M{"_id": id.ID1}
		update := bson.M{"$inc": bson.M{"due": order.Rent * order.Count}}
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			helpers.ErrHandler(err, w)
		}
	}
	var payment model.PaymentStatus

	payment.Due = order.Due
	payment.LastPaymentDate = order.PayDates[len(order.PayDates)-1]
	payment.NextPaymentDate = order.PayDates[len(order.PayDates)-1].AddDate(0, 1, 0)

	json.NewEncoder(w).Encode(payment)

	query.Endconn(client)
}

// HighDemanding
func AdminHighDemanding(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	Check("highdemanding/"+params["location"], "GET", w, r)

	var res model.ResponseResult

	productCollection, client, err := db.GetDBCollection("products")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println("params[location]: ", params["location"])
	if params["location"] == "all" {

		pipeline := []bson.M{

			{
				"$group": bson.M{
					"_id":       "$_id",
					"name":      bson.M{"$first": "$name"},
					"demand":    bson.M{"$first": "$demand"},
					"createdat": bson.M{"$first": "$createdat"},
				},
			},
			{
				"$sort": bson.M{"demand": -1},
			},
			{
				"$limit": 10,
			},
		}

		result, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		var demand model.AdminHighLowDemanding

		var demandList []model.AdminHighLowDemanding
		for result.Next(context.TODO()) {
			if err = result.Decode(&demand); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}
			demandList = append(demandList, demand)
		}

		json.NewEncoder(w).Encode(demandList)

	} else {

		LocID, err := primitive.ObjectIDFromHex(params["location"])
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		pipeline := []bson.M{
			{
				"$match": bson.M{

					"locationid": LocID,
				},
			},
			{
				"$group": bson.M{
					"_id":       "$_id",
					"name":      bson.M{"$first": "$name"},
					"demand":    bson.M{"$first": "$demand"},
					"createdat": bson.M{"$first": "$createdat"},
				},
			},
			{
				"$sort": bson.M{"demand": -1},
			},
			{
				"$limit": 10,
			},
		}

		result, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		var demand model.AdminHighLowDemanding

		var demandList []model.AdminHighLowDemanding
		for result.Next(context.TODO()) {
			if err = result.Decode(&demand); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}
			demandList = append(demandList, demand)
		}

		json.NewEncoder(w).Encode(demandList)
	}

	client.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}

// Low Demanding
func AdminLowDemanding(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	Check("lowdemanding/"+params["location"], "GET", w, r)

	var res model.ResponseResult
	productCollection, client, err := db.GetDBCollection("products")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println("params[location]: ", params["location"])
	if params["location"] == "all" {

		pipeline := []bson.M{

			{
				"$group": bson.M{
					"_id":       "$_id",
					"name":      bson.M{"$first": "$name"},
					"demand":    bson.M{"$first": "$demand"},
					"createdat": bson.M{"$first": "$createdat"},
				},
			},
			{
				"$sort": bson.M{"demand": 1},
			},
			{
				"$limit": 10,
			},
		}

		result, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		var demand model.AdminHighLowDemanding

		var demandList []model.AdminHighLowDemanding
		for result.Next(context.TODO()) {
			if err = result.Decode(&demand); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}
			demandList = append(demandList, demand)
		}

		json.NewEncoder(w).Encode(demandList)

	} else {

		LocID, err := primitive.ObjectIDFromHex(params["location"])
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		pipeline := []bson.M{
			{
				"$match": bson.M{

					"locationid": LocID,
				},
			},
			{
				"$group": bson.M{
					"_id":       "$_id",
					"name":      bson.M{"$first": "$name"},
					"demand":    bson.M{"$first": "$demand"},
					"createdat": bson.M{"$first": "$createdat"},
				},
			},
			{
				"$sort": bson.M{"demand": 1},
			},
			{
				"$limit": 10,
			},
		}

		result, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		var demand model.AdminHighLowDemanding

		var demandList []model.AdminHighLowDemanding
		for result.Next(context.TODO()) {
			if err = result.Decode(&demand); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}
			demandList = append(demandList, demand)
		}

		json.NewEncoder(w).Encode(demandList)
	}

	client.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}

//Admin High Profit
func AdminHighProfit(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	Check("highprofit/"+params["location"], "GET", w, r)

	var res model.ResponseResult

	productCollection, client1, err := db.GetDBCollection("products")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	if params["location"] == "all" {

		pipeline := []bson.M{
			{
				"$group": bson.M{
					"_id":   "$_id",
					"name":  bson.M{"$first": "$name"},
					"price": bson.M{"$first": "$price"},

					"collected": bson.M{"$first": "$collected"},
					"createdat": bson.M{"$first": "$createdat"},
				},
			},

			{
				"$sort": bson.M{"collected": -1, "createdat": -1},
			},
			{
				"$limit": 10,
			},
		}

		productResult, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		defer productResult.Close(context.TODO())

		var prodResult model.AdminProfitLossOutput
		var prodResultList []model.AdminProfitLossOutput
		for productResult.Next(context.TODO()) {
			if err = productResult.Decode(&prodResult); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			prodResultList = append(prodResultList, prodResult) // contains list of all category id's & Name

		}
		json.NewEncoder(w).Encode(prodResultList)

	} else {

		LocID, err := primitive.ObjectIDFromHex(params["location"])
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		pipeline := []bson.M{
			{
				"$match": bson.M{

					"locationid": LocID,
				},
			},
			{
				"$group": bson.M{
					"_id":   "$_id",
					"name":  bson.M{"$first": "$name"},
					"price": bson.M{"$first": "$price"},

					"collected": bson.M{"$first": "$collected"},
					"createdat": bson.M{"$first": "$createdat"},
				},
			},

			{
				"$sort": bson.M{"collected": -1, "createdat": -1},
			},
			{
				"$limit": 10,
			},
		}

		productResult, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		defer productResult.Close(context.TODO())

		var prodResult model.AdminProfitLossOutput
		var prodResultList []model.AdminProfitLossOutput
		for productResult.Next(context.TODO()) {
			if err = productResult.Decode(&prodResult); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			prodResultList = append(prodResultList, prodResult) // contains list of all category id's & Name

		}
		json.NewEncoder(w).Encode(prodResultList)

	}
	client1.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}

//Admin Least Profitable
func AdminLeastProfit(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	Check("leastprofit/"+params["location"], "GET", w, r)

	var res model.ResponseResult
	productCollection, client1, err := db.GetDBCollection("products")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println("Least Profit Started")

	if params["location"] == "all" {

		pipeline := []bson.M{
			{
				"$group": bson.M{
					"_id":   "$_id",
					"name":  bson.M{"$first": "$name"},
					"price": bson.M{"$first": "$price"},

					"collected": bson.M{"$first": "$collected"},
					"createdat": bson.M{"$first": "$createdat"},
				},
			},

			{
				"$sort": bson.M{"collected": 1, "createdat": 1},
			},
			{
				"$limit": 10,
			},
		}

		productResult, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		defer productResult.Close(context.TODO())

		var prodResult model.AdminProfitLossOutput
		var prodResultList []model.AdminProfitLossOutput
		prodResultList = nil
		for productResult.Next(context.TODO()) {
			if err = productResult.Decode(&prodResult); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			prodResultList = append(prodResultList, prodResult) // contains list of all category id's & Name

		}
		json.NewEncoder(w).Encode(prodResultList)

	} else {

		LocID, err := primitive.ObjectIDFromHex(params["location"])
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		pipeline := []bson.M{
			{
				"$match": bson.M{

					"locationid": LocID,
				},
			},
			{
				"$group": bson.M{
					"_id":   "$_id",
					"name":  bson.M{"$first": "$name"},
					"price": bson.M{"$first": "$price"},

					"collected": bson.M{"$first": "$collected"},
					"createdat": bson.M{"$first": "$createdat"},
				},
			},

			{
				"$sort": bson.M{"collected": 1, "createdat": 1},
			},
			{
				"$limit": 10,
			},
		}

		productResult, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		defer productResult.Close(context.TODO())

		var prodResult model.AdminProfitLossOutput
		var prodResultList []model.AdminProfitLossOutput
		prodResultList = nil
		for productResult.Next(context.TODO()) {
			if err = productResult.Decode(&prodResult); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			prodResultList = append(prodResultList, prodResult) // contains list of all category id's & Name

		}
		json.NewEncoder(w).Encode(prodResultList)

	}
	client1.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println("Least Profit Ended")
}

// swagger:route POST /numberchange User NumberChange
// Allows users to change their phone number.
// responses:
//   200: numberchangeResponse

// Allows users to change their phone number
// swagger:response numberchangeResponse
type numberchangeResponseWrapper struct {
	// The generated response
	// in:body
	Body models.ResponseResult
}

// consumes:
// - application/json
// swagger:parameters NumberChange
type numberchangeParamsWrapper struct {
	// details for number change
	// in:body
	// type:application/json
	Body models.User
}

// NumberChangeHandler : Handles the user number change
func NumberChangeHandler(w http.ResponseWriter, r *http.Request) {

	Check("numberchange", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")

	var userDetails model.User
	var response model.ResponseResult

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &userDetails)
	if err != nil {
		json.NewEncoder(w).Encode("Error in Unmarshalling")
		return
	}

	collection, client := query.Connection("user")
	err1 := collection.FindOne(context.TODO(), bson.M{"email": userDetails.Email}).Err()
	if err1 == nil {

		err1 := collection.FindOne(context.TODO(), bson.M{"phone": userDetails.Phone}).Err()
		if err1 != nil {

			otp = helpers.Resend(userDetails.Email, userDetails.Phone)
			filter := bson.M{"_id": userDetails.ID}
			update := bson.M{"$set": bson.M{"otp": otp, "expiry": time.Now().Add(time.Minute * 5)}}
			query.UpdateOne("user", filter, update)

			response.Result = "Email sent successfully!!"
			json.NewEncoder(w).Encode(response)

		} else {
			response.Result = "Number already taken!!"
			json.NewEncoder(w).Encode(response)

		}

	} else {
		response.Result = "Email not registered!!"
		json.NewEncoder(w).Encode(response)
	}
	defer query.Endconn(client)

}

func CategoryUpload(w http.ResponseWriter, r *http.Request) {
	Check("categoryupload", "POST", w, r)
	NoofSub := r.FormValue("noofsub")
	Count, _ := strconv.Atoi(NoofSub)
	var res model.ResponseResult

	for i := 0; i <= Count; i++ {

		_, _, err := r.FormFile("img" + strconv.Itoa(i+1))
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "Image missing")
			return
		}
	}
	collection, client := query.Connection("category")
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":          "$_id",
				"categoryName": bson.M{"$first": "$categoryName"},
			},
		},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
		json.NewEncoder(w).Encode(res)
		return
	}
	var category model.Category
	var catlist []model.Category
	catlist = nil

	for cursor.Next(context.TODO()) {

		if err = cursor.Decode(&category); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		catlist = append(catlist, category)

	}
	cursor.Close(context.TODO())
	query.Endconn(client)
	//json.NewEncoder(w).Encode(catlist)

	Categoryname := r.FormValue("categoryname")
	for i := 0; i < len(catlist); i++ {
		if Categoryname == catlist[i].Catname {
			res.Error = "Name already exists"
			json.NewEncoder(w).Encode(res)
			return
		}
	}
	var Subcategoryname []string

	for i := 0; i < Count; i++ {
		Subcategoryname = append(Subcategoryname, r.FormValue("subcategoryname"+strconv.Itoa(i)))
	}

	var UploadCategory model.Category
	var UploadSubcategory model.Subcategory

	UploadCategory.Catname = Categoryname

	categoryresult := query.InsertOne("category", UploadCategory)
	var Sidarray []primitive.ObjectID
	Cid, _ := categoryresult.InsertedID.(primitive.ObjectID)
	for i := 0; i < len(Subcategoryname); i++ {
		UploadSubcategory.SubName = Subcategoryname[i]
		UploadSubcategory.CategoryId = Cid
		subcategoryresult := query.InsertOne("subcategory", UploadSubcategory)
		Sid, _ := subcategoryresult.InsertedID.(primitive.ObjectID)
		Sidarray = append(Sidarray, Sid)

	}
	helpers.CategoryImageUpload(w, r, Count, Subcategoryname, Cid, Sidarray)

}

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	Check("category", "GET", w, r)
	var res model.ResponseResult
	collection, client := query.Connection("category")
	cursor, err := collection.Find(context.TODO(), bson.M{"archived": false}, options.Find().SetSort(bson.M{"categoryName": 1}))
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var category model.Category
	var catlist []model.Category
	catlist = nil

	for cursor.Next(context.TODO()) {

		if err = cursor.Decode(&category); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		catlist = append(catlist, category)

	}
	json.NewEncoder(w).Encode(catlist)
	query.Endconn(client)
}

func CategoryUpdate(w http.ResponseWriter, r *http.Request) {
	Check("categoryupdate", "POST", w, r)

	Id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	Name := r.FormValue("name")
	if Name == "" {
		fmt.Fprintf(w, "Name Field is Empty")
		return
	}

	helpers.UpdateImage(w, r, Id, Name)
}

// swagger:route POST /tokengenerate User tokengeneration
// Allows users to get a new access token.
// responses:
//   200: tokengenerateResponse

// Allows users to generate a new access token
// swagger:response tokengenerateResponse
type tokengenerateResponseWrapper struct {
	// The generated response
	// in:header

}

// consumes:
// - application/json

// swagger:parameters tokengeneration
type tokengenerateParamsWrapper struct {
	// The Refresh token
	// in:header
	// type:string
	Token1 string
	// Privilege
	// in:header
	// type:string
	Admin string
}

// TokenGenerator : Generates a new access token
func TokenGenerator(w http.ResponseWriter, r *http.Request) {

	Check("tokengenerate", "POST", w, r)
	w.Header().Set("Content-Type", "application/json")
	helpers.LoadEnv()

	var mySigningKey = []byte(helpers.GetEnvWithKey("MY_SIGNING_KEY"))
	var adminKey = []byte(helpers.GetEnvWithKey("ADMIN_KEY"))

	admin := r.Header["Admin"]
	check := admin[len(admin)-1]

	var token *jwt.Token
	var err error
	var tokenacc string

	var data models.User

	if r.Header["Token1"] != nil {

		if check == "true" {

			token, err = jwt.Parse(r.Header["Token1"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return adminKey, nil
			})
		} else {

			token, err = jwt.Parse(r.Header["Token1"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

		}

		if err != nil {
			fmt.Fprintf(w, err.Error())
			fmt.Println("trello")
			fmt.Println(err.Error())

		}

		if token.Valid {
			claims, _ := token.Claims.(jwt.MapClaims)
			fmt.Println(claims["id"])

			str := fmt.Sprintf("%v", claims["id"])
			oid, _ := primitive.ObjectIDFromHex(str)

			collection, client := query.Connection("user")
			err = collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: oid}}).Decode(&data)
			defer query.Endconn(client)

			tokenacc, err = helpers.GenerateJWTAccess(check, data)

			fmt.Println(tokenacc)
			w.Header().Add("Token", tokenacc)
			//w.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
		}
	} else {
		fmt.Println("hiko")
		fmt.Fprintf(w, "Not Authorized")
	}

}

func LocationlistHandler(w http.ResponseWriter, r *http.Request) {
	Check("locationlist", "GET", w, r)
	collection, client := query.Connection("location")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	var location model.LocationTable
	var locationlist []model.LocationTable
	locationlist = nil

	for cursor.Next(context.TODO()) {

		if err = cursor.Decode(&location); err != nil {
			helpers.ErrHandler(err, w)
		}
		locationlist = append(locationlist, location)

	}
	cursor.Close(context.TODO())
	query.Endconn(client)
	json.NewEncoder(w).Encode(locationlist)
}

func CategoryDelete(w http.ResponseWriter, r *http.Request) {
	Check("categorydelete", "PATCH", w, r)

	Id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	From := r.FormValue("from")
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	if From == "category" {
		query.UpdateOne("category", bson.M{"_id": Id}, bson.M{"$set": bson.M{"archived": true}})
	} else if From == "subcategory" {
		query.UpdateOne("subcategory", bson.M{"_id": Id}, bson.M{"$set": bson.M{"archived": true}})
	}

}

func AdminProductSumInput(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	Check("adminProductInput", "POST", w, r)

	var res model.ResponseResult
	productCollection, client1, err := db.GetDBCollection("products")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return

	}

	userCollection, client2, err := db.GetDBCollection("user")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var decodprod model.StockData
	var productResultList []model.StockData
	productResultList = nil

	pipeline1 := []bson.M{
		{
			"$group": bson.M{
				"_id": "$_id",
			},
		},
	}
	productResult, err := productCollection.Aggregate(context.TODO(), pipeline1) //.Decode(&sd)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	for productResult.Next(context.TODO()) {
		if err = productResult.Decode(&decodprod); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		productResultList = append(productResultList, decodprod)
	}

	mappp := make(map[string]int) //,len(catIdResultList))
	for i := 0; i < len(productResultList); i++ {

		prodId := productResultList[i].Id.Hex()
		mappp[prodId] = 0 //append(mappp[prodId], 0) // sum of rent
	}
	fmt.Println(" Initial mapp is : ", mappp)

	var uResult model.User
	var uResultList []model.User
	uResultList = nil

	pipeline2 := []bson.M{

		{
			"$match": bson.M{
				"currentorder": bson.M{"$exists": true, "$not": bson.M{"$size": 0}},
			},
		},
		{

			"$unwind": "$currentorder",
		},
		{
			"$group": bson.M{
				"_id": "$_id",
				"currentorder": bson.M{
					"$push": "$currentorder",
				},
			},
		},
	}

	userResult, err := userCollection.Aggregate(context.TODO(), pipeline2) //.count() //.Decode(&user)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	defer userResult.Close(context.TODO())

	for userResult.Next(context.TODO()) {
		uResult.CurrentOrder = nil
		if err = userResult.Decode(&uResult); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		uResultList = append(uResultList, uResult)
	}
	fmt.Println("\nuResultList  : ", uResultList)
	var sumSlice []int
	sumSlice = nil
	var prodKey string
	var prodKeySlice []string
	prodKeySlice = nil
	var adminProd model.AdminReportRentOutput
	var adminProdList []model.AdminReportRentOutput
	adminProdList = nil

	var order model.Order

	fmt.Println("\nlength of uResultList", len(uResultList))
	for outer := 0; outer < len(uResultList); outer++ {
		fmt.Println("111")

		for inner := 0; inner < len(uResultList[outer].CurrentOrder); inner++ {

			err = query.FindoneID("orders", uResultList[outer].CurrentOrder[inner], "_id").Decode(&order)
			if err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			prodId := order.P_id.Hex()

			fmt.Println("prodId : ", prodId)
			//fmt.Println("rent ", uResultList[outer].CurrentOrder[inner].Rent, " * ", "count  ", uResultList[outer].CurrentOrder[inner].Count)

			mappp[prodId] = mappp[prodId] + (order.Rent * order.Count)
			fmt.Println("mappp[prodCatId][0]  ", mappp)
		}
	}

	for key, value := range mappp {
		prodKey = key
		sumSlice = append(sumSlice, value)
		prodKeySlice = append(prodKeySlice, prodKey) //contains key/id of category

	}

	fmt.Println("\nsumSlice[n]", sumSlice)

	adminProd.SumArray = make([]model.Sum, len(productResultList))
	for n := 0; n < len(productResultList); n++ {
		for m := 0; m < len(productResultList); m++ {
			if prodKeySlice[m] == productResultList[n].Id.Hex() {
				adminProd.SumArray[n].SumRent = sumSlice[m]
				adminProd.SumArray[n].ReportDate = time.Now()
				break
			}
		}
		adminProdList = append(adminProdList, adminProd)
	}

	for i := 0; i < len(adminProdList); i++ {

		fmt.Println("Document Updated of product  rent : ", i, ":::", adminProd.SumArray[i].SumRent)
		update := bson.M{"$push": bson.M{"sumarray": adminProd.SumArray[i]}}
		////update := bson.M{"$pull": bson.M{"sumarray": bson.M{"sumrent": 3250 & 0}}}
		_, err = productCollection.UpdateOne(context.TODO(), bson.M{"_id": productResultList[i].Id}, update)
		fmt.Println("Updated")
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

	}

	response := "Report has been generated for the month"

	json.NewEncoder(w).Encode(response)

	client1.Disconnect(context.TODO()) //products
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	client2.Disconnect(context.TODO()) //User
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}

//City Rent
func AdminCityRent(w http.ResponseWriter, r *http.Request) {

	Check("cityRent", "GET", w, r)

	w.Header().Set("Content-Type", "application/json")

	var res model.ResponseResult

	locationCollection, client1, err := db.GetDBCollection("location")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	locationResult, err := locationCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	defer locationResult.Close(context.TODO())

	var locn model.LocationTable

	var locationList []model.LocationTable
	locationList = nil
	for locationResult.Next(context.TODO()) {
		if err = locationResult.Decode(&locn); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		locationList = append(locationList, locn) // contains list of all location id's & Name

	}

	productCollection, client1, err := db.GetDBCollection("products")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var locationResponseList []model.LocationResponse
	for i := 0; i < len(locationList); i++ {

		fmt.Println("\nlocation : ", i, "  :", locationList[i])
		pipeline := []bson.M{

			{
				"$match": bson.M{
					"locationid": locationList[i].LocationID,
				},
			},

			{
				"$group": bson.M{
					"_id":        "_id",
					"locationid": bson.M{"$first": "$locationid"},
					//	"name":      bson.M{"$first": "$name"},
					// "piii": bson.M{
					// 	"$sum": "$price"},

					"cityRent": bson.M{"$sum": "$collected"}, //"total": bson.M{"$sum": "$quantity"},
				},
			},
		}

		cityRentResult, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		fmt.Println("cityResult ", cityRentResult)

		defer cityRentResult.Close(context.TODO())
		//var sd bson.M
		var locn2 model.LocationResponse

		for cityRentResult.Next(context.TODO()) {
			if err = cityRentResult.Decode(&locn2); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			locationResponseList = append(locationResponseList, locn2)
		}

	}
	fmt.Println("\nlocationResponseList ", locationResponseList)

	var locn3 model.LocationTable
	var locn3List []model.LocationTable
	locn3List = nil
	for i := 0; i < len(locationResponseList); i++ {

		locn3.CityName = locationList[i].CityName
		locn3.LocationID = locationList[i].LocationID
		locn3.CityRent = locationResponseList[i].CityRent
		locn3List = append(locn3List, locn3)
	}

	sort.SliceStable(locn3List, func(i, j int) bool {
		return locn3List[j].CityRent < locn3List[i].CityRent
	})

	json.NewEncoder(w).Encode(locn3List)

	client1.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}

//OverallRent Collected
func AdminOverallRentCollected(w http.ResponseWriter, r *http.Request) {

	Check("overallRent", "GET", w, r)

	w.Header().Set("Content-Type", "application/json")

	var res model.ResponseResult

	productCollection, client1, err := db.GetDBCollection("products")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	pipeline := []bson.M{

		{
			"$group": bson.M{
				"_id":         "_id",
				"overallRent": bson.M{"$sum": "$collected"}, //"total": bson.M{"$sum": "$quantity"},
			},
		},
	}

	overallRentResult, err := productCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var totalRent model.LocationResponse

	for overallRentResult.Next(context.TODO()) {
		if err = overallRentResult.Decode(&totalRent); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
	}
	fmt.Println("overallRentResult  ", overallRentResult)
	fmt.Println("totalRent", totalRent)

	json.NewEncoder(w).Encode(totalRent)

	client1.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}

//SubCategory Level Sum
func AdminSubCategoryLevelSum(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	Check("subCategoryLevelSum/"+params["location"], "GET", w, r)

	var res model.ResponseResult
	LocID, err := primitive.ObjectIDFromHex(params["location"])
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	//Accessing  Sub Category Collection
	subCategoryCollection, client2, err := db.GetDBCollection("subcategory")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	subCategoryResult, err := subCategoryCollection.Find(context.TODO(), bson.M{}) //bson.M{"categoryid": avgInput.CategoryID})
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	defer subCategoryResult.Close(context.TODO())

	var subCatResult model.Subcategory

	var subCatResultList []model.Subcategory
	subCatResultList = nil

	for subCategoryResult.Next(context.TODO()) {

		if err = subCategoryResult.Decode(&subCatResult); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		subCatResultList = append(subCatResultList, subCatResult) // contains list of all subCategory id's & Name

	}

	fmt.Println("\nsubCatResultList  :", subCatResultList)

	//Product Collection
	productCollection, client3, err := db.GetDBCollection("products")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	start, err := time.Parse(time.RFC3339, "2021-01-01T17:42:45.872Z")
	end := start.AddDate(1, 0, 0)
	fmt.Println("\nLocID :", LocID)
	fmt.Println("\nstart", start)
	fmt.Println("\nend", end)

	var subCatRentList []model.SubcategoryRentSum
	subCatRentList = nil
	for i := 0; i < len(subCatResultList); i++ {
		fmt.Println("subCatResultList[i].SubId  :", i, subCatResultList[i].SubId)
		pipeline := []bson.M{

			{
				"$match": bson.M{
					"locationid":    LocID,
					"subcategoryid": subCatResultList[i].SubId,
				},
			},

			{
				"$group": bson.M{
					"_id": "_id",
					//"LocationID": bson.M{"$first": "$locationid"},
					"SubCatSum": bson.M{"$sum": "$collected"},
					"NoOfProd":  bson.M{"$sum": 1},
				},
			},
		}

		subCatTotalRent, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		defer subCatTotalRent.Close(context.TODO())

		var subCatRent model.SubcategoryRentSum

		for subCatTotalRent.Next(context.TODO()) {
			if err = subCatTotalRent.Decode(&subCatRent); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}
			subCatRentList = append(subCatRentList, subCatRent)
		}
		fmt.Println("totalRent", subCatRentList)
	}

	var subCatSubOutput model.SubcategoryRentSum
	var subCatSubOutputList []model.SubcategoryRentSum
	for i := 0; i < len(subCatRentList); i++ {
		subCatSubOutput.SubCatName = subCatResultList[i].SubName
		subCatSubOutput.SubCatSum = subCatRentList[i].SubCatSum
		subCatSubOutput.NoOfProd = subCatRentList[i].NoOfProd
		subCatSubOutputList = append(subCatSubOutputList, subCatSubOutput)
	}

	fmt.Println("\ncatSubOUtPutList", subCatSubOutputList)

	json.NewEncoder(w).Encode(subCatSubOutputList)

	client2.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	client3.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}

//Category LevelSum
func AdminCategoryLevelSum(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	Check("categoryLevelSum/"+params["location"], "GET", w, r)

	var res model.ResponseResult
	LocID, err := primitive.ObjectIDFromHex(params["location"])
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	//Accessing Category Collection
	categoryCollection, client2, err := db.GetDBCollection("category")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	categoryResult, err := categoryCollection.Find(context.TODO(), bson.M{}) //bson.M{"categoryid": avgInput.CategoryID})
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	defer categoryResult.Close(context.TODO())

	var catResult model.Category

	var catResultList []model.Category
	catResultList = nil

	for categoryResult.Next(context.TODO()) {

		if err = categoryResult.Decode(&catResult); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		catResultList = append(catResultList, catResult) // contains list of all subCategory id's & Name

	}

	fmt.Println("\nCatResultList  :", catResultList)

	//Product Collection
	productCollection, client3, err := db.GetDBCollection("products")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	start, err := time.Parse(time.RFC3339, "2021-01-01T17:42:45.872Z")
	end := start.AddDate(1, 0, 0)
	fmt.Println("\nLocID :", LocID)
	fmt.Println("\nstart", start)
	fmt.Println("\nend", end)

	var catRentList []model.CategoryRentSum
	catRentList = nil
	for i := 0; i < len(catResultList); i++ {
		fmt.Println("catResultList[i].categoryID  :", i, catResultList[i].Catid)
		pipeline := []bson.M{

			{
				"$match": bson.M{
					"locationid": LocID,
					"categoryid": catResultList[i].Catid,
					//"createdat":     bson.M{"$gte": start, "$lte": end},
				},
			},

			{
				"$group": bson.M{
					"_id": "_id",
					//"LocationID": bson.M{"$first": "$locationid"},
					"CatSum":   bson.M{"$sum": "$collected"},
					"NoOfProd": bson.M{"$sum": 1},
				},
			},
		}

		catTotalRent, err := productCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		defer catTotalRent.Close(context.TODO())

		var catRent model.CategoryRentSum

		for catTotalRent.Next(context.TODO()) {
			if err = catTotalRent.Decode(&catRent); err != nil {
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			catRentList = append(catRentList, catRent)
		}

		fmt.Println("totalRent", catRentList)
	}

	var catOutput model.CategoryRentSum
	var catOutputList []model.CategoryRentSum
	for i := 0; i < len(catRentList); i++ {
		catOutput.CatName = catResultList[i].Catname
		catOutput.CatSum = catRentList[i].CatSum
		catOutput.NoOfProd = catRentList[i].NoOfProd
		catOutputList = append(catOutputList, catOutput)
	}

	fmt.Println("\ncatSubOUtPutList", catOutputList)

	json.NewEncoder(w).Encode(catOutputList)

	client2.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	client3.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

}

//Subcategory List
func AdminSubcategoryList(w http.ResponseWriter, r *http.Request) {

	Check("subcategoryList", "GET", w, r)

	w.Header().Set("Content-Type", "application/json")

	var res model.ResponseResult

	subCategoryCollection, client, err := db.GetDBCollection("subcategory")
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	pipeline := []bson.M{

		{
			"$group": bson.M{
				"_id":  "$_id",
				"name": bson.M{"$first": "$name"},
			},
		},
	}
	subCategoryResult, err := subCategoryCollection.Aggregate(context.TODO(), pipeline) //bson.M{"categoryid": avgInput.CategoryID})
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	defer subCategoryResult.Close(context.TODO())

	var subCatResult model.Subcategory

	var subCatResultList []model.Subcategory
	subCatResultList = nil

	for subCategoryResult.Next(context.TODO()) {

		if err = subCategoryResult.Decode(&subCatResult); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}
		subCatResultList = append(subCatResultList, subCatResult) // contains list of all subCategory id's & Name

	}

	fmt.Println("\nsubCatResultList  :", subCatResultList)

	json.NewEncoder(w).Encode(subCatResultList)

	client.Disconnect(context.TODO())
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}

func AdminSearchEngine(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	Check("adminSearchEngine/"+params["status"]+"/"+params["search"], "GET", w, r)

	var res model.ResponseResult

	collection, client := query.Connection("user")

	fmt.Println("Searching..............")

	fmt.Println("\nStatus is :", params["status"])
	fmt.Println("\nSearch query is : ", params["search"])
	_, err := strconv.Atoi(params["search"])
	if err == nil {
		params["search"] = "+91" + params["search"]
		fmt.Println("+91 added to the cellphone number", params["search"])
	}

	var projection bson.M
	if params["status"] == "intransit" {
		projection = bson.M{"currentorder": 0, "pastorder": 0, "isadmin": 0}
	} else if params["status"] == "currentorder" {
		projection = bson.M{"intransit": 0, "pastorder": 0, "isadmin": 0}
	} else if params["status"] == "pastorder" {
		projection = bson.M{"intransit": 0, "currentorder": 0, "isadmin": 0}
	}

	filter := bson.M{"$text": bson.M{"$search": params["search"]}}
	result, err := collection.Find(context.TODO(), filter, options.Find().SetProjection(projection)) //, abc)
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	var show []model.User
	var user model.User
	defer result.Close(context.TODO())
	for result.Next(context.TODO()) {
		user.CurrentOrder = nil
		user.InTransit = nil
		user.PastOrder = nil
		if err = result.Decode(&user); err != nil {
			res.Error = err.Error()
			json.NewEncoder(w).Encode(res)
			return
		}

		show = append(show, user)
	}

	json.NewEncoder(w).Encode(show)
	query.Endconn(client)
	fmt.Println("Search Complete")
}

func SubcategoryUploadhandler(w http.ResponseWriter, r *http.Request) {
	Check("subcategoryupload", "POST", w, r)
	Subname := r.FormValue("name")
	if r.FormValue("name") == " " {
		json.NewEncoder(w).Encode("Name field must be filled")
	}
	collection, client := query.Connection("subcategory")
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":  "$_id",
				"name": bson.M{"$first": "$name"},
			},
		},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	var subcategory model.Subcategory
	var subcatlist []model.Subcategory
	subcatlist = nil

	for cursor.Next(context.TODO()) {

		if err = cursor.Decode(&subcategory); err != nil {
			helpers.ErrHandler(err, w)
		}
		subcatlist = append(subcatlist, subcategory)

	}
	cursor.Close(context.TODO())
	query.Endconn(client)

	var res model.ResponseResult

	for i := 0; i < len(subcatlist); i++ {
		if Subname == subcatlist[i].SubName {
			res.Error = "Name already exists"
			json.NewEncoder(w).Encode(res)
			return
		}
	}
	_, _, err = r.FormFile("img")
	if err != nil {

		fmt.Fprintf(w, "No Image Inserted")
		return

	}
	var CatId primitive.ObjectID
	CatId, err = primitive.ObjectIDFromHex(r.FormValue("categoryid"))
	if err != nil {
		helpers.ErrHandler(err, w)
	}
	var upload model.Subcategory
	upload.Archived = false
	upload.SubName = Subname
	upload.CategoryId = CatId
	result := query.InsertOne("subcategory", upload)
	Sid := result.InsertedID.(primitive.ObjectID)
	helpers.SubcategoryUpdateImage(w, r, Sid)
}
