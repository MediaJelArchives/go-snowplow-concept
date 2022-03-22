package main

import (
	"github.com/gin-gonic/gin"
	sp "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
)

func TrackerMiddleWare() gin.HandlerFunc {
	// Instantiate a subject
	subject := sp.InitSubject()
	// Initialize tracker config
	emitter := sp.InitEmitter(sp.RequireCollectorUri("collector.dmp.mediajel.ninja"))
	tracker := sp.InitTracker(sp.RequireEmitter(emitter),
		sp.OptionSubject(subject),
		sp.OptionNamespace("sp1"),
		sp.OptionAppId("golang-snowplow"),
		sp.OptionPlatform("srv"),
		sp.OptionBase64Encode(false))
	
		return func(context *gin.Context) {
			context.Set("tracker", tracker)
			context.Set("subject", subject)
			context.Next()
		}
}



func PageviewTrackingDemo (context *gin.Context) {
	tracker := context.MustGet("tracker").(*sp.Tracker)
	subject := context.MustGet("subject").(*sp.Subject)

	 tracker.TrackPageView(sp.PageViewEvent{
	 	PageUrl: sp.NewString("localhost:3000"),
	 	Subject: subject,
	 })
}

func main() {
	router := gin.Default()
	router.Use(TrackerMiddleWare())
	router.GET("/pageview", PageviewTrackingDemo)
	router.Run("localhost:3000")
}
