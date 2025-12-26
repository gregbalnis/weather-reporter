package models

// Location represents a geographical location.
type Location struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	Region    string  `json:"admin1"`
}

// Weather represents current weather conditions.
type Weather struct {
	Temperature   float64 `json:"temperature_2m"`
	Humidity      int     `json:"relative_humidity_2m"`
	ApparentTemp  float64 `json:"apparent_temperature"`
	Precipitation float64 `json:"precipitation"`
	CloudCover    int     `json:"cloud_cover"`
	Pressure      float64 `json:"surface_pressure"`
	WindSpeed     float64 `json:"wind_speed_10m"`
	WindDirection int     `json:"wind_direction_10m"`
	WindGusts     float64 `json:"wind_gusts_10m"`
}

// WeatherUnits represents the units for weather data.
type WeatherUnits struct {
	Temperature   string `json:"temperature_2m"`
	Humidity      string `json:"relative_humidity_2m"`
	ApparentTemp  string `json:"apparent_temperature"`
	Precipitation string `json:"precipitation"`
	CloudCover    string `json:"cloud_cover"`
	Pressure      string `json:"surface_pressure"`
	WindSpeed     string `json:"wind_speed_10m"`
	WindDirection string `json:"wind_direction_10m"`
	WindGusts     string `json:"wind_gusts_10m"`
}

// WeatherResponse is the full response from the weather service.
type WeatherResponse struct {
	Current      Weather      `json:"current"`
	CurrentUnits WeatherUnits `json:"current_units"`
}
