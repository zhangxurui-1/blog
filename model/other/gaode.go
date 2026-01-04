package other

// IPResponse 表示 IP 定位查询的响应结果
type IPResponse struct {
	Status    string `json:"status"`    // 返回结果状态，0表示失败，1表示成功
	Info      string `json:"info"`      // 返回状态说明
	InfoCode  string `json:"info_code"` // 状态码：10000代表正确
	Province  string `json:"province"`  // 省份名称
	City      string `json:"city"`      // 城市名称
	Adcode    string `json:"adcode"`    // 城市的 adcode 编码
	Rectangle string `json:"rectangle"` // 所在城市矩形区域范围
}

// Cast 表示天气预报中的每日数据
type Cast struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	DayWeather   string `json:"day_weather"`
	NightWeather string `json:"night_weather"`
	DayTemp      string `json:"day_temp"`
	NightTemp    string `json:"night_temp"`
	DayWind      string `json:"day_wind"`
	NightWind    string `json:"night_wind"`
	DayPower     string `json:"day_power"`
	NightPower   string `json:"night_power"`
}

// Live 表示实况天气数据
type Live struct {
	Province         string `json:"province"`
	City             string `json:"city"`
	Adcode           string `json:"adcode"`
	Weather          string `json:"weather"`
	Temperature      string `json:"temperature"`
	WindDirection    string `json:"wind_direction"`
	WindPower        string `json:"wind_power"`
	Humidity         string `json:"humidity"`          // 空气湿度
	ReportTime       string `json:"report_time"`       // 数据发布时间
	TemperatureFloat string `json:"temperature_float"` // 浮点型气温
	HumidityFloat    string `json:"humidity_float"`    // 浮点型湿度
}

// Forecast 表示天气预报信息
type Forecast struct {
	City       string `json:"city"`
	Adcode     string `json:"adcode"`
	Province   string `json:"province"`
	ReportTime string `json:"report_time"`
	Casts      []Cast `json:"casts"`
}

type WeatherResponse struct {
	Status   string   `json:"status"`    // 返回状态
	Count    string   `json:"count"`     // 返回结果总数目
	Info     string   `json:"info"`      // 返回的状态信息
	InfoCode string   `json:"info_code"` // 返回状态说明, 10000代表正确
	Lives    []Live   `json:"lives"`     // 实况天气数据信息
	Forecast Forecast `json:"forecast"`  // 预报天气信息数据
}
