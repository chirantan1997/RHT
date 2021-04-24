package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is ...
type User struct {
	ID           primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string               `json:"name,omitempty" bson:"name,omitempty"`
	Phone        string               `json:"phone,omitempty" bson:"phone,omitempty"`
	Email        string               `json:"email,omitempty" bson:"email,omitempty"`
	Address      string               `json:"address,omitempty" bson:"address,omitempty"`
	CurrentOrder []primitive.ObjectID `json:"currentorder,omitempty" bson:"currentorder,omitempty"`
	PastOrder    []primitive.ObjectID `json:"pastorder,omitempty" bson:"pastorder,omitempty"`
	InTransit    []primitive.ObjectID `json:"intransit,omitempty" bson:"intransit,omitempty"`

	Isadmin bool `json:"isadmin" bson:"isadmin"`
}

type Result struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Phone        string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	Address      string             `json:"address,omitempty" bson:"address,omitempty"`
	CurrentOrder []Order            `json:"currentorder,omitempty" bson:"currentorder,omitempty"`
	PastOrder    []Order            `json:"pastorder,omitempty" bson:"pastorder,omitempty"`
	InTransit    []Order            `json:"intransit,omitempty" bson:"intransit,omitempty"`
}

//login...

type Login struct {
	Contact string `json:"contact"`
}

// ResponseResult is ...
type ResponseResult struct {
	Error  string `json:"error,omitempty"`
	Status int    `json:"status,omitempty"`
	Result string `json:"result,omitempty"`
}

// OtpContainer ...
type OtpContainer struct {
	OtpEntered   string             `json:"otpentered,omitempty" bson:"otp,omitempty"`
	Number       string             `json:"phone,omitempty" bson:"phone,omitempty"`
	From         string             `json:"from,omitempty" bson:"from,omitempty"`
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	Address      string             `json:"address,omitempty" bson:"address,omitempty"`
	CurrentOrder []Product          `json:"currentorder,omitempty" bson:"currentorder,omitempty"`
	PastOrder    []Product          `json:"pastorder,omitempty" bson:"pastorder,omitempty"`
	InTransit    []Product          `json:"intransit,omitempty" bson:"intransit,omitempty"`
}

type Carousel struct {
	Carousel [7]string `json:"carousel"`
}

type Id struct {
	ID1      primitive.ObjectID `json:"id" bson:"id"`
	UserId   primitive.ObjectID `json:"userid" bson:"userid"`
	PID      primitive.ObjectID `josn:"pid,omitempty" bson:"pid,omitempty"`
	Duration int                `json:"duration,omitempty" bson:"duration,omitempty"`
	Sub      primitive.ObjectID `json:"sub,omitempty" bson:"sub,omitempty"`
	From     string             `json:"from,omitempty" bson:"from,omitempty"`
	Exist    bool               `json:"exist,omitempty" bson:"exist,omitempty"`
	Date     time.Time          `json:"date,omitempty" bson:"date,omitempty"`
}
type Items struct {
	Id            primitive.ObjectID   `json:"_id" bson:"_id"`
	Subcategoryid primitive.ObjectID   `json:"subcategoryid,omitempty" bson:"subcategoryid,omitempty"`
	Name          string               `json:"name,omitempty" bson:"name,omitempty"`
	Img           []string             `json:"img,omitempty" bson:"img,omitempty"`
	Details       string               `json:"details,omitempty" bson:"details,omitempty"`
	Price         int                  `json:"price,omitempty" bson:"price,omitempty"`
	Rent          int                  `json:"rent,omitempty" bson:"rent,omitempty"`
	Duration      int                  `json:"duration,omitempty" bson:"duration,omitempty"`
	Itemsid       []primitive.ObjectID `json:"itemsid,omitempty" bson:"itemsid,omitempty"`
	LocationID    primitive.ObjectID   `json:"locationid,omitempty" bson:"locationid,omitempty"`
	Stock         int                  `json:"stock" bson:"stock"`
	Deposit       int                  `json:"deposit,omitempty" bson:"deposit,omitempty"`
	Demand        int                  `json:"demand,omitempty" bson:"demand,omitempty"`
}

type Wishlist struct {
	Userid  primitive.ObjectID   `json:"userid" bson:"userid"`
	ItemsId []primitive.ObjectID `json:"itemsId,omitempty" bson:"itemsId,omitempty"`
}

type Wishlistarray struct {
	Wisharr []primitive.ObjectID `json:"itemsid" bson:"itemsid"`
}

type Account struct {
	ID    primitive.ObjectID `json:"id"`
	Exist bool               `json:"exist"`
}

