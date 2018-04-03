curl "https://api.untappd.com/v4/beer/info/7936?access_token=$1" > testdata/v4/beer/info/7936_access_token=accesstoken
exit
curl "https://api.untappd.com/v4/beer/info/7936?client_id=$1&client_secret=$2&compact=true" > testdata/v4/beer/info/7936_client_id=testid_client_secret=testsecret_compact=true
exit
curl "https://api.untappd.com/v4/venue/info/2194560?client_id=$1&client_secret=$2" > testdata/v4/venue/info/2194560_client_id=testid_client_secret=testsecret
