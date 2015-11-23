package main
import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "gopkg.in/mgo.v2"
    "net/url"
    "bytes"
    "math/rand"
    "gopkg.in/mgo.v2/bson"
    "strconv"
    "time"
   "strings"
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

 type TripJsonObjResp struct{
     Id int `json:"id"`
     Status string `json:"status"`
    Starting string `json:"starting_from_location_id"`
    BestRoute[] string `json:"best_route_location_ids"`  
    TotalCost float64 `json:"total_uber_costs"`
    TotalDuration float64 `json:"total_uber_duration"`
    TotalDistance float64 `json:"total_distance"` 
}


type TripJsonObjPutResp struct{
     Id int `json:"id"`
     Status string `json:"status"`
    Starting string `json:"starting_from_location_id"`
    Next string `json:"next_destination_location_id"`
    BestRoute[] string `json:"best_route_location_ids"`  
    TotalCost float64 `json:"total_uber_costs"`
    TotalDuration float64 `json:"total_uber_duration"`
    TotalDistance float64 `json:"total_distance"` 
    Eta float64 `json:"uber_wait_time_eta"`
}
type TripJsonObj struct{
    LocationIds[10] string `json:"location_ids"`
    Starting string `json:"starting_from_location_id"` 
}

type Best struct{
   ID int
   Duration float64
   Distance float64
   Estimate float64
}


type Test struct {
        Name string
}
 var comb string
 var ctr int
func connectToDb() *mgo.Session{
    session, err := mgo.Dial("mongodb://admin:admin@ds041404.mongolab.com:41404/parikshithdb")
        if err != nil {
                panic(err)
        }
    return session
            
}

func getData(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

    
         session :=connectToDb()

         c := session.DB("parikshithdb").C("coll273")
         var result JsonObj
                 id,_ := strconv.Atoi(p.ByName("id"))
        err := c.Find(bson.M{"id":id }).One(&result)
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
            default:
            fmt.Println(k, "is of a unknown type",vv)
    }

}           
        id,_ := strconv.Atoi(p.ByName("id"))
       if err := c.Update(bson.M{"id": id}, bson.M{"id": id,"name":jsonObj.Name,"address":jsonObj.Address,"city":jsonObj.City,"state":jsonObj.State,"zip":jsonObj.Zip,"coordinate":jsonObj.Coordinate});err != nil {
                rw.WriteHeader(404)
                return 
            }
              
        resultResp := JsonObj{id,jsonObj.Name,jsonObj.Address,jsonObj.City,jsonObj.State,jsonObj.Zip,jsonObj.Coordinate}
        
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
        fmt.Println(k, "is of a unknwon type:",vv)
    }

}           
        t := time.Now()
        seed,_:=strconv.Atoi(t.Format("20060102150405"))
        id := rand.Intn(seed)
       
        if err := c.Insert(&JsonObj{id,jsonObj.Name,jsonObj.Address,jsonObj.City,jsonObj.State,jsonObj.Zip,jsonObj.Coordinate}) ;err != nil {
        panic(err)
        }
              
        resultResp := JsonObj{id,jsonObj.Name,jsonObj.Address,jsonObj.City,jsonObj.State,jsonObj.Zip,jsonObj.Coordinate}
        js, err4 := json.Marshal(resultResp)
        if err4 != nil {
        panic(err4)
        }
        rw.Write(js)
}