type Cart struct {
	Userid  primitive.ObjectID `json:"userid,omitempty" bson:"userid,omitempty"`
	Product []Product          `json:"product,omitempty" bson:"product,omitempty"`
	Order   []CartOrder        `json:"order,omitempty" bson:"order,omitempty"`
}

type Product struct {
	P_id     primitive.ObjectID `json:"p_id,omitempty" bson:"p_id,omitempty"`
	Img      string             `json:"img,omitempty" bson:"img,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Count    int                `json:"count,omitempty" bson:"count,omitempty"`
	Duration int                `json:"duration,omitempty" bson:"duration,omitempty"`
	Rent     int                `json:"_rent,omitempty" bson:"_rent,omitempty"`
	Deposit  int                `json:"deposit,omitempty" bson:"deposit,omitempty"`
}
type CartOrder struct {
	P_id     primitive.ObjectID `json:"p_id,omitempty" bson:"p_id,omitempty"`
	Pstr     string             `json:"pstr,omitempty" bson:"pstr,omitempty"`
	Img      string             `json:"img,omitempty" bson:"img,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Count    int                `json:"count,omitempty" bson:"count,omitempty"`
	Duration int                `json:"duration,omitempty" bson:"duration,omitempty"`
	Rent     int                `json:"_rent,omitempty" bson:"_rent,omitempty"`
	Deposit  int                `json:"deposit,omitempty" bson:"deposit,omitempty"`
}
type Order struct {
	ID           primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	P_id         primitive.ObjectID   `json:"p_id,omitempty" bson:"p_id,omitempty"`
	Img          string               `json:"img,omitempty" bson:"img,omitempty" `
	Name         string               `json:"name,omitempty" bson:"name,omitempty"`
	Count        int                  `json:"count,omitempty" bson:"count,omitempty"`
	Items_count  []primitive.ObjectID `json:"items_count,omitempty" bson:"items_count,omitempty"`
	Rent         int                  `json:"_rent,omitempty" bson:"_rent,omitempty"`
	Duration     int                  `json:"duration,omitempty" bson:"duration,omitempty"`
	Date         time.Time            `json:"checkoutdate,omitempty" bson:"checkoutdate,omitempty"`
	IsCancellled bool                 `json:"iscancelled" bson:"iscancelled"`
	Due          int                  `json:"due" bson:"due"`
	PayDates     []time.Time          `json:"paydates,omitempty" bson:"paydates"`
	Deposit      int                  `json:"deposit,omitempty" bson:"deposit,omitempty"`
}

type Cartproduct struct {
	Status    bool               `json:"status" bson:"status"`
	Userid    primitive.ObjectID `json:"userid" bson:"userid"`
	Productid primitive.ObjectID `json:"productid" bson:"productid"`
}

//Id ...
type CartContainer struct {
	UserID primitive.ObjectID `json:"UserID" bson:"UserID"`
	ItemID primitive.ObjectID `json:"ItemID" bson:"ItemID"`
	Status bool               `json:"Status" bson:"Status"`
}

type Total struct {
	CartID primitive.ObjectID `json:"CartID" bson:"CartID"`
	UserID primitive.ObjectID `json:"UserID,omitempty" bson:"UserID,omitempty"`
	Rent   int                `json:"Rent" bson:"Rent"`
}

//SearchProduct ...
type SearchProduct struct {
	Search     string             `json:"Search,omitempty" bson:"Search,omitempty"`
	LocationID primitive.ObjectID `json:"locationid,omitempty" bson:"locationid,omitempty"`
}

type CartInput struct {
	Userid  primitive.ObjectID `json:"userid,omitempty" bson:"userid,omitempty"`
	Value   int                `json:"value,omitempty" bson:"value,omitempty"`
	Status  bool               `json:"status,omitempty" bson:"status,omitempty"`
	Product Product            `json:"product,omitempty" bson:"product,omitempty"`
	Order   CartOrder          `json:"order,omitempty" bson:"order,omitempty"`
}

type RemoveCartProduct struct {
	UserId    primitive.ObjectID `json:"userid" bson:"userid"`
	ProductId primitive.ObjectID `json:"p_id" bson:"p_id"`

	Count    int `json:"count" bson:"count"`
	Duration int `json:"duration" bson:"duration"`
	Rent     int `json:"_rent,omitempty" bson:"_rent,omitempty"`
}

//ProductStock

type StockId struct {
	ProductId primitive.ObjectID `json:"_id" bson:"_id"`
}

