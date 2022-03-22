package main

import (
	"github.com/gin-gonic/gin"
	sp "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
)

func main() {
	router := gin.Default()
	router.Use(TrackerMiddleWare())
	router.GET("/pageview", PageviewTrackingDemo)
	router.Run("localhost:3000")
}

func TrackerMiddleWare() gin.HandlerFunc {
	// Instantiate a subject & override values
	subject := sp.InitSubject()
	subject.SetUserId("{{USER_ID}}")
	subject.SetDomainUserId("{{DOMAIN_USER_ID")
	subject.SetIpAddress("{{IP_ADDRESS}}")
	subject.SetNetworkUserId("{{NETWORK_USER_ID}}")
	subject.SetUseragent("{{USER_AGENT}}")
	subject.SetTimeZone("{{TIME_ZONE}}")
	subject.SetViewPort(420, 420)
	subject.SetScreenResolution(420, 420)

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

func PageviewTrackingDemo(context *gin.Context) {
	tracker := context.MustGet("tracker").(*sp.Tracker)
	subject := context.MustGet("subject").(*sp.Subject)

	tracker.TrackPageView(sp.PageViewEvent{
		PageUrl: sp.NewString("{{DOMAIN_URL}}"),
		Subject: subject,
	})

}

func EcommerceTrackingDemo(context *gin.Context) {
	tracker := context.MustGet("tracker").(*sp.Tracker)
	// subject := context.MustGet("subject").(*sp.Subject)

	items := []sp.EcommerceTransactionItemEvent{
		sp.EcommerceTransactionItemEvent{
			Sku:      sp.NewString("pbz0026"),
			Price:    sp.NewFloat64(20),
			Quantity: sp.NewInt64(1),
		},
		sp.EcommerceTransactionItemEvent{
			Sku:      sp.NewString("pbz0038"),
			Price:    sp.NewFloat64(15),
			Quantity: sp.NewInt64(1),
			Name:     sp.NewString("red hat"),
			Category: sp.NewString("menswear"),
		},
	}
	tracker.TrackEcommerceTransaction(sp.EcommerceTransactionEvent{
		OrderId:     sp.NewString("6a8078be"),
		TotalValue:  sp.NewFloat64(35),
		Affiliation: sp.NewString("some-affiliation"),
		TaxValue:    sp.NewFloat64(6.12),
		Shipping:    sp.NewFloat64(30),
		City:        sp.NewString("Dijon"),
		State:       sp.NewString("Bourgogne"),
		Country:     sp.NewString("France"),
		Currency:    sp.NewString("EUR"),
		Items:       items,
	})

}
