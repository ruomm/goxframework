/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/6/20 13:40
 * @version 1.0
 */
package encryptx

import (
	"crypto"
	"fmt"
	"testing"
)

const (
	//PUB_KETY = "-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEAyYTKmKEUE4scVNQzmFUSKMlvSKOUxjrymSZ6K0SyuiHY6dplz89e\nYkkbz95MwjihQ6oA2NDB++uCfbm6rGOSJ1U8/prlfuIvbtiySupQonFj0FTWWQpr\n5ZydRhxSsbsYJH8hdcYPMTdwXfJunWhpJ9gRZWyRdR5EymAUdSS0zXjJZdXWhbId\npyOX7mpA1Lhk+8cU6jtqAuP+6CirZB/NGrlhnzbv8TqPTrB132YvsjdPm+CcKSK8\n2176aAFhdau8VNYdf8jCtz0qyQO/yfy+FdoYAnKF4yJBOFAcYRN5LdIGSg5aLeWV\n5wEx4xu4KXzNoKsCvAltvVfzL45xzdlVKQIDAQAB\n-----END RSA PUBLIC KEY-----"
	//PRI_
	//PRI_KEY = "308204a50201000282010100ca98cf666dd78c03b55475fa911d6f7a8b23a545d3bcd668a9f55e61e868a314e0fb4aa89d8740b12b33fc1e045db48f56ce87f946e1588b80f4ce396bdb0c7140235fdd8637bb81a2cf6913da10dfc16e4de6101e1972f2e573296fe0857e2b2523f194757b2ef6b7f0d6d538bf226cb75f452a643b821fa01fd73cfd5e7d175b0361b8b5a35b213bb9b9005f8f9852fe9ebfb6f29c3e986f7b2f7f8872c30df6dded1fadad9912c89d4e8fe4dfea6b043a3876e2edb0da84504b1d443522034dd145ff9ca4c89be135addabde29d3c550869396923a603d26d865cdf4c34e65f0cc6269340c80f186ec4c1ec8211ec3dda081c71a3e368e47629d8ac68b6d1020301000102820101009aee7d3cf1a732b5eb47a9e9726c36425a8169f49a5601098d5fcd4bc675aeb15ad411661d01bfe13d0ec6311659aaa92b5349fcc1cdb6ef08312e3c0f553690ace1e080021240dc846b6509ff6b8411e28ec3ef737536c8b5db79c6cac68b93e71533bbba93f77849766e7424af565e19654bf97d295cfb5e72bb213309bd52701af1465a34c2278d07d48cdd4954757e7a5db33a1c57ab5f935b2c9e2b069dd91b2cabfc7de70206472e16c738251af239c94b001a8b8e86519d00c4fdbdb6a517ecb33626b5678767262f87e516c00c67384bfd5a153407ff3c4851cf6c679d767189d2f2e9f1eb306c1274e3e6ebabe6b9e1d47f895d8236c5503f7784d902818100f1fd4a7abbeacc28da918a45707e4a4f38ba49f91fd0f638d7011b142310f34882f90b88bb768f3a6f7b79fcd75ad91410ec8c7e3e65e6de70d554c53e7877b6a141c869faa74a6892093194174c5a39859ec7a4081ff8b8c8fd34777c3f4077e3b8f317d54637351f31818f7d5f08f0fe06dc3f0812d8c8f4d197d812692fcf02818100d653a72acaed2dbe7844ba56470bbd06d0045f9b11e65b1efbe4e9e94ccdf15b4f9f57c03b9b44e4b14f3f936c92edd577a8202ee701ee929127cef30c6f9d06efafabf2d9328168315c3f31130922db132095808cf13c022e64b3f88a8e88359712d0cc6b05a0048ecef47fe03995200d8f887579553a68cafa021b385db75f02818100dc4dcc509063e21a0f62108fc72a325c8d388bbfd1c75b61c3dbaddb5751472aec91ee7e3cac6318c00599b92615ff2ad57d852a29847bfa669ed0de01518b2b2903ca813140bbed1786672c7b49779a869b57056ea02cbf8dbb76d890b4c4ec60d52ffab29f8a0342b2bf50f2c2625183f29af087592123523ebe0e68347ae10281805c9cb32c83996f5dd19c692464f8c68a8e1285b97d479bf24c88832703c02dddf60ef59d246498a57594b7f51d427430fcce927191f2bcc36aa3d802204a8e80f3cb6632bef5db0349e90189541f6b875cc184b892ae9eee965b7f85136239ab668783e00112e22d76042994a4305da7071511b32965d1a27cab0183ff9c476302818100a332b4299537028045cac54147a76e8469b991f62ebba565cebb90748b769e73fc59146cce02ab1ffa8f8a7ac8e19ed0ad3902ff764b1a1408da0d7d68ee5cc2977a82a19c408b7fd8018b1521c5505a65aac5c1b51e2801264e7c9085264a40fcaf4016e2e99e42e8033f84386667fb2b77be62f2b352e4b5fac0218521b9eb"
	//PUB_KEY = "3082010a0282010100ca98cf666dd78c03b55475fa911d6f7a8b23a545d3bcd668a9f55e61e868a314e0fb4aa89d8740b12b33fc1e045db48f56ce87f946e1588b80f4ce396bdb0c7140235fdd8637bb81a2cf6913da10dfc16e4de6101e1972f2e573296fe0857e2b2523f194757b2ef6b7f0d6d538bf226cb75f452a643b821fa01fd73cfd5e7d175b0361b8b5a35b213bb9b9005f8f9852fe9ebfb6f29c3e986f7b2f7f8872c30df6dded1fadad9912c89d4e8fe4dfea6b043a3876e2edb0da84504b1d443522034dd145ff9ca4c89be135addabde29d3c550869396923a603d26d865cdf4c34e65f0cc6269340c80f186ec4c1ec8211ec3dda081c71a3e368e47629d8ac68b6d10203010001"
	PRI_KEY = "MIIEpQIBAAKCAQEAypjPZm3XjAO1VHX6kR1veosjpUXTvNZoqfVeYehooxTg+0qonYdAsSsz/B4EXbSPVs6H+UbhWIuA9M45a9sMcUAjX92GN7uBos9pE9oQ38FuTeYQHhly8uVzKW/ghX4rJSPxlHV7Lva38NbVOL8ibLdfRSpkO4IfoB/XPP1efRdbA2G4taNbITu5uQBfj5hS/p6/tvKcPphvey9/iHLDDfbd7R+trZkSyJ1Oj+Tf6msEOjh24u2w2oRQSx1ENSIDTdFF/5ykyJvhNa3aveKdPFUIaTlpI6YD0m2GXN9MNOZfDMYmk0DIDxhuxMHsghHsPdoIHHGj42jkdinYrGi20QIDAQABAoIBAQCa7n088acytetHqelybDZCWoFp9JpWAQmNX81LxnWusVrUEWYdAb/hPQ7GMRZZqqkrU0n8wc227wgxLjwPVTaQrOHggAISQNyEa2UJ/2uEEeKOw+9zdTbItdt5xsrGi5PnFTO7upP3eEl2bnQkr1ZeGWVL+X0pXPtecrshMwm9UnAa8UZaNMInjQfUjN1JVHV+el2zOhxXq1+TWyyeKwad2Rssq/x95wIGRy4WxzglGvI5yUsAGouOhlGdAMT9vbalF+yzNia1Z4dnJi+H5RbADGc4S/1aFTQH/zxIUc9sZ512cYnS8unx6zBsEnTj5uur5rnh1H+JXYI2xVA/d4TZAoGBAPH9Snq76swo2pGKRXB+Sk84ukn5H9D2ONcBGxQjEPNIgvkLiLt2jzpve3n811rZFBDsjH4+ZebecNVUxT54d7ahQchp+qdKaJIJMZQXTFo5hZ7HpAgf+LjI/TR3fD9Ad+O48xfVRjc1HzGBj31fCPD+Btw/CBLYyPTRl9gSaS/PAoGBANZTpyrK7S2+eES6VkcLvQbQBF+bEeZbHvvk6elMzfFbT59XwDubROSxTz+TbJLt1XeoIC7nAe6SkSfO8wxvnQbvr6vy2TKBaDFcPzETCSLbEyCVgIzxPAIuZLP4io6INZcS0MxrBaAEjs70f+A5lSANj4h1eVU6aMr6Ahs4XbdfAoGBANxNzFCQY+IaD2IQj8cqMlyNOIu/0cdbYcPbrdtXUUcq7JHufjysYxjABZm5JhX/KtV9hSophHv6Zp7Q3gFRiyspA8qBMUC77ReGZyx7SXeahptXBW6gLL+Nu3bYkLTE7GDVL/qyn4oDQrK/UPLCYlGD8prwh1khI1I+vg5oNHrhAoGAXJyzLIOZb13RnGkkZPjGio4Shbl9R5vyTIiDJwPALd32DvWdJGSYpXWUt/UdQnQw/M6ScZHyvMNqo9gCIEqOgPPLZjK+9dsDSekBiVQfa4dcwYS4kq6e7pZbf4UTYjmrZoeD4AES4i12BCmUpDBdpwcVEbMpZdGifKsBg/+cR2MCgYEAozK0KZU3AoBFysVBR6duhGm5kfYuu6VlzruQdIt2nnP8WRRszgKrH/qPinrI4Z7QrTkC/3ZLGhQI2g19aO5cwpd6gqGcQIt/2AGLFSHFUFplqsXBtR4oASZOfJCFJkpA/K9AFuLpnkLoAz+EOGZn+yt3vmLys1LktfrAIYUhues="
	PUB_KEY = "MIIBCgKCAQEAypjPZm3XjAO1VHX6kR1veosjpUXTvNZoqfVeYehooxTg+0qonYdAsSsz/B4EXbSPVs6H+UbhWIuA9M45a9sMcUAjX92GN7uBos9pE9oQ38FuTeYQHhly8uVzKW/ghX4rJSPxlHV7Lva38NbVOL8ibLdfRSpkO4IfoB/XPP1efRdbA2G4taNbITu5uQBfj5hS/p6/tvKcPphvey9/iHLDDfbd7R+trZkSyJ1Oj+Tf6msEOjh24u2w2oRQSx1ENSIDTdFF/5ykyJvhNa3aveKdPFUIaTlpI6YD0m2GXN9MNOZfDMYmk0DIDxhuxMHsghHsPdoIHHGj42jkdinYrGi20QIDAQAB"
)