func postTrips(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    comb =" "
    var tripJsonObj TripJsonObj
    var result1 JsonObj
        session :=connectToDb()
        c := session.DB("parikshithdb").C("coll273")
           

        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
        panic(err)
        }
        if err := req.Body.Close(); err != nil {
        panic(err)
        }

      
       if err := json.Unmarshal(body, &tripJsonObj);
        err != nil {
                    fmt.Println("failed")
                    rw.WriteHeader(400)
                    return
        }
        if len(tripJsonObj.Starting)<=0 {
            rw.WriteHeader(400)
            return
        }
        
        //tmp values
        var tdur,test,tdis float64
        id,_ := strconv.Atoi(tripJsonObj.Starting)
        if err := c.Find(bson.M{"id":id}).One(&result1);err!=nil{
            panic(err)
        }

        m := make(map[int]int)
        var tmpPerm string
        
       var src,dest Coordinate
       for i:=0;tripJsonObj.LocationIds[i]!="";i++{  
            
            m[i+1],_ = strconv.Atoi(tripJsonObj.LocationIds[i]);
            tmpPerm = tmpPerm+strconv.Itoa(i+1)
        }
  
    perm(tmpPerm,m)
    combFinArr := strings.Split(comb,"\n")
    combLen := len(combFinArr)
    var minDistance float64 = 10000
    var prod float64 = 1000000
    var minInd int
    var best Best

    var dur,est,dis float64
    var i int
    var result2,result3 JsonObj
    for i=0;i<combLen-1;i++{  
         
        combArr := strings.Split(combFinArr[i],",")
        dur =0
        est =0
        dis =0
        
        //starts from 1 to eliminate the first empty id before ,
        if len(combArr)==2{
             id,_ := strconv.Atoi(combArr[1])
            if err := c.Find(bson.M{"id":id}).One(&result2);err!=nil{
            rw.WriteHeader(404)
            return 
            }
             src.Lat =  result1.Coordinate.Lat
            src.Lng =  result1.Coordinate.Lng
            dest.Lat = result2.Coordinate.Lat
                dest.Lng = result2.Coordinate.Lng
            
                tdur,test,tdis = invokeUberApi(src,dest,rw)
                dur = dur +tdur
                est = est + test
                dis = dis + tdis
                fmt.Println("1Between ",tripJsonObj.Starting," and ",combArr[1],"attr ar",tdur,test,tdis)
                prod = dis*est
                return
        }
       
        for j:=1 ;j < len(combArr)-1;j++{
            
            id,_ := strconv.Atoi(combArr[j])
            if err := c.Find(bson.M{"id":id}).One(&result2);err!=nil{
            rw.WriteHeader(404)
            return 
            }
            id_next ,_ := strconv.Atoi(combArr[j+1])
            if err := c.Find(bson.M{"id":id_next}).One(&result3);err!=nil{
            rw.WriteHeader(404)
            return 
            }
            if j==1{
                src.Lat =  result1.Coordinate.Lat
                src.Lng =  result1.Coordinate.Lng
                dest.Lat = result2.Coordinate.Lat
                dest.Lng = result2.Coordinate.Lng
                
                tdur,test,tdis = invokeUberApi(src,dest,rw)
                dur = dur +tdur
                est = est + test
                dis = dis + tdis
                fmt.Println("1Between ",tripJsonObj.Starting," and ",combArr[j],"attr ar",tdur,test,tdis)
              
                 src.Lat =result2.Coordinate.Lat
                 src.Lng = result2.Coordinate.Lng
                dest.Lat = result3.Coordinate.Lat
                dest.Lng = result3.Coordinate.Lng
                tdur,test,tdis = invokeUberApi(src,dest,rw)
                dur = dur +tdur
                est = est + test
                dis = dis + tdis
                fmt.Println("2Between ",combArr[j]," and ",combArr[j+1],"attr ar",tdur,test,tdis)
        
           }else{
                src.Lat =  result2.Coordinate.Lat
                src.Lng =  result2.Coordinate.Lng
                dest.Lat = result3.Coordinate.Lat
                dest.Lng = result3.Coordinate.Lng
               
                tdur,test,tdis = invokeUberApi(src,dest,rw)
                dur = dur +tdur
                est = est + test
                dis = dis + tdis
                 fmt.Println("3Between ",combArr[j]," and ",combArr[j+1],"attr ar",tdur,test,tdis)
            
            }
        }
        if prod > dis*est{
                minInd = i
                prod = dis*est
                minDistance = dis
                best.Distance = dis
                best.Duration = dur
                best.Estimate =est
                }
        fmt.Println("Total is dur",dur,"est is",est,"dis is ",dis,"and minInd",minInd,"i is ",i,"minDistance",minDistance)
    
   }
           
        t := time.Now()
        seed,_:=strconv.Atoi(t.Format("20060102150405"))
        id = rand.Intn(seed)
       
        var tripJsonObjResp TripJsonObjResp 
    

        tripJsonObjResp.Status = "planning"
        tripJsonObjResp.Starting= tripJsonObj.Starting
        splitArr := strings.Split(combFinArr[minInd],",")
        tripJsonObjResp.BestRoute = make([]string, len(splitArr)-1)
        var l int
        for i=0;i<len(tripJsonObj.LocationIds[i]);i++{
        tripJsonObjResp.BestRoute[i] = splitArr[i+1]
        l = i 
        }

            id,_ = strconv.Atoi(tripJsonObjResp.BestRoute[l-1])
            if err := c.Find(bson.M{"id":id}).One(&result2);err!=nil{
            rw.WriteHeader(404)
            return 
            }
            id_next ,_ := strconv.Atoi(tripJsonObjResp.Starting)
            if err := c.Find(bson.M{"id":id_next}).One(&result3);err!=nil{
            rw.WriteHeader(404)
            return 
            }
            fmt.Println("Going back to ",id_next)
            src.Lat =  result2.Coordinate.Lat
            src.Lng =  result2.Coordinate.Lng
            dest.Lat = result3.Coordinate.Lat
            dest.Lng = result3.Coordinate.Lng
            tdur,test,tdis = invokeUberApi(src,dest,rw)
                // dur = dur +tdur
                // est = est + test
                // dis = dis + tdis
        tripJsonObjResp.TotalCost = best.Estimate + test
        tripJsonObjResp.TotalDistance = best.Distance + tdis
        tripJsonObjResp.TotalDuration = best.Duration + tdur


        if err := c.Insert(&TripJsonObjResp{id,tripJsonObjResp.Status,tripJsonObjResp.Starting,tripJsonObjResp.BestRoute,tripJsonObjResp.TotalCost,tripJsonObjResp.TotalDuration,tripJsonObjResp.TotalDistance}) ;err != nil {
        panic(err)
        }
              
        resultResp := TripJsonObjResp{id,tripJsonObjResp.Status,tripJsonObjResp.Starting,tripJsonObjResp.BestRoute,tripJsonObjResp.TotalCost,tripJsonObjResp.TotalDuration,tripJsonObjResp.TotalDistance}
        js, err4 := json.Marshal(resultResp)
        if err4 != nil {
        panic(err4)
        }
        rw.Write(js)
}     
func getTripData(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

         session :=connectToDb()
         c := session.DB("parikshithdb").C("coll273")
         var result TripJsonObjResp
                 id,_:= strconv.Atoi(p.ByName("id"))
                          fmt.Println("id is ",id)
        err := c.Find(bson.M{"id":id }).One(&result)
       
        if err != nil {
            fmt.Println(err)
        rw.WriteHeader(404)
        return 
            }
        js, err := json.Marshal(result)
        if err != nil {
        panic(err)
        }
        rw.Write(js)

}

func updateTrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    var tripJsonObjPutResp TripJsonObjPutResp
    var result TripJsonObjResp
    var tripJsonObjStr, tripJsonObjEnd JsonObj
     
    session :=connectToDb()
         c := session.DB("parikshithdb").C("coll273")
        
        var status string
        var eta float64 
        id:= p.ByName("id")
        idd,_ := strconv.Atoi(id)
        if err := c.Find(bson.M{"id":idd}).One(&result);err!=nil{
            panic(err)
        }
         if ctr ==len(result.BestRoute){
        fmt.Println("Count is greater ",ctr," best",len(result.BestRoute))
        return 
        }
       
        starting_id,_ :=strconv.Atoi(result.Starting)
        if err := c.Find(bson.M{"id":starting_id}).One(&tripJsonObjStr);err!=nil{
            panic(err)
        }
        tripJsonObjPutResp.Next = result.BestRoute[ctr]
        
        ending_id,_ :=strconv.Atoi(result.BestRoute[ctr])
        if err := c.Find(bson.M{"id":ending_id}).One(&tripJsonObjEnd);err!=nil{
            panic(err)
        }
        ctr++

        status,eta = invokeSandBoxApi(getAuthTokenFromHeader(req.Header,rw), rw,tripJsonObjStr.Coordinate.Lat,tripJsonObjStr.Coordinate.Lng,tripJsonObjEnd.Coordinate.Lat,tripJsonObjEnd.Coordinate.Lng)
        fmt.Println(status)
        tripJsonObjPutResp.Status = "requesting"
        if ctr == len(result.BestRoute) {
        tripJsonObjPutResp.Status = "finished"
        }
        tripJsonObjPutResp.Starting= result.Starting
        
        tripJsonObjPutResp.BestRoute = result.BestRoute
        
        
        tripJsonObjPutResp.TotalCost = result.TotalCost
        tripJsonObjPutResp.TotalDistance = result.TotalDistance
        tripJsonObjPutResp.TotalDuration = result.TotalDuration
        tripJsonObjPutResp.Eta = eta

        resultResp := TripJsonObjPutResp{idd,tripJsonObjPutResp.Status,tripJsonObjPutResp.Starting,tripJsonObjPutResp.Next,tripJsonObjPutResp.BestRoute,tripJsonObjPutResp.TotalCost,tripJsonObjPutResp.TotalDuration,tripJsonObjPutResp.TotalDistance,tripJsonObjPutResp.Eta}
        js, err4 := json.Marshal(resultResp)
        if err4 != nil {
        panic(err4)
        }
        rw.Write(js)


}