type StockData struct {
	Id            primitive.ObjectID   `json:"_id" bson:"_id"`
	Subcategoryid primitive.ObjectID   `json:"subcategoryid" bson:"subcategoryid"`
	Name          string               `json:"name" bson:"name"`
	Img           []string             `json:"img" bson:"img"`
	Details       string               `json:"details" bson:"details"`
	Price         int                  `json:"price" bson:"price"`
	Rent          int                  `json:"rent" bson:"rent"`
	Duration      int                  `json:"duration" bson:"duration"`
	Itemsid       []primitive.ObjectID `json:"itemsid,omitempty" bson:"itemsid,omitempty"`
	LocationID    primitive.ObjectID   `json:"locationid" bson:"locationid"`
	Stock         int                  `json:"stock" bson:"stock"`
}

//Newlogin
type Newlogin struct {
	Userid  primitive.ObjectID `json:"userid,omitempty" bson:"userid,omitempty"`
	Contact string             `json:"contact" bson:"contact"`
	Name    string             `json:"name,omitempty" default:"guest" bson:"name,omitempty"`
	Address string             `json:"address,omitempty" bson:"address,omitempty"`
	Email   string             `json:"email,omitempty" bson:"email,omitempty"`
	Isadmin bool               `json:"isadmin" bson:"isadmin"`
}

//----------------------------------------------------------------------------
// ProductUpload ...
type ProductUpload struct {
	Id            primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Subcategoryid primitive.ObjectID   `json:"subcategoryid,omitempty" bson:"subcategoryid,omitempty"`
	Name          string               `json:"name,omitempty" bson:"name,omitempty"`
	Img           []string             `json:"img,omitempty" bson:"img,omitempty"`
	Details       string               `json:"details,omitempty" bson:"details,omitempty"`
	Price         int                  `json:"price,omitempty" bson:"price,omitempty"`
	Rent          int                  `json:"rent,omitempty" bson:"rent,omitempty"`
	Duration      int                  `json:"duration,omitempty" bson:"duration,omitempty"`
	Itemsid       []primitive.ObjectID `json:"itemsid,omitempty" bson:"itemsid,omitempty"`
	LocationID    primitive.ObjectID   `json:"locationid,omitempty" bson:"locationid,omitempty"`
	Stock         int                  `json:"stock,omitempty" bson:"stock,omitempty"`
	Deposit       int                  `json:"deposit,omitempty" bson:"deposit,omitempty"`
	Demand        int                  `json:"demand" bson:"demand"`
	Createdat     time.Time            `json:"createdat,omitempty" bson:"createdat,omitempty"`
}

// ProductDetails ...
type ProductDetails struct {
	Subcategoryid string `json:"subcategoryid" bson:"subcategoryid"`
	Name          string `json:"name" bson:"name"`
	Details       string `json:"details" bson:"details"`
	Price         string `json:"price" bson:"price"`
	Rent          string `json:"rent" bson:"rent"`
	Deposit       string `json:"deposit" bson:"deposit"`
	LocationID    string `json:"locationid" bson:"locationid"`
}

// ProductItems ...
type ProductItems struct {
	Productid primitive.ObjectID `json:"productid" bson:"productid"`
	Createdat time.Time          `json:"createdat,omitempty" bson:"createdat,omitempty"`
}

// ProductStock ...
type ProductStock struct {
	Productid string `json:"productid" bson:"productid"`
	Quantity  string `json:"quantity" bson:"quantity"`
}

// Delete ...
type Delete struct {
	Productid primitive.ObjectID `json:"productid" bson:"productid"`
}
type Deleteditems struct {
	Itemsid []primitive.ObjectID `json:"iemsid" bson:"itemsid"`
}

// ProductUpdate ...
type ProductUpdate struct {
	Productid primitive.ObjectID `json:"id" bson:"id"`
	Name      string             `json:"name" bson:"name"`
	Details   string             `json:"details" bson:"details"`
	Price     int                `json:"price" bson:"price"`
	Rent      int                `json:"rent" bson:"rent"`
	Deposit   int                `json:"deposit" bson:"deposit"`
}

//------------------------------------------------------------------
type UserReportId struct {
	Id primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	//Index           int                `json:"index,omitempty"`
	CategoryID      primitive.ObjectID `json:"categoryid,omitempty" bson:"categoryid,omitempty"`
	CategoryName    string             `json:"categoryname,omitempty" bson:"categoryname,omitempty"`
	SubCategoryID   primitive.ObjectID `json:"subcategoryid,omitempty" bson:"subcategoryid,omitempty"`
	SubCategoryName string             `json:"name,omitempty" bson:"name,omitempty"`
}

