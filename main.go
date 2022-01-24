package main
 
import (
    "fmt"
    "net/http"
    "os"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/gin-gonic/gin"
)

/*数据库链接*/
func InitMongoSession() *mgo.Session {
    mHost := "127.0.0.1"
    mPort := "27017"
    mUsername := "root" //mongodb的账号
    mPassword := "123456" //mongodb的密码
    session, err := mgo.Dial(mHost + ":" + mPort)
    if err != nil {
        fmt.Println("mgo.Dial-error:", err)
        os.Exit(0)
    }
    session.SetMode(mgo.Eventual, true)
    myDB := session.DB("admin") //这里的关键是连接mongodb后，选择admin数据库，然后登录，确保账号密码无误之后，该连接就一直能用了
    err = myDB.Login(mUsername, mPassword)
    if err != nil {
        fmt.Println("Login-error:", err)
        os.Exit(0)
    }
    session.SetPoolLimit(10)
    return session
}


/*定义结构体*/
type Surl struct {
    Urlname string `bson:"urlname"` 
    Nexturl string `bson:"nexturl"` 
}

/*执行插入*/
func insert(u,n string) string{
    notice := "daydream is inserting......"
    session := InitMongoSession()
 
    myDB := session.DB("admin").C("s_url")
    reurl := Surl{Urlname:u,Nexturl:n}
    err := myDB.Insert(reurl)
    if err != nil { panic(err) }
    return notice
}


/*获取参数并跳转*/
func main() {
    ///*
    r := gin.Default()
    r.LoadHTMLGlob("./html/*")// 指明html加载文件目录
    r.Handle("GET", "/", func(context *gin.Context) {
        // 返回HTML文件，响应状态码200，html文件名为index.html，模板参数为nil
        context.HTML(http.StatusOK, "index.html", nil)
    })
    r.POST("/insert", func(c *gin.Context) {
        //types := c.DefaultPostForm("type", "post")
        api := c.PostForm("API")
        url := c.PostForm("URL")
        //c.String(http.StatusOK, fmt.Sprintf("api:%s,url:%s,type:%s", api, url, types))
        c.String(http.StatusOK,insert(api,url))
        })
    /*所有定向跳转的查询、输出*/
    r.GET("/show", func(c *gin.Context) {
        session := InitMongoSession()
        myDB := session.DB("admin").C("s_url")
        var result []Surl
        myDB.Find(bson.M{}).All(&result)
        
        ///*
        for id,res :=range result{
            //fmt.Println(res)
            
            U := res.Urlname
            N := res.Nexturl
            c.String(http.StatusOK, fmt.Sprintf("id:%d,U:%s,N:%s\n",id,U,N))
            }
        //*/
        })
    /*查询并执行跳转*/
     r.GET("/:name", func(c *gin.Context) {
        name := c.Param("name")//获取跳转api
        
        session := InitMongoSession()
        myDB := session.DB("admin").C("s_url")
        result := Surl{}
        err := myDB.Find(bson.M{"urlname": name}).One(&result)
        if err != nil { panic(err) }
        n := result.Nexturl
        
        //c.String(http.StatusOK, name)//输出获取name的值
        //n="http://www.baidu.com"
        c.Redirect(http.StatusMovedPermanently,n)           //通过变量来实现参数传递重定向
    })
    r.Run()
    //*/
    
}
