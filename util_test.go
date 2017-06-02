package kazaam

import (
	"testing"
)

type IsJsonTest struct {
	json    string
	valid   bool
	isArray bool // encoding/json IsJson() does not support arrays
}

var jsonTests = []IsJsonTest{
	{
		json:  `{"key":"value"}`,
		valid: true,
	},
	{
		json:  `{"bad-data"}`,
		valid: false,
	},
	{
		json:  `{"text": "RT @PostGradProblem: In preparation for the NFL lockout, I will be spending twice as much time analyzing my fantasy baseball team during ...","truncated": true,"in_reply_to_user_id": null,"in_reply_to_status_id": null,"favorited": false,"source": "","in_reply_to_screen_name": null,"in_reply_to_status_id_str": null,"id_str": "54691802283900928","entities": {"user_mentions": [{"indices": [3,19],"screen_name": "PostGradProblem","id_str": "271572434","name": "PostGradProblems","id": 271572434}],"urls": [ ],"hashtags": [ ]},"contributors": null,"retweeted": false,"in_reply_to_user_id_str": null,"place": null,"retweet_count": 4,"created_at": "Sun Apr 03 23:48:36 +0000 2011","retweeted_status": {"text": "In preparation for the NFL lockout, I will be spending twice as much time analyzing my fantasy baseball team during company time. #PGP","truncated": false,"in_reply_to_user_id": null,"in_reply_to_status_id": null,"favorited": false,"source": "","in_reply_to_screen_name": null,"in_reply_to_status_id_str": null,"id_str": "54640519019642881","entities": {"user_mentions": [ ],"urls": [ ],"hashtags": [{"text": "PGP","indices": [130,134]}]},"contributors": null,"retweeted": false,"in_reply_to_user_id_str": null,"place": null,"retweet_count": 4,"created_at": "Sun Apr 03 20:24:49 +0000 2011","user": {"notifications": null,"profile_use_background_image": true,"statuses_count": 31,"profile_background_color": "C0DEED","followers_count": 3066,"profile_image_url": "http://a2.twimg.com/profile_images/1285770264/PGP_normal.jpg","listed_count": 6,"profile_background_image_url": "http://a3.twimg.com/a/1301071706/images/themes/theme1/bg.png","description": "","screen_name": "PostGradProblem","default_profile": true,"verified": false,"time_zone": null,"profile_text_color": "333333","is_translator": false,"profile_sidebar_fill_color": "DDEEF6","location": "","id_str": "271572434","default_profile_image": false,"profile_background_tile": false,"lang": "en","friends_count": 21,"protected": false,"favourites_count": 0,"created_at": "Thu Mar 24 19:45:44 +0000 2011","profile_link_color": "0084B4","name": "PostGradProblems","show_all_inline_media": false,"follow_request_sent": null,"geo_enabled": false,"profile_sidebar_border_color": "C0DEED","url": null,"id": 271572434,"contributors_enabled": false,"following": null,"utc_offset": null},"id": 54640519019642880,"coordinates": null,"geo": null},"user": {"notifications": null,"profile_use_background_image": true,"statuses_count": 351,"profile_background_color": "C0DEED","followers_count": 48,"profile_image_url": "http://a1.twimg.com/profile_images/455128973/gCsVUnofNqqyd6tdOGevROvko1_500_normal.jpg","listed_count": 0,"profile_background_image_url": "http://a3.twimg.com/a/1300479984/images/themes/theme1/bg.png","description": "watcha doin in my waters?","screen_name": "OldGREG85","default_profile": true,"verified": false,"time_zone": "Hawaii","profile_text_color": "333333","is_translator": false,"profile_sidebar_fill_color": "DDEEF6","location": "Texas","id_str": "80177619","default_profile_image": false,"profile_background_tile": false,"lang": "en","friends_count": 81,"protected": false,"favourites_count": 0,"created_at": "Tue Oct 06 01:13:17 +0000 2009","profile_link_color": "0084B4","name": "GG","show_all_inline_media": false,"follow_request_sent": null,"geo_enabled": false,"profile_sidebar_border_color": "C0DEED","url": null,"id": 80177619,"contributors_enabled": false,"following": null,"utc_offset": -36000},"id": 54691802283900930,"coordinates": null,"geo": null}`,
		valid: true,
	},
	{
		// ,"profile_text_color": "333333","is_translator: false,"profile_sidebar_fill_color": "DDEEF6",
		json:  `{"text": "RT @PostGradProblem: In preparation for the NFL lockout, I will be spending twice as much time analyzing my fantasy baseball team during ...","truncated": true,"in_reply_to_user_id": null,"in_reply_to_status_id": null,"favorited": false,"source": "","in_reply_to_screen_name": null,"in_reply_to_status_id_str": null,"id_str": "54691802283900928","entities": {"user_mentions": [{"indices": [3,19],"screen_name": "PostGradProblem","id_str": "271572434","name": "PostGradProblems","id": 271572434}],"urls": [ ],"hashtags": [ ]},"contributors": null,"retweeted": false,"in_reply_to_user_id_str": null,"place": null,"retweet_count": 4,"created_at": "Sun Apr 03 23:48:36 +0000 2011","retweeted_status": {"text": "In preparation for the NFL lockout, I will be spending twice as much time analyzing my fantasy baseball team during company time. #PGP","truncated": false,"in_reply_to_user_id": null,"in_reply_to_status_id": null,"favorited": false,"source": "","in_reply_to_screen_name": null,"in_reply_to_status_id_str": null,"id_str": "54640519019642881","entities": {"user_mentions": [ ],"urls": [ ],"hashtags": [{"text": "PGP","indices": [130,134]}]},"contributors": null,"retweeted": false,"in_reply_to_user_id_str": null,"place": null,"retweet_count": 4,"created_at": "Sun Apr 03 20:24:49 +0000 2011","user": {"notifications": null,"profile_use_background_image": true,"statuses_count": 31,"profile_background_color": "C0DEED","followers_count": 3066,"profile_image_url": "http://a2.twimg.com/profile_images/1285770264/PGP_normal.jpg","listed_count": 6,"profile_background_image_url": "http://a3.twimg.com/a/1301071706/images/themes/theme1/bg.png","description": "","screen_name": "PostGradProblem","default_profile": true,"verified": false,"time_zone": null,"profile_text_color": "333333","is_translator": false,"profile_sidebar_fill_color": "DDEEF6","location": "","id_str": "271572434","default_profile_image": false,"profile_background_tile": false,"lang": "en","friends_count": 21,"protected": false,"favourites_count": 0,"created_at": "Thu Mar 24 19:45:44 +0000 2011","profile_link_color": "0084B4","name": "PostGradProblems","show_all_inline_media": false,"follow_request_sent": null,"geo_enabled": false,"profile_sidebar_border_color": "C0DEED","url": null,"id": 271572434,"contributors_enabled": false,"following": null,"utc_offset": null},"id": 54640519019642880,"coordinates": null,"geo": null},"user": {"notifications": null,"profile_use_background_image": true,"statuses_count": 351,"profile_background_color": "C0DEED","followers_count": 48,"profile_image_url": "http://a1.twimg.com/profile_images/455128973/gCsVUnofNqqyd6tdOGevROvko1_500_normal.jpg","listed_count": 0,"profile_background_image_url": "http://a3.twimg.com/a/1300479984/images/themes/theme1/bg.png","description": "watcha doin in my waters?","screen_name": "OldGREG85","default_profile": true,"verified": false,"time_zone": "Hawaii","profile_text_color": "333333","is_translator: false,"profile_sidebar_fill_color": "DDEEF6","location": "Texas","id_str": "80177619","default_profile_image": false,"profile_background_tile": false,"lang": "en","friends_count": 81,"protected": false,"favourites_count": 0,"created_at": "Tue Oct 06 01:13:17 +0000 2009","profile_link_color": "0084B4","name": "GG","show_all_inline_media": false,"follow_request_sent": null,"geo_enabled": false,"profile_sidebar_border_color": "C0DEED","url": null,"id": 80177619,"contributors_enabled": false,"following": null,"utc_offset": -36000},"id": 54691802283900930,"coordinates": null,"geo": null}`,
		valid: false,
	},
	{
		// "hashtags": [{"text": "PGP","indices": [130,134]}},
		json:  `{"text": "RT @PostGradProblem: In preparation for the NFL lockout, I will be spending twice as much time analyzing my fantasy baseball team during ...","truncated": true,"in_reply_to_user_id": null,"in_reply_to_status_id": null,"favorited": false,"source": "","in_reply_to_screen_name": null,"in_reply_to_status_id_str": null,"id_str": "54691802283900928","entities": {"user_mentions": [{"indices": [3,19],"screen_name": "PostGradProblem","id_str": "271572434","name": "PostGradProblems","id": 271572434}],"urls": [ ],"hashtags": [ ]},"contributors": null,"retweeted": false,"in_reply_to_user_id_str": null,"place": null,"retweet_count": 4,"created_at": "Sun Apr 03 23:48:36 +0000 2011","retweeted_status": {"text": "In preparation for the NFL lockout, I will be spending twice as much time analyzing my fantasy baseball team during company time. #PGP","truncated": false,"in_reply_to_user_id": null,"in_reply_to_status_id": null,"favorited": false,"source": "","in_reply_to_screen_name": null,"in_reply_to_status_id_str": null,"id_str": "54640519019642881","entities": {"user_mentions": [ ],"urls": [ ],"hashtags": [{"text": "PGP","indices": [130,134]}},"contributors": null,"retweeted": false,"in_reply_to_user_id_str": null,"place": null,"retweet_count": 4,"created_at": "Sun Apr 03 20:24:49 +0000 2011","user": {"notifications": null,"profile_use_background_image": true,"statuses_count": 31,"profile_background_color": "C0DEED","followers_count": 3066,"profile_image_url": "http://a2.twimg.com/profile_images/1285770264/PGP_normal.jpg","listed_count": 6,"profile_background_image_url": "http://a3.twimg.com/a/1301071706/images/themes/theme1/bg.png","description": "","screen_name": "PostGradProblem","default_profile": true,"verified": false,"time_zone": null,"profile_text_color": "333333","is_translator": false,"profile_sidebar_fill_color": "DDEEF6","location": "","id_str": "271572434","default_profile_image": false,"profile_background_tile": false,"lang": "en","friends_count": 21,"protected": false,"favourites_count": 0,"created_at": "Thu Mar 24 19:45:44 +0000 2011","profile_link_color": "0084B4","name": "PostGradProblems","show_all_inline_media": false,"follow_request_sent": null,"geo_enabled": false,"profile_sidebar_border_color": "C0DEED","url": null,"id": 271572434,"contributors_enabled": false,"following": null,"utc_offset": null},"id": 54640519019642880,"coordinates": null,"geo": null},"user": {"notifications": null,"profile_use_background_image": true,"statuses_count": 351,"profile_background_color": "C0DEED","followers_count": 48,"profile_image_url": "http://a1.twimg.com/profile_images/455128973/gCsVUnofNqqyd6tdOGevROvko1_500_normal.jpg","listed_count": 0,"profile_background_image_url": "http://a3.twimg.com/a/1300479984/images/themes/theme1/bg.png","description": "watcha doin in my waters?","screen_name": "OldGREG85","default_profile": true,"verified": false,"time_zone": "Hawaii","profile_text_color": "333333","is_translator": false,"profile_sidebar_fill_color": "DDEEF6","location": "Texas","id_str": "80177619","default_profile_image": false,"profile_background_tile": false,"lang": "en","friends_count": 81,"protected": false,"favourites_count": 0,"created_at": "Tue Oct 06 01:13:17 +0000 2009","profile_link_color": "0084B4","name": "GG","show_all_inline_media": false,"follow_request_sent": null,"geo_enabled": false,"profile_sidebar_border_color": "C0DEED","url": null,"id": 80177619,"contributors_enabled": false,"following": null,"utc_offset": -36000},"id": 54691802283900930,"coordinates": null,"geo": null}`,
		valid: false,
	},
	{
		json:  `{"bad-key":falre}`,
		valid: false,
	},
	{
		json:  `{"boolean":true}`,
		valid: true,
	},
	{
		json:  `{"nums":773,"type":"integer"}`,
		valid: true,
	},
	{
		json:  `{"empty":null, "foo":["bar", "baz", 83, null]}`,
		valid: true,
	},
	{
		json:  `this is in no way json`,
		valid: false,
	},
	{
		json:    `["foo", false, [1, 2, 3], {"key":"value", "bool":false, "empty":null}, "bar"]`,
		valid:   true,
		isArray: true,
	},
	{
		json:    `["simple", "array"]`,
		valid:   true,
		isArray: true,
	},
	{
		json:    `["invalid", "array]`,
		valid:   false,
		isArray: true,
	},
	{
		json:    `[33, false, 77, random, "array"]`,
		valid:   false,
		isArray: true,
	},
	{
		json:  `{}`,
		valid: true,
	},
	{
		json:    `[]`,
		valid:   true,
		isArray: true,
	},
	{
		json:  ``,
		valid: false,
	},
	{
		json:  `"just a string"`,
		valid: false,
	},
}

func TestIsJson(t *testing.T) {
	for _, jt := range jsonTests {
		if !jt.isArray && IsJson([]byte(jt.json)) != jt.valid {
			t.Error("IsJson() failed")
			t.Log("Json: ", jt.json)
			t.Log("Expected: ", jt.valid)
		}
		if IsJsonFast([]byte(jt.json)) != jt.valid {
			t.Error("IsJsonFast() failed")
			t.Log("Json: ", jt.json)
			t.Log("Expected: ", jt.valid)
		}
	}
}
