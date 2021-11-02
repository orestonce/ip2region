package ip2region

import (
	_ "embed"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
	"sync"
)

//go:embed ip2region.db
var gIp2RegionDb []byte

const (
	cIndexBlockLength = 12
)

var gIp2Region struct {
	firstIndexPtr uint32 // super block index info
	lastIndexPtr  uint32
	totalBlocks   uint32
	initOnce      sync.Once
}

type IpInfo struct {
	Country  string
	Region   string
	Province string
	City     string
	ISP      string
}

func (ip IpInfo) String() string {
	return ip.Country + "|" + ip.Region + "|" + ip.Province + "|" + ip.City + "|" + ip.ISP
}

func GetIpInfo(ipStr string) (ipInfo IpInfo, err error) {
	ip, err := ip2uint32(ipStr)
	if err != nil {
		return ipInfo, err
	}
	return GetIpInfoFromUint32(ip)
}

func GetIpInfoFromUint32(u32ip uint32) (ipInfo IpInfo, err error) {
	gIp2Region.initOnce.Do(func() {
		gIp2Region.firstIndexPtr = getUint32(gIp2RegionDb, 0)
		gIp2Region.lastIndexPtr = getUint32(gIp2RegionDb, 4)
		gIp2Region.totalBlocks = (gIp2Region.lastIndexPtr-gIp2Region.firstIndexPtr)/cIndexBlockLength + 1
	})
	h := gIp2Region.totalBlocks
	var dataPtr, l uint32
	for l <= h {
		m := (l + h) / 2
		p := gIp2Region.firstIndexPtr + m*cIndexBlockLength
		sip := getUint32(gIp2RegionDb, p)
		eip := getUint32(gIp2RegionDb, p+4)
		if u32ip < sip {
			h = m - 1
		} else if u32ip <= eip {
			dataPtr = getUint32(gIp2RegionDb, p+8)
			break
		} else {
			l = m + 1
		}
	}
	if dataPtr <= 0 {
		return ipInfo, errors.New("ip2region internal error: invalid dataPtr " + strconv.Itoa(int(dataPtr)))
	}
	dataLen := (dataPtr >> 24) & 0xFF
	dataPtr = dataPtr & 0x00FFFFFF
	line := gIp2RegionDb[(dataPtr)+4 : dataPtr+dataLen]
	lineSlice := strings.Split(string(line), "|")
	length := len(lineSlice)
	if length < 5 {
		return ipInfo, errors.New(`ip2region internal error: invalid line ` + string(line))
	}
	getInfo := func(s string) string {
		if s == `0` {
			return ``
		}
		return s
	}
	ipInfo.Country = getInfo(lineSlice[0])
	ipInfo.Region = getInfo(lineSlice[1])
	ipInfo.Province = getInfo(lineSlice[2])
	ipInfo.City = getInfo(lineSlice[3])
	ipInfo.ISP = getInfo(lineSlice[4])
	return ipInfo, nil
}

func ip2uint32(IpStr string) (uint32, error) {
	ErrInvalidIp := errors.New("ip2region: invalid ipv4 " + strconv.Quote(IpStr))

	bits := strings.Split(IpStr, ".")
	if len(bits) != 4 {
		return 0, ErrInvalidIp
	}
	var sum uint32
	for i, n := range bits {
		bit, err := strconv.Atoi(n)
		if err != nil {
			return 0, ErrInvalidIp
		}
		if bit < 0 || bit > 255 {
			return 0, ErrInvalidIp
		}
		sum += uint32(bit) << uint32(24-8*i)
	}
	return sum, nil
}

func getUint32(b []byte, offset uint32) uint32 {
	return binary.LittleEndian.Uint32(b[offset:])
}