type UserReport struct {
	Name            string    `json:"name"`
	Paid            int       `json:"paid"`
	Payable         int       `json:"payable"`
	TotalRent       int       `json:"totalRent"`
	Profit          int       `json:"profit"`
	TotalPrice      int       `json:"totalPrice"`
	LastPaymentDate time.Time `json:"lastpaymentdate"`
	NextPaymentDate time.Time `json:"nextpaymentdate"`
	Due             int       `json:"due"`
}

type List struct {
	Page       string  `json:"page,omitempty" bson:"page,omitempty"`
	IsLastpage string  `json:"islastpage,omitempty" bson:"islastpage,omitempty"`
	Products   []Items `json:"products,omitempty" bson:"products,omitempty"`
}

type OtpCreds struct {
	Otp    string    `json:"otp" bson:"otp"`
	Expiry time.Time `json:"expiry,omitempty" bson:"expiry,omitempty"`
}

type Category struct {
	Catid    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Catname  string             `json:"categoryName,omitempty" bson:"categoryName,omitempty"`
	ImageUrl string             `json:"img,omitempty" bson:"img,omitempty"`
	Archived bool               `json:"archived,omitempty" bson:"archived"`
}

type Subcategory struct {
	SubId      primitive.ObjectID `json:"subid,omitempty" bson:"_id,omitempty"`
	CategoryId primitive.ObjectID `json:"categoryid,omitempty" bson:"categoryid,omitempty"`
	SubName    string             `json:"name,omitempty" bson:"name,omitempty"`
	ImageUrl   string             `json:"img,omitempty" bson:"img,omitempty"`
	Archived   bool               `json:"archived,omitempty" bson:"archived"`
}

type CatResponse struct {
	CategoryName string        `json:"category,omitempty" bson:"category,omitempty"`
	SubArray     []Subcategory `json:"subcategory,omitempty" bson:"subcategory,omitempty"`
}

type PaymentStatus struct {
	LastPaymentDate time.Time `json:"lastpaymentdate"`
	NextPaymentDate time.Time `json:"nextpaymentdate"`
	Due             int       `json:"due"`
}

type AdminProfitLossOutput struct {
	Name      string    `json:"name" bson:"name"`
	Price     int       `json:"price" bson:"price"`
	Collected int       `json:"collected" bson:"collected"`
	CreatedAt time.Time `json:"createdat" bson:"createdat"`
}

type AdminHighLowDemanding struct {
	Name      string    `json:"name" bson:"name"`
	Demand    int       `json:"demand" bson:"demand"`
	CreatedAt time.Time `json:"createdat" bson:"createdat"`
}

type AdminProductInput struct {
	ProductId   primitive.ObjectID `json:"p_id,omitempty" bson:"p_id,omitempty"`
	LocationID  primitive.ObjectID `json:"locationid,omitempty" bson:"locationid,omitempty"`
	ProductName string             `json:"p_name,omitempty" bson:"p_name,omitempty"`
	SumArray    Sum                `json:"sumarray,omitempty" bson:"sumarray,omitempty"`
}

type AdminProductOutput struct {
	ProductId   primitive.ObjectID `json:"p_id" bson:"p_id"`
	LocationID  primitive.ObjectID `json:"locationid,omitempty" bson:"locationid,omitempty"`
	ProductName string             `json:"p_name" bson:"p_name"`
	SumArray    []Sum              `json:"sumarray" bson:"sumarray"`
}

type Sum struct {
	SumRent    int       `json:"sumrent" bson:"sumrent"`
	ReportDate time.Time `json:"reportdate" bson:"reportdate"`
}

type location struct {
	LocationId primitive.ObjectID `json:"locationid,omitempty" bson:"_id"`
	CityName   string             `json:"city,omitempty" bson:"city,omitempty"`
	Img        string             `json:"img,omitempty" bson:"img,omitempty"`
}

type AdminReportRentOutput struct {
	SumArray []Sum `json:"sumarray" bson:"sumarray"`
}

type SubcategoryRentSum struct {
	SubCatName string `json:"subCatName"`
	SubCatSum  int    `json:"subCatSum"`
	NoOfProd   int    `json:"noOfProducts"`
}

type CategoryRentSum struct {
	CatName  string `json:"catName"`
	CatSum   int    `json:"catSum"`
	NoOfProd int    `json:"noOfProducts"`
}

type LocationTable struct {
	LocationID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CityName    string             `json:"city,omitempty" bson:"city,omitempty"`
	CityRent    int                `json:"cityrent,omitempty"`
	OverallRent uint64             `json:"overallRent,omitempty"`
}

type LocationResponse struct {
	CityName    string `json:"city,omitempty" bson:"city,omitempty"`
	CityRent    int    `json:"cityrent,omitempty"`
	OverallRent int    `json:"overallRent,omitempty"`
}
