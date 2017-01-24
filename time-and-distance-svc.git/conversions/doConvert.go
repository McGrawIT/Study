package conversions

func	doConvert ( fromUnits, toUnits string, inputValue float64 ) ( convertedValue float64 ) {

//	Get primary unit for both Measurement System ( e.g., Metric is meters )

	Systems	= MeasurementSystems {}

	Systems.FromSystem = fromRatio [ fromUnits ].Unit
	Systems.ToSystem = toRatio [ toUnits ].Unit

	convertedValue = inputValue * fromRatio [ fromUnits ].Ratio * Factor [ Systems ] * toRatio [ toUnits ].Ratio
 
	return
}

