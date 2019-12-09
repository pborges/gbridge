package proto

import "errors"

type CommandStatus string

const (
	CommandStatusSuccess CommandStatus = "SUCCESS"
	CommandStatusError   CommandStatus = "ERROR"
)

type ErrorCode error

var (
	ErrorCodeAuthExpired     ErrorCode = errors.New("authExpired")
	ErrorCodeAuthFailure     ErrorCode = errors.New("authFailure")
	ErrorCodeDeviceOffline   ErrorCode = errors.New("deviceOffline")
	ErrorCodeTimeout         ErrorCode = errors.New("timeout")
	ErrorCodeDeviceTurnedOff ErrorCode = errors.New("deviceTurnedOff")
	ErrorCodeDeviceNotFound  ErrorCode = errors.New("deviceNotFound")
	ErrorCodeValueOutofRange ErrorCode = errors.New("valueOutOfRange")
	ErrorCodeNotSupported    ErrorCode = errors.New("notSupported")
	ErrorCodeProtocolError   ErrorCode = errors.New("protocolError")
	ErrorCodeUnknown         ErrorCode = errors.New("unknownError")
)

type DeviceError error

var (
	DeviceErrorAboveMaximumLightEffectsDuration DeviceError = errors.New("aboveMaximumLightEffectsDuration") // That's more than the maximum duration of 1 hour. Please try again.
	DeviceErrorAboveMaximumTimerDuration        DeviceError = errors.New("aboveMaximumTimerDuration")        // I can only set <device(s)> for up to <time period>
	DeviceErrorActionNotAvailable               DeviceError = errors.New("actionNotAvailable")               // Sorry, I can't seem to do that right now.
	DeviceErrorActionUnavailableWhileRunning    DeviceError = errors.New("actionUnavailableWhileRunning")    // <device(s)> <is/are> currently running, so I can't make any changes.
	DeviceErrorAlreadyArmed                     DeviceError = errors.New("alreadyArmed")                     // <device(s)> <is/are> already armed.
	DeviceErrorAlreadyAtMax                     DeviceError = errors.New("alreadyAtMax")                     // <device(s)> <is/are> already set to the maximum temperature.
	DeviceErrorAlreadyAtMin                     DeviceError = errors.New("alreadyAtMin")                     // <device(s)> <is/are> already set to the minimum temperature.
	DeviceErrorAlreadyClosed                    DeviceError = errors.New("alreadyClosed")                    // <device(s)> <is/are> already closed.
	DeviceErrorAlreadyDisarmed                  DeviceError = errors.New("alreadyDisarmed")                  // <device(s)> <is/are> already disarmed.
	DeviceErrorAlreadyDocked                    DeviceError = errors.New("alreadyDocked")                    // <device(s)> <is/are> already docked.
	DeviceErrorAlreadyInState                   DeviceError = errors.New("alreadyInState")                   // <device(s)> <is/are> already in that state.
	DeviceErrorAlreadyLocked                    DeviceError = errors.New("alreadyLocked")                    // <device(s)> <is/are> already locked.
	DeviceErrorAlreadyOff                       DeviceError = errors.New("alreadyOff")                       // <device(s)> <is/are> already off.
	DeviceErrorAlreadyOn                        DeviceError = errors.New("alreadyOn")                        // <device(s)> <is/are> already on.
	DeviceErrorAlreadyOpen                      DeviceError = errors.New("alreadyOpen")                      // <device(s)> <is/are> already open.
	DeviceErrorAlreadyPaused                    DeviceError = errors.New("alreadyPaused")                    // <device(s)> <is/are> already paused.
	DeviceErrorAlreadyStarted                   DeviceError = errors.New("alreadyStarted")                   // <device(s)> <is/are> already started.
	DeviceErrorAlreadyStopped                   DeviceError = errors.New("alreadyStopped")                   // <device(s)> <is/are> already stopped.
	DeviceErrorAlreadyUnlocked                  DeviceError = errors.New("alreadyUnlocked")                  // <device(s)> <is/are> already unlocked.
	DeviceErrorAmountAboveLimit                 DeviceError = errors.New("amountAboveLimit")                 // That's more than what <device(s)> can support.
	DeviceErrorAppLaunchFailed                  DeviceError = errors.New("appLaunchFailed")                  // Sorry, failed to launch <app name> on <device(s)>.
	DeviceErrorArmFailure                       DeviceError = errors.New("armFailure")                       // <device(s)> couldn't be armed.
	DeviceErrorArmLevelNeeded                   DeviceError = errors.New("armLevelNeeded")                   // I'm not sure which level to set <device(s)> to. Try saying "Set <device(s)> to <low security>" or "Set <device(s)> to <high security>"
	DeviceErrorAuthFailure                      DeviceError = errors.New("authFailure")                      // I can't seem to reach your <device(s)>. Try checking the app to make sure your device is fully set up.
	DeviceErrorBagFull                          DeviceError = errors.New("bagFull")                          // <device(s)> <has/have> <a full bag/full bags>. Please empty <it/them> and try again.
	DeviceErrorBelowMinimumLightEffectsDuration DeviceError = errors.New("belowMinimumLightEffectsDuration") // That's less than the minimum duration of 5 minutes. Please try again.
	DeviceErrorBelowMinimumTimerDuration        DeviceError = errors.New("belowMinimumTimerDuration")        // I can't set <device(s)> for such a short time. Please try again.
	DeviceErrorBinFull                          DeviceError = errors.New("binFull")                          // <device(s)> <has/have> <a full bin/full bins>.
	DeviceErrorCancelArmingRestricted           DeviceError = errors.New("cancelArmingRestricted")           // Sorry, I could not cancel arming <device(s)>.
	DeviceErrorCommandInsertFailed              DeviceError = errors.New("commandInsertFailed")              // Unable to process commands for <device(s)>.
	DeviceErrorDegreesOutOfRange                DeviceError = errors.New("degreesOutOfRange")                // The requested degrees are out of range for <device(s)>.
	DeviceErrorDeviceDoorOpen                   DeviceError = errors.New("deviceDoorOpen")                   // The door is open on <device(s)>. Please close it and try again.
	DeviceErrorDeviceHandleClosed               DeviceError = errors.New("deviceHandleClosed")               // The handle is closed on <device(s)>. Please open it and try again.
	DeviceErrorDeviceJammingDetected            DeviceError = errors.New("deviceJammingDetected")            // <device(s)> <is/are> jammed.
	DeviceErrorDeviceLidOpen                    DeviceError = errors.New("deviceLidOpen")                    // The lid is open on <device(s)>. Please close it and try again.
	DeviceErrorDeviceNotDocked                  DeviceError = errors.New("deviceNotDocked")                  // Sorry, it looks like <device(s)> <isn't/aren't> docked. Please dock <it/them> and try again.
	DeviceErrorDeviceNotFound                   DeviceError = errors.New("deviceNotFound")                   // <device(s)> <is/are>n't available. You might want to try setting <it/them> up again.
	DeviceErrorDeviceNotReady                   DeviceError = errors.New("deviceNotReady")                   // <device(s)> <is/are>n't ready.
	DeviceErrorDeviceStuck                      DeviceError = errors.New("deviceStuck")                      // <device(s)> <is/are> stuck.
	DeviceErrorDeviceTampered                   DeviceError = errors.New("deviceTampered")                   // <device(s)> <has/have> been tampered with.
	DeviceErrorDirectResponseOnlyUnreachable    DeviceError = errors.New("directResponseOnlyUnreachable")    // <device(s)> <doesn't/don't> support remote control.
	DeviceErrorDisarmFailure                    DeviceError = errors.New("disarmFailure")                    // <device(s)> couldn't be disarmed.
	DeviceErrorDoorClosedTooLong                DeviceError = errors.New("doorClosedTooLong")                // It's been a while since the door on <device(s)> has been opened. Please open the door, make sure there's something inside, and try again.
	DeviceErrorEmergencyHeatOn                  DeviceError = errors.New("emergencyHeatOn")                  // <device(s)> <is/are> in Emergency Heat Mode, so <it/they>'ll have to be adjusted by hand.
	DeviceErrorFloorUnreachable                 DeviceError = errors.New("floorUnreachable")                 // <device(s)> can't reach that room. Please move <it/them> to the right floor and try again.
	DeviceErrorFunctionNotSupported             DeviceError = errors.New("functionNotSupported")             // Actually, that functionality isn't supported yet.
	DeviceErrorHardError                        DeviceError = errors.New("hardError")                        // Sorry, something went wrong and I'm unable to control your home device.
	DeviceErrorInAutoMode                       DeviceError = errors.New("inAutoMode")                       // <device(s)> <is/are> currently set to auto mode. To change the temperature, you'll need to switch <it/them> to a different mode.
	DeviceErrorInAwayMode                       DeviceError = errors.New("inAwayMode")                       // <device(s)> <is/are> currently set to away mode. To control your thermostat, you'll need to manually switch it to home mode using the Nest app on a phone, tablet, or computer.
	DeviceErrorInDryMode                        DeviceError = errors.New("inDryMode")                        // <device(s)> <is/are> currently set to dry mode. To change the temperature, you'll need to switch <it/them> to a different mode.
	DeviceErrorInEcoMode                        DeviceError = errors.New("inEcoMode")                        // <device(s)> <is/are> currently set to eco mode. To change the temperature, you'll need to switch <it/them> to a different mode.
	DeviceErrorInFanOnlyMode                    DeviceError = errors.New("inFanOnlyMode")                    // <device(s)> <is/are> currently set to fan-only mode. To change the temperature, you'll need to switch <it/them> to a different mode.
	DeviceErrorInHeatOrCool                     DeviceError = errors.New("inHeatOrCool")                     // <device(s)> <is/are>n't in heat/cool mode.
	DeviceErrorInHumidifierMode                 DeviceError = errors.New("inHumidifierMode")                 // <device(s)> <is/are> currently set to humidifier mode. To change the temperature, you'll need to switch <it/them> to a different mode.
	DeviceErrorInOffMode                        DeviceError = errors.New("inOffMode")                        // <device(s)> <is/are> currently off. To change the temperature, you'll need to switch <it/them> to a different mode.
	DeviceErrorInPurifierMode                   DeviceError = errors.New("inPurifierMode")                   // <device(s)> <is/are> currently set to purifier mode. To change the temperature, you'll need to switch <it/them> to a different mode.
	DeviceErrorInSleepMode                      DeviceError = errors.New("inSleepMode")                      // <device(s)> <is/are> in sleep mode. Please try again later.
	DeviceErrorLockFailure                      DeviceError = errors.New("lockFailure")                      // <device(s)> couldn't be locked.
	DeviceErrorLockedToRange                    DeviceError = errors.New("lockedToRange")                    // That temperature is outside the locked range on <device(s)>.
	DeviceErrorLowBattery                       DeviceError = errors.New("lowBattery")                       // <device(s)> <has/have> low battery.
	DeviceErrorMaxSettingReached                DeviceError = errors.New("maxSettingReached")                // <device(s)> <is/are> already set to the highest setting.
	DeviceErrorMaxSpeedReached                  DeviceError = errors.New("maxSpeedReached")                  // <device(s)> <is/are> already set to the maximum speed.
	DeviceErrorMinSettingReached                DeviceError = errors.New("minSettingReached")                // <device(s)> <is/are> already set to the lowest setting.
	DeviceErrorMinSpeedReached                  DeviceError = errors.New("minSpeedReached")                  // <device(s)> <is/are> already set to the minimum speed.
	DeviceErrorMissingSubscription              DeviceError = errors.New("missingSubscription")              // Sorry, it looks like you don't have an active subscription to this service. Please activate your subscription and try again.
	DeviceErrorNeedsAttachment                  DeviceError = errors.New("needsAttachment")                  // Sorry, it looks like <device(s)> <is/are> missing a required attachment. Please replace it and try again.
	DeviceErrorNeedsBin                         DeviceError = errors.New("needsBin")                         // Sorry, it looks like <device(s)> <is/are> missing a bin. Please replace it and try again.
	DeviceErrorNeedsPads                        DeviceError = errors.New("needsPads")                        // <device(s)> <need(s)> new pads.
	DeviceErrorNeedsSoftwareUpdate              DeviceError = errors.New("needsSoftwareUpdate")              // <device(s)> <need(s)> a software update.
	DeviceErrorNeedsWater                       DeviceError = errors.New("needsWater")                       // <device(s)> <need(s)> water.
	DeviceErrorNoAvailableApp                   DeviceError = errors.New("noAvailableApp")                   // Sorry, it looks like <app name> isn't available.
	DeviceErrorNotSupported                     DeviceError = errors.New("notSupported")                     // Sorry, that mode isn't available for <device(s)>.
	DeviceErrorObstructionDetected              DeviceError = errors.New("obstructionDetected")              // <device(s)> detected an obstruction.
	DeviceErrorOffline                          DeviceError = errors.New("offline")                          // deviceOffline : Sorry, it looks like <device(s)> <is/are>n't available right now.
	DeviceErrorOnRequiresMode                   DeviceError = errors.New("onRequiresMode")                   // Please specify which mode you want to turn on.
	DeviceErrorPassphraseIncorrect              DeviceError = errors.New("passphraseIncorrect")              // Sorry, it looks like that security code is incorrect.
	DeviceErrorPinIncorrect                     DeviceError = errors.New("pinIncorrect")                     // (passphraseIncorrect)
	DeviceErrorRangeTooClose                    DeviceError = errors.New("rangeTooClose")                    // Those are too close for a Heat-Cool range for <device(s)>. Choose temperatures that are farther apart.
	DeviceErrorRelinkRequired                   DeviceError = errors.New("relinkRequired")                   // Sorry, it looks like something went wrong with your account. Please use your Google Home or Assistant App to re-link <device(s)>.
	DeviceErrorRoomsOnDifferentFloors           DeviceError = errors.New("roomsOnDifferentFloors")           // <device(s)> can't reach those rooms because they're on different floors.
	DeviceErrorSafetyShutOff                    DeviceError = errors.New("safetyShutOff")                    // <device(s)> <is/are> in Safety Shut-Off Mode, so <it/they>'ll have to be adjusted by hand.
	DeviceErrorSecurityRestriction              DeviceError = errors.New("securityRestriction")              // <device(s)> <has/have> a security restriction.
	DeviceErrorSoftwareUpdateNotAvailable       DeviceError = errors.New("softwareUpdateNotAvailable")       // Sorry, there's no software update available on <device(s)>.
	DeviceErrorStartRequiresTime                DeviceError = errors.New("startRequiresTime")                // To do that, you'll need to tell me how long you'd like to run <device(s)>.
	DeviceErrorStillWarmingUp                   DeviceError = errors.New("stillWarmingUp")                   // <device(s)> <is/are> still warming up.
	DeviceErrorTankEmpty                        DeviceError = errors.New("tankEmpty")                        // <device(s)> <has/have> <an empty tank/empty tanks>. Please fill <it/them> and try again.
	DeviceErrorTargetAlreadyReached             DeviceError = errors.New("targetAlreadyReached")             // Sorry, it looks like that's already the current temperature.
	DeviceErrorTimerValueOutOfRange             DeviceError = errors.New("timerValueOutOfRange")             // <device(s)> can't be set for that amount of time.
	DeviceErrorTransientError                   DeviceError = errors.New("transientError")                   // Sorry, something went wrong controlling <device(s)>. Please try again.
	DeviceErrorTurnedOff                        DeviceError = errors.New("turnedOff")                        // deviceTurnedOff : <device(s)> <is/are> off right now, so I can't make any adjustments.
	DeviceErrorUnableToLocateDevice             DeviceError = errors.New("unableToLocateDevice")             // I wasn't able to locate <device(s)>.
	DeviceErrorUnknownFoodPreset                DeviceError = errors.New("unknownFoodPreset")                // <device(s)> doesn't support that food preset.
	DeviceErrorUnlockFailure                    DeviceError = errors.New("unlockFailure")                    // <device(s)> couldn't be unlocked.
	DeviceErrorValueOutOfRange                  DeviceError = errors.New("valueOutOfRange")                  // <device(s)> can't be set to that temperature.
)
