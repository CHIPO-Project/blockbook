// +build unittest

package chipo

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/martinboehm/btcutil/chaincfg"
)

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

// Test getting the address details from the address hash

func Test_GetAddrDescFromAddress_Mainnet(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{address: "CeULMQZRLEMxMZDuR5v6q1HC5VPRZKm3Uv"},
			want:    "76a914f1684a035088c20e76ece8e4dd79bdead0e1569a88ac",
			wantErr: false,
		},
	}
	parser := NewChipoParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_GetAddressesFromAddrDesc(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want2   bool
		wantErr bool
	}{
		{
			name:    "Normal",
			args:    args{script: "76a914f1684a035088c20e76ece8e4dd79bdead0e1569a88ac"},
			want:    []string{"CeULMQZRLEMxMZDuR5v6q1HC5VPRZKm3Uv"},
			want2:   true,
			wantErr: false,
		},
	}

	parser := NewChipoParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.script)
			got, got2, err := parser.GetAddressesFromAddrDesc(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressesFromAddrDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

// Test the packing and unpacking of raw transaction data

var (
	// Block Height 10
	testTx1       bchain.Tx
	testTxPacked1 = "0a208e29a4aaf8bfdcf2e485ac2d6494d85c3960d37f178d5826233ee601b647d14b126201000000010000000000000000000000000000000000000000000000000000000000000000ffffffff035a0103ffffffff0100e1f50500000000232102f6e1b93078a37d5f6b517cb19ada666302593cfeec40c6e24b78b8b8fc3ed625ac0000000018c596e3ee052000280a32100a06356130313033180028ffffffff0f3a510a0405f5e10010001a232102f6e1b93078a37d5f6b517cb19ada666302593cfeec40c6e24b78b8b8fc3ed625ac222243546a48564a47544850464e4536686145693168465075436b4271555443714c43464000"

	// Block Height 211
	testTx2       bchain.Tx
	testTxPacked2 = "0a20a5d735b4e57e6b05e48b488fdf3e671d337ce20191f65abff24e92a30ed594bc12b1020100000002829ac3f7566042134f6302761be74162b32c164b68b80bd0379df328041f088a000000004847304402200c8a342270e34087ea91ce2276e70c4065a12f403adebe827bc7748d7abcf42102203d91373002e860e37afd7b5374d38e1567d3c48cf7cb5300f3abc32d44ad645e01ffffffffc5dc7bcf612eb708d0a93a20e65211d6c7cd1e40590a4466564fd160ce55d1850000000049483045022100e68477ea1d9c8291e645080e9d31a2db58cee2fb236cf6ac53a70372b194ab9b022018d1c8c9b8098b46842de9882c20358691c3f154aa8275e4441fac8cacd14c3601ffffffff020cd5f505000000001976a914742c95653246b74561aafda87abe6930585d9cf688ac0010a5d4e80000001976a9146597e04894ac86c00eac131e24fc8f80f9e1844e88ac0000000018f4dae3ee05200028d30132760a0012208a081f0428f39d37d00bb8684b162cb36241e71b7602634f13426056f7c39a821800224847304402200c8a342270e34087ea91ce2276e70c4065a12f403adebe827bc7748d7abcf42102203d91373002e860e37afd7b5374d38e1567d3c48cf7cb5300f3abc32d44ad645e0128ffffffff0f32770a00122085d155ce60d14f5666440a59401ecdc7d61152e6203aa9d008b72e61cf7bdcc518002249483045022100e68477ea1d9c8291e645080e9d31a2db58cee2fb236cf6ac53a70372b194ab9b022018d1c8c9b8098b46842de9882c20358691c3f154aa8275e4441fac8cacd14c360128ffffffff0f3a470a0405f5d50c10001a1976a914742c95653246b74561aafda87abe6930585d9cf688ac2222435434414e316f6d467a5143616e72555a5a4d584c48484339346f43684665485a6d3a480a05e8d4a5100010011a1976a9146597e04894ac86c00eac131e24fc8f80f9e1844e88ac222243526a346b7845526f6d38447844385a4e6d4c4c6e4a6b34516761337239697676744000"
)

func init() {
    testTx1 = bchain.Tx{
        Hex: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff035a0103ffffffff0100e1f50500000000232102f6e1b93078a37d5f6b517cb19ada666302593cfeec40c6e24b78b8b8fc3ed625ac00000000",
         Txid: "8e29a4aaf8bfdcf2e485ac2d6494d85c3960d37f178d5826233ee601b647d14b",
         LockTime: 0,
         Vin: []bchain.Vin{
             {
                 Coinbase: "5a0103",
                 Sequence: 4294967295,
             },
         },
         Vout: []bchain.Vout{
            {
                ValueSat: *big.NewInt(100000000),
                N: 0,
                ScriptPubKey: bchain.ScriptPubKey{
                    Hex: "2102f6e1b93078a37d5f6b517cb19ada666302593cfeec40c6e24b78b8b8fc3ed625ac",
                    Addresses: []string{
                        "CTjHVJGTHPFNE6haEi1hFPuCkBqUTCqLCF",
                    },
                },
            },
        },
        Blocktime: 1574488901,
        Time: 1574488901,
    }

	testTx2 = bchain.Tx{
		Hex:      "0100000002829ac3f7566042134f6302761be74162b32c164b68b80bd0379df328041f088a000000004847304402200c8a342270e34087ea91ce2276e70c4065a12f403adebe827bc7748d7abcf42102203d91373002e860e37afd7b5374d38e1567d3c48cf7cb5300f3abc32d44ad645e01ffffffffc5dc7bcf612eb708d0a93a20e65211d6c7cd1e40590a4466564fd160ce55d1850000000049483045022100e68477ea1d9c8291e645080e9d31a2db58cee2fb236cf6ac53a70372b194ab9b022018d1c8c9b8098b46842de9882c20358691c3f154aa8275e4441fac8cacd14c3601ffffffff020cd5f505000000001976a914742c95653246b74561aafda87abe6930585d9cf688ac0010a5d4e80000001976a9146597e04894ac86c00eac131e24fc8f80f9e1844e88ac00000000",
		Txid:     "a5d735b4e57e6b05e48b488fdf3e671d337ce20191f65abff24e92a30ed594bc",
		LockTime: 0,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "47304402200c8a342270e34087ea91ce2276e70c4065a12f403adebe827bc7748d7abcf42102203d91373002e860e37afd7b5374d38e1567d3c48cf7cb5300f3abc32d44ad645e01",
				},
				Txid:     "8a081f0428f39d37d00bb8684b162cb36241e71b7602634f13426056f7c39a82",
				Vout:     0,
				Sequence: 4294967295,
			},
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "483045022100e68477ea1d9c8291e645080e9d31a2db58cee2fb236cf6ac53a70372b194ab9b022018d1c8c9b8098b46842de9882c20358691c3f154aa8275e4441fac8cacd14c3601",
				},
				Txid:     "85d155ce60d14f5666440a59401ecdc7d61152e6203aa9d008b72e61cf7bdcc5",
				Vout:     0,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(99996940),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a914742c95653246b74561aafda87abe6930585d9cf688ac",
					Addresses: []string{
						"CT4AN1omFzQCanrUZZMXLHHC94oChFeHZm",
					},
				},
			},
			{
				ValueSat: *big.NewInt(1000000000000),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a9146597e04894ac86c00eac131e24fc8f80f9e1844e88ac",
					Addresses: []string{
						"CRj4kxERom8DxD8ZNmLLnJk4Qga3r9ivvt",
					},
				},
			},
		},
		Blocktime: 1574497652,
		Time:      1574497652,
	}
}

