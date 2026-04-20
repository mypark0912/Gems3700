package data

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Asset type constants
const (
	AssetCompressor         = 1
	AssetFan                = 2
	AssetPump               = 3
	AssetMotor              = 4
	AssetMotorFeed          = 5
	AssetPowerSupply        = 6
	AssetTransformer        = 7
	AssetPrimaryTransformer = 8
)

// Fixed total size for diagnosis data (1 byte asset_type + 17 bytes fields max)
const diagnosisFixedSize = 18

// BinaryPacker provides methods to pack diagnosis data into binary format.
// All values are packed as uint8 (1 byte each).
type BinaryPacker struct{}

// PackPQData packs PQ data from a map into binary (15 x uint8).
func (bp *BinaryPacker) PackPQData(dataDict map[string]int) []byte {
	fields := []string{
		"VoltagePhaseAngleL", "CurrentRMS", "CrestFactor",
		"Unbalance", "DCCurrent", "Harmonics",
		"ZeroSequence", "NegativeSequence", "CurrentPhaseAngle",
		"PhaseAngle", "PowerFactor", "TotalDemandDistortion",
		"Power", "VoltageRMS", "DCVoltage",
	}
	buf := make([]byte, len(fields))
	for i, f := range fields {
		buf[i] = uint8(dataDict[f])
	}
	return buf
}

// PackFaultData packs Fault data from a map into binary (9 x uint8).
func (bp *BinaryPacker) PackFaultData(dataDict map[string]int) []byte {
	fields := []string{
		"PhaseOrder", "NoLoad", "OverCurrent",
		"NoPower", "OverVoltage", "UnderVoltage",
		"LowFrequency", "CF", "VF",
	}
	buf := make([]byte, len(fields))
	for i, f := range fields {
		buf[i] = uint8(dataDict[f])
	}
	return buf
}

// PackEventData packs Event data from a map into binary (7 x uint8).
func (bp *BinaryPacker) PackEventData(dataDict map[string]int) []byte {
	fields := []string{
		"TransientCurrentEvent",
		"OverCurrentEvent",
		"UnderCurrentEvent",
		"SagEvent",
		"SwellEvent",
		"InterruptionEvent",
		"TransientVoltageEvent",
	}
	buf := make([]byte, len(fields))
	for i, f := range fields {
		buf[i] = uint8(dataDict[f])
	}
	return buf
}

// PackDiagnosisData packs diagnosis data with 1 byte asset_type prefix,
// followed by field values as uint8, padded to a fixed size per asset type.
func (bp *BinaryPacker) PackDiagnosisData(dataDict map[string]int, assetType int) ([]byte, error) {
	var fields []string

	switch assetType {
	case AssetCompressor:
		fields = []string{
			"Turbulence", "Blade", "MechanicalUnbalance", "SoftFoot",
			"TorqueRipple", "DCLink", "Rectifier", "Switching",
			"Load", "GroundFault", "VIUnbalance", "CableConnection",
			"NoiseVibration", "Heat", "Rotor", "Stator", "Bearing",
		}
	case AssetFan:
		fields = []string{
			"Turbulence", "Blade", "MechanicalUnbalance", "SoftFoot",
			"TorqueRipple", "DCLink", "Rectifier", "Switching",
			"Load", "GroundFault", "VIUnbalance", "CableConnection",
			"NoiseVibration", "Heat", "Rotor", "Stator", "Bearing",
		}
	case AssetPump:
		fields = []string{
			"Cavitation", "Vane", "MechanicalUnbalance", "SoftFoot",
			"TorqueRipple", "DCLink", "Rectifier", "Switching",
			"Load", "GroundFault", "VIUnbalance", "CableConnection",
			"NoiseVibration", "Heat", "Rotor", "Stator", "Bearing",
		}
	case AssetMotor:
		// 15 fields + 2 padding
		fields = []string{
			"TorqueRipple", "MechanicalUnbalance", "SoftFoot",
			"DCLink", "Rectifier", "Switching", "Load",
			"GroundFault", "VIUnbalance", "CableConnection",
			"NoiseVibration", "Heat", "Rotor", "Stator", "Bearing",
		}
	case AssetMotorFeed:
		// 9 fields + 8 padding
		fields = []string{
			"DCLink", "Rectifier", "Switching", "Load",
			"GroundFault", "VIUnbalance", "CableConnection",
			"NoiseVibration", "Heat",
		}
	case AssetPowerSupply:
		// 6 fields + 11 padding
		fields = []string{
			"Load", "GroundFault", "VIUnbalance",
			"CableConnection", "NoiseVibration", "Heat",
		}
	case AssetTransformer:
		// 14 fields + 3 padding
		fields = []string{
			"Capacitor", "Core", "TapChanger", "Bushings",
			"Stress", "Winding", "Load", "Rectifier",
			"Switching", "GroundFault", "VIUnbalance",
			"CableConnection", "NoiseVibration", "Heat",
		}
	case AssetPrimaryTransformer:
		// 10 fields + 7 padding
		fields = []string{
			"Capacitor", "TapChanger", "Bushings", "Stress",
			"Load", "GroundFault", "VIUnbalance",
			"CableConnection", "NoiseVibration", "Heat",
		}
	default:
		return nil, fmt.Errorf("unknown asset type: %d", assetType)
	}

	// 1 byte asset type + 17 bytes data area (fields + padding)
	buf := make([]byte, diagnosisFixedSize)
	buf[0] = uint8(assetType)

	for i, f := range fields {
		buf[1+i] = uint8(dataDict[f])
	}
	// Remaining bytes in buf are already zero (padding)

	return buf, nil
}

// SaveDataToRedis packs data and stores it in Redis via HSET.
func SaveDataToRedis(
	redisClient *redis.Client,
	redisKey string,
	dataType string,
	dataDict map[string]int,
	assetType int,
) error {
	packer := &BinaryPacker{}
	var binaryData []byte

	switch dataType {
	case "PQ":
		binaryData = packer.PackPQData(dataDict)
	case "Fault":
		binaryData = packer.PackFaultData(dataDict)
	case "Event":
		binaryData = packer.PackEventData(dataDict)
	case "Diagnosis":
		if assetType == 0 {
			assetType = AssetMotor
		}
		var err error
		binaryData, err = packer.PackDiagnosisData(dataDict, assetType)
		if err != nil {
			return fmt.Errorf("pack diagnosis data: %w", err)
		}
	default:
		return fmt.Errorf("unknown data type: %s", dataType)
	}

	ctx := context.Background()
	if err := redisClient.HSet(ctx, redisKey, dataType, binaryData).Err(); err != nil {
		log.Printf("Error saving data to Redis key=%s field=%s: %v", redisKey, dataType, err)
		return fmt.Errorf("save to Redis: %w", err)
	}

	log.Printf("%s data saved (size: %d bytes)", dataType, len(binaryData))
	return nil
}