func permRec(pre string,str string,m map[int]int){
     ln := len(str)
     if ln==0{  
        tmp,_ := strconv.Atoi(pre)
        for tmp>0 { 
                comb = comb +","+strconv.Itoa(m[tmp%10])
            tmp = tmp/10
        }
        comb = comb +"\n"
        
    }else{
     for i, r := range str {
        permRec(pre+string(r),str[0:i]+str[i+1:ln],m)
    }
    
    }
    
}
func perm(s string,m map[int]int){
   permRec("",s,m)
   fmt.Println(comb)
}
func invokeGoogleApi(params string) []byte{
    url:="http://maps.google.com/maps/api/geocode/json?address="
     par,err:=UrlEncoded(params)
     if err!=nil{
        panic(err)
     }
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
     mux.POST("/trips",postTrips)
     mux.GET("/trips/:id",getTripData)
     mux.PUT("/trips/:id/request",updateTrip)

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
func invokeUberApi(src Coordinate,dest Coordinate,rw http.ResponseWriter)(duration,estimate,distance float64) {
var lowF,durationF,distanceF float64
url := "https://api.uber.com/v1/estimates/price?start_latitude="+strconv.FormatFloat(src.Lat, 'f', 6, 64)+"&start_longitude="+strconv.FormatFloat(src.Lng, 'f', 6, 64)+"&end_latitude="+strconv.FormatFloat(dest.Lat, 'f', 6, 64)+"&end_longitude="+strconv.FormatFloat(dest.Lng, 'f', 6, 64)
client := &http.Client{}
req, _ := http.NewRequest("GET",url,nil)
req.Header.Set("Authorization", " Token PUSBQrdLW8_O39XWcQVkpKhoFLVmSdm34iF8QUmt")
response, err:= client.Do(req)
if err!=nil{
    fmt.Println(err)
}    
contents,_ := ioutil.ReadAll(response.Body)
var resp interface{}
if err := json.Unmarshal(contents, &resp);
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
            tmp1 :=first["duration"]
             switch tmp1 := tmp1.(type) {
            case float64:
            durationF = float64(tmp1)
             default: fmt.Println("No Match!")
         }
           
            low := first["low_estimate"]
             
            switch low := low.(type) {
            case float64:
            lowF = float64(low)
             default: fmt.Println("No Match!")
         }
            estimate = lowF

            
            switch tmp := first["distance"].(type) {
            case float64:
            distanceF = float64(tmp)
             default: fmt.Println("No Match!")
         }
          
            default:
        fmt.Println(k, "is of a unknwon type",vv)
    }

    }
    duration = durationF
    distance = distanceF

    return
}

