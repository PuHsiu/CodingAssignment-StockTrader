package rest

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"trader/services"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func Run(wg *sync.WaitGroup, services []services.Service) {
	wg.Add(1)

	var baseRouter = Router{
		muxRouter: mux.NewRouter(),
	}

	baseRouter.muxRouter.Use()

	for _, service := range services {
		if s, ok := service.(RestService); ok {
			s.RestRouter(&baseRouter)
		}
	}

	_ = baseRouter.muxRouter.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		for _, method := range methods {
			fmt.Printf("Add route [%s] %s\n", method, pathTemplate)
		}

		return nil
	})

	host := viper.GetString("rest.host")
	if len(host) == 0 {
		host = "0.0.0.0:8080"
	}

	srv := &http.Server{
		Addr:         host,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      baseRouter.muxRouter,
	}

	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("listen error: %s\n", err.Error())
			return
		}
	}()
}

// TODO: Graceful shutdown
// TODO: Persistent
