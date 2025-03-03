package daikin2mqtt

import (
	"log"
	"strconv"
	"strings"

	"github.com/home2mqtt/hass"
	"github.com/samthor/daikin-go/api"
)

type IHVAC[StateType any] interface {
	hass.IHVAC
	Stop() (StateType, error)
	Restart(StateType) error
	State() StateType
	IsOn(StateType) bool
}

type Daikin interface {
	IHVAC[DaikinState]
	ReadSensor() (api.SensorInfo, error)
	GetWeekPowerEx() (GetWeekPowerEx, error)
	GetMonthPowerEx() (GetMonthPowerEx, error)

	// performs get-and-set to the ControlInfo with a modification on it. If the provided operation returns false, set is omitted
	GetAndSet(func(*DaikinState) bool) error
}

type DaikinState = api.ControlInfo

var _ Daikin = &ac{}

//var _ hass.IHVAC = &ac{}

/* 100Wh */
type Power []int

func (p Power) Sum() int {
	var sum int
	for _, d := range p {
		sum += d
	}
	return sum
}

type GetWeekPowerEx struct {
	DayOfWeek int
	Heat      Power
	Cool      Power
}

type GetMonthPowerEx struct {
	CurrHeat Power
	CurrCool Power
	PrevHeat Power
	PrevCool Power
}

type ac struct {
	baseurl string
}

// Fan implements Daikin.
func (ac *ac) Fan() hass.IEnumField {
	panic("unimplemented")
}

// Mode implements Daikin.
func (ac *ac) Mode() hass.IEnumField {
	panic("unimplemented")
}

// Power implements Daikin.
func (ac *ac) Power() hass.IField[bool] {
	panic("unimplemented")
}

// Swing implements Daikin.
func (ac *ac) Swing() hass.IEnumField {
	panic("unimplemented")
}

// TargetTemp implements Daikin.
func (ac *ac) TargetTemp() hass.IField[float64] {
	panic("unimplemented")
}

// Temp implements Daikin.
func (ac *ac) Temp() hass.ISensor[float64] {
	panic("unimplemented")
}

func parsePowerValues(values string) []int {
	data := strings.Split(values, "/")
	idata := make([]int, len(data))
	for i, d := range data {
		var err error
		idata[i], err = strconv.Atoi(d)
		if err != nil {
			log.Panicf("Could not parse integer '%s': %v", d, err)
		}
	}
	return idata
}

func New(baseurl string) Daikin {
	return &ac{baseurl: baseurl}
}

func (ac *ac) GetAndSet(stateop func(*api.ControlInfo) bool) error {
	ci, err := ac.GetControlInfo()
	if err != nil {
		return err
	}
	if stateop(&ci) {
		return ac.SetControlInfo(ci)
	}
	return nil
}

func (ac *ac) Stop() (api.ControlInfo, error) {
	ci, err := ac.GetControlInfo()
	if err != nil {
		return ci, err
	}
	if ci.Power {
		newci := ci
		newci.Power = false
		return ci, ac.SetControlInfo(newci)
	}
	return ci, nil
}

func (ac *ac) Restart(state api.ControlInfo) error {
	if state.Power {
		return ac.SetControlInfo(state)
	}
	return nil
}

func (ac *ac) State() api.ControlInfo {
	ci, _ := ac.GetControlInfo()
	return ci
}

func (ac *ac) IsOn(ci api.ControlInfo) bool {
	return ci.Power
}

func (ac *ac) ReadSensor() (api.SensorInfo, error) {
	values, err := api.Get(ac.baseurl, "aircon/get_sensor_info")
	if err != nil {
		return api.SensorInfo{}, err
	}
	result := api.ParseSensorInfo(values)
	return result, nil
}

func (ac *ac) GetControlInfo() (api.ControlInfo, error) {
	values, err := api.Get(ac.baseurl, "aircon/get_control_info")
	if err != nil {
		log.Println("Failed to get control info")
		return api.ControlInfo{}, err
	}
	result := api.ParseControlInfo(values)
	return result, nil
}

func (ac *ac) SetControlInfo(ci api.ControlInfo) error {
	values := ci.Values()
	_, err := api.Get(ac.baseurl, "aircon/set_control_info?"+values.Encode())
	if err != nil {
		log.Println("Failed to set control info")
	}
	return err
}

func (ac *ac) GetWeekPowerEx() (GetWeekPowerEx, error) {
	var result GetWeekPowerEx
	v, err := api.Get(ac.baseurl, "aircon/get_week_power_ex")
	if err != nil {
		return result, err
	}
	result.Heat = parsePowerValues(v.Get("week_heat"))
	result.Cool = parsePowerValues(v.Get("week_cool"))
	return result, nil
}

func (ac *ac) GetMonthPowerEx() (GetMonthPowerEx, error) {
	var result GetMonthPowerEx
	v, err := api.Get(ac.baseurl, "aircon/get_month_power_ex")
	if err != nil {
		return result, err
	}
	result.CurrHeat = parsePowerValues(v.Get("curr_month_heat"))
	result.CurrCool = parsePowerValues(v.Get("curr_month_cool"))
	result.PrevHeat = parsePowerValues(v.Get("prev_month_heat"))
	result.PrevCool = parsePowerValues(v.Get("prev_month_cool"))
	return result, nil
}