func TestRsaHelperCommon(t *testing.T) {

	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	var xHelper RsaHelper
	xHelper = &XRsa{
		ModePadding: MODE_PADDING_PKCS5,
	}
	//xHelper.GenrateKeyPair(2048)
	xHelper.LoadPrivateKey(MODE_KEY_BASE64, PRI_KEY)
	xHelper.LoadPulicKey(MODE_KEY_BASE64, PUB_KEY)
	//fmt.Println(xHelper.FormatPrivateKey(MODE_KEY_HEX_UPPER))
	//fmt.Println(xHelper.FormatPublicKey(MODE_KEY_HEX_UPPER))
	//fmt.Println(xHelper.SizeOfPrivateKey())
	//fmt.Println(xHelper.SizeOfPublicKey())
	origStr := generateToken(1024) + "      中华人民共和国      "
	origStr = ""
	//encStr, _ := xHelper.EncryptPKCS1v15StringBig(MODE_ENCODE_BASE64, origStr)
	//fmt.Println(encStr)
	//decStr, _ := xHelper.DecryptPKCS1v15StringBig(MODE_ENCODE_BASE64, encStr)
	//fmt.Println(decStr)
	//if origStr == decStr {
	//	fmt.Println("加密解密验证通过")
	//} else {
	//	fmt.Println("加密解密验证不通过通过")
	//}

	//sigStr, _ := xHelper.SignPSSByString(MODE_ENCODE_BASE64, crypto.MD5SHA1, origStr, nil)
	//fmt.Println(sigStr)
	//verifyErr := xHelper.VerifyPSSByString(MODE_ENCODE_BASE64, crypto.MD5SHA1, origStr, sigStr, nil)
	//if verifyErr == nil && len(sigStr) > 0 {
	//	fmt.Println("签名验证通过")
	//} else {
	//	fmt.Printf("签名验证不通过:%v", verifyErr)
	//}

	sigData, _ := xHelper.SignPSS(crypto.SHA384, []byte(origStr), nil)
	fmt.Println(sigData)
	verifyDataErr := xHelper.VerifyPSS(crypto.SHA384, []byte(origStr), sigData, nil)
	if verifyDataErr == nil && len(sigData) > 0 {
		fmt.Println("签名验证通过")
	} else {
		fmt.Printf("签名验证不通过:%v", verifyDataErr)
	}
}

