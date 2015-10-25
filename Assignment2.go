package main
import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "gopkg.in/mgo.v2"
    "net/url"
    "math/rand"
    "gopkg.in/mgo.v2/bson"
    "strconv"
   // "strings"
)


type Coordinate struct{
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
}
type JsonObj struct{
    Id int `json:"id"`
    Name string `json:"name"`
    Address string `json:"address"`
    City string `json:"city"`
    State string `json:"state"`
    Zip string `json:"zip"`
    Coordinate Coordinate `bson:"coordinate"` 
}


type Test struct {
        Name string
}

func connectToDb() *mgo.Session{
    session, err := mgo.Dial("mongodb://admin:admin@ds041404.mongolab.com:41404/parikshithdb")
        if err != nil {
                panic(err)
        }
           session.SetMode(mgo.Monotonic, true)
           fmt.Println("session is ",session)
    return session
            
}

func getData(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

    
         session :=connectToDb()

         c := session.DB("parikshithdb").C("coll273")
         var result JsonObj
                 fmt.Println("Name",p.ByName("id"))
                 id,_ := strconv.Atoi(p.ByName("id"))
        err := c.Find(bson.M{"id":id }).One(&result)
       // err := c.Find(nil).All(&result)
        if err != nil {
        rw.WriteHeader(404)
        return 
            }
        js, err := json.Marshal(result)
        if err != nil {
        panic(err)
        }
        rw.Write(js)

}
func deleteData(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    session :=connectToDb()

    c := session.DB("parikshithdb").C("coll273")
    fmt.Println("Name",p.ByName("id"))
    id,_ := strconv.Atoi(p.ByName("id"))
    err := c.Remove(bson.M{"id":id})
    if err != nil {
        rw.WriteHeader(404)
        return 
    }
}

func updateData(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    var jsonObj JsonObj

    session :=connectToDb()
         c := session.DB("parikshithdb").C("coll273")
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
        panic(err)
        }
        if err := req.Body.Close(); err != nil {
        panic(err)
        }

        var resp interface{}
        if err := json.Unmarshal(body, &jsonObj);
        err != nil {
                    rw.WriteHeader(400)
                    return
            }
        googBody :=(invokeGoogleApi(jsonObj.Address+","+jsonObj.City+","+jsonObj.State))
        json.Unmarshal(googBody, &resp)

        m := resp.(map[string]interface{})
        
        for k, v := range m {
                switch vv := v.(type) {
        case []interface{}:
            tmp,_ := v.([]interface{})
            first := tmp[0].(map[string]interface{})
          
            second :=first["geometry"].(map[string]interface{})
            third := second["location"].(map[string]interface{})
                       
            jsonObj.Coordinate.Lat = third["lat"].(float64)
            jsonObj.Coordinate.Lng = third["lng"].(float64)
        fmt.Println(k, "is an array:")
        for i, u := range vv {
            fmt.Println("wow!...............",i,u)
                }
        default:
        fmt.Println(k, "is of a type I don't know how to handle")
    }

}           
        id,_ := strconv.Atoi(p.ByName("id"))
       if err := c.Update(bson.M{"id": id}, bson.M{"id": id,"name":jsonObj.Name,"address":jsonObj.Address,"city":jsonObj.City,"state":jsonObj.State,"zip":jsonObj.Zip,"coordinate":jsonObj.Coordinate});err != nil {
                rw.WriteHeader(404)
                return 
            }
              
        resultResp := JsonObj{id,jsonObj.Name,jsonObj.Address,jsonObj.City,jsonObj.State,jsonObj.Zip,jsonObj.Coordinate}
        fmt.Println("Name:")
        js, err4 := json.Marshal(resultResp)
        if err4 != nil {
        panic(err4)
        }
        rw.Write(js)
}
func postData(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    var jsonObj JsonObj

        session :=connectToDb()
         c := session.DB("parikshithdb").C("coll273")
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
        panic(err)
        }
        if err := req.Body.Close(); err != nil {
        panic(err)
        }

        var resp interface{}
       if err := json.Unmarshal(body, &jsonObj);
       err != nil {
                    rw.WriteHeader(400)
                    return
            }
        if len(jsonObj.Address)<=0 && len(jsonObj.City)<=0 && len(jsonObj.State) <=0{
            rw.WriteHeader(400)
            return
        }
        googBody :=(invokeGoogleApi(jsonObj.Address+","+jsonObj.City+","+jsonObj.State))
        if err := json.Unmarshal(googBody, &resp);
        err != nil {
                    rw.WriteHeader(400)
                    return
            }

        m := resp.(map[string]interface{})
        
        for k, v := range m {
                switch vv := v.(type) {
        case []interface{}:
            tmp,_ := v.([]interface{})
            first := tmp[0].(map[string]interface{})
            second :=first["geometry"].(map[string]interface{})
            third := second["location"].(map[string]interface{})
           
            
            jsonObj.Coordinate.Lat = third["lat"].(float64)
            jsonObj.Coordinate.Lng = third["lng"].(float64)
        default:
        fmt.Println(k, "is of a type I don't know how to handle:",vv)
    }

}           
        id := rand.Int()
        if err := c.Insert(&JsonObj{id,jsonObj.Name,jsonObj.Address,jsonObj.City,jsonObj.State,jsonObj.Zip,jsonObj.Coordinate}) ;err != nil {
        panic(err)
        }
              
        resultResp := JsonObj{id,jsonObj.Name,jsonObj.Address,jsonObj.City,jsonObj.State,jsonObj.Zip,jsonObj.Coordinate}
        fmt.Println("Name:")
        js, err4 := json.Marshal(resultResp)
        if err4 != nil {
        panic(err4)
        }
        rw.Write(js)
}

func invokeGoogleApi(params string) []byte{
    url:="http://maps.google.com/maps/api/geocode/json?address="
     par,err:=UrlEncoded(params)
     if err!=nil{
        panic(err)
     }
    fmt.Println(url+par)
    response,_ := http.Get(url+par)
    contents,_ := ioutil.ReadAll(response.Body)
    return contents
}   

func main() {
     mux := httprouter.New()
     mux.GET("/locations/:id", getData)
     mux.POST("/locations",postData) 
     mux.DELETE("/locations/:id",deleteData)
     mux.PUT("/locations/:id",updateData)
     server := http.Server{
             Addr:        "127.0.0.1:8080",
             Handler: mux,
     }

     server.ListenAndServe()

}
func UrlEncoded(str string) (string, error) {
    u, err := url.Parse(str)
    if err != nil {
        return "", err
    }
    return u.String(), nil
}