func invokeSandBoxApi(authToken string,rw http.ResponseWriter,start_latitude float64,start_longitude float64,end_latitude float64,end_longitude float64) (status string,eta float64){
     url := "https://sandbox-api.uber.com/v1/requests"
     
     var payLoad = "{\"product_id\":\"a1111c8c-c720-46c3-8534-2fcdd730040d\",\"start_latitude\":\""+strconv.FormatFloat(start_latitude, 'f', 6, 64)+"\""+","+"\"start_longitude\":\""+strconv.FormatFloat(start_longitude, 'f', 6, 64)+"\""+"," +"\"end_latitude\":\""+strconv.FormatFloat(end_latitude, 'f', 6, 64)+"\""+","+
"\"end_longitude\":\""+strconv.FormatFloat(end_longitude,'f', 6, 64)+"\""+"}"
     fmt.Println("\""+strconv.FormatFloat(start_longitude, 'f', 6, 64)+"\"")
    var jsonStr = []byte(payLoad)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Authorization","Bearer "+authToken) // eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsInJlcXVlc3RfcmVjZWlwdCIsInJlcXVlc3QiLCJoaXN0b3J5X2xpdGUiXSwic3ViIjoiMGY0YTRkZmItZTcyMS00ZTQ4LTk1YWYtNzc1ZDgxMGUzNzVlIiwiaXNzIjoidWJlci11czEiLCJqdGkiOiJjNzM5NGFkOC00MDRiLTRkZjctYWQyOC1lNzAyZWVlMWI4OGIiLCJleHAiOjE0NTA2NDUzNjcsImlhdCI6MTQ0ODA1MzM2NiwidWFjdCI6Im9IUEVhMHZBOFhDek9iSkhRdVNvd0ltRUJodncxNSIsIm5iZiI6MTQ0ODA1MzI3NiwiYXVkIjoiekFRNUNuZnNLR1FMRWxJdGhCV1BNTG1YTHBJaHdYcTEifQ.BARKfSrr1YDLaq_uh1AbgjsYoAYeGIixNFbsRaoHGGFY2hDpL5gSTq7h9CYTQeWCkwbTvYPWnqRVUSFx8h7fy3DLZ_6KqLOufYs5LXrNtPqdO09MbC9-c2qhLy62NgORHhosIhuL7DjL_D7mVcVa0P7Kqtc0Z2azun9WVBWllx90mhnAtyffX9kLKZDgcs1txF3K68uQG1xe2rc81Io9al0ccEz2u6_fNWx7Uw0G1mNCYsKjCsmRRgeSBpZ-bTSN0UZ3P7boC3B7Sw7Y1rZpNYc_71U6tyGslTwD7eIeS_BtwI_lgMSIUWKn6IiC8BsbF0oFCQi63r1aJkUb6BpZBA")
    req.Header.Set("Content-Type", "application/json")
     
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    
    defer resp.Body.Close()
    contents,_ := ioutil.ReadAll(resp.Body)

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
  
    var respI interface{}
    if err := json.Unmarshal(contents, &respI);
    err != nil {
            rw.WriteHeader(400)
            return
    }
     m := respI.(map[string]interface{}) 
     
     status = m["status"].(string)
     
     eta = m["eta"].(float64)
     return

}
func getAuthTokenFromHeader(header http.Header,rw http.ResponseWriter) (authToken string){
   
    if header.Get("Authorization")==""{
        rw.WriteHeader(401)
        return
    }
     authToken = strings.Split(header.Get("Authorization")," ")[1]

    return authToken
}