func TestRsaHelperFile(t *testing.T) {

	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	var xHelper RsaHelper
	xHelper = &XRsa{
		ModePadding: MODE_PADDING_PKCS5,
	}
	//xHelper.GenrateKeyPair(2048)
	xHelper.LoadPrivateKey(MODE_KEY_BASE64, PRI_KEY)
	xHelper.LoadPulicKey(MODE_KEY_BASE64, PUB_KEY)
	origFile := "/Users/qx/Downloads/文本bom测试.txt"
	encFile := "/Users/qx/Downloads/文本bom测试_ENC.txt"
	decFile := "/Users/qx/Downloads/文本bom测试_DEC.txt"
	err := xHelper.EncryptPKCS1v15File(origFile, encFile, false)
	if err == nil {
		fmt.Println("文件加密通过")
	} else {
		fmt.Printf("文件加密不通过:%v", err)
	}
	err = xHelper.DecryptPKCS1v15File(encFile, decFile)
	if err == nil {
		fmt.Println("文件解密通过")
	} else {
		fmt.Printf("文件解密不通过:%v", err)
	}
	sha, _ := SumFileByString(MODE_ENCODE_HEX_LOWER, crypto.MD5, origFile)
	fmt.Println("sha:", sha)
	shaStr, _ := SumByString(MODE_ENCODE_HEX_LOWER, crypto.MD5, origFile)
	fmt.Println("shaStr:", shaStr)

	sigData, _ := xHelper.SignPSSFileByString(MODE_ENCODE_HEX_LOWER, crypto.MD5, origFile, nil)
	fmt.Println(sigData)
	verifyDataErr := xHelper.VerifyPSSFileByString(MODE_ENCODE_HEX_LOWER, crypto.MD5, origFile, sigData, nil)
	if verifyDataErr == nil && len(sigData) > 0 {
		fmt.Println("文件签名验证通过")
	} else {
		fmt.Printf("文件签名验证不通过:%v", verifyDataErr)
	}
}