func Test_PackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *ChipoParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "chipo_1",
			args: args{
				tx:        testTx1,
				height:    10,
				blockTime: 1574488901,
				parser:    NewChipoParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
		{
			name: "chipo_2",
			args: args{
				tx:        testTx2,
				height:    211,
				blockTime: 1574497652,
				parser:    NewChipoParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_UnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *ChipoParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "chipo_1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewChipoParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx1,
			want1:   10,
			wantErr: false,
		},
		{
			name: "chipo_2",
			args: args{
				packedTx: testTxPacked2,
				parser:   NewChipoParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx2,
			want1:   211,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

// Block test - looks for size, time, and transaction hashes

type testBlock struct {
	size int
	time int64
	tx   []string
}

var testParseBlockTxs = map[int]testBlock{
	10: {
		size: 179,
		time: 1574488901,
		tx: []string{
			"8e29a4aaf8bfdcf2e485ac2d6494d85c3960d37f178d5826233ee601b647d14b",
		},
	},
	211: {
		size: 1705,
		time: 1574497652,
		tx: []string{
            "d8ae3828e85db67d167a1bdbc90430eb6bbb3779dbb7e1deef531ad75d00d816",
            "a5d735b4e57e6b05e48b488fdf3e671d337ce20191f65abff24e92a30ed594bc",
            "a272455f6b98fd8ac791d3e0b0dae05b75c6d4e9cb3974460f780cd5386c3403",
            "6527145741c381244c0004f2b5f03b38d0831fdc3f55fd2525e6e08b1cd0110d",
            "c7910fae75843870f219044d1224f59ba152630a78cd7d8672e765c4e03fc3d1",
            "82488b33ee38a06013c396256af6a2555b7d5535b50382a0a6b4834a7679fc5b",
		},
	},
}

func helperLoadBlock(t *testing.T, height int) []byte {
	name := fmt.Sprintf("block_dump.%d", height)
	path := filepath.Join("testdata", name)

	d, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	d = bytes.TrimSpace(d)

	b := make([]byte, hex.DecodedLen(len(d)))
	_, err = hex.Decode(b, d)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func TestParseBlock(t *testing.T) {
	p := NewChipoParser(GetChainParams("main"), &btc.Configuration{})

	for height, tb := range testParseBlockTxs {
		b := helperLoadBlock(t, height)

		blk, err := p.ParseBlock(b)
		if err != nil {
			t.Fatal(err)
		}

		if blk.Size != tb.size {
			t.Errorf("ParseBlock() block size: got %d, want %d", blk.Size, tb.size)
		}

		if blk.Time != tb.time {
			t.Errorf("ParseBlock() block time: got %d, want %d", blk.Time, tb.time)
		}

		if len(blk.Txs) != len(tb.tx) {
			t.Errorf("ParseBlock() number of transactions: got %d, want %d", len(blk.Txs), len(tb.tx))
		}

		for ti, tx := range tb.tx {
			if blk.Txs[ti].Txid != tx {
				t.Errorf("ParseBlock() transaction %d: got %s, want %s", ti, blk.Txs[ti].Txid, tx)
			}
		}
	}
}
