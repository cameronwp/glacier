package baseconv

import "testing"

func TestErrors(t *testing.T) {
	tests := []struct {
		val, from, to string
	}{
		{"", DigitsHex, DigitsDec},
		{"0", "", DigitsDec},
		{"0", DigitsDec, ""},
		{"bad", DigitsBin, DigitsDec},
		{"bad", DigitsHex, DigitsDec},
	}

	for i, test := range tests {
		_, err := Convert(test.val, test.from, test.to)
		if err == nil {
			t.Errorf("test %d Convert(%s, %s, %s) should produce error", i, test.val, test.from, test.to)
		}
	}
}

var (
	DigitsJapanese = `〇一二三四五六七八九`
	DigitsThai     = `๐๑๒๓๔๕๖๗๘๙`
)

func TestConvert(t *testing.T) {
	tests := []struct {
		from, to, val, exp string
	}{
		{DigitsDec, DigitsBin, "0", "0"},
		{DigitsDec, DigitsBin, "8", "1000"},
		{DigitsDec, DigitsBin, "15", "1111"},
		{DigitsDec, DigitsBin, "16", "10000"},
		{DigitsDec, DigitsBin, "88", "1011000"},
		{DigitsDec, DigitsBin, "10000", "10011100010000"},

		{DigitsDec, DigitsHex, "0", "0"},
		{DigitsDec, DigitsHex, "8", "8"},
		{DigitsDec, DigitsHex, "15", "F"},
		{DigitsDec, DigitsHex, "16", "10"},
		{DigitsDec, DigitsHex, "88", "58"},
		{DigitsDec, DigitsHex, "10000", "2710"},

		{DigitsDec, Digits62, "16571982744576742462", "JkBr7U8j5pu"},
		{DigitsDec, Digits62, "46394851265279874948", "tHEuTuu3MIe"},
		{DigitsDec, Digits62, "21901407667833273510", "Q5sg7u76TLS"},
		{DigitsDec, Digits62, "8232087098322120342", "9o72rlp5Ff4"},
		{DigitsDec, Digits62, "6354358749246709610", "7ZP1tBlfpP8"},
		{DigitsDec, Digits62, "18089061068", "JkBr7U"},
		{DigitsDec, Digits62, "50642057182", "tHEuTu"},
		{DigitsDec, Digits62, "23906366962", "Q5sg7u"},
		{DigitsDec, Digits62, "8985691605", "9o72rl"},
		{DigitsDec, Digits62, "6936067049", "7ZP1tB"},
		{DigitsDec, Digits62, "799310853702667", "3eyJa0O7P"},

		{DigitsDec, Digits64, "20100203105211888256765428281344829", "zy4Emq2QfCp-XiH3uCz"},
		{DigitsDec, Digits64, "20110423215600563210173308035411215", "z-5EW8kNBfN70ADf2qF"},

		{DigitsHex, DigitsBin, "70B1D707EAC2EDF4C6389F440C7294B51FFF57BB", "111000010110001110101110000011111101010110000101110110111110100110001100011100010011111010001000000110001110010100101001011010100011111111111110101011110111011"},
		{DigitsHex, DigitsBin, "8FC60E7C3B3C48E9A6A7A5FE4F1FBC31", "10001111110001100000111001111100001110110011110001001000111010011010011010100111101001011111111001001111000111111011110000110001"},

		{DigitsHex, Digits36, "ABCDEF00001234567890", "3O47RE02JZQISVIO"},
		{DigitsHex, Digits36, "ABCDEF01234567890123456789ABCDEF", "A65XA07491KF5ZYFPVBO76G33"},

		{Digits62, DigitsHex, "CbAIDLj84gGC5ja7iycGV", "6AD547FFE02477B9473F7977E4D5E17"},
		{Digits62, DigitsHex, "4NIPilGjLxpUTo1HSISijR", "8FC60E7C3B3C48E9A6A7A5FE4F1FBC31"},
		{Digits62, DigitsHex, "4VQYD6oOarxQJ9NrunHTlq", "941532A06BE1443AA9D5D57BDF180A52"},
		{Digits62, DigitsHex, "5fy8kWtSqAuj2kZhjgETFe", "BA86B8F06FDF494487A08A491A19490E"},
		{Digits62, DigitsHex, "7n42DGM5Tflk9n8mt7Fhc7", "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"},

		{DigitsDec, "Christopher", "355927353784509896715106760", "iihtspiphoeCrCeshhorsrrtrh"},

		// unicode
		{DigitsDec, DigitsJapanese, "9876543210", `九八七六五四三二一〇`},
		{DigitsDec, DigitsJapanese, "98765432100123456789", `九八七六五四三二一〇〇一二三四五六七八九`},

		{DigitsDec, DigitsThai, "9876543210", `๙๘๗๖๕๔๓๒๑๐`},

		{DigitsHex, DigitsJapanese, "2710", "一〇〇〇〇"},
		{DigitsHex, DigitsThai, "2710", "๑๐๐๐๐"},
		{DigitsHex, "0一23456789", "2710", "一0000"},
		{ASCII, Digits62, "7n42-&DG$M5Tf@ lk 9n':8mt\\7Fh[c7", "tbWhtdOcCFBgQuGmYPp3uTZzvMSWkZkUI99"},
	}

	for i, test := range tests {
		v0, err := Convert(test.val, test.from, test.to)
		if err != nil {
			t.Fatal(err)
		}
		if test.exp != v0 {
			t.Errorf("test %d (%d->%d) expected %s, got: %s ", i, len(test.from), len(test.to), test.exp, v0)
		}

		v1, err := Convert(test.exp, test.to, test.from)
		if err != nil {
			t.Fatal(err)
		}
		if test.val != v1 {
			t.Errorf("test %d (%d->%d) expected %s, got: %s ", i, len(test.to), len(test.from), test.val, v1)
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	v0 := "1627734050041231452076"

	var tests = []struct {
		encode func(string, ...string) (string, error)
		decode func(string, ...string) (string, error)
		exp    string
	}{
		{EncodeBin, DecodeBin, "10110000011110101011001000001100110100001101011100000100010011110101100"},
		{EncodeOct, DecodeOct, "260365310146415340423654"},
		{EncodeHex, DecodeHex, "583D5906686B8227AC"},
		{Encode36, Decode36, "9JIRD8FBZKUI7G"},
		{Encode62, Decode62, "VHOZDWl3wc8a"},
		{Encode64, Decode64, "M3rP1cXhWYUi"},
	}

	for i, test := range tests {
		v1, err := test.encode(v0)
		if err != nil {
			t.Fatal(err)
		}
		if test.exp != v1 {
			t.Errorf("test %d values %s / %s should match", i, test.exp, v1)
		}

		v2, err := test.decode(v1)
		if err != nil {
			t.Fatal(err)
		}
		if v0 != v2 {
			t.Errorf("test %d values %s / %s should match", i, v0, v2)
		}

		v3, err := test.encode(v0, DigitsDec)
		if err != nil {
			t.Fatal(err)
		}
		if test.exp != v3 {
			t.Errorf("test %d values %s / %s should match", i, test.exp, v3)
		}

		v4, err := test.decode(v1, DigitsDec)
		if err != nil {
			t.Fatal(err)
		}
		if v0 != v4 {
			t.Errorf("test %d values %s / %s should match", i, v0, v4)
		}
	}
}

func BenchmarkConvert(b *testing.B) {
	tests := []struct {
		from, to, val, exp string
	}{
		{DigitsDec, DigitsBin, "0", "0"},
		{DigitsDec, DigitsBin, "8", "1000"},
		{DigitsDec, DigitsBin, "15", "1111"},
		{DigitsDec, DigitsBin, "16", "10000"},
		{DigitsDec, DigitsBin, "88", "1011000"},
		{DigitsDec, DigitsBin, "10000", "10011100010000"},

		{DigitsDec, DigitsHex, "0", "0"},
		{DigitsDec, DigitsHex, "8", "8"},
		{DigitsDec, DigitsHex, "15", "f"},
		{DigitsDec, DigitsHex, "16", "10"},
		{DigitsDec, DigitsHex, "88", "58"},
		{DigitsDec, DigitsHex, "10000", "2710"},

		{DigitsDec, Digits62, "16571982744576742462", "JkBr7U8j5pu"},
		{DigitsDec, Digits62, "46394851265279874948", "tHEuTuu3MIe"},
		{DigitsDec, Digits62, "21901407667833273510", "Q5sg7u76TLS"},
		{DigitsDec, Digits62, "8232087098322120342", "9o72rlp5Ff4"},
		{DigitsDec, Digits62, "6354358749246709610", "7ZP1tBlfpP8"},
		{DigitsDec, Digits62, "18089061068", "JkBr7U"},
		{DigitsDec, Digits62, "50642057182", "tHEuTu"},
		{DigitsDec, Digits62, "23906366962", "Q5sg7u"},
		{DigitsDec, Digits62, "8985691605", "9o72rl"},
		{DigitsDec, Digits62, "6936067049", "7ZP1tB"},
		{DigitsDec, Digits62, "799310853702667", "3eyJa0O7P"},

		{DigitsDec, Digits64, "20100203105211888256765428281344829", "zy4Emq2QfCp-XiH3uCz"},
		{DigitsDec, Digits64, "20110423215600563210173308035411215", "z-5EW8kNBfN70ADf2qF"},

		{DigitsHex, DigitsBin, "70B1D707EAC2EDF4C6389F440C7294B51FFF57BB", "111000010110001110101110000011111101010110000101110110111110100110001100011100010011111010001000000110001110010100101001011010100011111111111110101011110111011"},
		{DigitsHex, DigitsBin, "8FC60E7C3B3C48E9A6A7A5FE4F1FBC31", "10001111110001100000111001111100001110110011110001001000111010011010011010100111101001011111111001001111000111111011110000110001"},

		{DigitsHex, Digits36, "ABCDEF00001234567890", "3O47RE02JZQISVIO"},
		{DigitsHex, Digits36, "ABCDEF01234567890123456789ABCDEF", "A65XA07491KF5ZYFPVBO76G33"},

		{Digits62, DigitsHex, "CbAIDLj84gGC5ja7iycGV", "6AD547FFE02477B9473F7977E4D5E17"},
		{Digits62, DigitsHex, "4NIPilGjLxpUTo1HSISijR", "8FC60E7C3B3C48E9A6A7A5FE4F1FBC31"},
		{Digits62, DigitsHex, "4VQYD6oOarxQJ9NrunHTlq", "941532A06BE1443AA9D5D57BDF180A52"},
		{Digits62, DigitsHex, "5fy8kWtSqAuj2kZhjgETFe", "BA86B8F06FDF494487A08A491A19490E"},
		{Digits62, DigitsHex, "7n42DGM5Tflk9n8mt7Fhc7", "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"},

		{DigitsDec, "Christopher", "355927353784509896715106760", "iihtspiphoeCrCeshhorsrrtrh"},
	}

	for n := 0; n < b.N; n++ {
		for i, test := range tests {
			v0, err := Convert(test.val, test.from, test.to)
			if err != nil {
				b.Fatal(err)
			}
			if test.exp != v0 {
				b.Errorf("test %d (%d->%d) expected %s, got: %s ", i, len(test.from), len(test.to), test.exp, v0)
			}
		}
	}
}
