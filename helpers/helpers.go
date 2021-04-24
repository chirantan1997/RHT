package helpers

import (
	"Newton/models"
	"Newton/query"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dgrijalva/jwt-go"
	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/sulat"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mySigningKey = []byte(GetEnvWithKey("MY_SIGNING_KEY"))
var adminKey = []byte(GetEnvWithKey("ADMIN_KEY"))

//GetEnvWithKey : get env value
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

//LoadEnv : loading the env file
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

// GenerateJWT ...
func GenerateJWTAccess(str string, user models.User) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * 3).Unix()
	claims["id"] = user.ID
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["phone"] = user.Phone
	claims["isadmin"] = user.Isadmin
	claims["address"] = user.Address

	if str == "true" {
		tokenString, err := token.SignedString(adminKey)
		if err != nil {
			fmt.Println("Something went wrong")
		}
		fmt.Println(tokenString)
		return tokenString, err
	} else {
		tokenString, err := token.SignedString(mySigningKey)
		if err != nil {
			fmt.Println("Something went wrong")
		}
		fmt.Println(tokenString)
		return tokenString, err
	}
}

func GenerateJWTRefresh(str string, user models.User) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().AddDate(0, 0, 25).Unix()
	claims["id"] = user.ID
	claims["isadmin"] = user.Isadmin

	if str == "true" {
		tokenString, err := token.SignedString(adminKey)
		if err != nil {
			fmt.Println("Something went wrong")
		}
		fmt.Println(tokenString)
		return tokenString, err
	} else {
		tokenString, err := token.SignedString(mySigningKey)
		if err != nil {
			fmt.Println("Something went wrong")
		}
		fmt.Println(tokenString)
		return tokenString, err
	}
}

// Reverse ...
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ProductUploadHandler : handles the product upload
func ProductImageHandler(w http.ResponseWriter, r *http.Request, id primitive.ObjectID, sub string, cat string, pos string, from string) {

	w.Header().Set("Content-Type", "application/json")

	var imageName []string
	var imageURL []string
	imageName = nil
	imageURL = nil
	key := []string{
		"img1",
		"img2",
		"img3",
	}

	var count int
	var k int = 0
	var resarr []int
	pos = Reverse(pos)
	num, _ := strconv.Atoi(pos)

	for i := 0; i < len(pos); i++ {
		fmt.Println("ho")
		resarr = append(resarr, num%10)
		num = num / 10
	}

	for i := 0; i < 3; i++ {

		_, _, err := r.FormFile(key[i])
		if err != nil {
			break
		}
		count++
	}

	LoadEnv()
	awsAccessKeyID := GetEnvWithKey("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := GetEnvWithKey("AWS_SECRET_ACCESS_KEY")

	for i := 0; i < count; i++ {

		file, fileHeader, err := r.FormFile(key[i])
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "Could not get uploaded file")
			return
		}
		imageName = append(imageName, fileHeader.Filename)
		newString := cat + "/" + sub + "/" + imageName[i]

		s, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				awsAccessKeyID,     // id
				awsSecretAccessKey, // secret
				""),                // token can be left blank for now
		})
		if err != nil {
			fmt.Fprintf(w, "Could not upload file")
		}
		uploader := s3manager.NewUploader(s)

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:             aws.String("ckrht"),
			ACL:                aws.String("public-read"),
			Key:                aws.String(newString),
			Body:               file,
			ContentType:        aws.String("image/jpeg"),
			ContentDisposition: aws.String("inline; filename=" + fmt.Sprintf("%s", imageName[i])),
		})
		if err != nil {
			fmt.Printf("failed to upload file, %v", err)
			return
		}
		imageURL = append(imageURL, aws.StringValue(&result.Location))
		file.Close()
	}
	if from == "creation" {
		for i := 0; i < len(imageURL); i++ {

			filter := bson.M{"_id": id}
			update := bson.M{"$push": bson.M{"img": imageURL[i]}}

			query.UpdateOne("products", filter, update)
		}
	} else if from == "updation" {
		for i := 0; i < len(resarr); i++ {

			if resarr[i] == 1 {
				key := "img." + strconv.Itoa(i)
				filter := bson.M{"_id": id}
				update := bson.M{"$set": bson.M{key: imageURL[k]}}
				query.UpdateOne("products", filter, update)
				k++
			}
		}

	}

}

type Context struct {
	Otp string
}

