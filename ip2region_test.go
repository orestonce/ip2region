package ip2region

import (
	"crypto/md5"
	"encoding/hex"
	"testing"
)

func TestGetIpInfo(t *testing.T) {
	info, err := GetIpInfo("1.1.1.1")
	if err != nil {
		t.Fatal(err)
	}
	if info.Country != "澳大利亚" {
		t.Fatal(info.String())
	}
}

func BenchmarkGetIpInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetIpInfo("1.1.1.1")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetIpInfoFromUint32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetIpInfoFromUint32(0x12345678)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestGetIpInfoFromUint32(t *testing.T) {
	//ticker := time.NewTicker(time.Second * 5)
	//for i := uint32(0); i < math.MaxUint32; i++ {
	//	_, err := GetIpInfoFromUint32(i)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	select {
	//	case <-ticker.C:
	//		fmt.Println(i, strconv.FormatFloat(float64(i)*100/float64(math.MaxUint32), 'f', 2, 64)+"%")
	//	default:
	//	}
	//}
	//_, err := GetIpInfoFromUint32(math.MaxUint32)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//ticker.Stop()
	tmp := md5.Sum(gIp2RegionDb)
	if hex.EncodeToString(tmp[:]) != `b986251e5e8f59df915f5c7367ac626a` {
		t.Fatal("untested db data")
	}
}
