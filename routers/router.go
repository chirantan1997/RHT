package routers

import (
	controller "Newton/controllers"
	"Newton/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/account", controller.AccountHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth", controller.AuthHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/resend", controller.Resendotp).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/loginNew", controller.NewLoginHandler).Methods("POST")
	r.HandleFunc("/api/signupNew", controller.NewSignupHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/carousel", controller.Carousel).Methods("GET")

	r.HandleFunc("/api/productslist/{location}/{sub}/{page}", controller.ProductsList).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/usercreation", controller.UserCreationHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/wishlist", controller.WishlistHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/wishlistproducts", controller.WishlistProductsHandler).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/productdetails/{productid}", controller.ProductDetailsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/checkout", middleware.IsAuthorized(controller.CheckoutHandler)).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/searchengine", controller.SearchEngine).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/cartproducts", controller.CartProducts).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/removecartproduct", controller.RemoveCartProduct).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/cartfirsttime", controller.CartFirstTime).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/stock/{productid}", controller.ProductStock).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/cartupdate", controller.CartUpdate).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/stockcheck", controller.StockCheckHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/intransit", middleware.IsAuthorized(controller.IntransitHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/currentorder", middleware.IsAuthorized(controller.CurrentOrderHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/pastorder", middleware.IsAuthorized(controller.PastOrderHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/cancelorder", middleware.IsAuthorized(controller.CancelHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/payment", controller.PaymentHandler).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/editprofile", middleware.IsAuthorized(controller.ProfileHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/values", middleware.IsAuthorized(controller.ValueHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/details", controller.DetailsHandler).Methods("POST")
	r.HandleFunc("/api/stocker", controller.StockHandler).Methods("POST")
	r.HandleFunc("/api/delete", middleware.IsAuthorized(controller.DeleteHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/adminstock", middleware.IsAuthorized(controller.AdminStockHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/adminupdate", middleware.IsAuthorized(controller.AdminUpdateHandler)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/userlist/{status}", middleware.IsAuthorized(controller.UserList)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/categoryupload", middleware.IsAuthorized(controller.CategoryUpload)).Methods("POST")
	r.HandleFunc("/api/locationlist", controller.LocationlistHandler).Methods("GET")

	r.HandleFunc("/api/userreport", middleware.IsAuthorized(controller.UserReport)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/highdemanding/{location}", middleware.IsAuthorized(controller.AdminHighDemanding)).Methods("GET")
	r.HandleFunc("/api/lowdemanding/{location}", middleware.IsAuthorized(controller.AdminLowDemanding)).Methods("GET")
	r.HandleFunc("/api/highprofit/{location}", middleware.IsAuthorized(controller.AdminHighProfit)).Methods("GET")
	r.HandleFunc("/api/leastprofit/{location}", middleware.IsAuthorized(controller.AdminLeastProfit)).Methods("GET")
	r.HandleFunc("/api/categorylist", controller.CategoryListResponseHandler).Methods("GET")
	r.HandleFunc("/api/category", controller.CategoryHandler).Methods("GET")
	r.HandleFunc("/api/categoryupdate", middleware.IsAuthorized(controller.CategoryUpdate)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/paymentstatus", controller.PaymentStatusHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/numberchange", controller.NumberChangeHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/tokengenerate", controller.TokenGenerator).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/categorydelete", controller.CategoryDelete).Methods("PATCH")
	r.HandleFunc("/api/cityRent", controller.AdminCityRent).Methods("GET")
	r.HandleFunc("/api/overallRent", controller.AdminOverallRentCollected).Methods("GET")
	r.HandleFunc("/api/subCategoryLevelSum/{location}", controller.AdminSubCategoryLevelSum).Methods("GET")
	r.HandleFunc("/api/categoryLevelSum/{location}", controller.AdminCategoryLevelSum).Methods("GET")

	//r.HandleFunc("/api/subCategoryAverage/{location}/{start}/{end}", controller.AdminSubCategoryAverage).Methods("GET")
	//r.HandleFunc("/api/categoryAverage/{location}/{start}/{end}", controller.AdminCategoryAverage).Methods("GET")
	r.HandleFunc("/api/adminProductInput", controller.AdminProductSumInput).Methods("POST")

	r.HandleFunc("/api/subcategoryList", controller.AdminSubcategoryList).Methods("GET")

	r.HandleFunc("/api/subcategoryupload", controller.SubcategoryUploadhandler).Methods("POST")

	r.HandleFunc("/api/adminSearchEngine/{status}/{search}", controller.AdminSearchEngine).Methods("GET")

	return r
}