var FullHTML = `<div width="100%" style="background: #f8f8f8; padding: 0px 0px; font-family: arial; line-height: 100%; height: 100%; width: 100%; color: #514d6a;">
<div style="max-width: 700px; padding: 0px 0; margin: 0px auto; font-size: 14px;">
	<table border="0" cellpadding="0" cellspacing="0" style="width: 99.2711%; margin-bottom: 0px;" height="193">
		<tbody>
			<tr>
				<td style="vertical-align: top; padding-bottom: 0px; width: 100%;" align="center">
					<img style="-webkit-user-select: none; margin: auto; cursor: zoom-in; background-color: hsl(0, 0%, 90%); transition: background-color 300ms;" src="https://ckrht.s3.amazonaws.com/RHT.png" width="399" height="193" />
				</td>
			</tr>
		</tbody>
	</table>
	<div style="padding: 20px; background: #fff;">
		<table border="0" cellpadding="0" cellspacing="0" style="width: 98.9116%;" height="121">
			<tbody>
				<tr>
					<td style="border-bottom: 0px solid #f6f6f6;">
						<h1 style="font-size: 14px; font-family: arial; margin: 0px; font-weight: bold;">Greetings of the day  !</h1>
					</td>
				</tr>
				<tr>
					<td style="padding: 10px 0 0px 0;">
						<p>A request to reset your phone number has been made. If you did not make this request, simply ignore this email. If you did make this request, please reset your number.</p>
						<p>OTP for new number registration is <b style="background-color: transparent;">- {{.Otp}}</b>. </p>
						<p><b style="background-color: transparent;">- Thanks (RHT Team)</b>
						</p>
					</td>
				</tr>
				<tr>
					<td style="border-top: 1px solid #f6f6f6; padding-top: 20px; color: #777;">If you continue to have problems, please feel free to contact us at <a href="mailto:support@RHT.com">support@RHT.com</a>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
	<div style="text-align: center; font-size: 12px; color: #b2b2b5; margin-top: 20px;">
		<p>Powered by RHT.com</p>
	</div>
</div>
</div>`
var SGC = sulat.SGC{}

func SendGrid() {
	LoadEnv()
	SGC = sulat.SGC{
		SendGridAPIKey:   GetEnvWithKey("SEND_GRID_API_KEY"),
		SendGridEndPoint: GetEnvWithKey("SEND_GRID_END_POINT"),
		SendGridHost:     GetEnvWithKey("SEND_GRID_HOST"),
	}
}

func Resend(email string, phone string) (OTP string) {
	SendGrid()
	fmt.Println("Hello USER!")

	max := 9999
	min := 1000
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(rand.Intn(max-min+1) + min)

	templates := template.New("template")
	templates.New("FullHTML").Parse(FullHTML)
	var tpl bytes.Buffer
	context := Context{
		Otp: otp,
	}
	templates.Lookup("FullHTML").Execute(&tpl, context)
	LoadEnv()

	mailOpt := &sulat.SendMail{
		Subject: "RHT Support",
		From:    sulat.NewEmail("Team RHT", GetEnvWithKey("SENDER_EMAIL")),
		To:      sulat.NewEmail("RHT User", email),
	}

	htmlContent, err := sulat.SetHTML(&sulat.EmailHTMLFormat{
		IsFullHTML:       true,
		FullHTMLTemplate: tpl.String(),
	})

	isSend, err := sulat.SendEmailSG(mailOpt, htmlContent, &SGC)
	if err != nil {
		itrlog.Fatal(err)
	}
	fmt.Println("is sent: ", isSend)
	return otp
}

func CategoryImageUpload(w http.ResponseWriter, r *http.Request, count int, SubCategoryName []string, Cid primitive.ObjectID, Sidarray []primitive.ObjectID) {

	w.Header().Set("Content-Type", "application/json")

	var imageName []string
	var imageURL []string
	imageName = nil
	imageURL = nil
	var key []string
	for i := 1; i <= count+1; i++ {
		key = append(key, "img"+strconv.Itoa(i))
	}

	var newString string
	LoadEnv()
	awsAccessKeyID := GetEnvWithKey("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := GetEnvWithKey("AWS_SECRET_ACCESS_KEY")

	for i := 0; i <= count; i++ {

		file, fileHeader, err := r.FormFile(key[i])
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "Could not get uploaded file")
			return
		}
		imageName = append(imageName, fileHeader.Filename)
		if i == 0 {
			newString = "category" + "/" + imageName[i]
		} else {

			newString = "subcategory" + "/" + imageName[i]

		}

		s, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				awsAccessKeyID,     // id
				awsSecretAccessKey, // secret
				""),                // token can be left blank for now
		})
		if err != nil {
			fmt.Fprintf(w, "Could not upload file")
		}
		uploader := s3manager.NewUploader(s)

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:             aws.String("ckrht"),
			ACL:                aws.String("public-read"),
			Key:                aws.String(newString),
			Body:               file,
			ContentType:        aws.String("image/jpeg"),
			ContentDisposition: aws.String("inline; filename=" + fmt.Sprintf("%s", imageName[i])),
		})
		if err != nil {
			fmt.Printf("failed to upload file, %v", err)
			return
		}
		imageURL = append(imageURL, aws.StringValue(&result.Location))
		file.Close()
	}
	for i := 0; i < len(imageURL); i++ {
		if i == 0 {
			filter := bson.M{"_id": Cid}
			update := bson.M{"$set": bson.M{"img": imageURL[i], "archived": false}}

			query.UpdateOne("category", filter, update)
		} else {
			filter := bson.M{"_id": Sidarray[i-1]}
			update := bson.M{"$set": bson.M{"img": imageURL[i], "archived": false}}

			query.UpdateOne("subcategory", filter, update)

		}

	}

}

