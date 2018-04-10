mkdir -p testdata/v4/user/checkins/
curl "https://api.untappd.com/v4/beer/info/2407538?client_id=$1&client_secret=$2&compact=true" > testdata/v4/beer/info/2407538_client_id=testid_client_secret=testsecret_compact=true
exit
curl "https://api.untappd.com/v4/user/checkins/blahdybblah?client_id=$1&client_secret=$2&min_id=0" > testdata/v4/user/checkins/blahdyblah_client_id=testid_client_secret=testsecret_min_id=0
exit
curl "https://api.untappd.com/v4/user/checkins/brotherlogic?client_id=$1&client_secret=$2&min_id=579454692" > testdata/v4/user/checkins/brotherlogic_client_id=testid_client_secret=testsecret_min_id=579454692
exit
curl "https://api.untappd.com/v4/user/checkins/brotherlogic?client_id=$1&client_secret=$2&max_id=8471" > testdata/v4/user/checkins/brotherlogic_client_id=testid_client_secret=testsecret_max_id=8471
curl "https://api.untappd.com/v4/user/checkins/brotherlogic?client_id=$1&client_secret=$2&min_id=0" > testdata/v4/user/checkins/brotherlogic_client_id=testid_client_secret=testsecret_min_id=0
exit
curl "https://api.untappd.com/v4/beer/info/7936?access_token=$1" > testdata/v4/beer/info/7936_access_token=accesstoken
exit
curl "https://api.untappd.com/v4/beer/info/7936?client_id=$1&client_secret=$2&compact=true" > testdata/v4/beer/info/7936_client_id=testid_client_secret=testsecret_compact=true
exit
curl "https://api.untappd.com/v4/venue/info/2194560?client_id=$1&client_secret=$2" > testdata/v4/venue/info/2194560_client_id=testid_client_secret=testsecret
