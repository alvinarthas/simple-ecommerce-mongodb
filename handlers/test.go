package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/alvinarthas/simple-ecommerce-mongodb/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"

	"github.com/alvinarthas/simple-ecommerce-mongodb/config"
	"github.com/gin-gonic/gin"
)

/*
	Just a Test Handler, to try some features
*/

// Student sadas
type Student struct {
	Name         string  `json:"name"`
	Age          int64   `json:"age"`
	AverageScore float64 `json:"average_score"`
}

// TestFuncxx only for testing
func TestFuncxx(c *gin.Context) {
	var err error
	// INSERT

	// store := models.Store{
	// 	ID:       primitive.NewObjectID(),
	// 	Name:     "Toko Bagus",
	// 	UserName: "bagus_toko",
	// 	Adress:   "Cipa Cipa",
	// 	Email:    "tokok@bagus.com",
	// 	Phone:    "password",
	// 	Avatar:   "https://avatars3.githubusercontent.com/u/28726294?v=4",
	// }

	// users := models.User{
	// 	ID:          primitive.NewObjectID(),
	// 	UserName:    "cipa_)ipa",
	// 	FullName:    "Cipa Cipa",
	// 	Email:       "cipa@alvin.com",
	// 	Password:    "password",
	// 	SocialID:    "2872629224",
	// 	Provider:    "github",
	// 	Avatar:      "https://avatars3.githubusercontent.com/u/28726294?v=4",
	// 	CreatedDate: time.Now(),
	// 	LastUpdate:  primitive.Timestamp{T: uint32(time.Now().Unix())},
	// 	Store:       store,
	// }

	// collection := config.DB.Collection("users")

	// // INSERT
	// insertResult, err := collection.InsertOne(config.CTX, users)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// FIND ONE

	// filter := bson.M{
	// 	"store.verification_token": "axax000111222",
	// }

	// var users models.User
	// err = collection.FindOne(config.CTX, filter).Decode(&users)

	// if err != nil {
	// 	fmt.Println("Error calling FindOne():", err)
	// } else {
	// 	fmt.Println("FindOne() result:", users)
	// }

	// // Pass these options to the Find method
	// findOptions := options.Find()

	// // Here's an array in which you can store the decoded documents
	// var results []*models.User

	// // Passing bson.D{{}} as the filter matches all documents in the collection
	// cur, err := collection.Find(config.CTX, filter, findOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Finding multiple documents returns a cursor
	// // Iterating through the cursor allows us to decode documents one at a time
	// for cur.Next(context.TODO()) {

	// 	// create a value into which the single document can be decoded
	// 	var elem models.User
	// 	err := cur.Decode(&elem)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	results = append(results, &elem)
	// }

	// if err := cur.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	// // Close the cursor once finished
	// cur.Close(config.CTX)
	collection := config.DB.Collection("users")

	var users models.User
	objID, _ := primitive.ObjectIDFromHex("5ed6117736d571d792d0bed")
	err = collection.FindOne(config.CTX, bson.M{"_id": objID}).Decode(&users)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": err})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "successfuly register user, please check your email",
		"data":   users,
	})

}

// Te1stFunc only for testing
func Te1stFunc(c *gin.Context) {
	// var err error

	collection := config.DB.Collection("users")

	var result models.User
	objID, _ := primitive.ObjectIDFromHex("5ed09aadecbc356d40af95c2")
	err = collection.FindOne(config.CTX, bson.M{"_id": objID}).Decode(&result)

	if result.Store.UserName != "" {
		c.JSON(200, gin.H{
			"status": "successfuly register user, please check your email",
			"data":   result,
		})
	} else {
		c.JSON(200, gin.H{
			"status": "KOSONG",
		})
	}

}

// TestFunc test, test
func TestFunc(c *gin.Context) {

	esclient, err := config.GetESClient()

	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	//creating student object
	newStudent := Student{
		Name:         "Alvin Khair Arthas",
		Age:          50,
		AverageScore: 59.9,
	}

	dataJSON, err := json.Marshal(newStudent)
	js := string(dataJSON)
	ind, err := esclient.Index().
		Index("belajar").
		Type("_doc").
		// Id("4").
		BodyJson(js).
		Do(config.CTX)

	if err != nil {
		panic(err)
	}

	fmt.Println(ind)
	fmt.Println("[Elastic][InsertProduct]Insertion Successful")

	// // Search Index Data
	// var students []Student

	// searchSource := elastic.NewSearchSource()
	// // searchSource.Query(elastic.NewMatchQuery("match_all"))

	// /* this block will basically print out the es query */
	// queryStr, err1 := searchSource.Source()
	// queryJs, err2 := json.Marshal(queryStr)

	// if err1 != nil || err2 != nil {
	// 	fmt.Println("[esclient][GetResponse]err during query marshal=", err1, err2)
	// }
	// fmt.Println("[esclient]Final ESQuery=\n", string(queryJs))
	// /* until this block */

	// searchService := esclient.Search().Index("belajar").SearchSource(searchSource)

	// searchResult, err := searchService.Do(config.CTX)
	// if err != nil {
	// 	fmt.Println("[ProductsES][GetPIds]Error=", err)
	// 	return
	// }

	// for _, hit := range searchResult.Hits.Hits {
	// 	var student Student
	// 	err := json.Unmarshal(hit.Source, &student)
	// 	if err != nil {
	// 		fmt.Println("[Getting Students][Unmarshal] Err=", err)
	// 	}

	// 	students = append(students, student)
	// }

	// if err != nil {
	// 	fmt.Println("Fetching student fail: ", err)
	// } else {
	// 	// for _, s := range students {
	// 	// 	fmt.Printf("Student found Name: %s, Age: %d, Score: %f \n", s.Name, s.Age, s.AverageScore)
	// 	// }
	// 	c.JSON(200, gin.H{
	// 		"status": "successfuly query elastic search",
	// 		"data":   students,
	// 	})
	// }
}
