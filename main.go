package main
import (
    "clip"
    "github.com/gin-gonic/gin"

//    "github.com/DeanThompson/ginpprof"
)


func main() {
    clip.Init()
    router := gin.Default()

    router.POST("/", clip.PostHandler)

    // automatically add routers for net/http/pprof
    // e.g. /debug/pprof, /debug/pprof/heap, etc.
    //    ginpprof.Wrapper(router)

    clip.Logr.Infof("[CLIP] Started on port %s", clip.Config.Port)
    // Listen and server on 0.0.0.0:3001
    router.Run(clip.Config.Port)

}

