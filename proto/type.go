package proto

// DeviceType identifies what kind of device this is.
type DeviceType string

// A list of possible device types can be found here
// https://developers.google.com/assistant/smarthome/guides
const (
	DeviceTypeCamera              DeviceType = "action.devices.types.CAMERA"
	DeviceTypeLight               DeviceType = "action.devices.types.LIGHT"
	DeviceTypeFan                 DeviceType = "action.devices.types.FAN"
	DeviceTypeOutlet              DeviceType = "action.devices.types.OUTLET"
	DeviceTypeSwitch              DeviceType = "action.devices.types.SWITCH"
	DeviceTypeThermostat          DeviceType = "action.devices.types.THERMOSTAT"
	DeviceTypeBlinds              DeviceType = "action.devices.types.BLINDS"
	DeviceTypeAirConditioningUnit DeviceType = "action.devices.types.AC_UNIT"
	DeviceTypeAirFreshener        DeviceType = "action.devices.types.AIRFRESHENER"
	DeviceTypeAirPurifier         DeviceType = "action.devices.types.AIRPURIFIER"
	DeviceTypeAwning              DeviceType = "action.devices.types.AWING"
	DeviceTypeBathtub             DeviceType = "action.devices.types.BATHTUB"
	DeviceTypeBed                 DeviceType = "action.devices.types.BED"
	DeviceTypeBlender             DeviceType = "action.devices.types.BLENDER"
	DeviceTypeBoiler              DeviceType = "action.devices.types.BOILER"
	DeviceTypeCloset              DeviceType = "action.devices.types.CLOSET"
	DeviceTypeCoffeeMaker         DeviceType = "action.devices.types.COFFEE_MAKER"
	DeviceTypeCooktop             DeviceType = "action.devices.types.COOKTOP"
	DeviceTypeCurtain             DeviceType = "action.devices.types.CURTAIN"
	DeviceTypeDehumidifier        DeviceType = "action.devices.types.DEHUMIDIFIER"
	DeviceTypeDehydrator          DeviceType = "action.devices.types.DEHYDRATOR"
	DeviceTypeDishwasher          DeviceType = "action.devices.types.DISHWASHER"
	DeviceTypeDoor                DeviceType = "action.devices.types.DOOR"
	DeviceTypeDrawer              DeviceType = "action.devices.types.DRAWER"
	DeviceTypeDryer               DeviceType = "action.devices.types.DRYER"
	DeviceTypeFaucet              DeviceType = "action.devices.types.FAUCET"
	DeviceTypeFireplace           DeviceType = "action.devices.types.FIREPLACE"
	DeviceTypeFryer               DeviceType = "action.devices.types.FRYER"
	DeviceTypeGarage              DeviceType = "action.devices.types.GARAGE"
	DeviceTypeGate                DeviceType = "action.devices.types.GATE"
	DeviceTypeGrill               DeviceType = "action.devices.types.GRILL"
	DeviceTypeHeater              DeviceType = "action.devices.types.HEATER"
	DeviceTypeHood                DeviceType = "action.devices.types.HOOD"
	DeviceTypeHumidifier          DeviceType = "action.devices.types.HUMIDIFIER"
	DeviceTypeKettle              DeviceType = "action.devices.types.KETTLE"
	DeviceTypeLock                DeviceType = "action.devices.types.LOCK"
	DeviceTypeMop                 DeviceType = "action.devices.types.MOP"
	DeviceTypeMower               DeviceType = "action.devices.types.MOWER"
	DeviceTypeMicrowave           DeviceType = "action.devices.types.MICROWAVE"
	DeviceTypeMulticooker         DeviceType = "action.devices.types.MULTICOOKER"
	DeviceTypeOven                DeviceType = "action.devices.types.OVEN"
	DeviceTypePergola             DeviceType = "action.devices.types.PERGOLA"
	DeviceTypePetFeeder           DeviceType = "action.devices.types.PETFEEDER"
	DeviceTypePressureCooker      DeviceType = "action.devices.types.PRESSURECOOKER"
	DeviceTypeRadiator            DeviceType = "action.devices.types.RADIATOR"
	DeviceTypeRefrigerator        DeviceType = "action.devices.types.REFRIGERATOR"
	DeviceTypeScene               DeviceType = "action.devices.types.SCENE"
	DeviceTypeSecuritySystem      DeviceType = "action.devices.types.SECURITYSYSTEM"
	DeviceTypeShutter             DeviceType = "action.devices.types.SHUTTER"
	DeviceTypeShower              DeviceType = "action.devices.types.SHOWER"
	DeviceTypeSousVide            DeviceType = "action.devices.types.SOUSVIDE"
	DeviceTypeSprinkler           DeviceType = "action.devices.types.SPRINKLER"
	DeviceTypeStandMixer          DeviceType = "action.devices.types.STANDMIXER"
	DeviceTypeVacuum              DeviceType = "action.devices.types.VACUUM"
	DeviceTypeValve               DeviceType = "action.devices.types.VALVE"
	DeviceTypeWasher              DeviceType = "action.devices.types.WASHER"
	DeviceTypeWaterHeater         DeviceType = "action.devices.types.WATERHEATER"
	DeviceTypeWindow              DeviceType = "action.devices.types.WINDOW"
	DeviceTypeYogurtMaker         DeviceType = "action.devices.types.YOGURTMAKER" // Who doesn't own one of these.
)
