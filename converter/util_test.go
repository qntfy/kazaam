package converter

func getTestData() []byte {
	data := []byte(`
{
  "tests": {
    "test_int": 500,
    "test_float": 500.01,
    "test_float2": 500.0,
    "test_fraction": 0.5,
    "test_trim": "    blah   ",
    "test_money": "$6,000,000",
    "test_chars": "abcdefghijklmnopqrstuvwxyz",
	"test_mapped": "Texas",
	"test_null": null
  },
  "test_bool": true
}
`)
	return data
}