func TestAesCommon(t *testing.T) {

	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	var xHelper EncryptHelper
	xHelper = &XAes{
		ModeKey:     MODE_KEY_PEM,
		ModeEncode:  MODE_ENCODE_BASE64,
		ModePadding: MODE_PADDING_PKCS7,
	}
	keyStr, _ := xHelper.GenKeyIvString(24)
	ivStr, _ := xHelper.GenIVString()
	fmt.Println(keyStr)
	fmt.Println(ivStr)
	xHelper.SetKeyString(keyStr)
	xHelper.SetIVString(ivStr)
	//xHelper.SetBlockSize(16)
	origStr := "      中华人民共和国      " + generateToken(1024) + "      中华人民共和国      "
	//origStr = "as"
	//origStr = ""
	encStr, _ := xHelper.EncStringCBC(origStr)
	fmt.Println(encStr)
	decStr, _ := xHelper.DecStringCBC(encStr)
	fmt.Println(decStr)
	if origStr == decStr {
		fmt.Println("加密解密验证通过")
	} else {
		fmt.Println("加密解密验证不通过通过")
	}
}

func TestDesCommon(t *testing.T) {

	//time, _ := TimeParseByString(TIME_PATTERN_STANDARD, "2023-01-01 00:50:11")
	var xHelper EncryptHelper
	//xencrypt.
	xHelper = &XDes{
		ModeKey:     MODE_KEY_PEM,
		ModeEncode:  MODE_ENCODE_BASE64,
		ModePadding: MODE_PADDING_PKCS7,
	}
	//xHelper.SetAutoFillKey(true)
	keyStr, _ := xHelper.GenKeyIvString(24)
	ivStr, _ := xHelper.GenIVString()
	fmt.Println(keyStr)
	fmt.Println(ivStr)
	xHelper.SetKeyString(keyStr)
	xHelper.SetIVString(ivStr)
	//xHelper.SetBlockSize(16)
	origStr := "      中华人民共和国      " + generateToken(1024) + "      中华人民共和国      "
	//origStr = ""
	encStr, _ := xHelper.EncStringCBC(origStr)
	fmt.Println(encStr)
	decStr, _ := xHelper.DecStringCBC(encStr)
	fmt.Println(decStr)
	if origStr == decStr {
		fmt.Println("加密解密验证通过")
	} else {
		fmt.Println("加密解密验证不通过通过")
	}
}