func UpdateImage(w http.ResponseWriter, r *http.Request, Id primitive.ObjectID, name string) {
	w.Header().Set("Content-Type", "application/json")
	flag := false
	var imageName string
	var imageURL string
	key := "img"
	From := r.FormValue("from")
	var newString string
	LoadEnv()
	awsAccessKeyID := GetEnvWithKey("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := GetEnvWithKey("AWS_SECRET_ACCESS_KEY")

	file, fileHeader, err := r.FormFile(key)
	if err != nil {

		//fmt.Fprintf(w, "No Image Inserted")
		flag = true

	}
	if flag == false {
		imageName = fileHeader.Filename
		if From == "category" {
			newString = "category" + "/" + imageName
		} else if From == "subcategory" {
			newString = "subcategory" + "/" + imageName
		}

		s, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				awsAccessKeyID,     // id
				awsSecretAccessKey, // secret
				""),                // token can be left blank for now
		})
		if err != nil {
			fmt.Fprintf(w, "Could not upload file")
		}
		uploader := s3manager.NewUploader(s)

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:             aws.String("ckrht"),
			ACL:                aws.String("public-read"),
			Key:                aws.String(newString),
			Body:               file,
			ContentType:        aws.String("image/jpeg"),
			ContentDisposition: aws.String("inline; filename=" + fmt.Sprintf("%s", imageName)),
		})
		if err != nil {
			fmt.Printf("failed to upload file, %v", err)
			return
		}

		imageURL = aws.StringValue(&result.Location)
		file.Close()
	}

	filter := bson.M{"_id": Id}
	if From == "category" && flag == false {

		update := bson.M{"$set": bson.M{"img": imageURL, "categoryName": name}}

		query.UpdateOne("category", filter, update)
	} else if From == "subcategory" && flag == false {
		update := bson.M{"$set": bson.M{"img": imageURL, "name": name}}

		query.UpdateOne("subcategory", filter, update)
	} else if From == "category" && flag == true {
		update := bson.M{"$set": bson.M{"categoryName": name}}

		query.UpdateOne("category", filter, update)
	} else if From == "subcategory" && flag == true {
		update := bson.M{"$set": bson.M{"name": name}}

		query.UpdateOne("subcategory", filter, update)
	}

}

func ErrHandler(err error, w http.ResponseWriter) {
	var res models.ResponseResult
	res.Error = err.Error()
	json.NewEncoder(w).Encode(res)
	return
}

func SubcategoryUpdateImage(w http.ResponseWriter, r *http.Request, Id primitive.ObjectID) {
	w.Header().Set("Content-Type", "application/json")

	var imageName string
	var imageURL string

	var newString string
	LoadEnv()
	awsAccessKeyID := GetEnvWithKey("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := GetEnvWithKey("AWS_SECRET_ACCESS_KEY")

	file, fileHeader, err := r.FormFile("img")
	if err != nil {

		ErrHandler(err, w)

	}
	imageName = fileHeader.Filename
	newString = "subcategory" + "/" + imageName

	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKeyID,     // id
			awsSecretAccessKey, // secret
			""),                // token can be left blank for now
	})
	if err != nil {
		fmt.Fprintf(w, "Could not upload file")
	}
	uploader := s3manager.NewUploader(s)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String("ckrht"),
		ACL:                aws.String("public-read"),
		Key:                aws.String(newString),
		Body:               file,
		ContentType:        aws.String("image/jpeg"),
		ContentDisposition: aws.String("inline; filename=" + fmt.Sprintf("%s", imageName)),
	})
	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
		return
	}

	imageURL = aws.StringValue(&result.Location)
	file.Close()

	filter := bson.M{"_id": Id}
	update := bson.M{"$set": bson.M{"img": imageURL}}

	query.UpdateOne("subcategory", filter, update)
}
