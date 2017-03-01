package tenant

/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

 
import (
    "database/sql"
    "encoding/json"
    "net/http"
    "fmt"
    "time"

    "github.com/gorilla/mux"
  _ "github.com/lib/pq"
)


//TODO(nirs) polymorphisem

/////BASE - COMMON TO ALL SEEDERS (needs to move to another package)
/// Defintions
type EndpointSeeder struct{
    myDb        *sql.DB
}

/// Interface
func NewEndpointSeeder(aDb *sql.DB) *EndpointSeeder {
    endPointSeeder := new(EndpointSeeder)
    endPointSeeder.myDb = aDb
    return endPointSeeder
}

func (aEndpointSeeder *EndpointSeeder) Seed (aRouter *mux.Router, aApiPrefix string) {

    aRouter.HandleFunc(aApiPrefix+"/tenants", aEndpointSeeder.getListEndpoint).Methods("GET")
    aRouter.HandleFunc(aApiPrefix+"/tenants/{id}", aEndpointSeeder.getEndpoint).Methods("GET")
    aRouter.HandleFunc(aApiPrefix+"/tenants", aEndpointSeeder.createEndpoint).Methods("POST")
    aRouter.HandleFunc(aApiPrefix+"/tenants/{id}", aEndpointSeeder.updateEndpoint).Methods("PUT")
    aRouter.HandleFunc(aApiPrefix+"/tenants/{id}", aEndpointSeeder.deleteEndpoint).Methods("DELETE")
}

///Common utilities
func (aEndpointSeeder *EndpointSeeder) checkErr(aError error) {
    //TODO log
    if aError != nil {
        panic(aError)
    }
}


/////Class specific

type Tenant struct {
    Id        int       `json:"id"`
    Name      string    `json:"name"`
    Active    bool      `json:"active"`
    ParentId  int       `json:"parentId"`
}

func (aEndpointSeeder *EndpointSeeder) getEndpoint(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
    params := mux.Vars(aRequest)

    var id int
    var name string
    var active bool
    var parent_id sql.NullInt64
    var updated time.Time;
    
    err := aEndpointSeeder.myDb.QueryRow("SELECT * FROM tenant WHERE id="+params["id"]).Scan(&id, &name, &active, &parent_id, &updated)
    if err == sql.ErrNoRows{//TODO make it work
        aResponseWriter.WriteHeader(http.StatusNotFound)
        fmt.Fprint(aResponseWriter, params["id"]+" not found")
        return
    }
    aEndpointSeeder.checkErr(err)
    
    var parent_id1 int
    if parent_id.Valid {
        parent_id1 = int(parent_id.Int64)
    } else {
       parent_id1 = 0        
    }

    json.NewEncoder(aResponseWriter).Encode(Tenant{Id: id, Name: name, Active: active, ParentId: parent_id1})
    return
}
    

 
func (aEndpointSeeder *EndpointSeeder) getListEndpoint(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
    var tenants []Tenant
    rows, err := aEndpointSeeder.myDb.Query("SELECT * FROM tenant")
    aEndpointSeeder.checkErr(err)
    
    for rows.Next() {
        var id int
        var name string
        var active bool
        var parent_id sql.NullInt64
        var updated time.Time;
            
        err = rows.Scan(&id, &name, &active, &parent_id, &updated)
        aEndpointSeeder.checkErr(err)

        var parent_id1 int
        if parent_id.Valid {
            parent_id1 = int(parent_id.Int64)
        } else {
            parent_id1 = 0        
        }
    
        tenants = append(tenants, Tenant{Id: id, Name: name, Active: active, ParentId: parent_id1})
    }
    
    json.NewEncoder(aResponseWriter).Encode(tenants)

}
 
func (aEndpointSeeder *EndpointSeeder) createEndpoint(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
    var tenant Tenant
    _ = json.NewDecoder(aRequest.Body).Decode(&tenant)
    
    var tenantParent sql.NullInt64
    if tenant.ParentId==0{
        tenantParent.Valid = false
    }else{
    	tenantParent.Int64 = int64(tenant.ParentId)
    	tenantParent.Valid = true
    }
    

    var lastInsertId int
    err := aEndpointSeeder.myDb.QueryRow("INSERT INTO tenant(name, active, parent_id, last_updated) VALUES($1,$2,$3,$4) returning id;", tenant.Name, tenant.Active, tenantParent, "2012-12-09").Scan(&lastInsertId)
    aEndpointSeeder.checkErr(err)
    aEndpointSeeder.getListEndpoint(aResponseWriter, aRequest)
}



func (aEndpointSeeder *EndpointSeeder) updateEndpoint(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
    params := mux.Vars(aRequest)
    
    var tenant Tenant
    _ = json.NewDecoder(aRequest.Body).Decode(&tenant)

    var tenantParent sql.NullInt64
    if tenant.ParentId==0{
        tenantParent.Valid = false
    }else{
    	tenantParent.Int64 = int64(tenant.ParentId)
    	tenantParent.Valid = true
    }
    
    stmt, err := aEndpointSeeder.myDb.Prepare("update tenant set name=$1, active=$2, parent_id=$3, last_updated=$4 where id="+params["id"])
    if err == sql.ErrNoRows{//TODO make it work
        aResponseWriter.WriteHeader(http.StatusNotFound)
        return
    }
    aEndpointSeeder.checkErr(err)
    
    res, err := stmt.Exec(tenant.Name, tenant.Active, tenantParent, "2012-12-09")
    aEndpointSeeder.checkErr(err)
    
    _, err = res.RowsAffected()
    aEndpointSeeder.checkErr(err)
    
    aEndpointSeeder.getEndpoint(aResponseWriter, aRequest)
}
 
func (aEndpointSeeder *EndpointSeeder) deleteEndpoint(aResponseWriter http.ResponseWriter, aRequest *http.Request) {
    params := mux.Vars(aRequest)
    
    stmt, err := aEndpointSeeder.myDb.Prepare("delete from tenant where id="+params["id"])
    if err == sql.ErrNoRows{//TODO make it work
        aResponseWriter.WriteHeader(http.StatusNotFound)
        return
    }
    aEndpointSeeder.checkErr(err)
    
    res, err := stmt.Exec()
    aEndpointSeeder.checkErr(err)
    
    _, err = res.RowsAffected()
    aEndpointSeeder.checkErr(err)
    
    aEndpointSeeder.getListEndpoint(aResponseWriter, aRequest)
}



