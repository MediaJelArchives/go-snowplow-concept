package main

import (
	"github.com/gin-gonic/gin"
	sp "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
)

func main() {
	router := gin.Default()
	router.Use(TrackerMiddleWare())
	router.GET("/pageview", PageviewTrackingDemo)
	router.GET("/ecommerce", EcommerceTrackingDemo)
	router.Run("localhost:3000")
}

func TrackerMiddleWare() gin.HandlerFunc {
	// Instantiate a subject/user
	subject := sp.InitSubject()

	/*
	Optionally, if these functions are not implemented
	on a REST API request/response and the data resolution is lacking
	an option we have is manually overriding the subject/user values.

	Ex. 
	subject.SetUserId("{{USER_ID}}")
	subject.SetDomainUserId("{{DOMAIN_USER_ID}}") 
	subject.SetIpAddress("{{IP_ADDRESS}}")
	subject.SetNetworkUserId("{{NETWORK_USER_ID}}") 
	subject.SetUseragent("{{USER_AGENT}}")
	subject.SetTimeZone("{{TIME_ZONE}}")
	*/

	emitter := sp.InitEmitter(sp.RequireCollectorUri("collector.dmp.mediajel.ninja")) 
	tracker := sp.InitTracker(sp.RequireEmitter(emitter),
		sp.OptionSubject(subject), 
		sp.OptionNamespace("sp1"), 
		sp.OptionAppId("iHeartJane-golang"),
		sp.OptionPlatform("srv"))

	iHeartJaneContext := []sp.SelfDescribingJson{
		*sp.InitSelfDescribingJson(
			"iglu:com.mediajel.events/iheartjane/jsonschema/1-0-0",
			map[string]interface{}{
				"advertiserId": "{{ADVERTISER_ID}}",
				"advertiserName": "{{ADVERTISER_NAME}}",
				"storeId": "{{STORE_ID}}",
				"storeName": "{{STORE_NAME}}",
				"locationId": "{{LOCATION_ID}}",
				"locationName": "{{LOCATION_NAME}}",
				
			},
		),
	}

	return func(context *gin.Context) {
		context.Set("tracker", tracker)
		context.Set("subject", subject)
		context.Set("iHeartJaneContext", iHeartJaneContext)
		context.Next()
	}
}

func PageviewTrackingDemo(context *gin.Context) {
	tracker := context.MustGet("tracker").(*sp.Tracker)
	subject := context.MustGet("subject").(*sp.Subject)
	iHeartJaneContext := context.MustGet("iHeartJaneContext").([]sp.SelfDescribingJson)

	tracker.TrackPageView(sp.PageViewEvent{
		PageUrl: sp.NewString("{{PAGE_URL}}"),
		Subject: subject,
		Contexts: iHeartJaneContext,
	})

	context.JSON(200, gin.H{ "status": "success" })
}

func EcommerceTrackingDemo(context *gin.Context) {
	tracker := context.MustGet("tracker").(*sp.Tracker)
	iHeartJaneContext := context.MustGet("iHeartJaneContext").([]sp.SelfDescribingJson)

	// Example Cart Items
	items := []sp.EcommerceTransactionItemEvent{
		{
			Sku:      sp.NewString("{{PRODUCT_SKU}}"),
			Price:    sp.NewFloat64(420),
			Quantity: sp.NewInt64(1),
			Name:     sp.NewString("{{PRODUCT_NAME}}"),
			Category: sp.NewString("{{CATEGORY}}"),
		},
		{
			Sku:      sp.NewString("{{PRODUCT_SKU}}"),
			Price:    sp.NewFloat64(420),
			Quantity: sp.NewInt64(1),
			Name:     sp.NewString("{{PRODUCT_NAME}}"),
			Category: sp.NewString("{{CATEGORY}}"),
		},
	}

	// Example Transaction
	tracker.TrackEcommerceTransaction(sp.EcommerceTransactionEvent{
		OrderId:     sp.NewString("{{ORDER_ID}}"),
		TotalValue:  sp.NewFloat64(420),
		Affiliation: sp.NewString("{{STORE_ID}}"),
		TaxValue:    sp.NewFloat64(4.20),
		Shipping:    sp.NewFloat64(42),
		City:        sp.NewString("{{CITY}}"),
		State:       sp.NewString("{{STATE}}"),
		Country:     sp.NewString("{{COUNTRY}}"),
		Currency:    sp.NewString("USD"),
		Items:       items,
		Contexts: iHeartJaneContext,
	})

	context.JSON(200, gin.H{ "status": "success"})
}


