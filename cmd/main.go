package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"

	"gitlab-code-review-notifier/internal"
	"gitlab-code-review-notifier/internal/controller"
	"gitlab-code-review-notifier/internal/database"
	"gitlab-code-review-notifier/pkg/envutil"
	"gitlab-code-review-notifier/pkg/firingservice"
	"gitlab-code-review-notifier/pkg/gitlabservice"
	"gitlab-code-review-notifier/pkg/log"
	"gitlab-code-review-notifier/pkg/notifier"
	"gitlab-code-review-notifier/pkg/scheduler"
)

func main() {
	logger := log.NewLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.Fatalf("Failed to start: %v", r)
		}
	}()

	logLevel := envutil.GetEnvStr(internal.EnvLogLevel)
	if err := log.SetLevelString(envutil.GetEnvStr(internal.EnvLogLevel)); err != nil {
		panic(fmt.Errorf("set log level to %s: %v", logLevel, err))
	}

	logMode := envutil.GetEnvStr(internal.EnvLogMode)
	if err := log.SetModeString(logMode); err != nil {
		panic(fmt.Errorf("set log mode to %s: %v", logMode, err))
	}

	db, err := database.NewDbFromEnv()
	if err != nil {
		panic(fmt.Errorf("create DB instance: %v", err))
	}

	if err := db.Migrate(); err != nil {
		panic(fmt.Errorf("DB migration: %v", err))
	}

	envSchedulerFixedTimes := envutil.GetEnvStr(internal.EnvSchedulerFixedTimes)
	schedulerFixedTimes := make([]string, 0)
	for _, fixedTime := range regexp.MustCompile(`[ ;,]`).Split(envSchedulerFixedTimes, -1) {
		if len(fixedTime) > 0 {
			schedulerFixedTimes = append(schedulerFixedTimes, fixedTime)
		}
	}

	schedulerIntervalMinutes := envutil.GetEnvUintOrDefault(internal.EnvSchedulerIntervalMinutes, 0)
	timeZoneStr := envutil.GetEnvStrOrDefault(internal.EnvTimeZone, "UTC")
	timeZone, err := time.LoadLocation(timeZoneStr)
	if err != nil {
		panic(fmt.Errorf("parse location in %s %s: %v", internal.EnvTimeZone, timeZoneStr, err))
	}
	schedulerConf := scheduler.Config{
		RepeatInterval: schedulerIntervalMinutes,
		FixedTimes:     schedulerFixedTimes,
		TimeZone:       timeZone,
		WorkdayStartAt: int(envutil.GetEnvUintOrDefault(internal.EnvWorkdayStartsAt, 10)),
		WorkdayEndAt:   int(envutil.GetEnvUintOrDefault(internal.EnvWorkdayEndsAt, 19)),
	}
	sched := scheduler.NewScheduler(schedulerConf)

	clientRepository := database.NewClientRepository(db)
	gitlabUrl := envutil.MustGetEnvStr(internal.EnvGitlabUrl)
	gitlabClientFactory := gitlabservice.NewInstancedClientFactory(gitlabUrl)
	notifierFactory := notifier.NewFactory("pkg/notifier/templates")
	configuredClientFactory := firingservice.NewConfiguredClientFactory(gitlabClientFactory, notifierFactory)
	service := firingservice.NewFiringService()

	job := func() {
		logger.Infof("Starting firing job")
		clients, err := clientRepository.GetAll()
		if err != nil {
			logger.Errorf("Failed to get clients from repository: %v", err)
			return
		}
		for _, client := range clients {
			configuredClient, err := configuredClientFactory.MakeClient(*client)
			if err != nil {
				logger.Errorf("Failed to make configured client %d: %v", client.Id, err)
				return
			}
			logger.Infof("Start processing client %d", configuredClient.Config.Id)
			service.ProcessConfig(configuredClient)
		}
		logger.Infof("Ending firing job")
	}

	if err := sched.Submit(job); err != nil {
		panic(fmt.Errorf("submit scheduled job: %v", err))
	}

	go sched.Run()

	clientController := controller.NewClientController(clientRepository)

	r := mux.NewRouter()
	r.HandleFunc("/", RootHandler).Methods("GET")
	r.HandleFunc("/clients", clientController.GetAll).Methods("GET")
	r.HandleFunc("/clients", clientController.Create).Methods("POST")
	r.HandleFunc("/clients/{id:[0-9]+}", clientController.Get).Methods("GET")
	r.HandleFunc("/clients/{id:[0-9]+}", clientController.Update).Methods("PUT")
	r.HandleFunc("/clients/{id:[0-9]+}", clientController.Delete).Methods("DELETE")

	addr := ":8080"
	logger.Infof("Starting at %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(fmt.Errorf("start listening: %v", err))
	}

}

func RootHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "ok")
}
